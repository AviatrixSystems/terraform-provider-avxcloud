// Copyright Aviatrix Systems, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfresource

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

var resourcePrefixOverrides = map[string]string{
	"Fabric": "fb",
}

// GetResourcePrefix returns the prefix for a resource name.
func GetResourcePrefix(resourceName string) string {
	if prefix, ok := resourcePrefixOverrides[resourceName]; ok {
		return prefix
	}
	// if there are 2 or more capital letters, use those as the prefix
	caps := []rune{}
	for _, c := range resourceName {
		if c >= 'A' && c <= 'Z' {
			caps = append(caps, c)
		}
	}
	if len(caps) >= 2 {
		return strings.ToLower(string(caps[:min(3, len(caps))]))
	}
	// otherwise use the first 3 letters of the resource name
	return strings.ToLower(resourceName[:min(3, len(resourceName))])
}

// GenerateResourceID generates a unique resource ID with the given prefix
// Format: {prefix}-{base32_uuid}
// Example: net-abc123def456.
func GenerateResourceID(prefix string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to extended hexbase32 without padding, lowercase
	id := strings.ToLower(base32.HexEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(b))

	return fmt.Sprintf("%s-%s", prefix, id), nil
}
