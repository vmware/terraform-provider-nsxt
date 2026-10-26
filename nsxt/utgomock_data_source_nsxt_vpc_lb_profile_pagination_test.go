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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
)

func TestMockDataSourceNsxtVpcLBProfilePagination(t *testing.T) {
	includeMarkedForDelete := false
	type profile struct{ id, name string }
	type page struct {
		profiles []profile
		cursor   *string
		unnamed  bool
		err      error
	}
	specs := []struct {
		name          string
		newDataSource func() *schema.Resource
		typeName      string
		setup         func(t *testing.T, ctrl *gomock.Controller, pages []page)
	}{
		{
			name: "App", newDataSource: dataSourceNsxtVpcLBAppProfile, typeName: "TCP",
			setup: func(t *testing.T, ctrl *gomock.Controller, pages []page) {
				m := setupVpcLbAppProfileMock(ctrl)
				t.Cleanup(m.restore)
				mockSDK, ok := cliVpcLbAppProfilesClient(utl.SessionContext{}, nil).Client.(*vpcsmocks.MockVpcLbAppProfilesClient)
				require.True(t, ok)
				var cursor *string
				var calls []interface{}
				for _, pg := range pages {
					var results []*data.StructValue
					for _, p := range pg.profiles {
						item := lbAppProfileStructValue(t, p.id, p.name, model.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE)
						if pg.unnamed {
							item.SetField("display_name", data.NewOptionalValue(nil))
						}
						results = append(results, item)
					}
					calls = append(calls, mockSDK.EXPECT().List(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc,
						cursor, nil, gomock.Eq(&includeMarkedForDelete), nil, nil, nil, nil).
						Return(model.LBAppProfileListResult{Results: results, Cursor: pg.cursor}, pg.err))
					cursor = pg.cursor
				}
				gomock.InOrder(calls...)
			},
		},
		{
			name: "Monitor", newDataSource: dataSourceNsxtVpcLBMonitor, typeName: "TCP",
			setup: func(t *testing.T, ctrl *gomock.Controller, pages []page) {
				m := setupVpcLbMonitorProfileMock(ctrl)
				t.Cleanup(m.restore)
				mockSDK, ok := cliVpcLbMonitorProfilesClient(utl.SessionContext{}, nil).Client.(*vpcsmocks.MockVpcLbMonitorProfilesClient)
				require.True(t, ok)
				var cursor *string
				var calls []interface{}
				for _, pg := range pages {
					var results []*data.StructValue
					for _, p := range pg.profiles {
						item := lbMonitorProfileStructValue(t, p.id, p.name, model.LBMonitorProfile_RESOURCE_TYPE_LBTCPMONITORPROFILE)
						if pg.unnamed {
							item.SetField("display_name", data.NewOptionalValue(nil))
						}
						results = append(results, item)
					}
					calls = append(calls, mockSDK.EXPECT().List(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc,
						cursor, nil, gomock.Eq(&includeMarkedForDelete), nil, nil, nil, nil).
						Return(model.LBMonitorProfileListResult{Results: results, Cursor: pg.cursor}, pg.err))
					cursor = pg.cursor
				}
				gomock.InOrder(calls...)
			},
		},
		{
			name: "Persistence", newDataSource: dataSourceNsxtVpcLBPersistenceProfile, typeName: "COOKIE",
			setup: func(t *testing.T, ctrl *gomock.Controller, pages []page) {
				mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
				t.Cleanup(restore)
				var cursor *string
				var calls []interface{}
				for _, pg := range pages {
					var results []*data.StructValue
					for _, p := range pg.profiles {
						item := vpcLbDsPersistValue(t, p.id, p.name, model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE)
						if pg.unnamed {
							item.SetField("display_name", data.NewOptionalValue(nil))
						}
						results = append(results, item)
					}
					calls = append(calls, mockSDK.EXPECT().List(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc,
						cursor, nil, gomock.Eq(&includeMarkedForDelete), nil, nil, nil, nil).
						Return(model.LBPersistenceProfileListResult{Results: results, Cursor: pg.cursor}, pg.err))
					cursor = pg.cursor
				}
				gomock.InOrder(calls...)
			},
		},
		{
			name: "ClientSsl", newDataSource: dataSourceNsxtVpcLBClientSslProfile, typeName: "",
			setup: func(t *testing.T, ctrl *gomock.Controller, pages []page) {
				mockSDK, restore := setupVpcLbClientSslMock(t, ctrl)
				t.Cleanup(restore)
				var cursor *string
				var calls []interface{}
				for _, pg := range pages {
					var results []model.LBClientSslProfile
					for _, p := range pg.profiles {
						item := vpcLbDsClientSslItem(p.id, p.name)
						if pg.unnamed {
							item.DisplayName = nil
						}
						results = append(results, item)
					}
					calls = append(calls, mockSDK.EXPECT().List(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc,
						cursor, nil, gomock.Eq(&includeMarkedForDelete), nil, nil, nil, nil).
						Return(model.LBClientSslProfileListResult{Results: results, Cursor: pg.cursor}, pg.err))
					cursor = pg.cursor
				}
				gomock.InOrder(calls...)
			},
		},
		{
			name: "ServerSsl", newDataSource: dataSourceNsxtVpcLBServerSslProfile, typeName: "",
			setup: func(t *testing.T, ctrl *gomock.Controller, pages []page) {
				mockSDK, restore := setupVpcLbServerSslMock(t, ctrl)
				t.Cleanup(restore)
				var cursor *string
				var calls []interface{}
				for _, pg := range pages {
					var results []model.LBServerSslProfile
					for _, p := range pg.profiles {
						item := vpcLbDsServerSslItem(p.id, p.name)
						if pg.unnamed {
							item.DisplayName = nil
						}
						results = append(results, item)
					}
					calls = append(calls, mockSDK.EXPECT().List(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc,
						cursor, nil, gomock.Eq(&includeMarkedForDelete), nil, nil, nil, nil).
						Return(model.LBServerSslProfileListResult{Results: results, Cursor: pg.cursor}, pg.err))
					cursor = pg.cursor
				}
				gomock.InOrder(calls...)
			},
		},
	}
	next, last, empty := "next-page", "last-page", ""
	cases := []struct {
		name      string
		pages     []page
		inputName string
		typeOnly  bool
		wantID    string
		errMsg    string
	}{
		{
			name: "Skip unnamed profiles before a later name match", inputName: "target", wantID: "target-id",
			pages: []page{
				{profiles: []profile{{"unnamed", ""}}, unnamed: true, cursor: &next},
				{profiles: []profile{{"target-id", "target"}}},
			},
		},
		{
			name: "Type-only lookup matches unnamed profiles", typeOnly: true, wantID: "target-id",
			pages: []page{
				{cursor: &next},
				{profiles: []profile{{"target-id", ""}}, unnamed: true},
			},
		},
		{
			name: "Exact match on third page", inputName: "target", wantID: "target-id",
			pages: []page{
				{profiles: []profile{{"other", "other"}}, cursor: &next},
				{cursor: &last},
				{profiles: []profile{{"target-id", "target"}}, cursor: &empty},
			},
		},
		{
			name: "Prefix match on later page", inputName: "target", wantID: "target-id",
			pages: []page{
				{cursor: &next},
				{profiles: []profile{{"target-id", "target-long"}}},
			},
		},
		{
			name: "Later exact match beats earlier prefix", inputName: "target", wantID: "target-id",
			pages: []page{
				{profiles: []profile{{"prefix-id", "target-long"}}, cursor: &next},
				{profiles: []profile{{"target-id", "target"}}},
			},
		},
		{
			name: "Duplicate exact matches across pages", inputName: "target", errMsg: "multiple",
			pages: []page{
				{profiles: []profile{{"one", "target"}}, cursor: &next},
				{profiles: []profile{{"two", "target"}}},
			},
		},
		{
			name: "Duplicate prefix matches across pages", inputName: "target", errMsg: "multiple",
			pages: []page{
				{profiles: []profile{{"one", "target-one"}}, cursor: &next},
				{profiles: []profile{{"two", "target-two"}}},
			},
		},
		{
			name: "Later page error after a match", inputName: "target", errMsg: "Failed to read",
			pages: []page{
				{profiles: []profile{{"one", "target"}}, cursor: &next},
				{err: vapiErrors.InternalServerError{}},
			},
		},
		{
			name: "Not found after all pages", inputName: "target", errMsg: "not found",
			pages: []page{
				{profiles: []profile{{"one", "other"}}, cursor: &next},
				{},
			},
		},
		{
			name: "Type match on later page", typeOnly: true, wantID: "target-id",
			pages: []page{
				{cursor: &next},
				{profiles: []profile{{"target-id", "target"}}},
			},
		},
		{
			name: "Duplicate type matches across pages", typeOnly: true, errMsg: "multiple",
			pages: []page{
				{profiles: []profile{{"one", "first"}}, cursor: &next},
				{profiles: []profile{{"two", "second"}}},
			},
		},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			for _, tc := range cases {
				if tc.typeOnly && spec.typeName == "" {
					continue
				}
				t.Run(tc.name, func(t *testing.T) {
					vpcLbPrfSetNsxVersion(t, "9.0.0")
					ctrl := gomock.NewController(t)
					spec.setup(t, ctrl, tc.pages)
					raw := map[string]interface{}{"display_name": tc.inputName}
					if tc.typeOnly {
						raw["type"] = spec.typeName
					}
					ds := spec.newDataSource()
					d := schema.TestResourceDataRaw(t, ds.Schema, vpcLbWithContext(raw))
					err := ds.Read(d, newGoMockProviderClient())
					if tc.errMsg != "" {
						require.Error(t, err)
						assert.Contains(t, err.Error(), tc.errMsg)
						return
					}
					require.NoError(t, err)
					assert.Equal(t, tc.wantID, d.Id())
				})
			}
		})
	}
}
