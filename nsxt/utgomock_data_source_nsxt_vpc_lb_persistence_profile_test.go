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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
)

// setupVpcLbPersistenceMock, vpcLbPrf* helpers are defined in
// utgomock_resource_nsxt_vpc_lb_cookie_persistence_profile_test.go
// lbPersistenceProfileStructValue is defined in utgomock_resource_nsxt_policy_lb_cookie_persistence_profile_test.go

const (
	vpcLbDsPersistID   = "vpc-lb-persist-1"
	vpcLbDsPersistName = "Test VPC LB Persistence Profile"
)

func vpcLbDsPersistValue(t *testing.T, id, displayName, resourceType string) *data.StructValue {
	t.Helper()
	path := "/orgs/default/projects/" + vpcLbPrfProject + "/vpcs/" + vpcLbPrfVpc + "/vpc-lb-persistence-profiles/" + id
	return lbPersistenceProfileStructValue(t, nsxModel.LBPersistenceProfile{
		Id:           &id,
		DisplayName:  &displayName,
		Path:         &path,
		ResourceType: resourceType,
	}, nsxModel.LBPersistenceProfileBindingType())
}

func vpcLbDsPersistData(extra map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{"context": vpcLbPrfContext()}
	for k, v := range extra {
		res[k] = v
	}
	return res
}

func vpcLbDsPersistList(mockSDK *vpcsmocks.MockVpcLbPersistenceProfilesClient, results []*data.StructValue, err error) {
	mockSDK.EXPECT().List(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nsxModel.LBPersistenceProfileListResult{Results: results}, err)
}

