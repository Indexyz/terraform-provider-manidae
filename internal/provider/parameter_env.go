// Copyright (c) WANIX Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"crypto/sha256"
	"encoding/hex"
)

func ParameterEnvironmentVariable(name string) string {
	sum := sha256.Sum256([]byte(name))
	return "MANIDAE_PARAMETER_" + hex.EncodeToString(sum[:])
}
