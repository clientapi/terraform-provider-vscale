package provider

import (
	"context"
	"fmt"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ImagesDataSource{}

type ImagesDataSource struct {
	client *client.Client
}

type ImageModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Active      types.Bool   `tfsdk:"active"`
	Size        types.Int64  `tfsdk:"size"`
	Locations   types.List   `tfsdk:"locations"`
	RPlans      types.List   `tfsdk:"rplans"`
}

type ImagesDataSourceModel struct {
	Images []ImageModel `tfsdk:"images"`
}

func NewImagesDataSource() datasource.DataSource {
	return &ImagesDataSource{}
}

func (d *ImagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

// Let's write the correct method body
func (d *ImagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides list of available OS images in VScale.",
		Attributes: map[string]schema.Attribute{
			"images": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Identifier of the image (OS name).",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable description.",
						},
						"active": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the image is active.",
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "Size of the image in MB.",
						},
						"locations": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Locations where this image is available.",
						},
						"rplans": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Tariff plans compatible with this image.",
						},
					},
				},
			},
		},
	}
}

func (d *ImagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ImagesDataSourceModel

	images, err := d.client.GetImages()
	if err != nil {
		resp.Diagnostics.AddError("Error Fetching Images", err.Error())
		return
	}

	for _, img := range images {
		locs, diags := types.ListValueFrom(ctx, types.StringType, img.Locations)
		resp.Diagnostics.Append(diags...)

		rplans, diags := types.ListValueFrom(ctx, types.StringType, img.RPlans)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		state.Images = append(state.Images, ImageModel{
			ID:          types.StringValue(img.ID),
			Description: types.StringValue(img.Description),
			Active:      types.BoolValue(img.Active),
			Size:        types.Int64Value(int64(img.Size)),
			Locations:   locs,
			RPlans:      rplans,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