func TestMockDataSourceNsxtVpcLBPersistenceProfile(t *testing.T) {
	cookie := nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE
	sourceIP := nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBSOURCEIPPERSISTENCEPROFILE

	type testCase struct {
		name    string
		version string
		input   map[string]interface{}
		setup   func(t *testing.T, mockSDK *vpcsmocks.MockVpcLbPersistenceProfilesClient)
		errMsg  string // empty means success
		wantID  string
		wantTyp string
	}
	cases := []testCase{
		{
			name: "by ID", input: map[string]interface{}{"id": vpcLbDsPersistID},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				m.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbDsPersistID).
					Return(vpcLbDsPersistValue(t, vpcLbDsPersistID, vpcLbDsPersistName, cookie), nil)
			},
			wantID: vpcLbDsPersistID, wantTyp: "COOKIE",
		},
		{
			name: "by ID API error", input: map[string]interface{}{"id": vpcLbDsPersistID},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				m.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbDsPersistID).Return(nil, vapiErrors.InternalServerError{})
			},
			errMsg: "VPC LBPersistenceProfile",
		},
		{
			name: "by ID not found", input: map[string]interface{}{"id": vpcLbDsPersistID},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				m.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbDsPersistID).Return(nil, vapiErrors.NotFound{})
			},
			errMsg: "Failed to read",
		},
		{
			name: "by ID conversion error", input: map[string]interface{}{"id": vpcLbDsPersistID},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				m.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbDsPersistID).Return(&data.StructValue{}, nil)
			},
			errMsg: "resource_type",
		},
		{
			name: "by name perfect match", input: map[string]interface{}{"display_name": vpcLbDsPersistName},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				vpcLbDsPersistList(m, []*data.StructValue{
					vpcLbDsPersistValue(t, "other", "Other profile", sourceIP),
					vpcLbDsPersistValue(t, vpcLbDsPersistID, vpcLbDsPersistName, cookie),
				}, nil)
			},
			wantID: vpcLbDsPersistID, wantTyp: "COOKIE",
		},
		{
			name: "by name prefix match with matching type", input: map[string]interface{}{"display_name": "Test VPC LB", "type": "COOKIE"},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				vpcLbDsPersistList(m, []*data.StructValue{vpcLbDsPersistValue(t, vpcLbDsPersistID, vpcLbDsPersistName, cookie)}, nil)
			},
			wantID: vpcLbDsPersistID, wantTyp: "COOKIE",
		},
		{
			name: "by name prefix match filtered out by type", input: map[string]interface{}{"display_name": "Test VPC LB", "type": "SOURCE_IP"},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				vpcLbDsPersistList(m, []*data.StructValue{vpcLbDsPersistValue(t, vpcLbDsPersistID, vpcLbDsPersistName, cookie)}, nil)
			},
			errMsg: "not found",
		},
		{
			name: "by type only", input: map[string]interface{}{"type": "SOURCE_IP"},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				vpcLbDsPersistList(m, []*data.StructValue{
					vpcLbDsPersistValue(t, vpcLbDsPersistID, vpcLbDsPersistName, cookie),
					vpcLbDsPersistValue(t, "src-1", "source profile", sourceIP),
				}, nil)
			},
			wantID: "src-1", wantTyp: "SOURCE_IP",
		},
		{
			name: "multiple perfect matches", input: map[string]interface{}{"display_name": "dup"},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				vpcLbDsPersistList(m, []*data.StructValue{
					vpcLbDsPersistValue(t, "dup-1", "dup", cookie),
					vpcLbDsPersistValue(t, "dup-2", "dup", cookie),
				}, nil)
			},
			errMsg: "multiple",
		},
		{
			name: "multiple prefix matches", input: map[string]interface{}{"display_name": "prefix"},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				vpcLbDsPersistList(m, []*data.StructValue{
					vpcLbDsPersistValue(t, "pfx-1", "prefix-a", cookie),
					vpcLbDsPersistValue(t, "pfx-2", "prefix-b", sourceIP),
				}, nil)
			},
			errMsg: "multiple",
		},
		{
			name: "name not found", input: map[string]interface{}{"display_name": "nonexistent"},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				vpcLbDsPersistList(m, nil, nil)
			},
			errMsg: "not found",
		},
		{
			name: "list API error", input: map[string]interface{}{"display_name": "x"},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				vpcLbDsPersistList(m, nil, vapiErrors.InternalServerError{})
			},
			errMsg: "VPC LBPersistenceProfile",
		},
		{
			name: "list conversion error", input: map[string]interface{}{"display_name": "x"},
			setup: func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {
				vpcLbDsPersistList(m, []*data.StructValue{{}}, nil)
			},
			errMsg: "resource_type",
		},
		{
			name: "NSX version lower than 9.0.0", version: "4.2.0", input: map[string]interface{}{"id": vpcLbDsPersistID},
			setup:  func(t *testing.T, m *vpcsmocks.MockVpcLbPersistenceProfilesClient) {},
			errMsg: "9.0.0",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			version := tc.version
			if version == "" {
				version = "9.0.0"
			}
			vpcLbPrfSetNsxVersion(t, version)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()
			tc.setup(t, mockSDK)

			ds := dataSourceNsxtVpcLBPersistenceProfile()
			d := schema.TestResourceDataRaw(t, ds.Schema, vpcLbDsPersistData(tc.input))
			err := ds.Read(d, newGoMockProviderClient())
			if tc.errMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantID, d.Id())
			assert.Equal(t, tc.wantTyp, d.Get("type"))
			assert.NotEmpty(t, d.Get("display_name"))
			assert.Contains(t, d.Get("path"), "/vpcs/"+vpcLbPrfVpc+"/vpc-lb-persistence-profiles/")
		})
	}

	t.Run("none of id, name or type set", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		ds := dataSourceNsxtVpcLBPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, vpcLbDsPersistData(nil))
		require.NoError(t, d.Set("type", ""))
		err := ds.Read(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining")
	})

	t.Run("context is required", func(t *testing.T) {
		ds := dataSourceNsxtVpcLBPersistenceProfile()
		ctxSchema, ok := ds.Schema["context"]
		require.True(t, ok)
		assert.True(t, ctxSchema.Required)
	})
}
