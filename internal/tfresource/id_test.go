// Copyright Aviatrix Systems, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfresource_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/AviatrixSystems/terraform-provider-avxcloud/internal/tfresource"
)

func TestGetResourcePrefix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		resourceName string
		want         string
	}{
		{
			resourceName: "Fabric",
			want:         "fb",
		},
		{
			resourceName: "Network",
			want:         "net",
		},
		{
			resourceName: "CloudAccount",
			want:         "ca",
		},
		{
			resourceName: "DcfPolicyBlock",
			want:         "dpb",
		},
		{
			resourceName: "Ab",
			want:         "ab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.resourceName, func(t *testing.T) {
			t.Parallel()
			got := tfresource.GetResourcePrefix(tt.resourceName)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateResourceID(t *testing.T) {
	t.Parallel()
	// Test multiple generations to ensure uniqueness
	ids := make(map[string]bool)
	prefix := "net"

	for i := 0; i < 100; i++ {
		id, err := tfresource.GenerateResourceID(prefix)
		require.NoError(t, err)

		// Check prefix
		require.True(t, strings.HasPrefix(id, prefix+"-"), "ID should start with prefix and dash")

		// Check uniqueness
		require.False(t, ids[id], "Generated ID should be unique")
		ids[id] = true

		// Check that the random part is lowercase extended hex base32
		parts := strings.SplitN(id, "-", 2)
		require.Len(t, parts, 2)
		randomPart := parts[1]
		require.Len(t, randomPart, 26, "Random part should be 26 characters long")

		// Extended hexbase32 should only contain a-z and 0-9
		for _, c := range randomPart {
			require.True(t,
				(c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'),
				"Random part should be lowercase extended hexbase32")
		}
	}
}
