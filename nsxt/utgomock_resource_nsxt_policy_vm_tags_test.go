//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/infra/realized_state/VirtualMachinesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/VirtualMachinesClient.go VirtualMachinesClient
// mockgen -destination=mocks/infra/realized_state/enforcement_points/VirtualMachinesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/enforcement_points/VirtualMachinesClient.go VirtualMachinesClient
// mockgen -destination=mocks/infra/realized_state/enforcement_points/VifsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/enforcement_points/VifsClient.go VifsClient
// mockgen -destination=mocks/infra/realized_state/virtual_machines/TagsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/virtual_machines/TagsClient.go TagsClient

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	realizedstate "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state"
	realizedep "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state/enforcement_points"
	virtualmachines "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state/virtual_machines"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	rsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state/enforcement_points"
	vmmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state/virtual_machines"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vmExternalID   = "vm-external-uuid-1"
	vmDisplayName  = "test-vm"
	vmInstanceID   = "vm-external-uuid-1"
	vmResultCount1 = int64(1)
)

func vmAPIResponse() model.VirtualMachine {
	return model.VirtualMachine{
		ExternalId:  &vmExternalID,
		DisplayName: &vmDisplayName,
		ComputeIds:  []string{},
		Tags:        []model.Tag{},
	}
}

func vmListResponse() model.VirtualMachineListResult {
	return model.VirtualMachineListResult{
		Results:     []model.VirtualMachine{vmAPIResponse()},
		ResultCount: &vmResultCount1,
	}
}

func emptyVifListResponse() model.VirtualNetworkInterfaceListResult {
	total := int64(0)
	return model.VirtualNetworkInterfaceListResult{
		Results:     []model.VirtualNetworkInterface{},
		ResultCount: &total,
	}
}

func minimalVMTagsData() map[string]interface{} {
	return map[string]interface{}{
		"instance_id": vmInstanceID,
	}
}

func TestMockResourceNsxtPolicyVMTagsCreate(t *testing.T) {
	// Use version < 4.1.1 so VMs are found via list API and tags updated via Updatetags
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRealizedVMSDK := rsmocks.NewMockVirtualMachinesClient(ctrl)
	mockRealizedVMWrapper := &realizedstate.VirtualMachineClientContext{
		Client:     mockRealizedVMSDK,
		ClientType: utl.Local,
	}

	mockEPVMSDK := epmocks.NewMockVirtualMachinesClient(ctrl)
	mockEPVMWrapper := &realizedep.VirtualMachineClientContext{
		Client:     mockEPVMSDK,
		ClientType: utl.Local,
	}

	mockVifsSDK := epmocks.NewMockVifsClient(ctrl)
	mockVifsWrapper := &realizedep.VirtualNetworkInterfaceClientContext{
		Client:     mockVifsSDK,
		ClientType: utl.Local,
	}

	originalRealizedVM := cliRealizedVirtualMachinesClient
	originalEPVM := cliVirtualMachinesClient
	originalVifs := cliVifsClient
	defer func() {
		cliRealizedVirtualMachinesClient = originalRealizedVM
		cliVirtualMachinesClient = originalEPVM
		cliVifsClient = originalVifs
	}()
	cliRealizedVirtualMachinesClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedstate.VirtualMachineClientContext {
		return mockRealizedVMWrapper
	}
	cliVirtualMachinesClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedep.VirtualMachineClientContext {
		return mockEPVMWrapper
	}
	cliVifsClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedep.VirtualNetworkInterfaceClientContext {
		return mockVifsWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		// findNsxtPolicyVMByID → listAllPolicyVirtualMachines (Create)
		mockRealizedVMSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(vmListResponse(), nil)
		// updateNsxtPolicyVMTags
		mockEPVMSDK.EXPECT().Updatetags(gomock.Any(), gomock.Any()).Return(nil)
		// updateNsxtPolicyVMPortTags → listPolicyVifAttachmentsForVM → listAllPolicyVifs (empty portTags)
		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)
		// resourceNsxtPolicyVMTagsRead → findNsxtPolicyVMByID
		mockRealizedVMSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(vmListResponse(), nil)
		// setPolicyVMPortTagsInSchema → listPolicyVifAttachmentsForVM → listAllPolicyVifs
		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)

		res := resourceNsxtPolicyVMTags()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVMTagsData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVMTagsCreate(d, m)
		require.NoError(t, err)
		assert.Equal(t, vmExternalID, d.Id())
	})

	t.Run("Create fails when VM not found", func(t *testing.T) {
		mockRealizedVMSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.VirtualMachineListResult{Results: []model.VirtualMachine{}}, nil)

		res := resourceNsxtPolicyVMTags()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVMTagsData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVMTagsCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Cannot find VM")
	})
}

