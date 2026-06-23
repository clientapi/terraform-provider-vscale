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

var _ resource.Resource = &PTRRecordResource{}
var _ resource.ResourceWithImportState = &PTRRecordResource{}

type PTRRecordResource struct {
	client *client.Client
}

type PTRRecordResourceModel struct {
	ID      types.String `tfsdk:"id"`
	IP      types.String `tfsdk:"ip"`
	Content types.String `tfsdk:"content"`
}

func NewPTRRecordResource() resource.Resource {
	return &PTRRecordResource{}
}

func (r *PTRRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ptr_record"
}

func (r *PTRRecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DNS PTR (Reverse) record in VScale.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the PTR record.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip": schema.StringAttribute{
				Required:    true,
				Description: "IP address for which the PTR record will be created.",
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "Domain name (value of the reverse record).",
			},
		},
	}
}

func (r *PTRRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PTRRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PTRRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	record, err := r.client.CreatePTRRecord(data.IP.ValueString(), data.Content.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Creating PTR Record", err.Error())
		return
	}

	data.ID = types.StringValue(strconv.Itoa(record.ID))
	data.IP = types.StringValue(record.IP)
	data.Content = types.StringValue(record.Content)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PTRRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PTRRecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing PTR ID", err.Error())
		return
	}

	record, err := r.client.GetPTRRecord(id)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.IP = types.StringValue(record.IP)
	data.Content = types.StringValue(record.Content)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PTRRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PTRRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing PTR ID", err.Error())
		return
	}

	_, err = r.client.UpdatePTRRecord(id, data.IP.ValueString(), data.Content.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Updating PTR Record", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PTRRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PTRRecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing PTR ID", err.Error())
		return
	}

	err = r.client.DeletePTRRecord(id)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting PTR Record", err.Error())
		return
	}
}

func (r *PTRRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
