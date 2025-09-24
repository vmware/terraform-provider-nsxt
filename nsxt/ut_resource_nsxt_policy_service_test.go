package nsxt

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestResourceNsxtPolicyServiceCreate(t *testing.T) {
	mockProviderClient, m := newMockProviderClient(emptyBody)
	d := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": "foo",
		"description":  "This is a foo service",
	})

	err := resourceNsxtPolicyServiceCreate(d, mockProviderClient)
	
	log.Println("Cache Request body :", m.cache)
	
	require.Empty(t, err)
	
	assert.Equal(t, "foo", getGjsonString(m.cache[0], "display_name"))
}

func TestResourceNsxtPolicyServiceRead(t *testing.T) {
	mockProviderClient, _ := newMockProviderClient(serviceWithId)

	d := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": "foo",
	})
	d.SetId("foo")
	err := resourceNsxtPolicyServiceRead(d, mockProviderClient)

	require.Empty(t, err)
	assert.Equal(t, d.Get("path"), "/infra/services/my-http")
	assert.Equal(t, d.Id(), "foo")
}

func TestResourceNsxtPolicyServiceUpdate(t *testing.T) {
	mockProviderClient, m := newMockProviderClient(emptyBody)

	d := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": "foo",
	})
	d.SetId("foo")
	err := resourceNsxtPolicyServiceUpdate(d, mockProviderClient)

	log.Println("Cache Request body :", m.cache)

	require.Empty(t, err)
	assert.Equal(t, "foo", gjson.Get(m.cache[0], "display_name").String())
	assert.Equal(t, d.Id(), "foo")
}

func TestResourceNsxtPolicyServiceDelete(t *testing.T) {
	mockProviderClient, m := newMockProviderClient(emptyBody)

	d := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": "foo",
	})
	d.SetId("foo")
	err := resourceNsxtPolicyServiceDelete(d, mockProviderClient)
	
	log.Println("Cache Request body :", m.cache)

	require.Empty(t, err)
	assert.Equal(t, d.Id(), "foo")
}

