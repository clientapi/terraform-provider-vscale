package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &DomainRecordResource{}
var _ resource.ResourceWithImportState = &DomainRecordResource{}

type DomainRecordResource struct {
	client *client.Client
}

type DomainRecordResourceModel struct {
	ID       types.String `tfsdk:"id"`
	DomainID types.Int64  `tfsdk:"domain_id"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	TTL      types.Int64  `tfsdk:"ttl"`
	Content  types.String `tfsdk:"content"`
	Priority types.Int64  `tfsdk:"priority"`
	Weight   types.Int64  `tfsdk:"weight"`
	Port     types.Int64  `tfsdk:"port"`
	Target   types.String `tfsdk:"target"`
}

func NewDomainRecordResource() resource.Resource {
	return &DomainRecordResource{}
}

func (r *DomainRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_record"
}

func (r *DomainRecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DNS record for a Domain in VScale.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the domain record.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the domain to attach this record to.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the record (e.g., www, mail, or @ for root).",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Record type: SOA, NS, A, AAAA, CNAME, SRV, MX, TXT, SPF.",
			},
			"ttl": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(300),
				Description: "Time to Live (TTL) in seconds. Min 60, Max 604800.",
			},
			"content": schema.StringAttribute{
				Optional:    true,
				Description: "Value of the record (e.g., an IP address). Omitted for SRV records.",
			},
			"priority": schema.Int64Attribute{
				Optional:    true,
				Description: "Priority for MX and SRV records.",
			},
			"weight": schema.Int64Attribute{
				Optional:    true,
				Description: "Weight for SRV records.",
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Description: "Port for SRV records.",
			},
			"target": schema.StringAttribute{
				Optional:    true,
				Description: "Target hostname for SRV records.",
			},
		},
	}
}

func (r *DomainRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DomainRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DomainRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rec := &client.DomainRecord{
		Name:    data.Name.ValueString(),
		Type:    data.Type.ValueString(),
		TTL:     int(data.TTL.ValueInt64()),
		Content: data.Content.ValueString(),
	}

	if !data.Priority.IsNull() {
		val := int(data.Priority.ValueInt64())
		rec.Priority = &val
	}
	if !data.Weight.IsNull() {
		val := int(data.Weight.ValueInt64())
		rec.Weight = &val
	}
	if !data.Port.IsNull() {
		val := int(data.Port.ValueInt64())
		rec.Port = &val
	}
	if !data.Target.IsNull() {
		rec.Target = data.Target.ValueString()
	}

	created, err := r.client.CreateDomainRecord(int(data.DomainID.ValueInt64()), rec)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Domain Record", err.Error())
		return
	}

	data.ID = types.StringValue(strconv.Itoa(created.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DomainRecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	recordID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Record ID", err.Error())
		return
	}

	record, err := r.client.GetDomainRecord(int(data.DomainID.ValueInt64()), recordID)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Name = types.StringValue(record.Name)
	data.Type = types.StringValue(record.Type)
	data.TTL = types.Int64Value(int64(record.TTL))
	data.Content = types.StringValue(record.Content)

	if record.Priority != nil {
		data.Priority = types.Int64Value(int64(*record.Priority))
	} else {
		data.Priority = types.Int64Null()
	}
	if record.Weight != nil {
		data.Weight = types.Int64Value(int64(*record.Weight))
	} else {
		data.Weight = types.Int64Null()
	}
	if record.Port != nil {
		data.Port = types.Int64Value(int64(*record.Port))
	} else {
		data.Port = types.Int64Null()
	}
	if record.Target != "" {
		data.Target = types.StringValue(record.Target)
	} else {
		data.Target = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DomainRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	recordID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Record ID", err.Error())
		return
	}

	rec := &client.DomainRecord{
		Name:    data.Name.ValueString(),
		Type:    data.Type.ValueString(),
		TTL:     int(data.TTL.ValueInt64()),
		Content: data.Content.ValueString(),
	}

	if !data.Priority.IsNull() {
		val := int(data.Priority.ValueInt64())
		rec.Priority = &val
	}
	if !data.Weight.IsNull() {
		val := int(data.Weight.ValueInt64())
		rec.Weight = &val
	}
	if !data.Port.IsNull() {
		val := int(data.Port.ValueInt64())
		rec.Port = &val
	}
	if !data.Target.IsNull() {
		rec.Target = data.Target.ValueString()
	}

	_, err = r.client.UpdateDomainRecord(int(data.DomainID.ValueInt64()), recordID, rec)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Domain Record", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DomainRecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	recordID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Record ID", err.Error())
		return
	}

	err = r.client.DeleteDomainRecord(int(data.DomainID.ValueInt64()), recordID)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Domain Record", err.Error())
		return
	}
}

func (r *DomainRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import state needs both domain_id and record_id. Standard format: "domain_id,record_id"
	// ImportStatePassIDToState will put the raw ID into id. We will split it in Read or handle custom import.
	// Let's implement custom ImportState to split it.
	id := req.ID
	// Or split by comma
	parts := importSplit(id)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format 'domain_id,record_id'.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func importSplit(s string) []string {
	var result []string
	curr := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			result = append(result, curr)
			curr = ""
		} else {
			curr += string(s[i])
		}
	}
	result = append(result, curr)
	return result
}
