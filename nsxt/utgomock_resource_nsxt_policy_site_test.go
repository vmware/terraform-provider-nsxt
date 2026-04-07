// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	globalInfraMocks "github.com/vmware/terraform-provider-nsxt/mocks/global_infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	siteID          = "site-001"
	siteDisplayName = "Test Site"
	siteDescription = "Test site description"
	siteRevision    = int64(1)
	siteType        = gm_model.Site_SITE_TYPE_ONPREM_LM
	sitePath        = "/infra/sites/site-001"
	siteMaxRTT      = int64(250)
)

func siteAPIResponse() gm_model.Site {
	boolTrue := true
	return gm_model.Site{
		Id:                      &siteID,
		DisplayName:             &siteDisplayName,
		Description:             &siteDescription,
		Revision:                &siteRevision,
		Path:                    &sitePath,
		SiteType:                &siteType,
		MaximumRtt:              &siteMaxRTT,
		FailIfRtepMisconfigured: &boolTrue,
		FailIfRttExceeded:       &boolTrue,
	}
}

func minimalSiteData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":               siteDisplayName,
		"description":                siteDescription,
		"nsx_id":                     siteID,
		"site_type":                  siteType,
		"fail_if_rtep_misconfigured": true,
		"fail_if_rtt_exceeded":       true,
		"maximum_rtt":                250,
	}
}

func setupSiteMock(t *testing.T, ctrl *gomock.Controller) (*globalInfraMocks.MockGlobalSitesClient, func()) {
	mockSDK := globalInfraMocks.NewMockGlobalSitesClient(ctrl)
	mockWrapper := &infraapi.SiteClientContext{
		Client:     mockSDK,
		ClientType: utl.Global,
	}

	original := cliSitesClient
	cliSitesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.SiteClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliSitesClient = original }
}

func TestMockResourceNsxtPolicySiteGuard(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	t.Run("Create fails when not global manager", func(t *testing.T) {
		res := resourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSiteData())

		err := resourceNsxtPolicySiteCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}

func TestMockResourceNsxtPolicySiteCreate(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSiteMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(siteID).Return(gm_model.Site{}, notFoundErr),
			mockSDK.EXPECT().Patch(siteID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(siteID).Return(siteAPIResponse(), nil),
		)

		res := resourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSiteData())

		err := resourceNsxtPolicySiteCreate(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
		assert.Equal(t, siteID, d.Id())
		assert.Equal(t, siteDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(siteID).Return(siteAPIResponse(), nil)

		res := resourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSiteData())

		err := resourceNsxtPolicySiteCreate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicySiteRead(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSiteMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(siteID).Return(siteAPIResponse(), nil)

		res := resourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSiteData())
		d.SetId(siteID)

		err := resourceNsxtPolicySiteRead(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
		assert.Equal(t, siteDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(siteID).Return(gm_model.Site{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSiteData())
		d.SetId(siteID)

		err := resourceNsxtPolicySiteRead(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSiteData())

		err := resourceNsxtPolicySiteRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySiteUpdate(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSiteMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Update(siteID, gomock.Any()).Return(siteAPIResponse(), nil)

		res := resourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSiteData())
		d.SetId(siteID)

		err := resourceNsxtPolicySiteUpdate(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSiteData())

		err := resourceNsxtPolicySiteUpdate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySiteDelete(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSiteMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(siteID, nil).Return(nil)

		res := resourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSiteData())
		d.SetId(siteID)

		err := resourceNsxtPolicySiteDelete(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSiteData())

		err := resourceNsxtPolicySiteDelete(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
	})
}
