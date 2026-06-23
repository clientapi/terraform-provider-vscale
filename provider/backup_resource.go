package provider

import (
	"context"
	"fmt"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &BackupResource{}
var _ resource.ResourceWithImportState = &BackupResource{}

type BackupResource struct {
	client *client.Client
}

type BackupResourceModel struct {
	ID       types.String `tfsdk:"id"`
	ScaletID types.Int64  `tfsdk:"scalet_id"`
	Name     types.String `tfsdk:"name"`
	Template types.String `tfsdk:"template"`
	Size     types.Int64  `tfsdk:"size"`
	Location types.String `tfsdk:"location"`
	Created  types.String `tfsdk:"created"`
	Status   types.String `tfsdk:"status"`
}

func NewBackupResource() resource.Resource {
	return &BackupResource{}
}

func (r *BackupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup"
}

func (r *BackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a server backup in VScale. Creating a backup triggers a backup of the source Scalet.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the backup.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scalet_id": schema.Int64Attribute{
				Required:    true,
				Description: "ID of the source Scalet.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the backup.",
			},
			"template": schema.StringAttribute{
				Computed:    true,
				Description: "OS template from which server was built.",
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "Backup size in GB.",
			},
			"location": schema.StringAttribute{
				Computed:    true,
				Description: "Datacenter location where backup is stored.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the backup.",
			},
		},
	}
}

func (r *BackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BackupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backup, err := r.client.CreateBackup(int(data.ScaletID.ValueInt64()), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Backup", err.Error())
		return
	}

	data.ID = types.StringValue(backup.ID)
	data.Template = types.StringValue(backup.Template)
	data.Size = types.Int64Value(int64(backup.Size))
	data.Location = types.StringValue(backup.Location)
	data.Created = types.StringValue(backup.Created)
	data.Status = types.StringValue(backup.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BackupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backup, err := r.client.GetBackup(data.ID.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Template = types.StringValue(backup.Template)
	data.Size = types.Int64Value(int64(backup.Size))
	data.Location = types.StringValue(backup.Location)
	data.Created = types.StringValue(backup.Created)
	data.Status = types.StringValue(backup.Status)
	data.ScaletID = types.Int64Value(int64(backup.ScaletID))
	data.Name = types.StringValue(backup.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Recreations are handled by default (RequiresReplace not needed since we don't have update attributes)
}

func (r *BackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BackupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteBackup(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Backup", err.Error())
		return
	}
}

func (r *BackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
