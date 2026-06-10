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
	projects "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	svcmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	projectsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
)

func setupPolicyServicesDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*svcmocks.MockServicesClient, func()) {
	t.Helper()
	mockSDK := svcmocks.NewMockServicesClient(ctrl)
	wrapper := &cliinfra.ServiceClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliServicesClient
	cliServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.ServiceClientContext {
		return wrapper
	}
	return mockSDK, func() { cliServicesClient = orig }
}

func TestMockDataSourceNsxtPolicyServicesRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupPolicyServicesDataSourceMock(t, ctrl)
	defer restore()

	svcName := "svc-a"
	svcPath := "/infra/services/svc-a"

	t.Run("success", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: &svcName, Path: &svcPath},
			},
		}, nil)

		ds := dataSourceNsxtPolicyServices()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		err := dataSourceNsxtPolicyServicesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		require.Contains(t, items, svcName)
		assert.Equal(t, svcPath, items[svcName].(string))
		assert.NotEmpty(t, d.Id())
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyServices()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		err := dataSourceNsxtPolicyServicesRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("built_in_only", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: &svcName, Path: &svcPath},
			},
		}, nil)

		ds := dataSourceNsxtPolicyServices()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"built_in_only": true,
		})
		err := dataSourceNsxtPolicyServicesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		require.Contains(t, items, svcName)
		assert.Equal(t, svcPath, items[svcName].(string))
	})

	t.Run("include_shared_services", func(t *testing.T) {
		mockProjectSDK := svcmocks.NewMockServicesClient(ctrl)
		mockDefaultSDK := svcmocks.NewMockServicesClient(ctrl)

		wrapperProject := &cliinfra.ServiceClientContext{
			Client:     mockProjectSDK,
			ClientType: utl.Local,
		}
		wrapperDefault := &cliinfra.ServiceClientContext{
			Client:     mockDefaultSDK,
			ClientType: utl.Local,
		}

		originalCli := cliServicesClient
		cliServicesClient = func(sessCtx utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.ServiceClientContext {
			if sessCtx.ProjectID == "project1" {
				return wrapperProject
			}
			return wrapperDefault
		}
		defer func() { cliServicesClient = originalCli }()

		projSvcName := "proj-svc"
		projSvcPath := "/orgs/default/projects/project1/infra/services/proj-svc"
		mockProjectSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: &projSvcName, Path: &projSvcPath},
			},
		}, nil)

		mockDefaultSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: &svcName, Path: &svcPath},
			},
		}, nil)

		mockShared := projectsmocks.NewMockSharedWithMeClient(ctrl)
		originalShared := cliSharedWithMeClient
		cliSharedWithMeClient = func(_ vapiProtocolClient.Connector) projects.SharedWithMeClient {
			return mockShared
		}
		defer func() { cliSharedWithMeClient = originalShared }()

		sharedSvcPath := "/infra/services/shared-svc"
		sharedSvcName := "shared-svc"
		resourceType := "Service"
		mockShared.EXPECT().List("default", "project1", &resourceType).Return(nsxModel.SharedResourceListResult{
			Results: []nsxModel.SharedResource{
				{
					DisplayName: &sharedSvcName,
					ResourceObjects: []nsxModel.ResourceObject{
						{
							ResourcePath: &sharedSvcPath,
						},
					},
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyServices()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"include_shared_services": true,
			"context": []interface{}{
				map[string]interface{}{
					"project_id":  "project1",
					"from_global": false,
				},
			},
		})

		err := dataSourceNsxtPolicyServicesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		assert.Equal(t, projSvcPath, items[projSvcName])
		assert.Equal(t, svcPath, items[svcName])
		assert.Equal(t, sharedSvcPath, items[sharedSvcName])
	})
}
