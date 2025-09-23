package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestResourceNsxtPolicySegmentCreate(t *testing.T) {
    mockProviderClient, m := newMockProviderClient(emptyBody)
    d := schema.TestResourceDataRaw(t, resourceNsxtPolicySegment().Schema, map[string]interface{}{
        "display_name": "foo",
        "description": "This is a foo service",
        "transport_zone_path": "/infra/sites/default/enforcement-points/default/transport-zones/1b3a2f36-bfd1-443e-a0f6-4de01abc963e",
    })

    err := resourceNsxtPolicySegmentCreate(d, mockProviderClient)
    fmt.Printf("create results --%v--%v--%v--%v--%v--\n",  d.Id(), d.Get("nsx_id"), d.Get("description"), m.cache, gjson.Get(m.cache[0], "children.0.Segment.display_name"))
    
    assert.Equal(t, "foo", gjson.Get(m.cache[0], "children.0.Segment.display_name").String())
    fmt.Println(err)
    require.Empty(t, err)
}
// func TestResourceNsxtPolicySegmentRead(t *testing.T) {    
//     mockProviderClient, _ := newMockProviderClient(serviceWithId)
    
//     d := schema.TestResourceDataRaw(t, resourceNsxtPolicySegment().Schema, map[string]interface{}{
//         "display_name": "foo",
//     })
//     d.SetId("foo")
//     err := resourceNsxtPolicySegmentRead(d, mockProviderClient)
    
//     fmt.Printf("results after --%v--%v--", d.Get("path"), d.Get("description"))
//     fmt.Println(err)
//     require.Empty(t, err)
//     assert.Equal(t, d.Get("path"), "/infra/services/my-http")
//     assert.Equal(t, d.Id(), "foo")
// }