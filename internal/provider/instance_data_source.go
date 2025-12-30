// Copyright (c) WANIX Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type instanceDataSource struct{}

type instanceDataSourceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	ConnectionID types.String `tfsdk:"connection_id"`
	Identity     types.String `tfsdk:"identity"`
	Action       types.String `tfsdk:"action"`
	State        types.String `tfsdk:"state"`
	StartCount   types.Int64  `tfsdk:"start_count"`
}

func NewInstanceDataSource() datasource.DataSource {
	return &instanceDataSource{}
}

func (d *instanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

func (d *instanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads Manidae instance context from environment variables.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Instance ID from `MANIDAE_INSTANCE_ID` (must be a non-negative integer).",
			},
			"connection_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Connection ID from `MANIDAE_CONNECTION_ID`.",
			},
			"identity": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identity from `MANIDAE_IDENTITY`.",
			},
			"action": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Action from `MANIDAE_ACTION`.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Instance state from `MANIDAE_INSTANCE_STATE` (`on` or `off`).",
			},
			"start_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Derived from `state`: `1` when `on`, otherwise `0`.",
			},
		},
	}
}

func (d *instanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data instanceDataSourceModel

	id, idDiags := getRequiredUintEnvAsInt64("MANIDAE_INSTANCE_ID")
	resp.Diagnostics.Append(idDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	connectionID, connectionDiags := getRequiredEnvString("MANIDAE_CONNECTION_ID")
	resp.Diagnostics.Append(connectionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	identity, identityDiags := getRequiredEnvString("MANIDAE_IDENTITY")
	resp.Diagnostics.Append(identityDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	action, actionDiags := getRequiredEnvString("MANIDAE_ACTION")
	resp.Diagnostics.Append(actionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, stateDiags := getRequiredEnvString("MANIDAE_INSTANCE_STATE")
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	startCount, startCountDiags := deriveStartCount(state)
	resp.Diagnostics.Append(startCountDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.Int64Value(id)
	data.ConnectionID = types.StringValue(connectionID)
	data.Identity = types.StringValue(identity)
	data.Action = types.StringValue(action)
	data.State = types.StringValue(state)
	data.StartCount = types.Int64Value(startCount)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getRequiredEnvString(key string) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		diags.AddError("Missing environment variable", fmt.Sprintf("%q must be set", key))
		return "", diags
	}

	return strings.TrimSpace(value), diags
}

func getRequiredUintEnvAsInt64(key string) (int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	raw, rawDiags := getRequiredEnvString(key)
	diags.Append(rawDiags...)
	if diags.HasError() {
		return 0, diags
	}

	uintValue, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		diags.AddError("Invalid environment variable", fmt.Sprintf("%q must be a non-negative integer: %s", key, err))
		return 0, diags
	}

	if uintValue > (^uint64(0) >> 1) {
		diags.AddError("Invalid environment variable", fmt.Sprintf("%q is too large to fit into Terraform int64", key))
		return 0, diags
	}

	return int64(uintValue), diags
}

func deriveStartCount(state string) (int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch strings.ToLower(strings.TrimSpace(state)) {
	case "on":
		return 1, diags
	case "off":
		return 0, diags
	default:
		diags.AddError("Invalid instance state", "MANIDAE_INSTANCE_STATE must be either \"on\" or \"off\"")
		return 0, diags
	}
}
