package provider

import (
	"context"
	"fmt"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &PricesDataSource{}

type PricesDataSource struct {
	client *client.Client
}

type PricesDataSourceModel struct {
	Period       types.String  `tfsdk:"period"`
	BackupPrice  types.Float64 `tfsdk:"backup_price"`
	SmallHour    types.Float64 `tfsdk:"small_hour"`
	SmallMonth   types.Float64 `tfsdk:"small_month"`
	MediumHour   types.Float64 `tfsdk:"medium_hour"`
	MediumMonth  types.Float64 `tfsdk:"medium_month"`
	LargeHour    types.Float64 `tfsdk:"large_hour"`
	LargeMonth   types.Float64 `tfsdk:"large_month"`
	HugeHour     types.Float64 `tfsdk:"huge_hour"`
	HugeMonth    types.Float64 `tfsdk:"huge_month"`
	MonsterHour  types.Float64 `tfsdk:"monster_hour"`
	MonsterMonth types.Float64 `tfsdk:"monster_month"`
}

func NewPricesDataSource() datasource.DataSource {
	return &PricesDataSource{}
}

func (d *PricesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prices"
}

func (d *PricesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides price list for VScale resources.",
		Attributes: map[string]schema.Attribute{
			"period": schema.StringAttribute{
				Computed:    true,
				Description: "The pricing period start date.",
			},
			"backup_price": schema.Float64Attribute{
				Computed:    true,
				Description: "Backup storage cost per GB.",
			},
			"small_hour": schema.Float64Attribute{
				Computed: true,
			},
			"small_month": schema.Float64Attribute{
				Computed: true,
			},
			"medium_hour": schema.Float64Attribute{
				Computed: true,
			},
			"medium_month": schema.Float64Attribute{
				Computed: true,
			},
			"large_hour": schema.Float64Attribute{
				Computed: true,
			},
			"large_month": schema.Float64Attribute{
				Computed: true,
			},
			"huge_hour": schema.Float64Attribute{
				Computed: true,
			},
			"huge_month": schema.Float64Attribute{
				Computed: true,
			},
			"monster_hour": schema.Float64Attribute{
				Computed: true,
			},
			"monster_month": schema.Float64Attribute{
				Computed: true,
			},
		},
	}
}

func (d *PricesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PricesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state PricesDataSourceModel

	prices, err := d.client.GetPrices()
	if err != nil {
		resp.Diagnostics.AddError("Error Fetching Prices", err.Error())
		return
	}

	state.Period = types.StringValue(prices.Period)

	m := prices.Default
	if m != nil {
		if b, ok := m["backup"]; ok {
			if bnum, ok := b.(float64); ok {
				state.BackupPrice = types.Float64Value(bnum)
			}
		}

		sh, sm := getPriceDetail(m, "small")
		state.SmallHour = types.Float64Value(sh)
		state.SmallMonth = types.Float64Value(sm)

		mh, mm := getPriceDetail(m, "medium")
		state.MediumHour = types.Float64Value(mh)
		state.MediumMonth = types.Float64Value(mm)

		lh, lm := getPriceDetail(m, "large")
		state.LargeHour = types.Float64Value(lh)
		state.LargeMonth = types.Float64Value(lm)

		hh, hm := getPriceDetail(m, "huge")
		state.HugeHour = types.Float64Value(hh)
		state.HugeMonth = types.Float64Value(hm)

		monh, monm := getPriceDetail(m, "monster")
		state.MonsterHour = types.Float64Value(monh)
		state.MonsterMonth = types.Float64Value(monm)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getPriceDetail(m map[string]interface{}, key string) (float64, float64) {
	val, ok := m[key]
	if !ok {
		return 0, 0
	}
	inner, ok := val.(map[string]interface{})
	if !ok {
		return 0, 0
	}
	var hr, mo float64
	if h, ok := inner["hour"]; ok {
		if hnum, ok := h.(float64); ok {
			hr = hnum
		}
	}
	if moVal, ok := inner["month"]; ok {
		if monum, ok := moVal.(float64); ok {
			mo = monum
		}
	}
	return hr, mo
}
