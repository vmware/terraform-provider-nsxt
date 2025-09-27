package nsxt

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	servicesmock "github.com/vmware/terraform-provider-nsxt/mocks/servicemocks"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var (
	displayName = "fooname"
	description = "this is a mock"
	path        = "/infra/services/my-http"
	revision    = int64(1)
)

func newGoMockProviderClient() nsxtClients {
	mockProviderClient := constructMockProviderClient()
	mockProviderClient.NsxtClientConfig.HTTPClient = &http.Client{}
	mockProviderClient.PolicyHTTPClient = &http.Client{}
	return mockProviderClient
}

func TestMockResourceNsxtPolicyServiceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	nsxClients := newGoMockProviderClient()

	mockServicesSDK := servicesmock.NewMockServicesClient(ctrl)
	tmpServiceModel := model.Service{ // Replace this with the binding functions to convert the string literal json to struct
		DisplayName: &displayName,
		Description: &description,
	}

	mockServicesSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockServicesSDK.EXPECT().Get(gomock.Any()).Return(tmpServiceModel, nil).Times(2)
	mockWrapper := &cliinfra.ServiceClientContext{
		Client:     mockServicesSDK,
		ClientType: utl.Local,
	}

	cli = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ServiceClientContext {
		return mockWrapper
	}
	getNSXVersionFunc = func(connector client.Connector) (string, error) {
		return "9.1.0", nil
	}
	defer func() {
		cli = cliinfra.NewServicesClient
		getNSXVersionFunc = getNSXVersion
	}()

	// Positive test case
	d := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": displayName,
		"description":  description,
	})

	err := resourceNsxtPolicyServiceCreate(d, nsxClients)

	require.Empty(t, err)
	assert.Equal(t, displayName, d.Get("display_name"))
	assert.Equal(t, description, d.Get("description"))

	// Negative test case
	dn := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": displayName,
		"description":  description,
	})
	testFakeId := "foo"
	dn.Set("nsx_id", testFakeId) // set this to get an already exists error. This can be set for a negative test case

	err = resourceNsxtPolicyServiceCreate(dn, nsxClients)
	require.NotEmpty(t, err)
	require.EqualError(t, err, fmt.Sprintf("Resource with id %s already exists", testFakeId))
}

func TestMockResourceNsxtPolicyServiceRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	nsxClients := newGoMockProviderClient()
	mockServicesSDK := servicesmock.NewMockServicesClient(ctrl)

	tmpServiceModel := model.Service{ // Replace this with the binding functions to convert the string literal json to struct
		DisplayName: &displayName,
		Description: &description,
		Path:        &path,
		Revision:    &revision,
	}

	mockServicesSDK.EXPECT().Get("foo").Return(tmpServiceModel, nil)
	mockWrapper := &cliinfra.ServiceClientContext{
		Client:     mockServicesSDK,
		ClientType: utl.Local,
	}

	cli = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ServiceClientContext {
		return mockWrapper
	}
	getNSXVersionFunc = func(connector client.Connector) (string, error) {
		return "9.1.0", nil
	}
	defer func() {
		cli = cliinfra.NewServicesClient
		getNSXVersionFunc = getNSXVersion
	}()

	d := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": "foo",
	})
	d.SetId("foo")

	err := resourceNsxtPolicyServiceRead(d, nsxClients)

	require.Empty(t, err)
	assert.Equal(t, displayName, d.Get("display_name"))
	assert.Equal(t, description, d.Get("description"))
	assert.Equal(t, path, d.Get("path"))
}

func TestMockResourceNsxtPolicyServiceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	nsxClients := newGoMockProviderClient()

	mockServicesSDK := servicesmock.NewMockServicesClient(ctrl)
	tmpServiceModel := model.Service{ // Replace this with the binding functions to convert the string literal json to struct
		DisplayName: &displayName,
		Description: &description,
	}

	mockServicesSDK.EXPECT().Update("foo", gomock.Any()).Return(tmpServiceModel, nil).Times(1)
	mockServicesSDK.EXPECT().Get("foo").Return(tmpServiceModel, nil).Times(1)
	mockWrapper := &cliinfra.ServiceClientContext{
		Client:     mockServicesSDK,
		ClientType: utl.Local,
	}

	cli = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ServiceClientContext {
		return mockWrapper
	}
	getNSXVersionFunc = func(connector client.Connector) (string, error) {
		return "9.1.0", nil
	}
	defer func() {
		cli = cliinfra.NewServicesClient
		getNSXVersionFunc = getNSXVersion
	}()

	// Positive test case
	d := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": displayName,
		"description":  description,
	})
	d.SetId("foo")

	err := resourceNsxtPolicyServiceUpdate(d, nsxClients)

	require.Empty(t, err)
	assert.Equal(t, displayName, d.Get("display_name"))
	assert.Equal(t, description, d.Get("description"))

	// Negative test case
	dn := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": displayName,
		"description":  description,
	})

	err = resourceNsxtPolicyServiceUpdate(dn, nsxClients)
	require.NotEmpty(t, err)
	require.EqualError(t, err, "Error obtaining service id")
}

func TestMockResourceNsxtPolicyServiceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	nsxClients := newGoMockProviderClient()
	mockServicesSDK := servicesmock.NewMockServicesClient(ctrl)

	mockServicesSDK.EXPECT().Delete("foo").Return(nil).Times(1)
	mockWrapper := &cliinfra.ServiceClientContext{
		Client:     mockServicesSDK,
		ClientType: utl.Local,
	}

	cli = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ServiceClientContext {
		return mockWrapper
	}
	getNSXVersionFunc = func(connector client.Connector) (string, error) {
		return "9.1.0", nil
	}
	defer func() {
		cli = cliinfra.NewServicesClient
		getNSXVersionFunc = getNSXVersion
	}()

	d := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": "foo",
	})
	d.SetId("foo")

	err := resourceNsxtPolicyServiceDelete(d, nsxClients)
	require.Empty(t, err)

	// Negative test case
	dn := schema.TestResourceDataRaw(t, resourceNsxtPolicyService().Schema, map[string]interface{}{
		"display_name": "foo",
	})

	err = resourceNsxtPolicyServiceDelete(dn, nsxClients)
	require.NotEmpty(t, err)
	require.EqualError(t, err, "Error obtaining service id")

}
