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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services/custom_signature_versions"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	sigdiffmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/settings/firewall/security/intrusion_services/custom_signature_versions"
)

// dataSourceNsxtPolicyIdpsSignatureDiffRead reads through the cliIdsCustomSignaturesDiffClient
// package-level var, so it's mockable via setupIdpsSignatureDiffMock.

func setupIdpsSignatureDiffMock(ctrl *gomock.Controller) (*sigdiffmocks.MockCustomSignaturesDiffClient, func()) {
	mockClient := sigdiffmocks.NewMockCustomSignaturesDiffClient(ctrl)
	orig := cliIdsCustomSignaturesDiffClient
	cliIdsCustomSignaturesDiffClient = func(_ client.Connector) custom_signature_versions.CustomSignaturesDiffClient {
		return mockClient
	}
	return mockClient, func() { cliIdsCustomSignaturesDiffClient = orig }
}

func TestMockDataSourceNsxtPolicyIdpsSignatureDiffReadGuard(t *testing.T) {
	t.Run("Read fails for global manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIdpsSignatureDiff()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"signature_version_id": "default",
		})

		err := dataSourceNsxtPolicyIdpsSignatureDiffRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}

func TestMockDataSourceNsxtPolicyIdpsSignatureDiffReadMocked(t *testing.T) {
	t.Run("Read succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureDiffMock(ctrl)
		defer restore()
		mockClient.EXPECT().Get("default").Return(model.IdsCustomSignaturesDiff{
			NewlyAddedSignatures: []string{"sig-1"},
			DeletedSignatures:    []string{"sig-2"},
			ExistingSignatures:   []string{"sig-3"},
		}, nil)

		ds := dataSourceNsxtPolicyIdpsSignatureDiff()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"signature_version_id": "default"})

		err := dataSourceNsxtPolicyIdpsSignatureDiffRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "default", d.Id())
		assert.Equal(t, []interface{}{"sig-1"}, d.Get("newly_added_signatures"))
	})

	t.Run("Read propagates API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureDiffMock(ctrl)
		defer restore()
		mockClient.EXPECT().Get("default").Return(model.IdsCustomSignaturesDiff{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIdpsSignatureDiff()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"signature_version_id": "default"})

		err := dataSourceNsxtPolicyIdpsSignatureDiffRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
