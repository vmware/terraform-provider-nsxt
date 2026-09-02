//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/realized_state/VirtualMachinesClient.go -package=mocks -source=<sdk>/services/nsxt/infra/realized_state/VirtualMachinesClient.go VirtualMachinesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	realizedstate "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	realizedvmmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state"
)

func setupPolicyVMsListMock(t *testing.T, ctrl *gomock.Controller) (*realizedvmmocks.MockVirtualMachinesClient, func()) {
	mockSDK := realizedvmmocks.NewMockVirtualMachinesClient(ctrl)
	mockWrapper := &realizedstate.VirtualMachineClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliRealizedVirtualMachinesClient
	cliRealizedVirtualMachinesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *realizedstate.VirtualMachineClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliRealizedVirtualMachinesClient = original }
}

func policyVMsListResponse(vms ...model.VirtualMachine) model.VirtualMachineListResult {
	total := int64(len(vms))
	return model.VirtualMachineListResult{
		Results:     vms,
		ResultCount: &total,
	}
}

func makePolicyVM(displayName, externalID, biosID, instanceID, powerState, osName string) model.VirtualMachine {
	vm := model.VirtualMachine{
		DisplayName: &displayName,
		ExternalId:  &externalID,
		ComputeIds:  []string{"biosUuid:" + biosID, "instanceUuid:" + instanceID},
	}
	if powerState != "" {
		vm.PowerState = &powerState
	}
	if osName != "" {
		vm.GuestInfo = &model.GuestInfo{OsName: &osName}
	}
	return vm
}

func TestUnitNsxt_DataSourceNsxtPolicyVMsRead(t *testing.T) {
	t.Run("default value_type (bios_id) with no filters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupPolicyVMsListMock(t, ctrl)
		defer restore()

		vm := makePolicyVM("vm-1", "ext-1", "bios-1", "inst-1", model.VirtualMachine_POWER_STATE_VM_RUNNING, "Ubuntu")
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(policyVMsListResponse(vm), nil)

		ds := dataSourceNsxtPolicyVMs()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyVMsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		assert.Equal(t, "bios-1", items["vm-1"])
	})

	t.Run("value_type external_id skips VMs without external id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupPolicyVMsListMock(t, ctrl)
		defer restore()

		vmNoExt := makePolicyVM("vm-noext", "", "bios-2", "inst-2", "", "")
		vmNoExt.ExternalId = nil
		vmWithExt := makePolicyVM("vm-2", "ext-2", "bios-3", "inst-3", "", "")

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(policyVMsListResponse(vmNoExt, vmWithExt), nil)

		ds := dataSourceNsxtPolicyVMs()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"value_type": "external_id",
		})

		err := dataSourceNsxtPolicyVMsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		assert.Equal(t, "ext-2", items["vm-2"])
		_, hasNoExt := items["vm-noext"]
		assert.False(t, hasNoExt)
	})

	t.Run("value_type instance_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupPolicyVMsListMock(t, ctrl)
		defer restore()

		vm := makePolicyVM("vm-3", "ext-3", "bios-4", "inst-4", "", "")
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(policyVMsListResponse(vm), nil)

		ds := dataSourceNsxtPolicyVMs()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"value_type": "instance_id",
		})

		err := dataSourceNsxtPolicyVMsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		assert.Equal(t, "inst-4", items["vm-3"])
	})

	t.Run("filters by state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupPolicyVMsListMock(t, ctrl)
		defer restore()

		running := makePolicyVM("vm-running", "ext-r", "bios-r", "inst-r", model.VirtualMachine_POWER_STATE_VM_RUNNING, "")
		stopped := makePolicyVM("vm-stopped", "ext-s", "bios-s", "inst-s", model.VirtualMachine_POWER_STATE_VM_STOPPED, "")

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(policyVMsListResponse(running, stopped), nil)

		ds := dataSourceNsxtPolicyVMs()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"state": "running",
		})

		err := dataSourceNsxtPolicyVMsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		_, hasRunning := items["vm-running"]
		_, hasStopped := items["vm-stopped"]
		assert.True(t, hasRunning)
		assert.False(t, hasStopped)
	})

	t.Run("filters by guest_os prefix", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupPolicyVMsListMock(t, ctrl)
		defer restore()

		ubuntu := makePolicyVM("vm-ubuntu", "ext-u", "bios-u", "inst-u", "", "Ubuntu Linux")
		windows := makePolicyVM("vm-win", "ext-w", "bios-w", "inst-w", "", "Windows Server")

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(policyVMsListResponse(ubuntu, windows), nil)

		ds := dataSourceNsxtPolicyVMs()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"guest_os": "ubuntu",
		})

		err := dataSourceNsxtPolicyVMsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		_, hasUbuntu := items["vm-ubuntu"]
		_, hasWin := items["vm-win"]
		assert.True(t, hasUbuntu)
		assert.False(t, hasWin)
	})

	t.Run("filters by display_name regex", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupPolicyVMsListMock(t, ctrl)
		defer restore()

		match := makePolicyVM("prod-vm-1", "ext-p1", "bios-p1", "inst-p1", "", "")
		noMatch := makePolicyVM("dev-vm-1", "ext-d1", "bios-d1", "inst-d1", "", "")

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(policyVMsListResponse(match, noMatch), nil)

		ds := dataSourceNsxtPolicyVMs()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "^prod-.*",
		})

		err := dataSourceNsxtPolicyVMsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		_, hasProd := items["prod-vm-1"]
		_, hasDev := items["dev-vm-1"]
		assert.True(t, hasProd)
		assert.False(t, hasDev)
	})

	t.Run("invalid display_name regex errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupPolicyVMsListMock(t, ctrl)
		defer restore()

		vm := makePolicyVM("vm-1", "ext-1", "bios-1", "inst-1", "", "")
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(policyVMsListResponse(vm), nil)

		ds := dataSourceNsxtPolicyVMs()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "[invalid(",
		})

		err := dataSourceNsxtPolicyVMsRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("List API error is wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupPolicyVMsListMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.VirtualMachineListResult{}, errors.New("list failed"))

		ds := dataSourceNsxtPolicyVMs()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyVMsRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error reading Virtual Machines")
	})
}
