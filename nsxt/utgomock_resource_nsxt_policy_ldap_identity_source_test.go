//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/aaa/LdapIdentitySourcesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/aaa/LdapIdentitySourcesClient.go LdapIdentitySourcesClient

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	"github.com/vmware/terraform-provider-nsxt/api/aaa"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ldapmocks "github.com/vmware/terraform-provider-nsxt/mocks/aaa"
)

func minimalLDAPData() map[string]interface{} {
	return map[string]interface{}{
		"nsx_id":      "ldap-src-1",
		"description": "Test LDAP Source",
		"type":        activeDirectoryType,
		"domain_name": "corp.example.com",
		"base_dn":     "dc=corp,dc=example,dc=com",
	}
}

func setupLdapIdentitySourceMock(t *testing.T) *ldapmocks.MockLdapIdentitySourcesClient {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSDK := ldapmocks.NewMockLdapIdentitySourcesClient(ctrl)
	mockWrapper := &aaa.LdapIdentitySourcesClientContext{Client: mockSDK, ClientType: utl.Local}

	originalCli := cliLdapIdentitySourcesClient
	t.Cleanup(func() { cliLdapIdentitySourcesClient = originalCli })
	cliLdapIdentitySourcesClient = func(sessionContext utl.SessionContext, connector client.Connector) *aaa.LdapIdentitySourcesClientContext {
		return mockWrapper
	}

	return mockSDK
}

// ldapStructValue builds a *data.StructValue for an ActiveDirectoryIdentitySource, mimicking
// what the real SDK Get()/Update() calls would return.
func ldapStructValue(t *testing.T, id, description, domainName, baseDn string, revision int64, servers []nsxModel.IdentitySourceLdapServer) *data.StructValue {
	converter := bindings.NewTypeConverter()
	obj := nsxModel.ActiveDirectoryIdentitySource{
		Id:           &id,
		Description:  &description,
		Revision:     &revision,
		DomainName:   &domainName,
		BaseDn:       &baseDn,
		LdapServers:  servers,
		ResourceType: nsxModel.LdapIdentitySource_RESOURCE_TYPE_ACTIVEDIRECTORYIDENTITYSOURCE,
	}
	dataValue, errs := converter.ConvertToVapi(obj, nsxModel.ActiveDirectoryIdentitySourceBindingType())
	require.Nil(t, errs)
	return dataValue.(*data.StructValue)
}

func ldapEnabledServer(url string) nsxModel.IdentitySourceLdapServer {
	enabled := true
	useStarttls := false
	bindIdentity := "cn=admin"
	password := "secret"
	return nsxModel.IdentitySourceLdapServer{
		BindIdentity: &bindIdentity,
		Enabled:      &enabled,
		Password:     &password,
		Url:          &url,
		UseStarttls:  &useStarttls,
	}
}

func ldapDisabledServer(url string) nsxModel.IdentitySourceLdapServer {
	enabled := false
	useStarttls := false
	bindIdentity := "cn=admin"
	password := "secret"
	return nsxModel.IdentitySourceLdapServer{
		BindIdentity: &bindIdentity,
		Enabled:      &enabled,
		Password:     &password,
		Url:          &url,
		UseStarttls:  &useStarttls,
	}
}

func ldapDataWithServer(server map[string]interface{}) map[string]interface{} {
	data := minimalLDAPData()
	data["ldap_server"] = []interface{}{server}
	return data
}

