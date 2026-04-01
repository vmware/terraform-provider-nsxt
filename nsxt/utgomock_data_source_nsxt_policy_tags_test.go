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

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

func setupTagsMock(t *testing.T, ctrl *gomock.Controller) (*inframocks.MockTagsClient, func()) {
	mockSDK := inframocks.NewMockTagsClient(ctrl)
	mockWrapper := &infraapi.TagsClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTagsClient
	cliTagsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.TagsClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTagsClient = original }
}

func TestMockDataSourceNsxtTagsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTagsMock(t, ctrl)
	defer restore()

	t.Run("Read success returns tags", func(t *testing.T) {
		scope := "env"
		tag1 := "production"
		tag2 := "staging"
		listResult := nsxModel.TagInfoListResult{
			Results: []nsxModel.TagInfo{
				{Tag: &tag1},
				{Tag: &tag2},
			},
		}
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(listResult, nil)

		ds := dataSourceNsxtTags()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"scope": scope,
		})

		err := dataSourceNsxtTagsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").([]interface{})
		assert.Len(t, items, 2)
		assert.Equal(t, tag1, items[0])
		assert.Equal(t, tag2, items[1])
	})

	t.Run("Read with empty results returns empty items", func(t *testing.T) {
		listResult := nsxModel.TagInfoListResult{
			Results: []nsxModel.TagInfo{},
		}
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(listResult, nil)

		ds := dataSourceNsxtTags()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"scope": "empty-scope",
		})

		err := dataSourceNsxtTagsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").([]interface{})
		assert.Empty(t, items)
	})

	t.Run("Read List API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.TagInfoListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtTags()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"scope": "test-scope",
		})

		err := dataSourceNsxtTagsRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
