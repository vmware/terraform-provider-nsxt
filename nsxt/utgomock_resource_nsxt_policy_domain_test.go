//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	domainsapi "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	globalInfraMocks "github.com/vmware/terraform-provider-nsxt/mocks/global_infra"
	globalInfraDomainsMocks "github.com/vmware/terraform-provider-nsxt/mocks/global_infra/domains"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	domainID       = "domain-1"
	domainName     = "domain-fooname"
	domainPath     = "/infra/domains/domain-1"
	domainRevision = int64(1)
	domainSiteName = "default"
)

func newGoMockGlobalProviderClient() nsxtClients {
	c := newGoMockProviderClient()
	c.PolicyGlobalManager = true
	return c
}

func TestMockResourceNsxtPolicyDomainGuard(t *testing.T) {
	res := resourceNsxtPolicyDomain()

	t.Run("Read_fails_when_not_global_manager", func(t *testing.T) {
		util.NsxVersion = "3.2.0"
		defer func() { util.NsxVersion = "" }()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(domainID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDomainRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "local manager")
	})
}

func TestMockResourceNsxtPolicyDomainRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDomainsSDK := globalInfraMocks.NewMockDomainsClient(ctrl)
	mockDomainDeploymentMapsSDK := globalInfraDomainsMocks.NewMockDomainDeploymentMapsClient(ctrl)

	domainWrapper := &infraapi.DomainClientContext{
		Client:     mockDomainsSDK,
		ClientType: utl.Global,
	}
	deploymentMapWrapper := &domainsapi.DomainDeploymentMapClientContext{
		Client:     mockDomainDeploymentMapsSDK,
		ClientType: utl.Global,
	}

	originalDomains := cliDomainsClient
	originalDMs := cliDomainDeploymentMapsClient
	defer func() {
		cliDomainsClient = originalDomains
		cliDomainDeploymentMapsClient = originalDMs
	}()

	cliDomainsClient = func(sessionContext utl.SessionContext, connector client.Connector) *infraapi.DomainClientContext {
		return domainWrapper
	}
	cliDomainDeploymentMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *domainsapi.DomainDeploymentMapClientContext {
		return deploymentMapWrapper
	}

	res := resourceNsxtPolicyDomain()

	t.Run("Read_success", func(t *testing.T) {
		gmDomain := gm_model.Domain{
			DisplayName: &domainName,
			Path:        &domainPath,
			Revision:    &domainRevision,
		}
		mockDomainsSDK.EXPECT().Get(domainID).Return(gmDomain, nil)
		mockDomainDeploymentMapsSDK.EXPECT().List(domainID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			gm_model.DomainDeploymentMapListResult{
				Results: []gm_model.DomainDeploymentMap{
					{DisplayName: &domainSiteName},
				},
			}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(domainID)
		m := newGoMockGlobalProviderClient()
		err := resourceNsxtPolicyDomainRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, domainName, d.Get("display_name"))
	})

	t.Run("Read_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockGlobalProviderClient()
		err := resourceNsxtPolicyDomainRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Domain ID")
	})

	t.Run("Read_fails_when_Get_returns_not_found", func(t *testing.T) {
		mockDomainsSDK.EXPECT().Get(domainID).Return(gm_model.Domain{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(domainID)
		m := newGoMockGlobalProviderClient()
		err := resourceNsxtPolicyDomainRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read_fails_when_Get_returns_error", func(t *testing.T) {
		mockDomainsSDK.EXPECT().Get(domainID).Return(gm_model.Domain{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(domainID)
		m := newGoMockGlobalProviderClient()
		err := resourceNsxtPolicyDomainRead(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyDomainCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDomainsSDK := globalInfraMocks.NewMockDomainsClient(ctrl)
	domainWrapper := &infraapi.DomainClientContext{
		Client:     mockDomainsSDK,
		ClientType: utl.Global,
	}

	originalDomains := cliDomainsClient
	defer func() { cliDomainsClient = originalDomains }()

	cliDomainsClient = func(sessionContext utl.SessionContext, connector client.Connector) *infraapi.DomainClientContext {
		return domainWrapper
	}

	res := resourceNsxtPolicyDomain()

	t.Run("Create_fails_when_resource_already_exists", func(t *testing.T) {
		mockDomainsSDK.EXPECT().Get(domainID).Return(gm_model.Domain{
			DisplayName: &domainName,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":       domainID,
			"display_name": domainName,
		})
		m := newGoMockGlobalProviderClient()
		err := resourceNsxtPolicyDomainCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyDomainDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDomainsSDK := globalInfraMocks.NewMockDomainsClient(ctrl)
	domainWrapper := &infraapi.DomainClientContext{
		Client:     mockDomainsSDK,
		ClientType: utl.Global,
	}

	originalDomains := cliDomainsClient
	defer func() { cliDomainsClient = originalDomains }()

	cliDomainsClient = func(sessionContext utl.SessionContext, connector client.Connector) *infraapi.DomainClientContext {
		return domainWrapper
	}

	res := resourceNsxtPolicyDomain()

	t.Run("Delete_success", func(t *testing.T) {
		mockDomainsSDK.EXPECT().Delete(domainID).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(domainID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDomainDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDomainDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Domain ID")
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockDomainsSDK.EXPECT().Delete(domainID).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(domainID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDomainDelete(d, m)
		require.Error(t, err)
	})
}

func setupDomainMocks(ctrl *gomock.Controller) (*globalInfraMocks.MockDomainsClient, *globalInfraDomainsMocks.MockDomainDeploymentMapsClient, func()) {
	mockDomainsSDK := globalInfraMocks.NewMockDomainsClient(ctrl)
	mockDMSDK := globalInfraDomainsMocks.NewMockDomainDeploymentMapsClient(ctrl)

	domainWrapper := &infraapi.DomainClientContext{Client: mockDomainsSDK, ClientType: utl.Global}
	dmWrapper := &domainsapi.DomainDeploymentMapClientContext{Client: mockDMSDK, ClientType: utl.Global}

	originalDomains := cliDomainsClient
	originalDMs := cliDomainDeploymentMapsClient
	cliDomainsClient = func(sessionContext utl.SessionContext, connector client.Connector) *infraapi.DomainClientContext {
		return domainWrapper
	}
	cliDomainDeploymentMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *domainsapi.DomainDeploymentMapClientContext {
		return dmWrapper
	}
	return mockDomainsSDK, mockDMSDK, func() {
		cliDomainsClient = originalDomains
		cliDomainDeploymentMapsClient = originalDMs
	}
}

func TestUnitNsxt_createChildDomainDeploymentMap(t *testing.T) {
	m := newGoMockProviderClient()

	dataValue, err := createChildDomainDeploymentMap(m, "domain-1", "default")
	require.NoError(t, err)
	assert.NotNil(t, dataValue)
}

func TestUnitNsxt_setDomainStructWithChildrenCreate(t *testing.T) {
	m := newGoMockProviderClient()
	id := "domain-1"
	domain := &model.Domain{Id: &id}

	err := setDomainStructWithChildren(m, domain, []string{"default"}, false)
	require.NoError(t, err)
	assert.Len(t, domain.Children, 1)
}

func TestMockNsxtSetDomainStructWithChildrenUpdate(t *testing.T) {
	t.Run("no stale locations to remove", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockDM, restore := setupDomainMocks(ctrl)
		defer restore()
		mockDM.EXPECT().List(domainID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			gm_model.DomainDeploymentMapListResult{Results: []gm_model.DomainDeploymentMap{{DisplayName: &domainSiteName}}}, nil,
		)

		m := newGoMockProviderClient()
		id := domainID
		domain := &model.Domain{Id: &id}
		err := setDomainStructWithChildren(m, domain, []string{domainSiteName}, true)
		require.NoError(t, err)
		assert.Len(t, domain.Children, 1)
	})

	t.Run("removes stale locations", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockDM, restore := setupDomainMocks(ctrl)
		defer restore()
		staleName := "stale-site"
		mockDM.EXPECT().List(domainID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			gm_model.DomainDeploymentMapListResult{Results: []gm_model.DomainDeploymentMap{
				{DisplayName: &domainSiteName},
				{DisplayName: &staleName},
			}}, nil,
		)

		m := newGoMockProviderClient()
		id := domainID
		domain := &model.Domain{Id: &id}
		err := setDomainStructWithChildren(m, domain, []string{domainSiteName}, true)
		require.NoError(t, err)
		// one child for the new location, one delete-child for the stale location
		assert.Len(t, domain.Children, 2)
	})

	t.Run("List error propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockDM, restore := setupDomainMocks(ctrl)
		defer restore()
		mockDM.EXPECT().List(domainID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			gm_model.DomainDeploymentMapListResult{}, vapiErrors.InternalServerError{},
		)

		m := newGoMockProviderClient()
		id := domainID
		domain := &model.Domain{Id: &id}
		err := setDomainStructWithChildren(m, domain, []string{domainSiteName}, true)
		require.Error(t, err)
	})
}

func TestMockNsxtResourceNsxtPolicyDomainExists(t *testing.T) {
	t.Run("fails when not global manager", func(t *testing.T) {
		exists, err := resourceNsxtPolicyDomainExists(domainID, nil, false)
		require.Error(t, err)
		assert.False(t, exists)
	})

	t.Run("returns false on not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockDomainsSDK, _, restore := setupDomainMocks(ctrl)
		defer restore()
		mockDomainsSDK.EXPECT().Get(domainID).Return(gm_model.Domain{}, vapiErrors.NotFound{})

		exists, err := resourceNsxtPolicyDomainExists(domainID, nil, true)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("propagates other errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockDomainsSDK, _, restore := setupDomainMocks(ctrl)
		defer restore()
		mockDomainsSDK.EXPECT().Get(domainID).Return(gm_model.Domain{}, vapiErrors.InternalServerError{})

		exists, err := resourceNsxtPolicyDomainExists(domainID, nil, true)
		require.Error(t, err)
		assert.False(t, exists)
	})
}

func TestMockResourceNsxtPolicyDomainUpdateEmptyID(t *testing.T) {
	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtPolicyDomainUpdate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Domain ID")
	})
}

// Verify that the model conversion from gm_model to local model works correctly.
func TestDomainModelConversion(t *testing.T) {
	displayName := "test-domain"
	path := "/infra/domains/test-domain"
	revision := int64(2)

	gmDomain := gm_model.Domain{
		DisplayName: &displayName,
		Path:        &path,
		Revision:    &revision,
	}

	lmObj, err := convertModelBindingType(gmDomain, gm_model.DomainBindingType(), model.DomainBindingType())
	require.NoError(t, err)
	obj, ok := lmObj.(model.Domain)
	require.True(t, ok)
	assert.Equal(t, displayName, *obj.DisplayName)
	assert.Equal(t, path, *obj.Path)
}
