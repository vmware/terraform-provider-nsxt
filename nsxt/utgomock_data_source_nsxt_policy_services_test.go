//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	svcmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
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
}
