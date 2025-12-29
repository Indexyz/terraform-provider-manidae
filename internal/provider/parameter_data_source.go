// Copyright (c) WANIX Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	parameterTypeString = "string"
	parameterTypeNumber = "number"
)

type parameterDataSource struct{}

type parameterValidationModel struct {
	Min types.Number `tfsdk:"min"`
	Max types.Number `tfsdk:"max"`
}

type parameterOptionModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type parameterDataSourceModel struct {
	ID                  types.String             `tfsdk:"id"`
	Name                types.String             `tfsdk:"name"`
	DisplayName         types.String             `tfsdk:"display_name"`
	Description         types.String             `tfsdk:"description"`
	Type                types.String             `tfsdk:"type"`
	Default             types.Dynamic            `tfsdk:"default"`
	Value               types.Dynamic            `tfsdk:"value"`
	EnvironmentVariable types.String             `tfsdk:"environment_variable"`
	Validation          parameterValidationModel `tfsdk:"validation"`
	Options             []parameterOptionModel   `tfsdk:"option"`
}

func NewParameterDataSource() datasource.DataSource {
	return &parameterDataSource{}
}

func (d *parameterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_parameter"
}

func (d *parameterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads a parameter value from an environment variable derived from `name`, falling back to `default`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal identifier (same as `name`).",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Parameter name (used to derive the environment variable key).",
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Human-friendly display name.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Human-friendly description.",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Parameter type. Supported values: `string`, `number`. If unset, inferred from `default`.",
			},
			"default": schema.DynamicAttribute{
				Optional:            true,
				MarkdownDescription: "Default value used when the environment variable is not set.",
			},
			"value": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: "Resolved value (from environment variable if set, otherwise `default`).",
			},
			"environment_variable": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Environment variable key used to resolve the value.",
			},
		},
		Blocks: map[string]schema.Block{
			"validation": schema.SingleNestedBlock{
				MarkdownDescription: "Numeric validation (only valid when `type = \"number\"`).",
				Attributes: map[string]schema.Attribute{
					"min": schema.NumberAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum allowed value (inclusive).",
					},
					"max": schema.NumberAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum allowed value (inclusive).",
					},
				},
			},
			"option": schema.ListNestedBlock{
				MarkdownDescription: "Allowed values (enum) when `type = \"string\"`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Human-friendly option label.",
						},
						"value": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Allowed value.",
						},
					},
				},
			},
		},
	}
}

func (d *parameterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data parameterDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parameterName := data.Name.ValueString()
	parameterType, typeDiags := resolveParameterType(data.Type, data.Default)
	resp.Diagnostics.Append(typeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Type = types.StringValue(parameterType)

	envKey := ParameterEnvironmentVariable(parameterName)
	data.EnvironmentVariable = types.StringValue(envKey)

	value, valueDiags := resolveParameterValue(parameterType, envKey, data.Default)
	resp.Diagnostics.Append(valueDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateParameterValue(parameterType, value, data.Validation, data.Options)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(parameterName)
	data.Value = types.DynamicValue(value)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func resolveParameterType(typeAttr types.String, defaultValue types.Dynamic) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if typeAttr.IsUnknown() {
		diags.AddError("Invalid type", "`type` must be known")
		return "", diags
	}

	if !typeAttr.IsNull() {
		raw := strings.ToLower(strings.TrimSpace(typeAttr.ValueString()))
		switch raw {
		case parameterTypeString, parameterTypeNumber:
			return raw, diags
		default:
			diags.AddError("Invalid type", fmt.Sprintf("unsupported `type` %q (supported: %q, %q)", raw, parameterTypeString, parameterTypeNumber))
			return "", diags
		}
	}

	if defaultValue.IsUnknown() {
		diags.AddError("Invalid default", "`default` must be known to infer `type`")
		return "", diags
	}

	if defaultValue.IsNull() {
		diags.AddError("Missing type", "`type` is required when `default` is not set")
		return "", diags
	}

	switch defaultValue.UnderlyingValue().(type) {
	case types.String:
		return parameterTypeString, diags
	case types.Number:
		return parameterTypeNumber, diags
	default:
		diags.AddError("Invalid default", "unsupported `default` type (supported: string, number)")
		return "", diags
	}
}

func resolveParameterValue(parameterType string, envKey string, defaultValue types.Dynamic) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	rawEnv, hasEnv := os.LookupEnv(envKey)
	if hasEnv {
		return parseParameterValue(parameterType, rawEnv, true)
	}

	if defaultValue.IsUnknown() {
		diags.AddError("Missing value", fmt.Sprintf("environment variable %q is not set and `default` is unknown", envKey))
		return nil, diags
	}

	if defaultValue.IsNull() {
		diags.AddError("Missing value", fmt.Sprintf("environment variable %q is not set and `default` is not configured", envKey))
		return nil, diags
	}

	underlying := defaultValue.UnderlyingValue()
	switch parameterType {
	case parameterTypeString:
		switch v := underlying.(type) {
		case types.String:
			if v.IsUnknown() {
				diags.AddError("Invalid default", "`default` must be known")
				return nil, diags
			}
			if v.IsNull() {
				diags.AddError("Invalid default", "`default` must not be null")
				return nil, diags
			}
			return v, diags
		default:
			diags.AddError("Invalid default", "expected `default` to be a string")
			return nil, diags
		}
	case parameterTypeNumber:
		switch v := underlying.(type) {
		case types.Number:
			if v.IsUnknown() {
				diags.AddError("Invalid default", "`default` must be known")
				return nil, diags
			}
			if v.IsNull() {
				diags.AddError("Invalid default", "`default` must not be null")
				return nil, diags
			}
			return v, diags
		case types.String:
			if v.IsUnknown() {
				diags.AddError("Invalid default", "`default` must be known")
				return nil, diags
			}
			if v.IsNull() {
				diags.AddError("Invalid default", "`default` must not be null")
				return nil, diags
			}
			return parseParameterValue(parameterType, v.ValueString(), false)
		default:
			diags.AddError("Invalid default", "expected `default` to be a number")
			return nil, diags
		}
	default:
		diags.AddError("Invalid type", fmt.Sprintf("unsupported `type` %q", parameterType))
		return nil, diags
	}
}

func parseParameterValue(parameterType string, raw string, fromEnv bool) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch parameterType {
	case parameterTypeString:
		return types.StringValue(raw), diags
	case parameterTypeNumber:
		number, ok := new(big.Float).SetString(strings.TrimSpace(raw))
		if !ok {
			source := "`default`"
			if fromEnv {
				source = "environment variable"
			}
			diags.AddError("Invalid number", fmt.Sprintf("%s value %q cannot be parsed as a number", source, raw))
			return nil, diags
		}
		return types.NumberValue(number), diags
	default:
		diags.AddError("Invalid type", fmt.Sprintf("unsupported `type` %q", parameterType))
		return nil, diags
	}
}

