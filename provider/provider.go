package provider

import (
	"context"
	"os"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &VScaleProvider{}

type VScaleProvider struct {
	version string
}

type VScaleProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &VScaleProvider{
			version: version,
		}
	}
}

func (p *VScaleProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vscale"
	resp.Version = p.version
}

func (p *VScaleProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for VScale (Selectel VDS) API.",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The VScale API token. Can also be set via the VSCALE_TOKEN environment variable.",
			},
		},
	}
}

func (p *VScaleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data VScaleProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv("VSCALE_TOKEN")

	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"Missing VScale API Token",
			"The provider requires a VScale API Token to be configured. "+
				"Please set the 'token' attribute in the provider block or the VSCALE_TOKEN environment variable.",
		)
		return
	}

	c := client.NewClient(token)

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *VScaleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSSHKeyResource,
		NewScaletResource,
		NewDomainResource,
		NewDomainRecordResource,
		NewDomainTagResource,
		NewPTRRecordResource,
		NewBackupResource,
	}
}

func (p *VScaleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountDataSource,
		NewImagesDataSource,
		NewLocationsDataSource,
		NewRPlansDataSource,
		NewPricesDataSource,
		NewBackupsDataSource,
	}
}
