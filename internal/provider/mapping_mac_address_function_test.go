// Copyright (c) WANIX Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestMappingMacAddressFunction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fn := NewMappingMacAddressFunction()

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			basetypes.NewNumberValue(big.NewFloat(1)),
			basetypes.NewStringValue("test"),
		}),
	}
	resp := function.RunResponse{
		Result: function.NewResultData(basetypes.NewStringNull()),
	}

	fn.Run(ctx, req, &resp)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	got, ok := resp.Result.Value().(basetypes.StringValue)
	if !ok {
		t.Fatalf("expected basetypes.StringValue result, got %T", resp.Result.Value())
	}

	if got.ValueString() != "f9:cc:b0:a8:cd:2b" {
		t.Fatalf("expected %q, got %q", "f9:cc:b0:a8:cd:2b", got.ValueString())
	}

	macRe := regexp.MustCompile(`^[0-9a-f]{2}(:[0-9a-f]{2}){5}$`)
	if !macRe.MatchString(got.ValueString()) {
		t.Fatalf("expected MAC address format, got %q", got.ValueString())
	}
}

func TestMappingMacAddressFunction_UsesFirstSixBytes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fn := NewMappingMacAddressFunction()

	const (
		namespace = "ns"
		id        = int64(42)
	)

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			basetypes.NewNumberValue(big.NewFloat(float64(id))),
			basetypes.NewStringValue(namespace),
		}),
	}
	resp := function.RunResponse{
		Result: function.NewResultData(basetypes.NewStringNull()),
	}

	fn.Run(ctx, req, &resp)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	got, ok := resp.Result.Value().(basetypes.StringValue)
	if !ok {
		t.Fatalf("expected basetypes.StringValue result, got %T", resp.Result.Value())
	}

	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%d", namespace, id)))
	want := fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		sum[0], sum[1], sum[2], sum[3], sum[4], sum[5],
	)
	if got.ValueString() != want {
		t.Fatalf("expected %q, got %q", want, got.ValueString())
	}
}

func TestMappingMacAddressFunction_RejectsNonIntegerID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fn := NewMappingMacAddressFunction()

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			basetypes.NewNumberValue(big.NewFloat(1.5)),
			basetypes.NewStringValue("test"),
		}),
	}
	resp := function.RunResponse{
		Result: function.NewResultData(basetypes.NewStringNull()),
	}

	fn.Run(ctx, req, &resp)

	if resp.Error == nil {
		t.Fatalf("expected error, got none")
	}
}
