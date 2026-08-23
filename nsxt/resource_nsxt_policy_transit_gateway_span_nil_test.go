// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetSpanFromSchemaNilSpan guards against a regression where the
// Terraform SDK populates an optional "span" block that contains only
// empty/default attributes with a single nil element rather than an empty
// slice, or hands getSpanFromSchema a bare nil. Unchecked type assertions on
// that element previously triggered a panic at apply time; getSpanFromSchema
// must instead treat these shapes as "no span configured".
func TestGetSpanFromSchemaNilSpan(t *testing.T) {
	structVal, err := getSpanFromSchema(nil)
	require.NoError(t, err)
	require.Nil(t, structVal)

	structVal, err = getSpanFromSchema([]interface{}{})
	require.NoError(t, err)
	require.Nil(t, structVal)

	structVal, err = getSpanFromSchema([]interface{}{nil})
	require.NoError(t, err)
	require.Nil(t, structVal)
}
