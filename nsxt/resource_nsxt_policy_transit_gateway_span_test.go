// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetSpanFromSchemaZoneBasedEmptyZones guards against a regression where
// span.zone_based_span with an explicitly empty zone_external_ids list was
// dropped from the create/update payload entirely. The SDKv2 diff engine
// emits no attribute entries for the nested block's contents when the list
// stays at zero elements, so the reconstructed ResourceData can hand
// getSpanFromSchema a zone_based_span element that is a bare nil instead of
// a map. getSpanFromSchema must still recognize the block as present (via
// the list length) and produce a ZoneBasedSpan rather than panicking or
// silently dropping the span.
func TestGetSpanFromSchemaZoneBasedEmptyZones(t *testing.T) {
	span := []interface{}{
		map[string]interface{}{
			"cluster_based_span": []interface{}{},
			"zone_based_span":    []interface{}{nil},
		},
	}

	structVal, err := getSpanFromSchema(span)
	require.NoError(t, err)
	require.NotNil(t, structVal, "getSpanFromSchema must not drop a present-but-empty zone_based_span block")

	out, err := setSpanFromSchema(structVal)
	require.NoError(t, err)
	require.NotNil(t, out)

	spanList := out.([]interface{})
	require.Len(t, spanList, 1)
	spanMap := spanList[0].(map[string]interface{})
	zbs := spanMap["zone_based_span"].([]interface{})
	require.Len(t, zbs, 1)
}
