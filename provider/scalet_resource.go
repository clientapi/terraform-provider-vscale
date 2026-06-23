package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ScaletResource{}
var _ resource.ResourceWithImportState = &ScaletResource{}

type ScaletResource struct {
	client *client.Client
}

type AddressModel struct {
	Netmask types.String `tfsdk:"netmask"`
	Gateway types.String `tfsdk:"gateway"`
	Address types.String `tfsdk:"address"`
}

type ScaletResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	MakeFrom       types.String `tfsdk:"make_from"`
	RPlan          types.String `tfsdk:"rplan"`
	Location       types.String `tfsdk:"location"`
	DoStart        types.Bool   `tfsdk:"do_start"`
	Password       types.String `tfsdk:"password"`
	Keys           types.List   `tfsdk:"keys"`
	Status         types.String `tfsdk:"status"`
	PublicAddress  types.Object `tfsdk:"public_address"`
	PrivateAddress types.Object `tfsdk:"private_address"`
}

var addressAttrTypes = map[string]attr.Type{
	"netmask": types.StringType,
	"gateway": types.StringType,
	"address": types.StringType,
}

func mapAddressToObject(ctx context.Context, addr *client.Address) (types.Object, diag.Diagnostics) {
	if addr == nil {
		return types.ObjectNull(addressAttrTypes), nil
	}

	var netmask, gateway, address types.String
	if addr.Netmask != "" {
		netmask = types.StringValue(addr.Netmask)
	} else {
		netmask = types.StringNull()
	}

	if addr.Gateway != "" {
		gateway = types.StringValue(addr.Gateway)
	} else {
		gateway = types.StringNull()
	}

	if addr.Address != "" {
		address = types.StringValue(addr.Address)
	} else {
		address = types.StringNull()
	}

	model := AddressModel{
		Netmask: netmask,
		Gateway: gateway,
		Address: address,
	}
	return types.ObjectValueFrom(ctx, addressAttrTypes, model)
}

func NewScaletResource() resource.Resource {
	return &ScaletResource{}
}

func (r *ScaletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scalet"
}

