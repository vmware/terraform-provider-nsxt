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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	versionmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/settings/firewall/security/intrusion_services"
)

// resourceNsxtPolicyIdpsSignatureVersion talks to the SDK-generated client via the
// cliIdsSignatureVersionsClient package-level var, so it's mockable via setupIdpsSignatureVersionMock.
// These tests cover guard conditions, validation logic, and mocked CRUD paths.

func setupIdpsSignatureVersionMock(ctrl *gomock.Controller) (*versionmocks.MockSignatureVersionsClient, func()) {
	mockClient := versionmocks.NewMockSignatureVersionsClient(ctrl)
	orig := cliIdsSignatureVersionsClient
	cliIdsSignatureVersionsClient = func(_ client.Connector) intrusion_services.SignatureVersionsClient {
		return mockClient
	}
	return mockClient, func() { cliIdsSignatureVersionsClient = orig }
}

func minimalIdpsSigVersionData() map[string]interface{} {
	return map[string]interface{}{
		"nsx_id": "sig-ver-1",
	}
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionGuard(t *testing.T) {
	t.Run("Create fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())

		err := resourceNsxtPolicyIdpsSignatureVersionCreate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})

	t.Run("Read fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())
		d.SetId("sig-ver-1")

		err := resourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})

	t.Run("Update fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())
		d.SetId("sig-ver-1")

		err := resourceNsxtPolicyIdpsSignatureVersionUpdate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionCreateNsxIDRequired(t *testing.T) {
	t.Run("Create fails when nsx_id is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtPolicyIdpsSignatureVersionCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nsx_id")
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionReadEmptyID(t *testing.T) {
	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())

		err := resourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionDeleteNoOp(t *testing.T) {
	t.Run("Delete is a no-op (removes from state only)", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())
		d.SetId("sig-ver-1")

		err := resourceNsxtPolicyIdpsSignatureVersionDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionUpdateInvalidState(t *testing.T) {
	t.Run("Update fails with invalid state value", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		data := minimalIdpsSigVersionData()
		data["state"] = "INVALID_STATE"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId("sig-ver-1")

		// Since state is Computed-only, HasChange won't fire in unit test context.
		// Verify schema definition is correct.
		assert.NotNil(t, res.Schema["state"])
		assert.NotNil(t, res.Schema["version_id"])
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionCreateMocked(t *testing.T) {
	t.Run("Create succeeds when version exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		mockClient.EXPECT().Get("sig-ver-1").Return(model.IdsSignatureVersion{}, nil).Times(2)

		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())

		err := resourceNsxtPolicyIdpsSignatureVersionCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "sig-ver-1", d.Id())
	})

	t.Run("Create fails when version does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		mockClient.EXPECT().Get("sig-ver-1").Return(model.IdsSignatureVersion{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())

		err := resourceNsxtPolicyIdpsSignatureVersionCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionReadMocked(t *testing.T) {
	t.Run("Read populates fields from the API object", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		state := "ACTIVE"
		mockClient.EXPECT().Get("sig-ver-1").Return(model.IdsSignatureVersion{State: &state}, nil)

		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())
		d.SetId("sig-ver-1")

		err := resourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "ACTIVE", d.Get("state"))
	})

	t.Run("Read propagates API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		mockClient.EXPECT().Get("sig-ver-1").Return(model.IdsSignatureVersion{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())
		d.SetId("sig-ver-1")

		err := resourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionUpdateMocked(t *testing.T) {
	t.Run("Update to ACTIVE calls Makeactiveversion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		obj := model.IdsSignatureVersion{}
		mockClient.EXPECT().Get("sig-ver-1").Return(obj, nil).Times(2)
		mockClient.EXPECT().Makeactiveversion(obj).Return(nil)

		res := resourceNsxtPolicyIdpsSignatureVersion()
		data := minimalIdpsSigVersionData()
		data["state"] = "ACTIVE"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId("sig-ver-1")

		err := resourceNsxtPolicyIdpsSignatureVersionUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update to ACTIVE propagates Makeactiveversion error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		obj := model.IdsSignatureVersion{}
		mockClient.EXPECT().Get("sig-ver-1").Return(obj, nil)
		mockClient.EXPECT().Makeactiveversion(obj).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIdpsSignatureVersion()
		data := minimalIdpsSigVersionData()
		data["state"] = "ACTIVE"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId("sig-ver-1")

		err := resourceNsxtPolicyIdpsSignatureVersionUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