func TestMockResourceNsxtPolicyLdapIdentitySourceRead(t *testing.T) {
	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLDAPData())

		err := resourceNsxtPolicyLdapIdentitySourceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read success sets all fields", func(t *testing.T) {
		mockSDK := setupLdapIdentitySourceMock(t)

		servers := []nsxModel.IdentitySourceLdapServer{ldapEnabledServer("ldap://server1.corp.example.com")}
		structValue := ldapStructValue(t, "ldap-src-1", "Test LDAP Source", "corp.example.com", "dc=corp,dc=example,dc=com", 2, servers)

		mockSDK.EXPECT().Get("ldap-src-1").Return(structValue, nil)

		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("ldap-src-1")

		err := resourceNsxtPolicyLdapIdentitySourceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "ldap-src-1", d.Get("nsx_id"))
		assert.Equal(t, "Test LDAP Source", d.Get("description"))
		assert.Equal(t, 2, d.Get("revision"))
		assert.Equal(t, activeDirectoryType, d.Get("type"))
		assert.Equal(t, "corp.example.com", d.Get("domain_name"))
		assert.Equal(t, "dc=corp,dc=example,dc=com", d.Get("base_dn"))

		ldapServers := d.Get("ldap_server").([]interface{})
		require.Len(t, ldapServers, 1)
		serverMap := ldapServers[0].(map[string]interface{})
		assert.Equal(t, "ldap://server1.corp.example.com", serverMap["url"])
		assert.True(t, serverMap["enabled"].(bool))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK := setupLdapIdentitySourceMock(t)

		mockSDK.EXPECT().Get("ldap-src-1").Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("ldap-src-1")

		err := resourceNsxtPolicyLdapIdentitySourceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})
}

func TestMockResourceNsxtPolicyLdapIdentitySourceCreate(t *testing.T) {
	t.Run("Create success probes enabled server then updates", func(t *testing.T) {
		mockSDK := setupLdapIdentitySourceMock(t)

		servers := []nsxModel.IdentitySourceLdapServer{ldapEnabledServer("ldap://server1.corp.example.com")}
		structValue := ldapStructValue(t, "ldap-src-1", "Test LDAP Source", "corp.example.com", "dc=corp,dc=example,dc=com", 1, servers)

		successResult := nsxModel.IdentitySourceLdapServerProbeResult_RESULT_SUCCESS
		url := "ldap://server1.corp.example.com"

		// getOrGenerateID's presence check: nsx_id is set, so it calls Get(id) first to
		// verify the ID isn't already taken.
		mockSDK.EXPECT().
			Get("ldap-src-1").
			Return(nil, vapiErrors.NotFound{})
		mockSDK.EXPECT().
			Probeidentitysource(gomock.Any()).
			Return(nsxModel.LdapIdentitySourceProbeResults{
				Results: []nsxModel.IdentitySourceLdapServerProbeResult{
					{Result: &successResult, Url: &url},
				},
			}, nil)
		mockSDK.EXPECT().
			Update("ldap-src-1", gomock.Any()).
			Return(structValue, nil)
		mockSDK.EXPECT().
			Get("ldap-src-1").
			Return(structValue, nil)

		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, ldapDataWithServer(map[string]interface{}{
			"url":     "ldap://server1.corp.example.com",
			"enabled": true,
		}))

		err := resourceNsxtPolicyLdapIdentitySourceCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "ldap-src-1", d.Id())
		assert.Equal(t, "ldap-src-1", d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockSDK := setupLdapIdentitySourceMock(t)

		existing := ldapStructValue(t, "ldap-src-1", "existing", "corp.example.com", "dc=corp,dc=example,dc=com", 1, nil)
		mockSDK.EXPECT().Get("ldap-src-1").Return(existing, nil)

		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, ldapDataWithServer(map[string]interface{}{
			"url":     "ldap://server1.corp.example.com",
			"enabled": false,
		}))

		err := resourceNsxtPolicyLdapIdentitySourceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create does not probe when no server is enabled", func(t *testing.T) {
		mockSDK := setupLdapIdentitySourceMock(t)

		servers := []nsxModel.IdentitySourceLdapServer{ldapDisabledServer("ldap://server1.corp.example.com")}
		structValue := ldapStructValue(t, "ldap-src-1", "Test LDAP Source", "corp.example.com", "dc=corp,dc=example,dc=com", 1, servers)

		// No Probeidentitysource EXPECT() is set: if the code under test called it
		// unexpectedly, gomock would fail the test since the call is not registered.
		mockSDK.EXPECT().
			Get("ldap-src-1").
			Return(nil, vapiErrors.NotFound{})
		mockSDK.EXPECT().
			Update("ldap-src-1", gomock.Any()).
			Return(structValue, nil)
		mockSDK.EXPECT().
			Get("ldap-src-1").
			Return(structValue, nil)

		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, ldapDataWithServer(map[string]interface{}{
			"url":     "ldap://server1.corp.example.com",
			"enabled": false,
		}))

		err := resourceNsxtPolicyLdapIdentitySourceCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Create fails when probe reports failure", func(t *testing.T) {
		mockSDK := setupLdapIdentitySourceMock(t)

		failureResult := nsxModel.IdentitySourceLdapServerProbeResult_RESULT_FAILURE
		url := "ldap://server1.corp.example.com"
		errType := "CONNECTION_REFUSED"

		mockSDK.EXPECT().
			Get("ldap-src-1").
			Return(nil, vapiErrors.NotFound{})
		mockSDK.EXPECT().
			Probeidentitysource(gomock.Any()).
			Return(nsxModel.LdapIdentitySourceProbeResults{
				Results: []nsxModel.IdentitySourceLdapServerProbeResult{
					{
						Result: &failureResult,
						Url:    &url,
						Errors: []nsxModel.LdapProbeError{{ErrorType: &errType}},
					},
				},
			}, nil)

		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, ldapDataWithServer(map[string]interface{}{
			"url":     "ldap://server1.corp.example.com",
			"enabled": true,
		}))

		err := resourceNsxtPolicyLdapIdentitySourceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "probe failed")
	})
}

