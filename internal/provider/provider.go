// Copyright Aviatrix Systems, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	configtf "github.com/AviatrixSystems/terraform-provider-avxcloud/internal/gen-terraform/aviatrix/config/v1alpha"
	"github.com/AviatrixSystems/terraform-provider-avxcloud/internal/clients"
)

var _ provider.Provider = &AviatrixProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AviatrixProvider{
			version: version,
		}
	}
}

type AviatrixProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *AviatrixProvider) Metadata(ctx context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "avxcloud"
	resp.Version = p.version
}

func (p *AviatrixProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"addr": schema.StringAttribute{
				Description: "The API address of the Aviatrix controller. Defaults to https://api.cloud.aviatrix.com.",
				Optional:    true,
			},
			"api_access_key": schema.StringAttribute{
				Description: "PaaS API access key.",
				Required:    true,
				Sensitive:   true,
			},
			"insecure_tls": schema.BoolAttribute{
				Description: "Disable TLS verification",
				Optional:    true,
			},
		},
	}
}

type ProviderConfig struct {
	Addr         types.String `tfsdk:"addr"`
	InsecureTLS  types.Bool   `tfsdk:"insecure_tls"`
	APIAccessKey types.String `tfsdk:"api_access_key"`
}

func (p *AviatrixProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ProviderConfig

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clients, err := clients.New(ctx, &clients.Config{
		Addr:         config.Addr.ValueString(),
		InsecureTLS:  config.InsecureTLS.ValueBool(),
		APIAccessKey: config.APIAccessKey.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"cannot create client",
			err.Error(),
		)
		return
	}

	resp.ResourceData = clients
}

func (p *AviatrixProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *AviatrixProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		configtf.NewCloudAccountResource,
		configtf.NewNetworkResource,
		configtf.NewNetworkInspectionResource,
	}
}
