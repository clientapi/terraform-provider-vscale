package provider

import (
	"context"
	"fmt"

	"terraform-provider-vscale/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &AccountDataSource{}

type AccountDataSource struct {
	client *client.Client
}

type AccountDataSourceModel struct {
	ID         types.String  `tfsdk:"id"`
	Name       types.String  `tfsdk:"name"`
	Middlename types.String  `tfsdk:"middlename"`
	Surname    types.String  `tfsdk:"surname"`
	Email      types.String  `tfsdk:"email"`
	Mobile     types.String  `tfsdk:"mobile"`
	Country    types.String  `tfsdk:"country"`
	ActDate    types.String  `tfsdk:"actdate"`
	FaceID     types.String  `tfsdk:"face_id"`
	State      types.String  `tfsdk:"state"`
	Balance    types.Float64 `tfsdk:"balance"`
	Bonus      types.Float64 `tfsdk:"bonus"`
}

func NewAccountDataSource() datasource.DataSource {
	return &AccountDataSource{}
}

func (d *AccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (d *AccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches information about the current VScale account and billing balance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Account ID.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "User's first name.",
			},
			"middlename": schema.StringAttribute{
				Computed:    true,
				Description: "User's middle name.",
			},
			"surname": schema.StringAttribute{
				Computed:    true,
				Description: "User's last name.",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "Registered email address.",
			},
			"mobile": schema.StringAttribute{
				Computed:    true,
				Description: "Registered mobile number.",
			},
			"country": schema.StringAttribute{
				Computed:    true,
				Description: "User's country.",
			},
			"actdate": schema.StringAttribute{
				Computed:    true,
				Description: "Account activation date.",
			},
			"face_id": schema.StringAttribute{
				Computed:    true,
				Description: "Internal type identifier for physical/legal entity.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Account state (1 for active, 0 for inactive).",
			},
			"balance": schema.Float64Attribute{
				Computed:    true,
				Description: "Current balance amount in kopecks/rubles.",
			},
			"bonus": schema.Float64Attribute{
				Computed:    true,
				Description: "Current bonus amount.",
			},
		},
	}
}

func (d *AccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AccountDataSourceModel

	info, err := d.client.GetAccountInfo()
	if err != nil {
		resp.Diagnostics.AddError("Error Fetching Account Info", err.Error())
		return
	}

	balance, err := d.client.GetBalance()
	if err != nil {
		resp.Diagnostics.AddError("Error Fetching Balance Info", err.Error())
		return
	}

	state.ID = types.StringValue(info.ID)
	state.Name = types.StringValue(info.Name)
	state.Middlename = types.StringValue(info.Middlename)
	state.Surname = types.StringValue(info.Surname)
	state.Email = types.StringValue(info.Email)
	state.Mobile = types.StringValue(info.Mobile)
	state.Country = types.StringValue(info.Country)
	state.ActDate = types.StringValue(info.ActDate)
	state.FaceID = types.StringValue(info.FaceID)
	state.State = types.StringValue(info.State)

	state.Balance = types.Float64Value(balance.Balance)
	state.Bonus = types.Float64Value(balance.Bonus)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
