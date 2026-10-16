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

// setupVpcLbServerSslMock and the vpcLbServerSsl* fixtures are defined in
// utgomock_resource_nsxt_vpc_lb_server_ssl_profile_test.go

func vpcLbDsServerSslData(extra map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{"context": vpcLbPrfContext()}
	for k, v := range extra {
		res[k] = v
	}
	return res
}

func vpcLbDsServerSslItem(id, name string) nsxModel.LBServerSslProfile {
	path := "/orgs/default/projects/" + vpcLbPrfProject + "/vpcs/" + vpcLbPrfVpc + "/vpc-lb-server-ssl-profiles/" + id
	return nsxModel.LBServerSslProfile{Id: &id, DisplayName: &name, Description: &vpcLbServerSslDescription, Path: &path}
}

func TestMockDataSourceNsxtVpcLBServerSslProfile(t *testing.T) {
	type testCase struct {
		name    string
		version string
		input   map[string]interface{}
		setup   func(m *vpcsmocks.MockVpcLbServerSslProfilesClient)
		errMsg  string
		wantID  string
	}
	list := func(m *vpcsmocks.MockVpcLbServerSslProfilesClient, results []nsxModel.LBServerSslProfile, err error) {
		m.EXPECT().List(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBServerSslProfileListResult{Results: results}, err)
	}
	cases := []testCase{
		{
			name: "by ID", input: map[string]interface{}{"id": vpcLbServerSslID},
			setup: func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {
				m.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(vpcLbServerSslAPIResponse(), nil)
			},
			wantID: vpcLbServerSslID,
		},
		{
			name: "by ID not found", input: map[string]interface{}{"id": vpcLbServerSslID},
			setup: func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {
				m.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(nsxModel.LBServerSslProfile{}, vapiErrors.NotFound{})
			},
			errMsg: "Failed to read",
		},
		{
			name: "by ID API error", input: map[string]interface{}{"id": vpcLbServerSslID},
			setup: func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {
				m.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(nsxModel.LBServerSslProfile{}, vapiErrors.InternalServerError{})
			},
			errMsg: "LBServerSslProfile",
		},
		{
			name: "by name perfect match", input: map[string]interface{}{"display_name": vpcLbServerSslDisplayName},
			setup: func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {
				list(m, []nsxModel.LBServerSslProfile{vpcLbDsServerSslItem("other", vpcLbServerSslDisplayName+" extra"), vpcLbDsServerSslItem(vpcLbServerSslID, vpcLbServerSslDisplayName)}, nil)
			},
			wantID: vpcLbServerSslID,
		},
		{
			name: "by name prefix match", input: map[string]interface{}{"display_name": "Test VPC LB"},
			setup: func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {
				list(m, []nsxModel.LBServerSslProfile{vpcLbDsServerSslItem(vpcLbServerSslID, vpcLbServerSslDisplayName)}, nil)
			},
			wantID: vpcLbServerSslID,
		},
		{
			name: "multiple perfect matches", input: map[string]interface{}{"display_name": "dup"},
			setup: func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {
				list(m, []nsxModel.LBServerSslProfile{vpcLbDsServerSslItem("d1", "dup"), vpcLbDsServerSslItem("d2", "dup")}, nil)
			},
			errMsg: "multiple",
		},
		{
			name: "multiple prefix matches", input: map[string]interface{}{"display_name": "prefix"},
			setup: func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {
				list(m, []nsxModel.LBServerSslProfile{vpcLbDsServerSslItem("p1", "prefix-a"), vpcLbDsServerSslItem("p2", "prefix-b")}, nil)
			},
			errMsg: "multiple",
		},
		{
			name: "name not found", input: map[string]interface{}{"display_name": "nonexistent"},
			setup: func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {
				list(m, nil, nil)
			},
			errMsg: "not found",
		},
		{
			name: "list API error", input: map[string]interface{}{"display_name": "x"},
			setup: func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {
				list(m, nil, vapiErrors.InternalServerError{})
			},
			errMsg: "LBServerSslProfile",
		},
		{
			name: "neither id nor name", input: map[string]interface{}{},
			setup:  func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {},
			errMsg: "Error obtaining",
		},
		{
			name: "NSX version lower than 9.0.2", version: "9.0.1", input: map[string]interface{}{"id": vpcLbServerSslID},
			setup:  func(m *vpcsmocks.MockVpcLbServerSslProfilesClient) {},
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
			mockSDK, restore := setupVpcLbServerSslMock(t, ctrl)
			defer restore()
			tc.setup(mockSDK)

			ds := dataSourceNsxtVpcLBServerSslProfile()
			d := schema.TestResourceDataRaw(t, ds.Schema, vpcLbDsServerSslData(tc.input))
			err := ds.Read(d, newGoMockProviderClient())
			if tc.errMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantID, d.Id())
			assert.NotEmpty(t, d.Get("display_name"))
			assert.Equal(t, vpcLbServerSslDescription, d.Get("description"))
		})
	}

	t.Run("context is required", func(t *testing.T) {
		ds := dataSourceNsxtVpcLBServerSslProfile()
		ctxSchema, ok := ds.Schema["context"]
		require.True(t, ok)
		assert.True(t, ctxSchema.Required)
	})
}
