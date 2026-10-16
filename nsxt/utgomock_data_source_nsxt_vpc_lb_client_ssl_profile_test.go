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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
)

// setupVpcLbClientSslMock and the vpcLbClientSsl* fixtures are defined in
// utgomock_resource_nsxt_vpc_lb_client_ssl_profile_test.go

func vpcLbDsClientSslData(extra map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{"context": vpcLbPrfContext()}
	for k, v := range extra {
		res[k] = v
	}
	return res
}

func vpcLbDsClientSslItem(id, name string) nsxModel.LBClientSslProfile {
	path := "/orgs/default/projects/" + vpcLbPrfProject + "/vpcs/" + vpcLbPrfVpc + "/vpc-lb-client-ssl-profiles/" + id
	return nsxModel.LBClientSslProfile{Id: &id, DisplayName: &name, Description: &vpcLbClientSslDescription, Path: &path}
}

func TestMockDataSourceNsxtVpcLBClientSslProfile(t *testing.T) {
	type testCase struct {
		name    string
		version string
		input   map[string]interface{}
		setup   func(m *vpcsmocks.MockVpcLbClientSslProfilesClient)
		errMsg  string
		wantID  string
	}
	list := func(m *vpcsmocks.MockVpcLbClientSslProfilesClient, results []nsxModel.LBClientSslProfile, err error) {
		m.EXPECT().List(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBClientSslProfileListResult{Results: results}, err)
	}
	cases := []testCase{
		{
			name: "by ID", input: map[string]interface{}{"id": vpcLbClientSslID},
			setup: func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {
				m.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(vpcLbClientSslAPIResponse(), nil)
			},
			wantID: vpcLbClientSslID,
		},
		{
			name: "by ID not found", input: map[string]interface{}{"id": vpcLbClientSslID},
			setup: func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {
				m.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(nsxModel.LBClientSslProfile{}, vapiErrors.NotFound{})
			},
			errMsg: "Failed to read",
		},
		{
			name: "by ID API error", input: map[string]interface{}{"id": vpcLbClientSslID},
			setup: func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {
				m.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(nsxModel.LBClientSslProfile{}, vapiErrors.InternalServerError{})
			},
			errMsg: "LBClientSslProfile",
		},
		{
			name: "by name perfect match", input: map[string]interface{}{"display_name": vpcLbClientSslDisplayName},
			setup: func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {
				list(m, []nsxModel.LBClientSslProfile{vpcLbDsClientSslItem("other", vpcLbClientSslDisplayName+" extra"), vpcLbDsClientSslItem(vpcLbClientSslID, vpcLbClientSslDisplayName)}, nil)
			},
			wantID: vpcLbClientSslID,
		},
		{
			name: "by name prefix match", input: map[string]interface{}{"display_name": "Test VPC LB"},
			setup: func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {
				list(m, []nsxModel.LBClientSslProfile{vpcLbDsClientSslItem(vpcLbClientSslID, vpcLbClientSslDisplayName)}, nil)
			},
			wantID: vpcLbClientSslID,
		},
		{
			name: "multiple perfect matches", input: map[string]interface{}{"display_name": "dup"},
			setup: func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {
				list(m, []nsxModel.LBClientSslProfile{vpcLbDsClientSslItem("d1", "dup"), vpcLbDsClientSslItem("d2", "dup")}, nil)
			},
			errMsg: "multiple",
		},
		{
			name: "multiple prefix matches", input: map[string]interface{}{"display_name": "prefix"},
			setup: func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {
				list(m, []nsxModel.LBClientSslProfile{vpcLbDsClientSslItem("p1", "prefix-a"), vpcLbDsClientSslItem("p2", "prefix-b")}, nil)
			},
			errMsg: "multiple",
		},
		{
			name: "name not found", input: map[string]interface{}{"display_name": "nonexistent"},
			setup: func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {
				list(m, nil, nil)
			},
			errMsg: "not found",
		},
		{
			name: "list API error", input: map[string]interface{}{"display_name": "x"},
			setup: func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {
				list(m, nil, vapiErrors.InternalServerError{})
			},
			errMsg: "LBClientSslProfile",
		},
		{
			name: "neither id nor name", input: map[string]interface{}{},
			setup:  func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {},
			errMsg: "Error obtaining",
		},
		{
			name: "NSX version lower than 9.0.2", version: "9.0.1", input: map[string]interface{}{"id": vpcLbClientSslID},
			setup:  func(m *vpcsmocks.MockVpcLbClientSslProfilesClient) {},
			errMsg: "9.0.2",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			version := tc.version
			if version == "" {
				version = "9.0.2"
			}
			vpcLbPrfSetNsxVersion(t, version)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbClientSslMock(t, ctrl)
			defer restore()
			tc.setup(mockSDK)

			ds := dataSourceNsxtVpcLBClientSslProfile()
			d := schema.TestResourceDataRaw(t, ds.Schema, vpcLbDsClientSslData(tc.input))
			err := ds.Read(d, newGoMockProviderClient())
			if tc.errMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantID, d.Id())
			assert.NotEmpty(t, d.Get("display_name"))
			assert.Equal(t, vpcLbClientSslDescription, d.Get("description"))
		})
	}

	t.Run("context is required", func(t *testing.T) {
		ds := dataSourceNsxtVpcLBClientSslProfile()
		ctxSchema, ok := ds.Schema["context"]
		require.True(t, ok)
		assert.True(t, ctxSchema.Required)
	})
}
