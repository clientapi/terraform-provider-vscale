package provider

import (
	"context"
	"fmt"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &LocationsDataSource{}

type LocationsDataSource struct {
	client *client.Client
}

type LocationModel struct {
	ID                types.String `tfsdk:"id"`
	Description       types.String `tfsdk:"description"`
	Active            types.Bool   `tfsdk:"active"`
	PrivateNetworking types.Bool   `tfsdk:"private_networking"`
	Templates         types.List   `tfsdk:"templates"`
	RPlans            types.List   `tfsdk:"rplans"`
}

type LocationsDataSourceModel struct {
	Locations []LocationModel `tfsdk:"locations"`
}

func NewLocationsDataSource() datasource.DataSource {
	return &LocationsDataSource{}
}

func (d *LocationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_locations"
}

func (d *LocationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides list of available datacenters/locations in VScale.",
		Attributes: map[string]schema.Attribute{
			"locations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifier of the location (e.g. spb0).",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable description.",
						},
						"active": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the location is active.",
						},
						"private_networking": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether private networking is supported in this location.",
						},
						"templates": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Templates available in this location.",
						},
						"rplans": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Tariff plans available in this location.",
						},
					},
				},
			},
		},
	}
}

func (d *LocationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LocationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state LocationsDataSourceModel

	locations, err := d.client.GetLocations()
	if err != nil {
		resp.Diagnostics.AddError("Error Fetching Locations", err.Error())
		return
	}

	for _, loc := range locations {
		templates, diags := types.ListValueFrom(ctx, types.StringType, loc.Templates)
		resp.Diagnostics.Append(diags...)

		rplans, diags := types.ListValueFrom(ctx, types.StringType, loc.RPlans)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		state.Locations = append(state.Locations, LocationModel{
			ID:                types.StringValue(loc.ID),
			Description:       types.StringValue(loc.Description),
			Active:            types.BoolValue(loc.Active),
			PrivateNetworking: types.BoolValue(loc.PrivateNetworking),
			Templates:         templates,
			RPlans:            rplans,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
