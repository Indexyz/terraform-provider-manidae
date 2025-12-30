// Copyright (c) WANIX Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestDeriveStartCount(t *testing.T) {
	t.Run("on", func(t *testing.T) {
		got, diags := deriveStartCount("on")
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %#v", diags)
		}
		if got != 1 {
			t.Fatalf("expected 1, got %d", got)
		}
	})

	t.Run("off", func(t *testing.T) {
		got, diags := deriveStartCount("off")
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %#v", diags)
		}
		if got != 0 {
			t.Fatalf("expected 0, got %d", got)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		_, diags := deriveStartCount("maybe")
		if !diags.HasError() {
			t.Fatalf("expected error, got none")
		}
	})
}

func TestGetRequiredEnvString(t *testing.T) {
	t.Setenv("MANIDAE_CONNECTION_ID", "  abc123  ")

	got, diags := getRequiredEnvString("MANIDAE_CONNECTION_ID")
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if got != "abc123" {
		t.Fatalf("expected %q, got %q", "abc123", got)
	}
}

func TestInstanceDataSourceSchema_HasIdentity(t *testing.T) {
	t.Parallel()

	ds := NewInstanceDataSource()

	var resp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	if _, ok := resp.Schema.Attributes["identity"]; !ok {
		t.Fatalf("expected schema to include attribute %q", "identity")
	}
}

func TestInstanceDataSourceRead_RequiresIdentity(t *testing.T) {
	t.Setenv("MANIDAE_INSTANCE_ID", "1")
	t.Setenv("MANIDAE_CONNECTION_ID", "cid")
	t.Setenv("MANIDAE_ACTION", "action")
	t.Setenv("MANIDAE_INSTANCE_STATE", "on")

	ds := NewInstanceDataSource()

	var resp datasource.ReadResponse
	ds.Read(context.Background(), datasource.ReadRequest{}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected error, got none")
	}
}
