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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	apiprojects "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	cnamocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
)

var (
	cnaID          = "cna-test-id"
	cnaDisplayName = "test-cna"
	cnaDescription = "Test CNA"
	cnaRevision    = int64(1)
	cnaPath        = "/orgs/default/projects/project1/centralized-network-attachments/cna-test-id"
	cnaOrgID       = "default"
	cnaProjectID   = "project1"
	cnaSubnetPath  = "/orgs/default/projects/project1/vpcs/vpc1/subnets/sub1"
)

func cnaAPIResponse() nsxModel.CentralizedNetworkAttachment {
	return nsxModel.CentralizedNetworkAttachment{
		Id:          &cnaID,
		DisplayName: &cnaDisplayName,
		Description: &cnaDescription,
		Revision:    &cnaRevision,
		Path:        &cnaPath,
		SubnetPath:  &cnaSubnetPath,
	}
}

func minimalCNAData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": cnaDisplayName,
		"description":  cnaDescription,
		"nsx_id":       cnaID,
		"subnet_path":  cnaSubnetPath,
		"context": []interface{}{
			map[string]interface{}{
				"project_id":  cnaProjectID,
				"vpc_id":      "",
				"from_global": false,
			},
		},
	}
}

func setupCNAMock(t *testing.T, ctrl *gomock.Controller) (*cnamocks.MockCentralizedNetworkAttachmentsClient, func()) {
	mockSDK := cnamocks.NewMockCentralizedNetworkAttachmentsClient(ctrl)
	mockWrapper := &apiprojects.CentralizedNetworkAttachmentClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  cnaProjectID,
	}

	original := cliCNAClient
	cliCNAClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.CentralizedNetworkAttachmentClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliCNAClient = original }
}

func TestMockResourceNsxtPolicyCNACreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCNAMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(nsxModel.CentralizedNetworkAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(cnaAPIResponse(), nil),
		)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNACreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cnaID, d.Id())
		assert.Equal(t, cnaSubnetPath, d.Get("subnet_path"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(cnaAPIResponse(), nil)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNACreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when Patch returns error", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(nsxModel.CentralizedNetworkAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNACreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create with private_network_attachment and interface_subnet", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		var captured nsxModel.CentralizedNetworkAttachment
		gomock.InOrder(
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(nsxModel.CentralizedNetworkAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).DoAndReturn(
				func(_, _, _ string, obj nsxModel.CentralizedNetworkAttachment) error {
					captured = obj
					return nil
				}),
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(cnaAPIResponse(), nil),
		)

		data := minimalCNAData()
		data["private_network_attachment"] = false
		data["interface_subnet"] = []interface{}{
			map[string]interface{}{
				"interface_ip_address": []interface{}{"192.168.1.10"},
				"prefix_length":        24,
				"ha_vip_ip_address":    "192.168.1.100",
			},
		}

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyProjectCNACreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.NotNil(t, captured.PrivateNetworkAttachment)
		assert.False(t, *captured.PrivateNetworkAttachment)
		assert.False(t, d.Get("private_network_attachment").(bool))
		require.Len(t, captured.InterfaceSubnets, 1)
		assert.Equal(t, int64(24), *captured.InterfaceSubnets[0].PrefixLength)
		assert.Equal(t, "192.168.1.100", *captured.InterfaceSubnets[0].HaVipIpAddress)
	})
}

