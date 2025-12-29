// Copyright (c) WANIX Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParameterEnvironmentVariable(t *testing.T) {
	got := ParameterEnvironmentVariable("root_volume_size_gb")

	const prefix = "MANIDAE_PARAMETER_"
	if len(got) != len(prefix)+64 {
		t.Fatalf("expected env var length %d, got %d (%q)", len(prefix)+64, len(got), got)
	}
	if got[:len(prefix)] != prefix {
		t.Fatalf("expected env var prefix %q, got %q", prefix, got)
	}
}

func TestValidateParameterValue_NumberMin(t *testing.T) {
	value := types.NumberValue(new(big.Float).SetInt64(19))
	diags := validateParameterValue(parameterTypeNumber, value, &parameterValidationModel{
		Min: types.NumberValue(new(big.Float).SetInt64(20)),
	}, nil)

	if !diags.HasError() {
		t.Fatalf("expected validation error, got none")
	}
}

func TestValidateParameterValue_StringOptions(t *testing.T) {
	value := types.StringValue("SA2.MEDIUM8")
	diags := validateParameterValue(parameterTypeString, value, nil, []parameterOptionModel{
		{Value: types.StringValue("SA2.MEDIUM2")},
		{Value: types.StringValue("SA2.MEDIUM4")},
	})

	if !diags.HasError() {
		t.Fatalf("expected options validation error, got none")
	}
}
