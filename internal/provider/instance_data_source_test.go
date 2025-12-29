// Copyright (c) WANIX Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

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
