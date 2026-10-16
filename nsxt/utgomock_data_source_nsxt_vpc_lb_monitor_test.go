//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	apivpcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func TestMockDataSourceNsxtVpcLBMonitor(t *testing.T) {
	runVpcLbProfileDataSourceTests(t, vpcLbProfileDataSourceSpec{
		newDataSource: dataSourceNsxtVpcLBMonitor,
		setupMock:     setupVpcLbMonitorProfileMock,
		structValue:   lbMonitorProfileStructValue,
		resourceType:  nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBTCPMONITORPROFILE,
		typeName:      "TCP",
		otherTypeName: "HTTP",
		apiSegment:    "vpc-lb-monitor-profiles",
	})
}

func TestMockVpcLbMonitorProfileAPISelection(t *testing.T) {
	vpcCtx := utl.SessionContext{ClientType: utl.VPC, ProjectID: vpcLbTestProject, VPCID: vpcLbTestVPC}

	t.Run("VPC client unsupported returns nil", func(t *testing.T) {
		original := cliVpcLbMonitorProfilesClient
		defer func() { cliVpcLbMonitorProfilesClient = original }()
		cliVpcLbMonitorProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbMonitorProfileClientContext {
			return nil
		}
		assert.Nil(t, getLbMonitorProfileAPI(vpcCtx, nil))
	})

	t.Run("Infra client unsupported returns nil", func(t *testing.T) {
		original := cliLbMonitorProfilesClient
		defer func() { cliLbMonitorProfilesClient = original }()
		cliLbMonitorProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LBMonitorProfileClientContext {
			return nil
		}
		assert.Nil(t, getLbMonitorProfileAPI(utl.SessionContext{ClientType: utl.Local}, nil))
	})

	t.Run("VPC context selects the VPC adapter", func(t *testing.T) {
		api := getLbMonitorProfileAPI(vpcCtx, nil)
		require.NotNil(t, api)
		assert.IsType(t, vpcLbMonitorProfileAdapter{}, api)
	})

	t.Run("Local context selects the infra client", func(t *testing.T) {
		api := getLbMonitorProfileAPI(utl.SessionContext{ClientType: utl.Local}, nil)
		require.NotNil(t, api)
		assert.IsType(t, &infraapi.LBMonitorProfileClientContext{}, api)
	})

	t.Run("Unsupported client type for VPC wrapper", func(t *testing.T) {
		assert.Nil(t, apivpcs.NewVpcLbMonitorProfilesClient(utl.SessionContext{ClientType: utl.Global}, nil))
	})
}
