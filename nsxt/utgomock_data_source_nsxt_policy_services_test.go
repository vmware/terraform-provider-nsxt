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
	orgs "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs"
	projects "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	svcmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	orgsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs"
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
	isDefaultTrue := true
	isDefaultFalse := false

	t.Run("default project, built_in_only = false (default behavior)", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: &svcName, Path: &svcPath},
				{DisplayName: ptrString("svc-def"), Path: ptrString("/infra/services/svc-def"), IsDefault: &isDefaultTrue},
			},
		}, nil)

		ds := dataSourceNsxtPolicyServices()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		err := dataSourceNsxtPolicyServicesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		assert.Len(t, items, 2)
		assert.Equal(t, svcPath, items[svcName].(string))
		assert.Equal(t, "/infra/services/svc-def", items["svc-def"].(string))
		assert.NotEmpty(t, d.Id())
	})

	t.Run("default project, built_in_only = true", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: &svcName, Path: &svcPath},
				{DisplayName: ptrString("svc-def"), Path: ptrString("/infra/services/svc-def"), IsDefault: &isDefaultTrue},
			},
		}, nil)

		ds := dataSourceNsxtPolicyServices()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"built_in_only": true,
		})
		err := dataSourceNsxtPolicyServicesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		assert.Len(t, items, 2)
		assert.Equal(t, svcPath, items[svcName].(string))
		assert.Equal(t, "/infra/services/svc-def", items["svc-def"].(string))
	})

	t.Run("user project, built_in_only = false (default behavior)", func(t *testing.T) {
		mockProjectSDK := svcmocks.NewMockServicesClient(ctrl)
		wrapperProject := &cliinfra.ServiceClientContext{
			Client:     mockProjectSDK,
			ClientType: utl.Local,
		}

		originalCli := cliServicesClient
		cliServicesClient = func(sessCtx utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.ServiceClientContext {
			return wrapperProject
		}
		defer func() { cliServicesClient = originalCli }()

		mockProjectSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: ptrString("proj-svc"), Path: ptrString("/orgs/default/projects/project1/infra/services/proj-svc")},
				{DisplayName: ptrString("proj-svc-def"), Path: ptrString("/orgs/default/projects/project1/infra/services/proj-svc-def"), IsDefault: &isDefaultTrue},
				{DisplayName: ptrString("proj-svc-custom"), Path: ptrString("/orgs/default/projects/project1/infra/services/proj-svc-custom"), IsDefault: &isDefaultFalse},
			},
		}, nil)

		ds := dataSourceNsxtPolicyServices()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
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
		assert.Len(t, items, 3)
		assert.Contains(t, items, "proj-svc")
		assert.Contains(t, items, "proj-svc-custom")
		assert.Contains(t, items, "proj-svc-def")
	})

	t.Run("user project, built_in_only = true", func(t *testing.T) {
		mockDefaultSDK := svcmocks.NewMockServicesClient(ctrl)
		wrapperDefault := &cliinfra.ServiceClientContext{
			Client:     mockDefaultSDK,
			ClientType: utl.Local,
		}

		originalCli := cliServicesClient
		cliServicesClient = func(sessCtx utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.ServiceClientContext {
			return wrapperDefault
		}
		defer func() { cliServicesClient = originalCli }()

		mockDefaultSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: ptrString("proj-svc-custom"), Path: ptrString("/orgs/default/projects/project1/infra/services/proj-svc-custom"), IsDefault: &isDefaultFalse},
				{DisplayName: ptrString("proj-svc-built-in"), Path: ptrString("/orgs/default/projects/project1/infra/services/proj-svc-built-in"), IsDefault: &isDefaultTrue},
			},
		}, nil)

		ds := dataSourceNsxtPolicyServices()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"built_in_only": true,
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
		assert.Len(t, items, 1)
		assert.Equal(t, "/orgs/default/projects/project1/infra/services/proj-svc-built-in", items["proj-svc-built-in"].(string))
	})

	t.Run("user project, include_shared_services = true", func(t *testing.T) {
		mockProjectSDK := svcmocks.NewMockServicesClient(ctrl)
		wrapperProject := &cliinfra.ServiceClientContext{
			Client:     mockProjectSDK,
			ClientType: utl.Local,
		}

		wrapperDefault := &cliinfra.ServiceClientContext{
			Client:     mockSDK,
			ClientType: utl.Local,
		}

		originalCli := cliServicesClient
		cliServicesClient = func(sessCtx utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.ServiceClientContext {
			if sessCtx.ProjectID != "" {
				return wrapperProject
			}
			return wrapperDefault
		}
		defer func() { cliServicesClient = originalCli }()

		mockProjectSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: ptrString("proj-svc"), Path: ptrString("/orgs/default/projects/project1/infra/services/proj-svc")},
			},
		}, nil)

		resourceType := "Service"
		orgSharedSvcPath := "/infra/services/org-shared-svc"
		orgSharedSvcName := "org-shared-svc"
		sharedSvcPath := "/infra/services/shared-svc"
		sharedSvcName := "shared-svc"

		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: &orgSharedSvcName, Path: &orgSharedSvcPath},
				{DisplayName: &sharedSvcName, Path: &sharedSvcPath},
			},
		}, nil)

		mockOrgShared := orgsmocks.NewMockSharedWithMeClient(ctrl)
		originalOrgShared := cliOrgSharedWithMeClient
		cliOrgSharedWithMeClient = func(_ vapiProtocolClient.Connector) orgs.SharedWithMeClient {
			return mockOrgShared
		}
		defer func() { cliOrgSharedWithMeClient = originalOrgShared }()

		mockShared := projectsmocks.NewMockSharedWithMeClient(ctrl)
		originalShared := cliProjectSharedWithMeClient
		cliProjectSharedWithMeClient = func(_ vapiProtocolClient.Connector) projects.SharedWithMeClient {
			return mockShared
		}
		defer func() { cliProjectSharedWithMeClient = originalShared }()

		mockOrgShared.EXPECT().List("default", &resourceType).Return(nsxModel.SharedResourceListResult{
			Results: []nsxModel.SharedResource{
				{
					DisplayName: &orgSharedSvcName,
					ResourceObjects: []nsxModel.ResourceObject{
						{
							ResourcePath: &orgSharedSvcPath,
						},
					},
				},
			},
		}, nil)

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
		assert.Len(t, items, 3)
		assert.Equal(t, "/orgs/default/projects/project1/infra/services/proj-svc", items["proj-svc"])
		assert.Equal(t, orgSharedSvcPath, items[orgSharedSvcName])
		assert.Equal(t, sharedSvcPath, items[sharedSvcName])
	})

	t.Run("user project, include_shared_services = true, built_in_only = true", func(t *testing.T) {
		mockProjectSDK := svcmocks.NewMockServicesClient(ctrl)
		wrapperProject := &cliinfra.ServiceClientContext{
			Client:     mockProjectSDK,
			ClientType: utl.Local,
		}

		wrapperDefault := &cliinfra.ServiceClientContext{
			Client:     mockSDK,
			ClientType: utl.Local,
		}

		originalCli := cliServicesClient
		cliServicesClient = func(sessCtx utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.ServiceClientContext {
			if sessCtx.ProjectID != "" {
				return wrapperProject
			}
			return wrapperDefault
		}
		defer func() { cliServicesClient = originalCli }()

		// Project-level services query (since built_in_only=true, non-default custom services will be filtered out)
		mockProjectSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: ptrString("proj-svc-custom"), Path: ptrString("/orgs/default/projects/project1/infra/services/proj-svc-custom"), IsDefault: &isDefaultFalse},
			},
		}, nil)

		// Expect list call to fetch default space built-in services to filter shared services
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{
			Results: []nsxModel.Service{
				{DisplayName: ptrString("svc-def-1"), Path: ptrString("/infra/services/svc-def-1"), IsDefault: &isDefaultTrue},
				{DisplayName: ptrString("svc-custom"), Path: ptrString("/infra/services/svc-custom"), IsDefault: &isDefaultFalse},
			},
		}, nil)

		mockOrgShared := orgsmocks.NewMockSharedWithMeClient(ctrl)
		originalOrgShared := cliOrgSharedWithMeClient
		cliOrgSharedWithMeClient = func(_ vapiProtocolClient.Connector) orgs.SharedWithMeClient {
			return mockOrgShared
		}
		defer func() { cliOrgSharedWithMeClient = originalOrgShared }()

		mockShared := projectsmocks.NewMockSharedWithMeClient(ctrl)
		originalShared := cliProjectSharedWithMeClient
		cliProjectSharedWithMeClient = func(_ vapiProtocolClient.Connector) projects.SharedWithMeClient {
			return mockShared
		}
		defer func() { cliProjectSharedWithMeClient = originalShared }()

		resourceType := "Service"
		// This one is built-in (path is /infra/services/svc-def-1)
		orgSharedSvcPath := "/infra/services/svc-def-1"
		orgSharedSvcName := "svc-def-1"
		mockOrgShared.EXPECT().List("default", &resourceType).Return(nsxModel.SharedResourceListResult{
			Results: []nsxModel.SharedResource{
				{
					DisplayName: &orgSharedSvcName,
					ResourceObjects: []nsxModel.ResourceObject{
						{
							ResourcePath: &orgSharedSvcPath,
						},
					},
				},
			},
		}, nil)

		// This one is NOT built-in (path is /infra/services/svc-custom)
		sharedSvcPath := "/infra/services/svc-custom"
		sharedSvcName := "svc-custom"
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
			"built_in_only":           true,
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
		assert.Len(t, items, 1)
		assert.Contains(t, items, "svc-def-1")
		assert.NotContains(t, items, "svc-custom")
	})

	t.Run("user project, include_shared_services = true, shared client returns NotFound error", func(t *testing.T) {
		mockProjectSDK := svcmocks.NewMockServicesClient(ctrl)
		wrapperProject := &cliinfra.ServiceClientContext{
			Client:     mockProjectSDK,
			ClientType: utl.Local,
		}

		wrapperDefault := &cliinfra.ServiceClientContext{
			Client:     mockSDK,
			ClientType: utl.Local,
		}

		originalCli := cliServicesClient
		cliServicesClient = func(sessCtx utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.ServiceClientContext {
			if sessCtx.ProjectID != "" {
				return wrapperProject
			}
			return wrapperDefault
		}
		defer func() { cliServicesClient = originalCli }()

		mockProjectSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{}, nil)
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{}, nil)

		mockOrgShared := orgsmocks.NewMockSharedWithMeClient(ctrl)
		originalOrgShared := cliOrgSharedWithMeClient
		cliOrgSharedWithMeClient = func(_ vapiProtocolClient.Connector) orgs.SharedWithMeClient {
			return mockOrgShared
		}
		defer func() { cliOrgSharedWithMeClient = originalOrgShared }()

		mockShared := projectsmocks.NewMockSharedWithMeClient(ctrl)
		originalShared := cliProjectSharedWithMeClient
		cliProjectSharedWithMeClient = func(_ vapiProtocolClient.Connector) projects.SharedWithMeClient {
			return mockShared
		}
		defer func() { cliProjectSharedWithMeClient = originalShared }()

		resourceType := "Service"
		mockOrgShared.EXPECT().List("default", &resourceType).Return(nsxModel.SharedResourceListResult{}, vapiErrors.NotFound{})
		mockShared.EXPECT().List("default", "project1", &resourceType).Return(nsxModel.SharedResourceListResult{}, vapiErrors.NotFound{})

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
		assert.Len(t, items, 0)
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ServiceListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyServices()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		err := dataSourceNsxtPolicyServicesRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func ptrString(s string) *string {
	return &s
}
