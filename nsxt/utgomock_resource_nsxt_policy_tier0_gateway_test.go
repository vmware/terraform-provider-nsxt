// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run mockgen for Tier0sClient and LocaleServicesClient
// in api/infra and api/infra/tier_0s.

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apipkg "github.com/vmware/terraform-provider-nsxt/api"
	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	tier0localeservices "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	t0mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	localeServicesMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	t0GatewayID       = "t0-gw-1"
	t0DisplayName     = "tier0-fooname"
	t0Description     = "tier0 mock description"
	t0Path            = "/infra/tier-0s/t0-gw-1"
	t0Revision        = int64(1)
	t0FailoverMode    = "PREEMPTIVE"
	t0HaMode          = "ACTIVE_STANDBY"
	t0DisableFirewall = false
)

func TestMockResourceNsxtPolicyTier0GatewayRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTier0sSDK := t0mocks.NewMockTier0sClient(ctrl)
	mockLocaleServicesSDK := localeServicesMocks.NewMockLocaleServicesClient(ctrl)

	tier0Wrapper := &cliinfra.Tier0ClientContext{
		Client:     mockTier0sSDK,
		ClientType: utl.Local,
	}
	localeServicesWrapper := &tier0localeservices.LocaleServicesClientContext{
		Client:     mockLocaleServicesSDK,
		ClientType: utl.Local,
	}

	originalTier0s := cliTier0sClient
	originalLocaleServices := cliTier0LocaleServicesClient
	defer func() {
		cliTier0sClient = originalTier0s
		cliTier0LocaleServicesClient = originalLocaleServices
	}()

	cliTier0sClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}
	cliTier0LocaleServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0localeservices.LocaleServicesClientContext {
		return localeServicesWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(t0GatewayID).Return(model.Tier0{
			DisplayName:     &t0DisplayName,
			Description:     &t0Description,
			Path:            &t0Path,
			Revision:        &t0Revision,
			FailoverMode:    &t0FailoverMode,
			HaMode:          &t0HaMode,
			DisableFirewall: &t0DisableFirewall,
		}, nil)
		resultCount := int64(0)
		mockLocaleServicesSDK.EXPECT().List(t0GatewayID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.LocaleServicesListResult{Results: []model.LocaleServices{}, ResultCount: &resultCount}, nil)

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(t0GatewayID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, t0DisplayName, d.Get("display_name"))
		assert.Equal(t, t0Description, d.Get("description"))
		assert.Equal(t, t0Path, d.Get("path"))
		assert.Equal(t, int(t0Revision), d.Get("revision"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Tier0 ID")
	})
}

func TestMockResourceNsxtPolicyTier0GatewayCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTier0sSDK := t0mocks.NewMockTier0sClient(ctrl)
	mockLocaleServicesSDK := localeServicesMocks.NewMockLocaleServicesClient(ctrl)
	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)

	tier0Wrapper := &cliinfra.Tier0ClientContext{
		Client:     mockTier0sSDK,
		ClientType: utl.Local,
	}
	localeServicesWrapper := &tier0localeservices.LocaleServicesClientContext{
		Client:     mockLocaleServicesSDK,
		ClientType: utl.Local,
	}

	originalTier0s := cliTier0sClient
	originalLocaleServices := cliTier0LocaleServicesClient
	originalInfra := cliInfraClient
	defer func() {
		cliTier0sClient = originalTier0s
		cliTier0LocaleServicesClient = originalLocaleServices
		cliInfraClient = originalInfra
	}()

	cliTier0sClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}
	cliTier0LocaleServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0localeservices.LocaleServicesClientContext {
		return localeServicesWrapper
	}
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	t.Run("Create success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockTier0sSDK.EXPECT().Get(gomock.Any()).Return(model.Tier0{
			DisplayName:     &t0DisplayName,
			Description:     &t0Description,
			Path:            &t0Path,
			Revision:        &t0Revision,
			FailoverMode:    &t0FailoverMode,
			HaMode:          &t0HaMode,
			DisableFirewall: &t0DisableFirewall,
		}, nil)
		resultCount := int64(0)
		mockLocaleServicesSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.LocaleServicesListResult{Results: []model.LocaleServices{}, ResultCount: &resultCount}, nil)

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0DisplayName,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, t0DisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyTier0GatewayUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTier0sSDK := t0mocks.NewMockTier0sClient(ctrl)
	mockLocaleServicesSDK := localeServicesMocks.NewMockLocaleServicesClient(ctrl)
	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)

	tier0Wrapper := &cliinfra.Tier0ClientContext{
		Client:     mockTier0sSDK,
		ClientType: utl.Local,
	}
	localeServicesWrapper := &tier0localeservices.LocaleServicesClientContext{
		Client:     mockLocaleServicesSDK,
		ClientType: utl.Local,
	}

	originalTier0s := cliTier0sClient
	originalLocaleServices := cliTier0LocaleServicesClient
	originalInfra := cliInfraClient
	defer func() {
		cliTier0sClient = originalTier0s
		cliTier0LocaleServicesClient = originalLocaleServices
		cliInfraClient = originalInfra
	}()

	cliTier0sClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}
	cliTier0LocaleServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0localeservices.LocaleServicesClientContext {
		return localeServicesWrapper
	}
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	t.Run("Update success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockTier0sSDK.EXPECT().Get(t0GatewayID).Return(model.Tier0{
			DisplayName:     &t0DisplayName,
			Description:     &t0Description,
			Path:            &t0Path,
			Revision:        &t0Revision,
			FailoverMode:    &t0FailoverMode,
			HaMode:          &t0HaMode,
			DisableFirewall: &t0DisableFirewall,
		}, nil)
		resultCount := int64(0)
		mockLocaleServicesSDK.EXPECT().List(t0GatewayID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.LocaleServicesListResult{Results: []model.LocaleServices{}, ResultCount: &resultCount}, nil)

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0DisplayName,
		})
		d.SetId(t0GatewayID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, t0DisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0DisplayName,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Tier0 ID")
	})
}

func TestMockResourceNsxtPolicyTier0GatewayDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)

	originalInfra := cliInfraClient
	defer func() { cliInfraClient = originalInfra }()
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	t.Run("Delete success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0DisplayName,
		})
		d.SetId(t0GatewayID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Tier0 ID")
	})
}
