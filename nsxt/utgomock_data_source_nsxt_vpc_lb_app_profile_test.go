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

func TestMockDataSourceNsxtVpcLBAppProfile(t *testing.T) {
	runVpcLbProfileDataSourceTests(t, vpcLbProfileDataSourceSpec{
		newDataSource: dataSourceNsxtVpcLBAppProfile,
		setupMock:     setupVpcLbAppProfileMock,
		structValue:   lbAppProfileStructValue,
		resourceType:  nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE,
		typeName:      "TCP",
		otherTypeName: "HTTP",
		apiSegment:    "vpc-lb-app-profiles",
	})
}

func TestMockVpcLbAppProfileAPISelection(t *testing.T) {
	vpcCtx := utl.SessionContext{ClientType: utl.VPC, ProjectID: vpcLbTestProject, VPCID: vpcLbTestVPC}

	t.Run("VPC client unsupported returns nil", func(t *testing.T) {
		original := cliVpcLbAppProfilesClient
		defer func() { cliVpcLbAppProfilesClient = original }()
		cliVpcLbAppProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbAppProfileClientContext {
			return nil
		}
		assert.Nil(t, getLbAppProfileAPI(vpcCtx, nil))
	})

	t.Run("Infra client unsupported returns nil", func(t *testing.T) {
		original := cliLbAppProfilesClient
		defer func() { cliLbAppProfilesClient = original }()
		cliLbAppProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LBAppProfileClientContext {
			return nil
		}
		assert.Nil(t, getLbAppProfileAPI(utl.SessionContext{ClientType: utl.Local}, nil))
	})

	t.Run("VPC context selects the VPC adapter", func(t *testing.T) {
		api := getLbAppProfileAPI(vpcCtx, nil)
		require.NotNil(t, api)
		assert.IsType(t, vpcLbAppProfileAdapter{}, api)
	})

	t.Run("Local context selects the infra client", func(t *testing.T) {
		api := getLbAppProfileAPI(utl.SessionContext{ClientType: utl.Local}, nil)
		require.NotNil(t, api)
		assert.IsType(t, &infraapi.LBAppProfileClientContext{}, api)
	})

	t.Run("Unsupported client type for VPC wrapper", func(t *testing.T) {
		assert.Nil(t, apivpcs.NewVpcLbAppProfilesClient(utl.SessionContext{ClientType: utl.Global}, nil))
	})
}
