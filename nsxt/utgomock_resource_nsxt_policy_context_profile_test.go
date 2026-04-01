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

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	contextprofiles "github.com/vmware/terraform-provider-nsxt/api/infra/context_profiles"
	customattributes "github.com/vmware/terraform-provider-nsxt/api/infra/context_profiles/custom_attributes"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	cpAttrMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/context_profiles"
	caMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/context_profiles/custom_attributes"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	ctxProfileID       = "ctx-profile-1"
	ctxProfileName     = "ctx-profile-fooname"
	ctxProfileRevision = int64(1)
	ctxProfilePath     = "/infra/context-profiles/ctx-profile-1"
	ctxProfileDomain   = "example.com"
)

func TestMockResourceNsxtPolicyContextProfileCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContextProfilesSDK := inframocks.NewMockContextProfilesClient(ctrl)
	mockAttributesSDK := cpAttrMocks.NewMockAttributesClient(ctrl)
	mockCustomAttrSDK := caMocks.NewMockDefaultClient(ctrl)

	ctxProfileWrapper := &cliinfra.PolicyContextProfileClientContext{
		Client:     mockContextProfilesSDK,
		ClientType: utl.Local,
	}
	attrWrapper := &contextprofiles.AttributeClientContext{
		Client:     mockAttributesSDK,
		ClientType: utl.Local,
	}
	customAttrWrapper := &customattributes.PolicyCustomAttributesClientContext{
		Client:     mockCustomAttrSDK,
		ClientType: utl.Local,
	}

	originalCli := cliContextProfilesClient
	originalAttrCli := cliContextProfileAttributesClient
	originalCustomAttrCli := cliContextProfileCustomAttributesClient
	defer func() {
		cliContextProfilesClient = originalCli
		cliContextProfileAttributesClient = originalAttrCli
		cliContextProfileCustomAttributesClient = originalCustomAttrCli
	}()
	cliContextProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.PolicyContextProfileClientContext {
		return ctxProfileWrapper
	}
	cliContextProfileAttributesClient = func(sessionContext utl.SessionContext, connector client.Connector) *contextprofiles.AttributeClientContext {
		return attrWrapper
	}
	cliContextProfileCustomAttributesClient = func(sessionContext utl.SessionContext, connector client.Connector) *customattributes.PolicyCustomAttributesClientContext {
		return customAttrWrapper
	}

	res := resourceNsxtPolicyContextProfile()

	t.Run("Create_success", func(t *testing.T) {
		domainKey := model.PolicyAttributes_KEY_DOMAIN_NAME
		resultCount := int64(1)
		mockAttributesSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{
				ResultCount: &resultCount,
				Results: []model.PolicyContextProfile{
					{
						Attributes: []model.PolicyAttributes{
							{Key: &domainKey, Value: []string{ctxProfileDomain}},
						},
					},
				},
			}, nil).AnyTimes()
		mockCustomAttrSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{}, nil).AnyTimes()

		mockContextProfilesSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockContextProfilesSDK.EXPECT().Get(gomock.Any()).Return(model.PolicyContextProfile{
			DisplayName: &ctxProfileName,
			Path:        &ctxProfilePath,
			Revision:    &ctxProfileRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": ctxProfileName,
			"domain_name": []interface{}{
				map[string]interface{}{
					"value":         []interface{}{ctxProfileDomain},
					"sub_attribute": []interface{}{},
				},
			},
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create_fails_when_no_attributes", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": ctxProfileName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "At least one attribute")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		domainKey := model.PolicyAttributes_KEY_DOMAIN_NAME
		resultCount := int64(1)
		mockAttributesSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{
				ResultCount: &resultCount,
				Results: []model.PolicyContextProfile{
					{
						Attributes: []model.PolicyAttributes{
							{Key: &domainKey, Value: []string{ctxProfileDomain}},
						},
					},
				},
			}, nil).AnyTimes()
		mockCustomAttrSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{}, nil).AnyTimes()

		mockContextProfilesSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": ctxProfileName,
			"domain_name": []interface{}{
				map[string]interface{}{
					"value":         []interface{}{ctxProfileDomain},
					"sub_attribute": []interface{}{},
				},
			},
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyContextProfileRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContextProfilesSDK := inframocks.NewMockContextProfilesClient(ctrl)
	ctxProfileWrapper := &cliinfra.PolicyContextProfileClientContext{
		Client:     mockContextProfilesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliContextProfilesClient
	defer func() { cliContextProfilesClient = originalCli }()
	cliContextProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.PolicyContextProfileClientContext {
		return ctxProfileWrapper
	}

	res := resourceNsxtPolicyContextProfile()

	t.Run("Read_success", func(t *testing.T) {
		mockContextProfilesSDK.EXPECT().Get(ctxProfileID).Return(model.PolicyContextProfile{
			DisplayName: &ctxProfileName,
			Path:        &ctxProfilePath,
			Revision:    &ctxProfileRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ctxProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, ctxProfileName, d.Get("display_name"))
	})

	t.Run("Read_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ContextProfile ID")
	})

	t.Run("Read_fails_when_Get_returns_not_found", func(t *testing.T) {
		mockContextProfilesSDK.EXPECT().Get(ctxProfileID).Return(model.PolicyContextProfile{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ctxProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyContextProfileUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContextProfilesSDK := inframocks.NewMockContextProfilesClient(ctrl)
	mockAttributesSDK := cpAttrMocks.NewMockAttributesClient(ctrl)
	mockCustomAttrSDK := caMocks.NewMockDefaultClient(ctrl)

	ctxProfileWrapper := &cliinfra.PolicyContextProfileClientContext{
		Client:     mockContextProfilesSDK,
		ClientType: utl.Local,
	}
	attrWrapper := &contextprofiles.AttributeClientContext{
		Client:     mockAttributesSDK,
		ClientType: utl.Local,
	}
	customAttrWrapper := &customattributes.PolicyCustomAttributesClientContext{
		Client:     mockCustomAttrSDK,
		ClientType: utl.Local,
	}

	originalCli := cliContextProfilesClient
	originalAttrCli := cliContextProfileAttributesClient
	originalCustomAttrCli := cliContextProfileCustomAttributesClient
	defer func() {
		cliContextProfilesClient = originalCli
		cliContextProfileAttributesClient = originalAttrCli
		cliContextProfileCustomAttributesClient = originalCustomAttrCli
	}()
	cliContextProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.PolicyContextProfileClientContext {
		return ctxProfileWrapper
	}
	cliContextProfileAttributesClient = func(sessionContext utl.SessionContext, connector client.Connector) *contextprofiles.AttributeClientContext {
		return attrWrapper
	}
	cliContextProfileCustomAttributesClient = func(sessionContext utl.SessionContext, connector client.Connector) *customattributes.PolicyCustomAttributesClientContext {
		return customAttrWrapper
	}

	res := resourceNsxtPolicyContextProfile()

	t.Run("Update_success", func(t *testing.T) {
		domainKey := model.PolicyAttributes_KEY_DOMAIN_NAME
		resultCount := int64(1)
		mockAttributesSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{
				ResultCount: &resultCount,
				Results: []model.PolicyContextProfile{
					{
						Attributes: []model.PolicyAttributes{
							{Key: &domainKey, Value: []string{ctxProfileDomain}},
						},
					},
				},
			}, nil).AnyTimes()
		mockCustomAttrSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.PolicyContextProfileListResult{}, nil).AnyTimes()

		mockContextProfilesSDK.EXPECT().Patch(ctxProfileID, gomock.Any(), gomock.Any()).Return(nil)
		mockContextProfilesSDK.EXPECT().Get(ctxProfileID).Return(model.PolicyContextProfile{
			DisplayName: &ctxProfileName,
			Path:        &ctxProfilePath,
			Revision:    &ctxProfileRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": ctxProfileName,
			"domain_name": []interface{}{
				map[string]interface{}{
					"value":         []interface{}{ctxProfileDomain},
					"sub_attribute": []interface{}{},
				},
			},
		})
		d.SetId(ctxProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyContextProfileDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContextProfilesSDK := inframocks.NewMockContextProfilesClient(ctrl)
	ctxProfileWrapper := &cliinfra.PolicyContextProfileClientContext{
		Client:     mockContextProfilesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliContextProfilesClient
	defer func() { cliContextProfilesClient = originalCli }()
	cliContextProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.PolicyContextProfileClientContext {
		return ctxProfileWrapper
	}

	res := resourceNsxtPolicyContextProfile()

	t.Run("Delete_success", func(t *testing.T) {
		force := true
		mockContextProfilesSDK.EXPECT().Delete(ctxProfileID, &force, gomock.Any()).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ctxProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileDelete(d, m)
		require.Error(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		force := true
		mockContextProfilesSDK.EXPECT().Delete(ctxProfileID, &force, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ctxProfileID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyContextProfileDelete(d, m)
		require.Error(t, err)
	})
}
