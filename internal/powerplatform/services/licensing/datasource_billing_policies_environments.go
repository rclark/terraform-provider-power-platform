// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &BillingPoliciesEnvironmetsDataSource{}
	_ datasource.DataSourceWithConfigure = &BillingPoliciesEnvironmetsDataSource{}
)

func NewBillingPoliciesEnvironmetsDataSource() datasource.DataSource {
	return &BillingPoliciesEnvironmetsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_billing_policies_environments",
	}
}

type BillingPoliciesEnvironmetsDataSource struct {
	LicensingClient  LicensingClient
	ProviderTypeName string
	TypeName         string
}

type BillingPoliciesEnvironmetsListDataSourceModel struct {
	BillingPolicyId string   `tfsdk:"billing_policy_id"`
	Environments    []string `tfsdk:"environments"`
}

func (d *BillingPoliciesEnvironmetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *BillingPoliciesEnvironmetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the environments associated with a billing policy",
		MarkdownDescription: "Fetches the environments associated with a [billing policy](https://learn.microsoft.com/en-us/power-platform/admin/pay-as-you-go-overview#what-is-a-billing-policy).\n\nThis data source uses the [List Billing Policy Environments](https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/billing-policy-environment/list-billing-policy-environments) endpoint in the Power Platform API.",
		Attributes: map[string]schema.Attribute{
			"billing_policy_id": schema.StringAttribute{
				Required:            true,
				Description:         "The id of the billing policy",
				MarkdownDescription: "The id of the billing policy",
			},
			"environments": schema.SetAttribute{
				Description:         "The environments associated with the billing policy",
				MarkdownDescription: "The environments associated with the billing policy",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *BillingPoliciesEnvironmetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientApi := req.ProviderData.(*api.ProviderClient).Api

	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.LicensingClient = NewLicensingClient(clientApi)
}

func (d *BillingPoliciesEnvironmetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state BillingPoliciesEnvironmetsListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE START: %s", d.ProviderTypeName))

	if resp.Diagnostics.HasError() {
		return
	}

	diag := req.Config.GetAttribute(ctx, path.Root("billing_policy_id"), &state.BillingPolicyId)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	environments, err := d.LicensingClient.GetEnvironmentsForBillingPolicy(ctx, state.BillingPolicyId)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	state.Environments = environments

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}