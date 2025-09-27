package nsxt

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceNsxtPolicySegmentCreate(t *testing.T) {
	mockProviderClient, m := newMockProviderClient(emptyBody)
	d := schema.TestResourceDataRaw(t, resourceNsxtPolicySegment().Schema, map[string]interface{}{
		"display_name":        "foo",
		"description":         "This is a foo service",
		"transport_zone_path": "/infra/sites/default/enforcement-points/default/transport-zones/1b3a2f36-bfd1-443e-a0f6-4de01abc963e",
		// Add all the 'Required' arguments here depending on the resource.
	})

	err := resourceNsxtPolicySegmentCreate(d, mockProviderClient)

	log.Println("Cache Request body :", m.cache)

	require.Empty(t, err)

	// Test the value you get from the API cache with the actual value.
	// You can add multiple checks here based on the resource (like description, tags, ip_address)
	// The actual value in the below function (3rd argument) can be given in a jq query format. Eg.: "children.0.Segment.display_name"
	// This can be obtained easily by going through the logs above. Search for "Cache Request body".

	assert.Equal(t, "foo", getGjsonString(m.cache[0], "children.0.Segment.display_name"))

}

func TestResourceNsxtPolicySegmentRead(t *testing.T) {
	mockProviderClient, _ := newMockProviderClient(segmentWithId)

	d := schema.TestResourceDataRaw(t, resourceNsxtPolicySegment().Schema, map[string]interface{}{
		"display_name": "foo",
		// Add all the 'Required' arguments here depending on the resource.
	})
	d.SetId("foo")

	err := resourceNsxtPolicySegmentRead(d, mockProviderClient)

	require.Empty(t, err)

	// Test the value you get from the API cache with the actual value.
	// You can add multiple checks here based on the resource (like description, tags, ip_address)

	assert.Equal(t, d.Get("path"), "/infra/tier-1s/cgw/segments/web-tier")
	assert.Equal(t, d.Id(), "foo")
}

func TestResourceNsxtPolicySegmentUpdate(t *testing.T) {
	mockProviderClient, m := newMockProviderClient(emptyBody)

	d := schema.TestResourceDataRaw(t, resourceNsxtPolicySegment().Schema, map[string]interface{}{
		"display_name":        "foo",
		"transport_zone_path": "/infra/sites/default/enforcement-points/default/transport-zones/1b3a2f36-bfd1-443e-a0f6-4de01abc963e",
		// Add all the 'Required' arguments here depending on the resource.
	})
	d.SetId("foo")
	err := resourceNsxtPolicySegmentUpdate(d, mockProviderClient)

	log.Println("Cache Request body :", m.cache)

	require.Empty(t, err)

	// Test the value you get from the API cache with the actual value.
	// You can add multiple checks here based on the resource (like description, tags, ip_address)
	// The actual value in the below function (3rd argument) can be given in a jq query format. Eg.: "children.0.Segment.display_name"
	// This can be obtained easily by going through the logs above. Search for "Cache Request body".

	assert.Equal(t, "foo", getGjsonString(m.cache[0], "children.0.Segment.display_name"))
	assert.Equal(t, d.Id(), "foo")
}

func TestResourceNsxtPolicySegmentDelete(t *testing.T) {
	mockProviderClient, m := newMockProviderClient(emptyBody)

	d := schema.TestResourceDataRaw(t, resourceNsxtPolicySegment().Schema, map[string]interface{}{
		"display_name": "foo",
		// Add all the 'Required' arguments here depending on the resource.
	})
	d.SetId("foo")
	err := resourceNsxtPolicySegmentDelete(d, mockProviderClient)

	log.Println("Cache Request body :", m.cache)

	// Test the value you get from the API cache with the actual value.
	// You can add multiple checks here based on the resource (like description, tags, ip_address)

	require.Empty(t, err)
	assert.Equal(t, d.Id(), "foo")
}