func (r *ScaletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Scalet (Virtual Private Server) in VScale.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier (CTID) of the Scalet.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the Scalet.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"make_from": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the OS image or backup from which to build the Scalet. Changing this triggers rebuild.",
			},
			"rplan": schema.StringAttribute{
				Required:    true,
				Description: "Tariff plan ID (e.g. medium, large). Changing this triggers plan upgrade.",
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "Datacenter location ID (e.g. spb0).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"do_start": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to start the Scalet immediately after creation.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Root password for the Scalet (if keys are not used). Changing this triggers rebuild.",
			},
			"keys": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "List of SSH key IDs to load onto the Scalet. Can be updated dynamically.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the Scalet (e.g., started, stopped, defined).",
			},
			"public_address": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Public network interface details.",
				Attributes: map[string]schema.Attribute{
					"netmask": schema.StringAttribute{
						Computed: true,
					},
					"gateway": schema.StringAttribute{
						Computed: true,
					},
					"address": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"private_address": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Private network interface details.",
				Attributes: map[string]schema.Attribute{
					"netmask": schema.StringAttribute{
						Computed: true,
					},
					"gateway": schema.StringAttribute{
						Computed: true,
					},
					"address": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *ScaletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ScaletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ScaletResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var keys []int
	if !data.Keys.IsNull() && !data.Keys.IsUnknown() {
		var keyVals []int64
		diags := data.Keys.ElementsAs(ctx, &keyVals, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, k := range keyVals {
			keys = append(keys, int(k))
		}
	}

	createReq := &client.CreateScaletRequest{
		MakeFrom: data.MakeFrom.ValueString(),
		RPlan:    data.RPlan.ValueString(),
		DoStart:  data.DoStart.ValueBool(),
		Name:     data.Name.ValueString(),
		Password: data.Password.ValueString(),
		Location: data.Location.ValueString(),
		Keys:     keys,
	}

	scalet, err := r.client.CreateScalet(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Scalet", err.Error())
		return
	}

	data.ID = types.StringValue(strconv.Itoa(scalet.CTID))
	data.Status = types.StringValue(scalet.Status)

	pubObj, diags := mapAddressToObject(ctx, scalet.PublicAddress)
	resp.Diagnostics.Append(diags...)
	data.PublicAddress = pubObj

	privObj, diags := mapAddressToObject(ctx, scalet.PrivateAddress)
	resp.Diagnostics.Append(diags...)
	data.PrivateAddress = privObj

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScaletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ScaletResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctid, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Scalet ID", err.Error())
		return
	}

	scalet, err := r.client.GetScalet(ctid)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Name = types.StringValue(scalet.Name)
	data.MakeFrom = types.StringValue(scalet.MakeFrom)
	data.RPlan = types.StringValue(scalet.RPlan)
	data.Location = types.StringValue(scalet.Location)
	data.Status = types.StringValue(scalet.Status)

	pubObj, diags := mapAddressToObject(ctx, scalet.PublicAddress)
	resp.Diagnostics.Append(diags...)
	data.PublicAddress = pubObj

	privObj, diags := mapAddressToObject(ctx, scalet.PrivateAddress)
	resp.Diagnostics.Append(diags...)
	data.PrivateAddress = privObj

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScaletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ScaletResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctid, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Scalet ID", err.Error())
		return
	}

	// 1. Check if RPlan changed (upgrade plan)
	if !plan.RPlan.Equal(state.RPlan) {
		_, err := r.client.UpgradeScalet(ctid, plan.RPlan.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error Upgrading Scalet Plan", err.Error())
			return
		}
	}

	// 2. Check if MakeFrom or Password changed (rebuild server)
	makeFromChanged := !plan.MakeFrom.Equal(state.MakeFrom) && !state.MakeFrom.IsNull() && !state.MakeFrom.IsUnknown() && state.MakeFrom.ValueString() != ""
	passwordChanged := !plan.Password.Equal(state.Password) && !state.Password.IsNull() && !state.Password.IsUnknown() && state.Password.ValueString() != ""

	if makeFromChanged || passwordChanged {
		pass := plan.Password.ValueString()
		if pass == "" {
			pass = "RebuiltPwd123!"
		}
		_, err := r.client.RebuildScalet(ctid, pass)
		if err != nil {
			resp.Diagnostics.AddError("Error Rebuilding Scalet", err.Error())
			return
		}
	}

	// 3. Check if Keys changed (update SSH keys)
	if !plan.Keys.Equal(state.Keys) {
		var keys []int
		if !plan.Keys.IsNull() && !plan.Keys.IsUnknown() {
			var keyVals []int64
			diags := plan.Keys.ElementsAs(ctx, &keyVals, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			for _, k := range keyVals {
				keys = append(keys, int(k))
			}
		}

		_, err := r.client.UpdateScaletSSHKeys(ctid, keys)
		if err != nil {
			resp.Diagnostics.AddError("Error Updating Scalet SSH Keys", err.Error())
			return
		}
	}

	// Fetch fresh state from API
	scalet, err := r.client.GetScalet(ctid)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Scalet State after Update", err.Error())
		return
	}

	plan.Status = types.StringValue(scalet.Status)

	pubObj, diags := mapAddressToObject(ctx, scalet.PublicAddress)
	resp.Diagnostics.Append(diags...)
	plan.PublicAddress = pubObj

	privObj, diags := mapAddressToObject(ctx, scalet.PrivateAddress)
	resp.Diagnostics.Append(diags...)
	plan.PrivateAddress = privObj

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ScaletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ScaletResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctid, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Scalet ID", err.Error())
		return
	}

	err = r.client.DeleteScalet(ctid)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Scalet", err.Error())
		return
	}
}

func (r *ScaletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
