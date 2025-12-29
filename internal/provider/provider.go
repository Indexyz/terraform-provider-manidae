// Copyright (c) WANIX Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ManidaeProvider satisfies the provider interface.
var _ provider.Provider = &ManidaeProvider{}

// ManidaeProvider defines the provider implementation.
type ManidaeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ManidaeProviderModel describes the provider data model.
type ManidaeProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *ManidaeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "manidae"
	resp.Version = p.version
}

func (p *ManidaeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *ManidaeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ManidaeProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *ManidaeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *ManidaeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewParameterDataSource,
		NewInstanceDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ManidaeProvider{
			version: version,
		}
	}
}
