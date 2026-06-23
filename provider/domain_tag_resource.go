package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &DomainTagResource{}
var _ resource.ResourceWithImportState = &DomainTagResource{}

type DomainTagResource struct {
	client *client.Client
}

type DomainTagResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Domains types.List   `tfsdk:"domains"`
}

func NewDomainTagResource() resource.Resource {
	return &DomainTagResource{}
}

func (r *DomainTagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_tag"
}

func (r *DomainTagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DNS domain tag in VScale.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the domain tag.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the tag.",
			},
			"domains": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Domains attached to this tag.",
			},
		},
	}
}

func (r *DomainTagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DomainTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DomainTagResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var domains []string
	if !data.Domains.IsNull() && !data.Domains.IsUnknown() {
		resp.Diagnostics.Append(data.Domains.ElementsAs(ctx, &domains, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tag, err := r.client.CreateDomainTag(data.Name.ValueString(), domains)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Domain Tag", err.Error())
		return
	}

	data.ID = types.StringValue(strconv.Itoa(tag.ID))
	data.Name = types.StringValue(tag.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DomainTagResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Tag ID", err.Error())
		return
	}

	tag, err := r.client.GetDomainTag(id)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Name = types.StringValue(tag.Name)

	domainsList, diags := types.ListValueFrom(ctx, types.StringType, tag.Domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Domains = domainsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DomainTagResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Tag ID", err.Error())
		return
	}

	var domains []string
	if !data.Domains.IsNull() && !data.Domains.IsUnknown() {
		resp.Diagnostics.Append(data.Domains.ElementsAs(ctx, &domains, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	_, err = r.client.UpdateDomainTag(id, data.Name.ValueString(), domains)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Domain Tag", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DomainTagResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Tag ID", err.Error())
		return
	}

	err = r.client.DeleteDomainTag(id)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Domain Tag", err.Error())
		return
	}
}

func (r *DomainTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
