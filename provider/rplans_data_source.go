package provider

import (
	"context"
	"fmt"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &RPlansDataSource{}

type RPlansDataSource struct {
	client *client.Client
}

type RPlanModel struct {
	ID        types.String `tfsdk:"id"`
	Memory    types.Int64  `tfsdk:"memory"`
	CPUs      types.Int64  `tfsdk:"cpus"`
	Disk      types.Int64  `tfsdk:"disk"`
	Addresses types.Int64  `tfsdk:"addresses"`
	Locations types.List   `tfsdk:"locations"`
	Templates types.List   `tfsdk:"templates"`
}

type RPlansDataSourceModel struct {
	RPlans []RPlanModel `tfsdk:"rplans"`
}

func NewRPlansDataSource() datasource.DataSource {
	return &RPlansDataSource{}
}

func (d *RPlansDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rplans"
}

func (d *RPlansDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides list of available pricing/tariff plans (rplans) in VScale.",
		Attributes: map[string]schema.Attribute{
			"rplans": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifier of the plan (e.g. small, medium).",
						},
						"memory": schema.Int64Attribute{
							Computed:    true,
							Description: "RAM size in MB.",
						},
						"cpus": schema.Int64Attribute{
							Computed:    true,
							Description: "CPU core count.",
						},
						"disk": schema.Int64Attribute{
							Computed:    true,
							Description: "Disk size in GB or MB.",
						},
						"addresses": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of IP addresses included.",
						},
						"locations": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Locations where this plan is available.",
						},
						"templates": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Templates compatible with this plan.",
						},
					},
				},
			},
		},
	}
}

func (d *RPlansDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RPlansDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RPlansDataSourceModel

	rplans, err := d.client.GetRPlans()
	if err != nil {
		resp.Diagnostics.AddError("Error Fetching RPlans", err.Error())
		return
	}

	for _, rp := range rplans {
		locations, diags := types.ListValueFrom(ctx, types.StringType, rp.Locations)
		resp.Diagnostics.Append(diags...)

		templates, diags := types.ListValueFrom(ctx, types.StringType, rp.Templates)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		state.RPlans = append(state.RPlans, RPlanModel{
			ID:        types.StringValue(rp.ID),
			Memory:    types.Int64Value(int64(rp.Memory)),
			CPUs:      types.Int64Value(int64(rp.CPUs)),
			Disk:      types.Int64Value(int64(rp.Disk)),
			Addresses: types.Int64Value(int64(rp.Addresses)),
			Locations: locations,
			Templates: templates,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
