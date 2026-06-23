package provider

import (
	"context"
	"fmt"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &BackupsDataSource{}

type BackupsDataSource struct {
	client *client.Client
}

type BackupModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Template types.String `tfsdk:"template"`
	Size     types.Int64  `tfsdk:"size"`
	Location types.String `tfsdk:"location"`
	Created  types.String `tfsdk:"created"`
	Active   types.Bool   `tfsdk:"active"`
	Locked   types.Bool   `tfsdk:"locked"`
	ScaletID types.Int64  `tfsdk:"scalet_id"`
	Status   types.String `tfsdk:"status"`
}

type BackupsDataSourceModel struct {
	Backups []BackupModel `tfsdk:"backups"`
}

func NewBackupsDataSource() datasource.DataSource {
	return &BackupsDataSource{}
}

func (d *BackupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backups"
}

func (d *BackupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides list of all backups in VScale.",
		Attributes: map[string]schema.Attribute{
			"backups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifier of the backup.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
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
							Description: "Storage datacenter location.",
						},
						"created": schema.StringAttribute{
							Computed:    true,
							Description: "Creation timestamp.",
						},
						"active": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the backup is active.",
						},
						"locked": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the backup is locked.",
						},
						"scalet_id": schema.Int64Attribute{
							Computed:    true,
							Description: "ID of the source Scalet.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Current status of the backup.",
						},
					},
				},
			},
		},
	}
}

func (d *BackupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *BackupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state BackupsDataSourceModel

	backups, err := d.client.GetBackups()
	if err != nil {
		resp.Diagnostics.AddError("Error Fetching Backups", err.Error())
		return
	}

	for _, b := range backups {
		state.Backups = append(state.Backups, BackupModel{
			ID:       types.StringValue(b.ID),
			Name:     types.StringValue(b.Name),
			Template: types.StringValue(b.Template),
			Size:     types.Int64Value(int64(b.Size)),
			Location: types.StringValue(b.Location),
			Created:  types.StringValue(b.Created),
			Active:   types.BoolValue(b.Active),
			Locked:   types.BoolValue(b.Locked),
			ScaletID: types.Int64Value(int64(b.ScaletID)),
			Status:   types.StringValue(b.Status),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
