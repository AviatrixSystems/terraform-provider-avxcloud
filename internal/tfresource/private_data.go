// Copyright Aviatrix Systems, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfresource

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const generationKey = "generation"

type PrivateData interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

func SetGeneration(ctx context.Context, private PrivateData, gen uint32) diag.Diagnostics {
	genStr := strconv.FormatUint(uint64(gen), 10)
	return private.SetKey(ctx, generationKey, []byte(genStr))
}

func GetGeneration(ctx context.Context, private PrivateData) (uint32, diag.Diagnostics) {
	genBytes, diags := private.GetKey(ctx, "generation")
	if diags.HasError() {
		return 0, diags
	}
	gen, err := strconv.ParseUint(string(genBytes), 10, 32)
	if err != nil {
		diags.AddError("failed to parse generation", err.Error())
	}
	return uint32(gen), diags
}
