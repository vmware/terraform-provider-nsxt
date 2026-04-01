//go:build unittest

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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	policyHAProfileID       = "ha-profile-1"
	policyHAProfileSitePath = "/infra/sites/default"
	policyHAProfileSiteID   = "default"
	policyHAProfileEPID     = "default"
	policyHAProfileName     = "ha-profile-fooname"
	policyHAProfileRevision = int64(1)
	policyHAProfilePath     = "/infra/sites/default/enforcement-points/default/edge-cluster-high-availability-profiles/ha-profile-1"
)

func TestMockResourceNsxtPolicyEdgeHighAvailabilityProfileCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHAProfilesSDK := epmocks.NewMockEdgeClusterHighAvailabilityProfilesClient(ctrl)
	mockSitesSDK := inframocks.NewMockSitesClient(ctrl)

	haProfileWrapper := &enforcementpoints.PolicyEdgeHighAvailabilityProfileClientContext{
		Client:     mockHAProfilesSDK,
		ClientType: utl.Local,
	}
	siteWrapper := &cliinfra.SiteClientContext{
		Client:     mockSitesSDK,
		ClientType: utl.Local,
	}

	originalHAProfiles := cliEdgeClusterHighAvailabilityProfilesClient
	originalSites := cliSitesClient
	defer func() {
		cliEdgeClusterHighAvailabilityProfilesClient = originalHAProfiles
		cliSitesClient = originalSites
	}()
	cliEdgeClusterHighAvailabilityProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeHighAvailabilityProfileClientContext {
		return haProfileWrapper
	}
	cliSitesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SiteClientContext {
		return siteWrapper
	}

	res := resourceNsxtPolicyEdgeHighAvailabilityProfile()

	t.Run("Create_success", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(policyHAProfileSiteID).Return(model.Site{}, nil)
		mockHAProfilesSDK.EXPECT().Get(policyHAProfileSiteID, policyHAProfileEPID, gomock.Any()).Return(model.PolicyEdgeHighAvailabilityProfile{}, vapiErrors.NotFound{})
		mockHAProfilesSDK.EXPECT().Patch(policyHAProfileSiteID, policyHAProfileEPID, gomock.Any(), gomock.Any()).Return(nil)
		mockHAProfilesSDK.EXPECT().Get(policyHAProfileSiteID, policyHAProfileEPID, gomock.Any()).Return(model.PolicyEdgeHighAvailabilityProfile{
			DisplayName: &policyHAProfileName,
			Path:        &policyHAProfilePath,
			Revision:    &policyHAProfileRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":      policyHAProfileName,
			"site_path":         policyHAProfileSitePath,
			"enforcement_point": policyHAProfileEPID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeHighAvailabilityProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, policyHAProfileName, d.Get("display_name"))
	})

	t.Run("Create_fails_when_already_exists", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(policyHAProfileSiteID).Return(model.Site{}, nil)
		mockHAProfilesSDK.EXPECT().Get(policyHAProfileSiteID, policyHAProfileEPID, policyHAProfileID).Return(model.PolicyEdgeHighAvailabilityProfile{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":            policyHAProfileID,
			"site_path":         policyHAProfileSitePath,
			"enforcement_point": policyHAProfileEPID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeHighAvailabilityProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_site_path_invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         "/infra/no-sites/default",
			"enforcement_point": policyHAProfileEPID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeHighAvailabilityProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Site ID")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(policyHAProfileSiteID).Return(model.Site{}, nil)
		mockHAProfilesSDK.EXPECT().Get(policyHAProfileSiteID, policyHAProfileEPID, gomock.Any()).Return(model.PolicyEdgeHighAvailabilityProfile{}, vapiErrors.NotFound{})
		mockHAProfilesSDK.EXPECT().Patch(policyHAProfileSiteID, policyHAProfileEPID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         policyHAProfileSitePath,
			"enforcement_point": policyHAProfileEPID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeHighAvailabilityProfileCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEdgeHighAvailabilityProfileRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHAProfilesSDK := epmocks.NewMockEdgeClusterHighAvailabilityProfilesClient(ctrl)
	haProfileWrapper := &enforcementpoints.PolicyEdgeHighAvailabilityProfileClientContext{
		Client:     mockHAProfilesSDK,
		ClientType: utl.Local,
	}

	originalHAProfiles := cliEdgeClusterHighAvailabilityProfilesClient
	defer func() { cliEdgeClusterHighAvailabilityProfilesClient = originalHAProfiles }()
	cliEdgeClusterHighAvailabilityProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeHighAvailabilityProfileClientContext {
		return haProfileWrapper
	}

	res := resourceNsxtPolicyEdgeHighAvailabilityProfile()

	t.Run("Read_success", func(t *testing.T) {
		bfdProbeInterval := int64(500)
		mockHAProfilesSDK.EXPECT().Get(policyHAProfileSiteID, policyHAProfileEPID, policyHAProfileID).Return(model.PolicyEdgeHighAvailabilityProfile{
			DisplayName:      &policyHAProfileName,
			Path:             &policyHAProfilePath,
			Revision:         &policyHAProfileRevision,
			BfdProbeInterval: &bfdProbeInterval,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         policyHAProfileSitePath,
			"enforcement_point": policyHAProfileEPID,
		})
		d.SetId(policyHAProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeHighAvailabilityProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, policyHAProfileName, d.Get("display_name"))
	})

	t.Run("Read_clears_id_when_not_found", func(t *testing.T) {
		mockHAProfilesSDK.EXPECT().Get(policyHAProfileSiteID, policyHAProfileEPID, policyHAProfileID).Return(model.PolicyEdgeHighAvailabilityProfile{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         policyHAProfileSitePath,
			"enforcement_point": policyHAProfileEPID,
		})
		d.SetId(policyHAProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeHighAvailabilityProfileRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyEdgeHighAvailabilityProfileUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHAProfilesSDK := epmocks.NewMockEdgeClusterHighAvailabilityProfilesClient(ctrl)
	haProfileWrapper := &enforcementpoints.PolicyEdgeHighAvailabilityProfileClientContext{
		Client:     mockHAProfilesSDK,
		ClientType: utl.Local,
	}

	originalHAProfiles := cliEdgeClusterHighAvailabilityProfilesClient
	defer func() { cliEdgeClusterHighAvailabilityProfilesClient = originalHAProfiles }()
	cliEdgeClusterHighAvailabilityProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeHighAvailabilityProfileClientContext {
		return haProfileWrapper
	}

	res := resourceNsxtPolicyEdgeHighAvailabilityProfile()

	t.Run("Update_success", func(t *testing.T) {
		mockHAProfilesSDK.EXPECT().Update(policyHAProfileSiteID, policyHAProfileEPID, policyHAProfileID, gomock.Any()).Return(model.PolicyEdgeHighAvailabilityProfile{}, nil)
		mockHAProfilesSDK.EXPECT().Get(policyHAProfileSiteID, policyHAProfileEPID, policyHAProfileID).Return(model.PolicyEdgeHighAvailabilityProfile{
			DisplayName: &policyHAProfileName,
			Path:        &policyHAProfilePath,
			Revision:    &policyHAProfileRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":      policyHAProfileName,
			"site_path":         policyHAProfileSitePath,
			"enforcement_point": policyHAProfileEPID,
		})
		d.SetId(policyHAProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeHighAvailabilityProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_Update_returns_error", func(t *testing.T) {
		mockHAProfilesSDK.EXPECT().Update(policyHAProfileSiteID, policyHAProfileEPID, policyHAProfileID, gomock.Any()).Return(model.PolicyEdgeHighAvailabilityProfile{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         policyHAProfileSitePath,
			"enforcement_point": policyHAProfileEPID,
		})
		d.SetId(policyHAProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeHighAvailabilityProfileUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEdgeHighAvailabilityProfileDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHAProfilesSDK := epmocks.NewMockEdgeClusterHighAvailabilityProfilesClient(ctrl)
	haProfileWrapper := &enforcementpoints.PolicyEdgeHighAvailabilityProfileClientContext{
		Client:     mockHAProfilesSDK,
		ClientType: utl.Local,
	}

	originalHAProfiles := cliEdgeClusterHighAvailabilityProfilesClient
	defer func() { cliEdgeClusterHighAvailabilityProfilesClient = originalHAProfiles }()
	cliEdgeClusterHighAvailabilityProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeHighAvailabilityProfileClientContext {
		return haProfileWrapper
	}

	res := resourceNsxtPolicyEdgeHighAvailabilityProfile()

	t.Run("Delete_success", func(t *testing.T) {
		mockHAProfilesSDK.EXPECT().Delete(policyHAProfileSiteID, policyHAProfileEPID, policyHAProfileID).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         policyHAProfileSitePath,
			"enforcement_point": policyHAProfileEPID,
		})
		d.SetId(policyHAProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeHighAvailabilityProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockHAProfilesSDK.EXPECT().Delete(policyHAProfileSiteID, policyHAProfileEPID, policyHAProfileID).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         policyHAProfileSitePath,
			"enforcement_point": policyHAProfileEPID,
		})
		d.SetId(policyHAProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeHighAvailabilityProfileDelete(d, m)
		require.Error(t, err)
	})
}
