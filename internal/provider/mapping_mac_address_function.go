// Copyright (c) WANIX Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ function.Function = (*mappingMacAddressFunction)(nil)

type mappingMacAddressFunction struct{}

func NewMappingMacAddressFunction() function.Function {
	return &mappingMacAddressFunction{}
}

func (f *mappingMacAddressFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "mapping_mac_address"
}

func (f *mappingMacAddressFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Derive a deterministic MAC address from a namespace and numeric identifier.",
		Parameters: []function.Parameter{
			function.NumberParameter{
				Name: "id",
			},
			function.StringParameter{
				Name: "namespace",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *mappingMacAddressFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var id basetypes.NumberValue
	var namespace string

	if funcErr := req.Arguments.Get(ctx, &id, &namespace); funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	idFloat := id.ValueBigFloat()
	if idFloat == nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(0, "id is required"))
		return
	}

	if !idFloat.IsInt() {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(0, "id must be an integer number"))
		return
	}

	idInt := new(big.Int)
	idFloat.Int(idInt)

	sum := sha256.Sum256([]byte(namespace + "|" + idInt.String()))

	mac := fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		sum[0], sum[1], sum[2], sum[3], sum[4], sum[5],
	)

	if funcErr := resp.Result.Set(ctx, mac); funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	}
}