func validateParameterValue(parameterType string, value attr.Value, validation parameterValidationModel, options []parameterOptionModel) diag.Diagnostics {
	var diags diag.Diagnostics

	switch parameterType {
	case parameterTypeString:
		minSet := !validation.Min.IsNull() && !validation.Min.IsUnknown()
		maxSet := !validation.Max.IsNull() && !validation.Max.IsUnknown()
		if minSet || maxSet {
			diags.AddError("Invalid validation", "`validation` is only supported when `type = \"number\"`")
			return diags
		}

		stringValue, ok := value.(types.String)
		if !ok {
			diags.AddError("Invalid value", "expected resolved value to be a string")
			return diags
		}
		if stringValue.IsNull() || stringValue.IsUnknown() {
			diags.AddError("Invalid value", "resolved value must be known")
			return diags
		}

		if len(options) == 0 {
			return diags
		}

		allowed := make(map[string]struct{}, len(options))
		for i, opt := range options {
			if opt.Value.IsUnknown() || opt.Value.IsNull() {
				diags.AddError("Invalid option", fmt.Sprintf("option[%d].value must be set", i))
				continue
			}
			allowed[opt.Value.ValueString()] = struct{}{}
		}
		if diags.HasError() {
			return diags
		}

		if _, ok := allowed[stringValue.ValueString()]; !ok {
			diags.AddError("Invalid value", fmt.Sprintf("value %q is not one of the configured options", stringValue.ValueString()))
		}

		return diags
	case parameterTypeNumber:
		if len(options) > 0 {
			diags.AddError("Invalid option", "`option` blocks are only supported when `type = \"string\"`")
			return diags
		}

		numberValue, ok := value.(types.Number)
		if !ok {
			diags.AddError("Invalid value", "expected resolved value to be a number")
			return diags
		}
		if numberValue.IsNull() || numberValue.IsUnknown() {
			diags.AddError("Invalid value", "resolved value must be known")
			return diags
		}

		val := numberValue.ValueBigFloat()
		if validation.Min.IsUnknown() || validation.Max.IsUnknown() {
			diags.AddError("Invalid validation", "`validation.min` and `validation.max` must be known when set")
			return diags
		}

		if !validation.Min.IsNull() && !validation.Max.IsNull() {
			if validation.Min.ValueBigFloat().Cmp(validation.Max.ValueBigFloat()) > 0 {
				diags.AddError("Invalid validation", "`validation.min` must be <= `validation.max`")
				return diags
			}
		}

		if !validation.Min.IsNull() {
			if val.Cmp(validation.Min.ValueBigFloat()) < 0 {
				diags.AddError("Invalid value", fmt.Sprintf("value %s is less than validation.min %s", val.String(), validation.Min.ValueBigFloat().String()))
				return diags
			}
		}

		if !validation.Max.IsNull() {
			if val.Cmp(validation.Max.ValueBigFloat()) > 0 {
				diags.AddError("Invalid value", fmt.Sprintf("value %s is greater than validation.max %s", val.String(), validation.Max.ValueBigFloat().String()))
				return diags
			}
		}

		return diags
	default:
		diags.AddError("Invalid type", fmt.Sprintf("unsupported `type` %q", parameterType))
		return diags
	}
}
