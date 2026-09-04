//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	realizedep "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state/enforcement_points"
)

// setupPolicyVMVifsMock overrides cliVifsClient to return the given mock,
// reusing the mocks/wrapper types already established for
// resource_nsxt_policy_vm_tags.go (see epmocks.MockVifsClient /
// realizedep.VirtualNetworkInterfaceClientContext in
// utgomock_resource_nsxt_policy_vm_tags_test.go).
func setupPolicyVMVifsMock(t *testing.T, ctrl *gomock.Controller) (*epmocks.MockVifsClient, func()) {
	mockVifsSDK := epmocks.NewMockVifsClient(ctrl)
	mockWrapper := &realizedep.VirtualNetworkInterfaceClientContext{
		Client:     mockVifsSDK,
		ClientType: utl.Local,
	}

	original := cliVifsClient
	cliVifsClient = func(_ utl.SessionContext, _ client.Connector) *realizedep.VirtualNetworkInterfaceClientContext {
		return mockWrapper
	}
	return mockVifsSDK, func() { cliVifsClient = original }
}

func TestUnitNsxt_DataSourceNsxtPolicyVMIDRead(t *testing.T) {
	t.Run("by external_id success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVifsSDK, restoreVifs := setupPolicyVMVifsMock(t, ctrl)
		defer restoreVifs()

		vmSV := vmToStructValue(t, vmAPIResponse())
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{vmSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)

		ds := dataSourceNsxtPolicyVM()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"external_id": vmExternalID,
		})

		err := dataSourceNsxtPolicyVMIDRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vmExternalID, d.Id())
		assert.Equal(t, vmDisplayName, d.Get("display_name"))
	})

	t.Run("by external_id not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyVM()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"external_id": "no-such-vm",
		})

		err := dataSourceNsxtPolicyVMIDRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading Virtual Machine")
	})

	t.Run("by display_name perfect match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVifsSDK, restoreVifs := setupPolicyVMVifsMock(t, ctrl)
		defer restoreVifs()

		vmSV := vmToStructValue(t, vmAPIResponse())
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{vmSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)

		ds := dataSourceNsxtPolicyVM()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": vmDisplayName,
		})

		err := dataSourceNsxtPolicyVMIDRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vmExternalID, d.Id())
	})

	t.Run("by display_name prefix match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVifsSDK, restoreVifs := setupPolicyVMVifsMock(t, ctrl)
		defer restoreVifs()

		vm := vmAPIResponse()
		vmSV := vmToStructValue(t, vm)
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{vmSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)

		ds := dataSourceNsxtPolicyVM()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "test-v",
		})

		err := dataSourceNsxtPolicyVMIDRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vmExternalID, d.Id())
	})

	t.Run("by display_name multiple matches", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		vm1 := vmAPIResponse()
		name2 := "test-vm-2"
		ext2 := "vm-external-uuid-2"
		vm2 := vmAPIResponse()
		vm2.DisplayName = &name2
		vm2.ExternalId = &ext2

		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{vmToStructValue(t, vm1), vmToStructValue(t, vm2)}, ResultCount: i64(2)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyVM()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "test-v",
		})

		err := dataSourceNsxtPolicyVMIDRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Found 2 Virtual Machines")
	})

	t.Run("by display_name not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyVM()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtPolicyVMIDRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unable to find Virtual Machine")
	})

	t.Run("vif listing error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVifsSDK, restoreVifs := setupPolicyVMVifsMock(t, ctrl)
		defer restoreVifs()

		vmSV := vmToStructValue(t, vmAPIResponse())
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{vmSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.VirtualNetworkInterfaceListResult{}, errors.New("vif list failed"))

		ds := dataSourceNsxtPolicyVM()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"external_id": vmExternalID,
		})

		err := dataSourceNsxtPolicyVMIDRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error getting the VIF attachment ids")
	})
}
