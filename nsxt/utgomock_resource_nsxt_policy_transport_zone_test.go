//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/infra/SitesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/SitesClient.go SitesClient
// mockgen -destination=mocks/infra/sites/enforcement_points/TransportZonesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points/TransportZonesClient.go TransportZonesClient

package nsxt

import (
	"errors"
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
	tzID          = "tz-1"
	tzDisplayName = "test-transport-zone"
	tzDescription = "Transport zone description"
	tzPath        = "/infra/sites/default/enforcement-points/default/transport-zones/tz-1"
	tzParentPath  = "/infra/sites/default/enforcement-points/default"
	tzRevision    = int64(1)
	tzSiteID      = "default"
	tzEpID        = "default"
	tzSitePath    = "/infra/sites/default"
	tzType        = model.PolicyTransportZone_TZ_TYPE_OVERLAY_STANDARD
)

func tzAPIResponse() model.PolicyTransportZone {
	return model.PolicyTransportZone{
		Id:          &tzID,
		DisplayName: &tzDisplayName,
		Description: &tzDescription,
		Path:        &tzPath,
		ParentPath:  &tzParentPath,
		Revision:    &tzRevision,
		TzType:      &tzType,
	}
}

func minimalTZData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":   tzDisplayName,
		"description":    tzDescription,
		"transport_type": tzType,
		"site_path":      tzSitePath,
	}
}

func TestMockResourceNsxtPolicyTransportZoneCreate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSitesSDK := inframocks.NewMockSitesClient(ctrl)
	mockSitesWrapper := &cliinfra.SiteClientContext{
		Client:     mockSitesSDK,
		ClientType: utl.Local,
	}

	mockTZSDK := epmocks.NewMockTransportZonesClient(ctrl)
	mockTZWrapper := &enforcementpoints.PolicyTransportZoneClientContext{
		Client:     mockTZSDK,
		ClientType: utl.Local,
	}

	originalSites := cliSitesClient
	originalTZ := cliTransportZonesClient
	defer func() {
		cliSitesClient = originalSites
		cliTransportZonesClient = originalTZ
	}()
	cliSitesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SiteClientContext {
		return mockSitesWrapper
	}
	cliTransportZonesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyTransportZoneClientContext {
		return mockTZWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(gomock.Any()).Return(model.Site{}, nil)
		mockTZSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.PolicyTransportZone{}, vapiErrors.NotFound{})
		mockTZSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(model.PolicyTransportZone{}, nil)
		mockTZSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(tzAPIResponse(), nil)

		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTZData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, tzDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(gomock.Any()).Return(model.Site{}, nil)
		mockTZSDK.EXPECT().Get(gomock.Any(), gomock.Any(), tzID).Return(model.PolicyTransportZone{Id: &tzID}, nil)

		res := resourceNsxtPolicyTransportZone()
		data := minimalTZData()
		data["nsx_id"] = tzID
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when site not found", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(gomock.Any()).Return(model.Site{}, errors.New("site not found"))

		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTZData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneCreate(d, m)
		require.Error(t, err)
	})

	t.Run("Create fails when Patch returns error", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(gomock.Any()).Return(model.Site{}, nil)
		mockTZSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.PolicyTransportZone{}, vapiErrors.NotFound{})
		mockTZSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(model.PolicyTransportZone{}, errors.New("API error"))

		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTZData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtPolicyTransportZoneRead(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTZSDK := epmocks.NewMockTransportZonesClient(ctrl)
	mockTZWrapper := &enforcementpoints.PolicyTransportZoneClientContext{
		Client:     mockTZSDK,
		ClientType: utl.Local,
	}

	originalTZ := cliTransportZonesClient
	defer func() { cliTransportZonesClient = originalTZ }()
	cliTransportZonesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyTransportZoneClientContext {
		return mockTZWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockTZSDK.EXPECT().Get(tzSiteID, tzEpID, tzID).Return(tzAPIResponse(), nil)

		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         tzSitePath,
			"enforcement_point": tzEpID,
		})
		d.SetId(tzID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, tzDisplayName, d.Get("display_name"))
		assert.Equal(t, tzDescription, d.Get("description"))
		assert.Equal(t, tzSitePath, d.Get("site_path"))
		assert.Equal(t, int(tzRevision), d.Get("revision"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         tzSitePath,
			"enforcement_point": tzEpID,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining Resource ID")
	})

	t.Run("Read clears ID when not found", func(t *testing.T) {
		mockTZSDK.EXPECT().Get(tzSiteID, tzEpID, tzID).Return(model.PolicyTransportZone{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         tzSitePath,
			"enforcement_point": tzEpID,
		})
		d.SetId(tzID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyTransportZoneUpdate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTZSDK := epmocks.NewMockTransportZonesClient(ctrl)
	mockTZWrapper := &enforcementpoints.PolicyTransportZoneClientContext{
		Client:     mockTZSDK,
		ClientType: utl.Local,
	}

	originalTZ := cliTransportZonesClient
	defer func() { cliTransportZonesClient = originalTZ }()
	cliTransportZonesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyTransportZoneClientContext {
		return mockTZWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockTZSDK.EXPECT().Patch(tzSiteID, tzEpID, tzID, gomock.Any()).Return(model.PolicyTransportZone{}, nil)
		mockTZSDK.EXPECT().Get(tzSiteID, tzEpID, tzID).Return(tzAPIResponse(), nil)

		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTZData())
		d.Set("enforcement_point", tzEpID)
		d.SetId(tzID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, tzDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTZData())
		d.Set("enforcement_point", tzEpID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining Resource ID")
	})
}

func TestMockResourceNsxtPolicyTransportZoneDelete(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTZSDK := epmocks.NewMockTransportZonesClient(ctrl)
	mockTZWrapper := &enforcementpoints.PolicyTransportZoneClientContext{
		Client:     mockTZSDK,
		ClientType: utl.Local,
	}

	originalTZ := cliTransportZonesClient
	defer func() { cliTransportZonesClient = originalTZ }()
	cliTransportZonesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyTransportZoneClientContext {
		return mockTZWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockTZSDK.EXPECT().Delete(tzSiteID, tzEpID, tzID).Return(nil)

		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         tzSitePath,
			"enforcement_point": tzEpID,
		})
		d.SetId(tzID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         tzSitePath,
			"enforcement_point": tzEpID,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining Resource ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockTZSDK.EXPECT().Delete(tzSiteID, tzEpID, tzID).Return(errors.New("API error"))

		res := resourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         tzSitePath,
			"enforcement_point": tzEpID,
		})
		d.SetId(tzID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTransportZoneDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}
