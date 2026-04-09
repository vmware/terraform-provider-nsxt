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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	customattributes "github.com/vmware/terraform-provider-nsxt/api/infra/context_profiles/custom_attributes"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	caMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/context_profiles/custom_attributes"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	customAttrKey   = model.PolicyCustomAttributes_KEY_DOMAIN_NAME
	customAttrValue = "example.com"
	customAttrID    = model.PolicyCustomAttributes_KEY_DOMAIN_NAME + "~example.com"
)

func TestMockResourceNsxtPolicyContextProfileCustomAttributeCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDefaultSDK := caMocks.NewMockDefaultClient(ctrl)
	defaultWrapper := &customattributes.PolicyCustomAttributesClientContext{
		Client:     mockDefaultSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDefaultClient
	defer func() { cliDefaultClient = originalCli }()
	cliDefaultClient = func(sessionContext utl.SessionContext, connector client.Connector) *customattributes.PolicyCustomAttributesClientContext {
		return defaultWrapper
	}

	res := resourceNsxtPolicyContextProfileCustomAttribute()

	t.Run("Create_success", func(t *testing.T) {
		mockDefaultSDK.EXPECT().Create(gomock.Any(), "add").Return(nil)
		resultCount := int64(1)
		mockDefaultSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{
				ResultCount: &resultCount,
				Results: []model.PolicyContextProfile{
					{
						Attributes: []model.PolicyAttributes{
							{Key: &customAttrKey, Value: []string{customAttrValue}},
						},
					},
				},
			}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"key":       customAttrKey,
			"attribute": customAttrValue,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileCustomAttributeCreate(d, m)
		require.NoError(t, err)
		assert.Equal(t, customAttrID, d.Id())
		assert.Equal(t, customAttrKey, d.Get("key"))
		assert.Equal(t, customAttrValue, d.Get("attribute"))
	})

	t.Run("Create_fails_when_Create_returns_error", func(t *testing.T) {
		mockDefaultSDK.EXPECT().Create(gomock.Any(), "add").Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"key":       customAttrKey,
			"attribute": customAttrValue,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileCustomAttributeCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyContextProfileCustomAttributeRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDefaultSDK := caMocks.NewMockDefaultClient(ctrl)
	defaultWrapper := &customattributes.PolicyCustomAttributesClientContext{
		Client:     mockDefaultSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDefaultClient
	defer func() { cliDefaultClient = originalCli }()
	cliDefaultClient = func(sessionContext utl.SessionContext, connector client.Connector) *customattributes.PolicyCustomAttributesClientContext {
		return defaultWrapper
	}

	res := resourceNsxtPolicyContextProfileCustomAttribute()

	t.Run("Read_success", func(t *testing.T) {
		resultCount := int64(1)
		mockDefaultSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{
				ResultCount: &resultCount,
				Results: []model.PolicyContextProfile{
					{
						Attributes: []model.PolicyAttributes{
							{Key: &customAttrKey, Value: []string{customAttrValue}},
						},
					},
				},
			}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(customAttrID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileCustomAttributeRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, customAttrKey, d.Get("key"))
		assert.Equal(t, customAttrValue, d.Get("attribute"))
	})

	t.Run("Read_fails_when_not_found", func(t *testing.T) {
		resultCount := int64(0)
		mockDefaultSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{
				ResultCount: &resultCount,
				Results:     []model.PolicyContextProfile{},
			}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(customAttrID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileCustomAttributeRead(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyContextProfileCustomAttributeDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDefaultSDK := caMocks.NewMockDefaultClient(ctrl)
	defaultWrapper := &customattributes.PolicyCustomAttributesClientContext{
		Client:     mockDefaultSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDefaultClient
	defer func() { cliDefaultClient = originalCli }()
	cliDefaultClient = func(sessionContext utl.SessionContext, connector client.Connector) *customattributes.PolicyCustomAttributesClientContext {
		return defaultWrapper
	}

	res := resourceNsxtPolicyContextProfileCustomAttribute()

	t.Run("Delete_success", func(t *testing.T) {
		resultCount := int64(1)
		mockDefaultSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{
				ResultCount: &resultCount,
				Results: []model.PolicyContextProfile{
					{
						Attributes: []model.PolicyAttributes{
							{Key: &customAttrKey, Value: []string{customAttrValue}},
						},
					},
				},
			}, nil)
		mockDefaultSDK.EXPECT().Create(gomock.Any(), "remove").Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(customAttrID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileCustomAttributeDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_Create_remove_returns_error", func(t *testing.T) {
		resultCount := int64(1)
		mockDefaultSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{
				ResultCount: &resultCount,
				Results: []model.PolicyContextProfile{
					{
						Attributes: []model.PolicyAttributes{
							{Key: &customAttrKey, Value: []string{customAttrValue}},
						},
					},
				},
			}, nil)
		mockDefaultSDK.EXPECT().Create(gomock.Any(), "remove").Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(customAttrID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileCustomAttributeDelete(d, m)
		require.Error(t, err)
	})
}
