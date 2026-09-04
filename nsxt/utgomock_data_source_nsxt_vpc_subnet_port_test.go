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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

const dsVpcPortSubnetPath = "/orgs/default/projects/proj1/vpcs/vpc1/subnets/subnet1"

func dsVpcSubnetPortResponse(portID, lportAttachmentID string) model.VpcSubnetPort {
	displayName := "port-" + portID
	desc := "test port " + portID
	path := dsVpcPortSubnetPath + "/ports/" + portID
	return model.VpcSubnetPort{
		Id:          &portID,
		DisplayName: &displayName,
		Description: &desc,
		Path:        &path,
		Attachment: &model.PortAttachment{
			Id: &lportAttachmentID,
		},
	}
}

func dsVpcSubnetPortsListResponse(ports ...model.VpcSubnetPort) model.VpcSubnetPortListResult {
	total := int64(len(ports))
	return model.VpcSubnetPortListResult{
		Results:     ports,
		ResultCount: &total,
	}
}

func TestUnitNsxt_DataSourceNsxtVpcSubnetPortRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	t.Run("finds port matching VM's VIF attachment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPortsSDK, restorePorts := setupExtAddrMock(t, ctrl)
		defer restorePorts()
		mockVifsSDK, restoreVifs := setupPolicyVMVifsMock(t, ctrl)
		defer restoreVifs()

		lportAttachID := "lport-attach-1"
		vif := model.VirtualNetworkInterface{
			LportAttachmentId: &lportAttachID,
			OwnerVmId:         &vmExternalID,
		}
		total := int64(1)
		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.VirtualNetworkInterfaceListResult{Results: []model.VirtualNetworkInterface{vif}, ResultCount: &total}, nil)

		port := dsVpcSubnetPortResponse("port1", lportAttachID)
		mockPortsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(dsVpcSubnetPortsListResponse(port), nil)

		ds := dataSourceNsxtVpcSubnetPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"subnet_path": dsVpcPortSubnetPath,
			"vm_id":       vmExternalID,
		})

		err := dataSourceNsxtVpcSubnetPortRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "port1", d.Id())
		assert.Equal(t, "port-port1", d.Get("display_name"))
	})

	t.Run("no matching port errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPortsSDK, restorePorts := setupExtAddrMock(t, ctrl)
		defer restorePorts()
		mockVifsSDK, restoreVifs := setupPolicyVMVifsMock(t, ctrl)
		defer restoreVifs()

		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)

		mockPortsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(dsVpcSubnetPortsListResponse(dsVpcSubnetPortResponse("port1", "other-attach-id")), nil)

		ds := dataSourceNsxtVpcSubnetPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"subnet_path": dsVpcPortSubnetPath,
			"vm_id":       vmExternalID,
		})

		err := dataSourceNsxtVpcSubnetPortRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find port")
	})

	t.Run("vif listing error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVifsSDK, restoreVifs := setupPolicyVMVifsMock(t, ctrl)
		defer restoreVifs()

		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.VirtualNetworkInterfaceListResult{}, errors.New("vif fail"))

		ds := dataSourceNsxtVpcSubnetPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"subnet_path": dsVpcPortSubnetPath,
			"vm_id":       vmExternalID,
		})

		err := dataSourceNsxtVpcSubnetPortRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list port attachments")
	})

	t.Run("invalid subnet_path errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVifsSDK, restoreVifs := setupPolicyVMVifsMock(t, ctrl)
		defer restoreVifs()

		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)

		ds := dataSourceNsxtVpcSubnetPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"subnet_path": "/infra/tier-1s/t1",
			"vm_id":       vmExternalID,
		})

		err := dataSourceNsxtVpcSubnetPortRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("ports List API error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPortsSDK, restorePorts := setupExtAddrMock(t, ctrl)
		defer restorePorts()
		mockVifsSDK, restoreVifs := setupPolicyVMVifsMock(t, ctrl)
		defer restoreVifs()

		mockVifsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(emptyVifListResponse(), nil)

		mockPortsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.VpcSubnetPortListResult{}, errors.New("list failed"))

		ds := dataSourceNsxtVpcSubnetPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"subnet_path": dsVpcPortSubnetPath,
			"vm_id":       vmExternalID,
		})

		err := dataSourceNsxtVpcSubnetPortRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