func TestMockResourceNsxtPolicyCNARead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCNAMock(t, ctrl)
	defer restore()

	t.Run("Read success sets all fields", func(t *testing.T) {
		prefixLen := int64(24)
		haVip := "192.168.1.100"
		privateNetAttach := false
		apiResp := cnaAPIResponse()
		apiResp.PrivateNetworkAttachment = &privateNetAttach
		apiResp.InterfaceSubnets = []nsxModel.CentralizedNetworkAttachmentInterfaceSubnet{
			{
				InterfaceIpAddress: []string{"192.168.1.10"},
				PrefixLength:       &prefixLen,
				HaVipIpAddress:     &haVip,
			},
		}
		mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(apiResp, nil)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())
		d.SetId(cnaID)

		err := resourceNsxtPolicyProjectCNARead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cnaDisplayName, d.Get("display_name"))
		assert.Equal(t, cnaSubnetPath, d.Get("subnet_path"))
		assert.False(t, d.Get("private_network_attachment").(bool))

		subnets := d.Get("interface_subnet").([]interface{})
		require.Len(t, subnets, 1)
		subMap := subnets[0].(map[string]interface{})
		assert.Equal(t, 24, subMap["prefix_length"])
		assert.Equal(t, haVip, subMap["ha_vip_ip_address"])
	})

	for _, testCase := range []struct {
		name       string
		configured interface{}
		expected   bool
	}{
		{name: "Read missing field preserves false", configured: false, expected: false},
		{name: "Read missing field preserves true", configured: true, expected: true},
		{name: "Read missing field uses schema default", expected: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(cnaAPIResponse(), nil)

			res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
			data := minimalCNAData()
			if testCase.configured != nil {
				data["private_network_attachment"] = testCase.configured
			}
			d := schema.TestResourceDataRaw(t, res.Schema, data)
			d.SetId(cnaID)

			err := resourceNsxtPolicyProjectCNARead(d, newGoMockProviderClient())
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, d.Get("private_network_attachment"))
		})
	}

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(nsxModel.CentralizedNetworkAttachment{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())
		d.SetId(cnaID)

		err := resourceNsxtPolicyProjectCNARead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNARead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyCNAUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCNAMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).Return(cnaAPIResponse(), nil),
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(cnaAPIResponse(), nil),
		)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())
		d.SetId(cnaID)

		err := resourceNsxtPolicyProjectCNAUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update with private_network_attachment and interface_subnet", func(t *testing.T) {
		var captured nsxModel.CentralizedNetworkAttachment
		gomock.InOrder(
			mockSDK.EXPECT().Update(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).DoAndReturn(
				func(_, _, _ string, obj nsxModel.CentralizedNetworkAttachment) (nsxModel.CentralizedNetworkAttachment, error) {
					captured = obj
					return cnaAPIResponse(), nil
				}),
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(cnaAPIResponse(), nil),
		)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		data := minimalCNAData()
		data["private_network_attachment"] = false
		data["interface_subnet"] = []interface{}{
			map[string]interface{}{
				"interface_ip_address": []interface{}{"192.168.1.10"},
				"prefix_length":        24,
				"ha_vip_ip_address":    "192.168.1.100",
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(cnaID)

		err := resourceNsxtPolicyProjectCNAUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.NotNil(t, captured.PrivateNetworkAttachment)
		assert.False(t, *captured.PrivateNetworkAttachment)
		assert.False(t, d.Get("private_network_attachment").(bool))
		require.Len(t, captured.InterfaceSubnets, 1)
		assert.Equal(t, int64(24), *captured.InterfaceSubnets[0].PrefixLength)
		assert.Equal(t, "192.168.1.100", *captured.InterfaceSubnets[0].HaVipIpAddress)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNAUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyCNAPrivateNetworkAttachment(t *testing.T) {
	for _, operation := range []string{"create", "update"} {
		t.Run(operation, func(t *testing.T) {
			for _, testCase := range []struct {
				name       string
				configured interface{}
				expected   bool
			}{
				{name: "explicit false", configured: false, expected: false},
				{name: "explicit true", configured: true, expected: true},
				{name: "omitted defaults to true", expected: true},
			} {
				t.Run(testCase.name, func(t *testing.T) {
					ctrl := gomock.NewController(t)
					mockSDK, restore := setupCNAMock(t, ctrl)
					defer restore()

					var captured nsxModel.CentralizedNetworkAttachment
					response := cnaAPIResponse()
					response.PrivateNetworkAttachment = &testCase.expected
					data := minimalCNAData()
					if testCase.configured != nil {
						data["private_network_attachment"] = testCase.configured
					}
					res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
					d := schema.TestResourceDataRaw(t, res.Schema, data)

					if operation == "create" {
						gomock.InOrder(
							mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).
								Return(nsxModel.CentralizedNetworkAttachment{}, vapiErrors.NotFound{}),
							mockSDK.EXPECT().Patch(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).
								DoAndReturn(func(_, _, _ string, obj nsxModel.CentralizedNetworkAttachment) error {
									captured = obj
									return nil
								}),
							mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(response, nil),
						)
						require.NoError(t, resourceNsxtPolicyProjectCNACreate(d, newGoMockProviderClient()))
					} else {
						d.SetId(cnaID)
						gomock.InOrder(
							mockSDK.EXPECT().Update(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).
								DoAndReturn(func(
									_, _, _ string,
									obj nsxModel.CentralizedNetworkAttachment,
								) (nsxModel.CentralizedNetworkAttachment, error) {
									captured = obj
									return response, nil
								}),
							mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(response, nil),
						)
						require.NoError(t, resourceNsxtPolicyProjectCNAUpdate(d, newGoMockProviderClient()))
					}

					require.NotNil(t, captured.PrivateNetworkAttachment)
					assert.Equal(t, testCase.expected, *captured.PrivateNetworkAttachment)
					assert.Equal(t, testCase.expected, d.Get("private_network_attachment"))
				})
			}
		})
	}
}

func TestMockResourceNsxtPolicyCNADelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCNAMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(cnaOrgID, cnaProjectID, cnaID).Return(nil)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())
		d.SetId(cnaID)

		err := resourceNsxtPolicyProjectCNADelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNADelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestUnitNsxt_PolicyProjectCNASchemaValidation(t *testing.T) {
	res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
	subnetPathSchema := res.Schema["subnet_path"]
	validate := subnetPathSchema.ValidateFunc

	// Negative cases: CIDRs and invalid paths must fail validation
	_, errs := validate("10.0.0.0/8", "subnet_path")
	assert.NotEmpty(t, errs, "CIDR string should be rejected by validatePolicyPath()")

	_, errs = validate("invalid-path", "subnet_path")
	assert.NotEmpty(t, errs, "Non-path string should be rejected")

	// Positive cases: Valid infra and project paths must pass validation
	_, errs = validate("/infra/tier-1s/t1", "subnet_path")
	assert.Empty(t, errs, "Valid policy path should be accepted")

	_, errs = validate("/orgs/default/projects/project1/vpcs/vpc1/subnets/sub1", "subnet_path")
	assert.Empty(t, errs, "Project-scoped subnet path should be accepted")

	// Verify private_network_attachment schema
	pnaSchema := res.Schema["private_network_attachment"]
	assert.NotNil(t, pnaSchema)
	assert.Equal(t, schema.TypeBool, pnaSchema.Type)
	assert.True(t, pnaSchema.Optional)
	assert.Equal(t, true, pnaSchema.Default)
}