func TestMockResourceNsxtPolicyLdapIdentitySourceUpdate(t *testing.T) {
	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLDAPData())

		err := resourceNsxtPolicyLdapIdentitySourceUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		mockSDK := setupLdapIdentitySourceMock(t)

		servers := []nsxModel.IdentitySourceLdapServer{ldapEnabledServer("ldap://server2.corp.example.com")}
		structValue := ldapStructValue(t, "ldap-src-1", "updated description", "corp.example.com", "dc=corp,dc=example,dc=com", 2, servers)

		successResult := nsxModel.IdentitySourceLdapServerProbeResult_RESULT_SUCCESS
		url := "ldap://server2.corp.example.com"

		mockSDK.EXPECT().
			Probeidentitysource(gomock.Any()).
			Return(nsxModel.LdapIdentitySourceProbeResults{
				Results: []nsxModel.IdentitySourceLdapServerProbeResult{
					{Result: &successResult, Url: &url},
				},
			}, nil)
		mockSDK.EXPECT().
			Update("ldap-src-1", gomock.Any()).
			Return(structValue, nil)
		mockSDK.EXPECT().
			Get("ldap-src-1").
			Return(structValue, nil)

		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, ldapDataWithServer(map[string]interface{}{
			"url":     "ldap://server2.corp.example.com",
			"enabled": true,
		}))
		d.SetId("ldap-src-1")

		err := resourceNsxtPolicyLdapIdentitySourceUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "updated description", d.Get("description"))
	})

	t.Run("Update fails when Update API errors", func(t *testing.T) {
		mockSDK := setupLdapIdentitySourceMock(t)

		mockSDK.EXPECT().
			Update("ldap-src-1", gomock.Any()).
			Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, ldapDataWithServer(map[string]interface{}{
			"url":     "ldap://server1.corp.example.com",
			"enabled": false,
		}))
		d.SetId("ldap-src-1")

		err := resourceNsxtPolicyLdapIdentitySourceUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLdapIdentitySourceDelete(t *testing.T) {
	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLDAPData())

		err := resourceNsxtPolicyLdapIdentitySourceDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Delete success", func(t *testing.T) {
		mockSDK := setupLdapIdentitySourceMock(t)

		mockSDK.EXPECT().Delete("ldap-src-1").Return(nil)

		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLDAPData())
		d.SetId("ldap-src-1")

		err := resourceNsxtPolicyLdapIdentitySourceDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockSDK := setupLdapIdentitySourceMock(t)

		mockSDK.EXPECT().Delete("ldap-src-1").Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLDAPData())
		d.SetId("ldap-src-1")

		err := resourceNsxtPolicyLdapIdentitySourceDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
