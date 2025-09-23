package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	api "github.com/vmware/go-vmware-nsxt"
)

func TestResourceNsxtPolicyServiceCreate(t *testing.T) {
	mockProviderClient, m := newMockProviderClient(emptyBody)
	d := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": "foo",
		"description":  "This is a foo service",
	})

	err := resourceNsxtPolicyServiceCreate(d, mockProviderClient)
	fmt.Printf("create results --%v--%v--%v--%v--\n", d.Id(), d.Get("nsx_id"), d.Get("description"), m.cache)

	assert.Equal(t, "foo", gjson.Get(m.cache[0], "display_name").String())
	fmt.Println(err)
	require.Empty(t, err)
}

func TestResourceNsxtPolicyServiceRead(t *testing.T) {
	mockProviderClient, _ := newMockProviderClient(serviceWithId)

	d := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": "foo",
	})
	d.SetId("foo")
	err := resourceNsxtPolicyServiceRead(d, mockProviderClient)

	fmt.Printf("results after --%v--%v--", d.Get("path"), d.Get("description"))
	fmt.Println(err)
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

	fmt.Printf("results after --%v--%v--%v--", d.Get("path"), d.Get("description"), m.cache)
	fmt.Println(err)
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

	fmt.Printf("results after --%v--%v--%v--", d.Get("path"), d.Get("description"), m.cache)
	fmt.Println(err)
	require.Empty(t, err)
	// assert.Equal(t, "foo", gjson.Get(m.cache[0], "display_name").String())
	assert.Equal(t, d.Id(), "foo")
}

func ConstructMockProviderClient() nsxtClients {
	commonConfig := commonProviderConfig{
		RemoteAuth:             false,
		ToleratePartialSuccess: false,
		MaxRetries:             2,
		MinRetryInterval:       0,
		MaxRetryInterval:       0,
		RetryStatusCodes:       []int{404, 400},
		Username:               "username",
		Password:               "password",
	}

	nsxtClient := nsxtClients{
		CommonConfig: commonConfig,
	}
	nsxtClient.NsxtClientConfig = &api.Configuration{
		BasePath:   "/api/v1",
		Scheme:     "https",
		UserAgent:  "terraform-provider-nsxt",
		UserName:   "username",
		Password:   "password",
		RemoteAuth: true,
		Insecure:   true,
	}

	return nsxtClient
}