func TestMockResourceNsxtPolicyVMTagsRead(t *testing.T) {
	// Use version < 4.1.1 so VMs are found via list API
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRealizedVMSDK := rsmocks.NewMockVirtualMachinesClient(ctrl)
	mockRealizedVMWrapper := &realizedstate.VirtualMachineClientContext{
		Client:     mockRealizedVMSDK,
		ClientType: utl.Local,
	}

	mockVifsSDK := epmocks.NewMockVifsClient(ctrl)
	mockVifsWrapper := &realizedep.VirtualNetworkInterfaceClientContext{
		Client:     mockVifsSDK,
		ClientType: utl.Local,
	}

	originalRealizedVM := cliRealizedVirtualMachinesClient
	originalVifs := cliVifsClient
	defer func() {
		cliRealizedVirtualMachinesClient = originalRealizedVM
		cliVifsClient = originalVifs
	}()
	cliRealizedVirtualMachinesClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedstate.VirtualMachineClientContext {
		return mockRealizedVMWrapper
	}
	cliVifsClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedep.VirtualNetworkInterfaceClientContext {
		return mockVifsWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockRealizedVMSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(vmListResponse(), nil)
		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)

		res := resourceNsxtPolicyVMTags()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"instance_id": vmInstanceID,
		})
		d.SetId(vmExternalID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVMTagsRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, vmExternalID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyVMTags()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVMTagsRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Virtual Machine ID")
	})

	t.Run("Read clears ID when VM not found", func(t *testing.T) {
		mockRealizedVMSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.VirtualMachineListResult{Results: []model.VirtualMachine{}}, nil)

		res := resourceNsxtPolicyVMTags()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"instance_id": vmInstanceID,
		})
		d.SetId(vmExternalID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVMTagsRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyVMTagsDelete(t *testing.T) {
	// Use version < 4.1.1 so VMs are found via list API and tags cleared via Updatetags
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRealizedVMSDK := rsmocks.NewMockVirtualMachinesClient(ctrl)
	mockRealizedVMWrapper := &realizedstate.VirtualMachineClientContext{
		Client:     mockRealizedVMSDK,
		ClientType: utl.Local,
	}

	mockEPVMSDK := epmocks.NewMockVirtualMachinesClient(ctrl)
	mockEPVMWrapper := &realizedep.VirtualMachineClientContext{
		Client:     mockEPVMSDK,
		ClientType: utl.Local,
	}

	mockVifsSDK := epmocks.NewMockVifsClient(ctrl)
	mockVifsWrapper := &realizedep.VirtualNetworkInterfaceClientContext{
		Client:     mockVifsSDK,
		ClientType: utl.Local,
	}

	originalRealizedVM := cliRealizedVirtualMachinesClient
	originalEPVM := cliVirtualMachinesClient
	originalVifs := cliVifsClient
	defer func() {
		cliRealizedVirtualMachinesClient = originalRealizedVM
		cliVirtualMachinesClient = originalEPVM
		cliVifsClient = originalVifs
	}()
	cliRealizedVirtualMachinesClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedstate.VirtualMachineClientContext {
		return mockRealizedVMWrapper
	}
	cliVirtualMachinesClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedep.VirtualMachineClientContext {
		return mockEPVMWrapper
	}
	cliVifsClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedep.VirtualNetworkInterfaceClientContext {
		return mockVifsWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		// findNsxtPolicyVMByID
		mockRealizedVMSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(vmListResponse(), nil)
		// updateNsxtPolicyVMTags (clear tags)
		mockEPVMSDK.EXPECT().Updatetags(gomock.Any(), gomock.Any()).Return(nil)
		// updateNsxtPolicyVMPortTags → listPolicyVifAttachmentsForVM → listAllPolicyVifs
		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)

		res := resourceNsxtPolicyVMTags()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"instance_id": vmInstanceID,
		})
		d.SetId(vmExternalID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVMTagsDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete clears ID when VM not found", func(t *testing.T) {
		mockRealizedVMSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.VirtualMachineListResult{Results: []model.VirtualMachine{}}, nil)

		res := resourceNsxtPolicyVMTags()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"instance_id": vmInstanceID,
		})
		d.SetId(vmExternalID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVMTagsDelete(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyVMTagsUpdate(t *testing.T) {
	// Use version 4.1.1+ so tags are updated via the newer TagsClient.Create path
	util.NsxVersion = "4.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRealizedVMSDK := rsmocks.NewMockVirtualMachinesClient(ctrl)
	mockRealizedVMWrapper := &realizedstate.VirtualMachineClientContext{
		Client:     mockRealizedVMSDK,
		ClientType: utl.Local,
	}

	mockVMTagsSDK := vmmocks.NewMockTagsClient(ctrl)
	mockVMTagsWrapper := &virtualmachines.TagsClientContext{
		Client:     mockVMTagsSDK,
		ClientType: utl.Local,
	}

	mockVifsSDK := epmocks.NewMockVifsClient(ctrl)
	mockVifsWrapper := &realizedep.VirtualNetworkInterfaceClientContext{
		Client:     mockVifsSDK,
		ClientType: utl.Local,
	}

	originalRealizedVM := cliRealizedVirtualMachinesClient
	originalVMTags := cliVirtualMachineTagsClient
	originalVifs := cliVifsClient
	defer func() {
		cliRealizedVirtualMachinesClient = originalRealizedVM
		cliVirtualMachineTagsClient = originalVMTags
		cliVifsClient = originalVifs
	}()
	cliRealizedVirtualMachinesClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedstate.VirtualMachineClientContext {
		return mockRealizedVMWrapper
	}
	cliVirtualMachineTagsClient = func(sessionContext utl.SessionContext, connector client.Connector) *virtualmachines.TagsClientContext {
		return mockVMTagsWrapper
	}
	cliVifsClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedep.VirtualNetworkInterfaceClientContext {
		return mockVifsWrapper
	}

	t.Run("Update success (uses newer TagsClient)", func(t *testing.T) {
		// findNsxtPolicyVMByID → list (version < 4.1.2 uses list API)
		mockRealizedVMSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(vmListResponse(), nil)
		// updateNsxtPolicyVMTags via cliVirtualMachineTagsClient.Create (version >= 4.1.1)
		mockVMTagsSDK.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		// updateNsxtPolicyVMPortTags → listAllPolicyVifs
		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)
		// Read after update → findNsxtPolicyVMByID
		mockRealizedVMSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(vmListResponse(), nil)
		// setPolicyVMPortTagsInSchema → listAllPolicyVifs
		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)

		res := resourceNsxtPolicyVMTags()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVMTagsData())
		d.SetId(vmExternalID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVMTagsUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, vmExternalID, d.Id())
	})
}
