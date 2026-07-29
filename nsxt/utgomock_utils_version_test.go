//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/node/VersionClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/node/VersionClient.go VersionClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"

	nodeAPI "github.com/vmware/terraform-provider-nsxt/api/nsx/node"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	nodemocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/node"
)

func setupVersionMock(ctrl *gomock.Controller) (*nodemocks.MockVersionClient, func()) {
	mockSDK := nodemocks.NewMockVersionClient(ctrl)
	mockWrapper := &nodeAPI.NodeVersionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliVersionClient
	cliVersionClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *nodeAPI.NodeVersionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVersionClient = original }
}

func TestMockNsxtGetNSXVersion(t *testing.T) {
	t.Run("prefers ProductVersion when present", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		nodeVersion := "4.1.0.0.0"
		productVersion := "4.1.0"
		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{NodeVersion: &nodeVersion, ProductVersion: &productVersion}, nil)

		version, err := getNSXVersion(nil)
		require.NoError(t, err)
		assert.Equal(t, "4.1.0", version)
	})

	t.Run("falls back to NodeVersion when ProductVersion is absent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		nodeVersion := "4.1.0.0.0"
		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{NodeVersion: &nodeVersion}, nil)

		version, err := getNSXVersion(nil)
		require.NoError(t, err)
		assert.Equal(t, "4.1.0.0.0", version)
	})

	t.Run("Get error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{}, errors.New("connection error"))

		_, err := getNSXVersion(nil)
		require.Error(t, err)
	})
}

func TestMockNsxtInitNSXVersion(t *testing.T) {
	t.Run("success sets util.NsxVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		// NodeVersion must be non-nil: getNSXVersion unconditionally dereferences
		// it for a debug log before checking ProductVersion.
		nodeVersion := "4.1.0.0.0"
		productVersion := "4.1.0"
		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{NodeVersion: &nodeVersion, ProductVersion: &productVersion}, nil)

		err := initNSXVersion(nil)
		require.NoError(t, err)
	})

	t.Run("error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{}, errors.New("connection error"))

		err := initNSXVersion(nil)
		require.Error(t, err)
	})
}
