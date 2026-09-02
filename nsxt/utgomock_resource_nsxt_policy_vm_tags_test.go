//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/infra/realized_state/enforcement_points/VirtualMachinesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/enforcement_points/VirtualMachinesClient.go VirtualMachinesClient
// mockgen -destination=mocks/infra/realized_state/enforcement_points/VifsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/enforcement_points/VifsClient.go VifsClient
// mockgen -destination=mocks/infra/realized_state/virtual_machines/TagsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/virtual_machines/TagsClient.go TagsClient

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	realizedep "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state/enforcement_points"
	virtualmachines "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state/virtual_machines"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state/enforcement_points"
	vmmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state/virtual_machines"
)

var (
	vmExternalID  = "vm-external-uuid-1"
	vmDisplayName = "test-vm"
	vmInstanceID  = "vm-external-uuid-1"
)

func vmAPIResponse() model.VirtualMachine {
	return model.VirtualMachine{
		ExternalId:  &vmExternalID,
		DisplayName: &vmDisplayName,
		ComputeIds:  []string{},
		Tags:        []model.Tag{},
	}
}

// vmToStructValue converts a VirtualMachine into the *data.StructValue shape
// returned by the inventory search API, for use with seqQueryListClient.
func vmToStructValue(t *testing.T, vm model.VirtualMachine) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(vm, model.VirtualMachineBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	originalVMTags := cliVirtualMachineTagsClient
	originalVifs := cliVifsClient
	defer func() {
		cliVirtualMachineTagsClient = originalVMTags
		cliVifsClient = originalVifs
	}()
	cliVirtualMachineTagsClient = func(sessionContext utl.SessionContext, connector client.Connector) *virtualmachines.TagsClientContext {
		return mockVMTagsWrapper
	}
	cliVifsClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedep.VirtualNetworkInterfaceClientContext {
		return mockVifsWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		vmSV := vmToStructValue(t, vmAPIResponse())
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			// findNsxtPolicyVMByID (Create)
			{Results: []*data.StructValue{vmSV}, ResultCount: i64(1)},
			// resourceNsxtPolicyVMTagsRead → findNsxtPolicyVMByID
			{Results: []*data.StructValue{vmSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		// updateNsxtPolicyVMTags via cliVirtualMachineTagsClient.Create
		mockVMTagsSDK.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		// updateNsxtPolicyVMPortTags → listPolicyVifAttachmentsForVM → listAllPolicyVifs (empty portTags)
		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)
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
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		res := resourceNsxtPolicyVMTags()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVMTagsData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVMTagsCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Cannot find VM")
	})
}

func TestMockResourceNsxtPolicyVMTagsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVifsSDK := epmocks.NewMockVifsClient(ctrl)
	mockVifsWrapper := &realizedep.VirtualNetworkInterfaceClientContext{
		Client:     mockVifsSDK,
		ClientType: utl.Local,
	}

	originalVifs := cliVifsClient
	defer func() { cliVifsClient = originalVifs }()
	cliVifsClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedep.VirtualNetworkInterfaceClientContext {
		return mockVifsWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		vmSV := vmToStructValue(t, vmAPIResponse())
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{vmSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

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
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	originalVMTags := cliVirtualMachineTagsClient
	originalVifs := cliVifsClient
	defer func() {
		cliVirtualMachineTagsClient = originalVMTags
		cliVifsClient = originalVifs
	}()
	cliVirtualMachineTagsClient = func(sessionContext utl.SessionContext, connector client.Connector) *virtualmachines.TagsClientContext {
		return mockVMTagsWrapper
	}
	cliVifsClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedep.VirtualNetworkInterfaceClientContext {
		return mockVifsWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		vmSV := vmToStructValue(t, vmAPIResponse())
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			// findNsxtPolicyVMByID
			{Results: []*data.StructValue{vmSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		// updateNsxtPolicyVMTags (clear tags) via cliVirtualMachineTagsClient.Create
		mockVMTagsSDK.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
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
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	originalVMTags := cliVirtualMachineTagsClient
	originalVifs := cliVifsClient
	defer func() {
		cliVirtualMachineTagsClient = originalVMTags
		cliVifsClient = originalVifs
	}()
	cliVirtualMachineTagsClient = func(sessionContext utl.SessionContext, connector client.Connector) *virtualmachines.TagsClientContext {
		return mockVMTagsWrapper
	}
	cliVifsClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedep.VirtualNetworkInterfaceClientContext {
		return mockVifsWrapper
	}

	t.Run("Update success (uses newer TagsClient)", func(t *testing.T) {
		vmSV := vmToStructValue(t, vmAPIResponse())
		stub := &seqQueryListClient{responses: []model.SearchResponse{
			// findNsxtPolicyVMByID (Update/Create)
			{Results: []*data.StructValue{vmSV}, ResultCount: i64(1)},
			// Read after update → findNsxtPolicyVMByID
			{Results: []*data.StructValue{vmSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		// updateNsxtPolicyVMTags via cliVirtualMachineTagsClient.Create (version >= 4.1.1)
		mockVMTagsSDK.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		// updateNsxtPolicyVMPortTags → listAllPolicyVifs
		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)
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
