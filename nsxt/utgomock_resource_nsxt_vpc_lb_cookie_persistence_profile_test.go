//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	apivpcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
)

func setupVpcLbPersistenceMock(t *testing.T, ctrl *gomock.Controller) (*vpcsmocks.MockVpcLbPersistenceProfilesClient, func()) {
	mockSDK := vpcsmocks.NewMockVpcLbPersistenceProfilesClient(ctrl)
	mockWrapper := &apivpcs.VpcLbPersistenceProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  vpcLbPrfProject,
		VPCID:      vpcLbPrfVpc,
	}
	original := cliVpcLbPersistenceProfilesClient
	cliVpcLbPersistenceProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbPersistenceProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcLbPersistenceProfilesClient = original }
}

// vpcLbPersistCase describes one VPC persistence profile resource for the shared test suite.
type vpcLbPersistCase struct {
	name        string
	newRes      func() *schema.Resource
	id          string
	displayName string
	extraData   map[string]interface{}
	// newModel returns the typed API model and its binding type for the given display name
	newModel func(displayName string) (interface{}, bindings.BindingType)
	// checkRead is invoked on the resource data after a successful read of newModel("<displayName>")
	checkRead func(t *testing.T, d *schema.ResourceData)
	// cacheAware indicates that the shared read code has a cache-aware path
	cacheAware   bool
	cacheTypeKey string
}

func (c vpcLbPersistCase) data() map[string]interface{} {
	res := map[string]interface{}{
		"display_name": c.displayName,
		"nsx_id":       c.id,
		"context":      vpcLbPrfContext(),
	}
	for k, v := range c.extraData {
		res[k] = v
	}
	return res
}

func (c vpcLbPersistCase) structValue(t *testing.T) *data.StructValue {
	t.Helper()
	obj, bt := c.newModel(c.displayName)
	val, errs := bindings.NewTypeConverter().ConvertToVapi(obj, bt)
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func runVpcLbPersistenceResourceSuite(t *testing.T, c vpcLbPersistCase) {
	t.Run("Create", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			vpcLbPrfSetNsxVersion(t, "9.0.2")
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()

			gomock.InOrder(
				mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id).Return(nil, vapiErrors.NotFound{}),
				mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id, gomock.Any()).Return(nil),
				mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id).Return(c.structValue(t), nil),
			)

			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			require.NoError(t, res.Create(d, newGoMockProviderClient()))
			assert.Equal(t, c.id, d.Id())
			assert.Equal(t, c.displayName, d.Get("display_name"))
			if c.checkRead != nil {
				c.checkRead(t, d)
			}
		})

		t.Run("already exists", func(t *testing.T) {
			vpcLbPrfSetNsxVersion(t, "9.0.2")
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()

			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id).Return(c.structValue(t), nil)

			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			err := res.Create(d, newGoMockProviderClient())
			require.Error(t, err)
			assert.Contains(t, err.Error(), "already exists")
		})

		t.Run("NSX version lower than 9.0.2", func(t *testing.T) {
			vpcLbPrfSetNsxVersion(t, "9.0.1")
			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			err := res.Create(d, newGoMockProviderClient())
			require.Error(t, err)
			assert.Contains(t, err.Error(), "9.0.2")
		})

		t.Run("patch error", func(t *testing.T) {
			vpcLbPrfSetNsxVersion(t, "9.0.2")
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()

			gomock.InOrder(
				mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id).Return(nil, vapiErrors.NotFound{}),
				mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id, gomock.Any()).Return(vapiErrors.InternalServerError{}),
			)

			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			require.Error(t, res.Create(d, newGoMockProviderClient()))
			assert.Equal(t, "", d.Id())
		})
	})

	t.Run("Read", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()

			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id).Return(c.structValue(t), nil)

			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			d.SetId(c.id)
			require.NoError(t, res.Read(d, newGoMockProviderClient()))
			assert.Equal(t, c.displayName, d.Get("display_name"))
			if c.checkRead != nil {
				c.checkRead(t, d)
			}
		})

		t.Run("not found clears ID", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()

			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id).Return(nil, vapiErrors.NotFound{})

			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			d.SetId(c.id)
			require.NoError(t, res.Read(d, newGoMockProviderClient()))
			assert.Equal(t, "", d.Id())
		})

		t.Run("API error is propagated", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()

			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id).Return(nil, vapiErrors.InternalServerError{})

			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			d.SetId(c.id)
			require.Error(t, res.Read(d, newGoMockProviderClient()))
		})

		t.Run("empty ID", func(t *testing.T) {
			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			require.Error(t, res.Read(d, newGoMockProviderClient()))
		})

		if c.cacheAware {
			t.Run("cache is bypassed for VPC context", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
				defer restore()

				mc := newGoMockProviderClientCacheEnabled()
				res := c.newRes()
				d := schema.TestResourceDataRaw(t, res.Schema, c.data())
				d.SetId(c.id)
				require.True(t, isCacheEnabledForRead(d, mc))
				require.False(t, isLbCacheEnabledForRead(d, mc))

				// an empty cache entry would result in a tag patch if the cache path was used
				query := getCacheQueryKey(c.cacheTypeKey, d, mc)
				tc := gcache.getTypeCache(c.cacheTypeKey)
				tc.data[query] = map[string]*data.StructValue{}
				defer delete(gcache.byTyp, c.cacheTypeKey)

				// exactly one Get, and no Patch
				mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id).Return(c.structValue(t), nil).Times(1)
				require.NoError(t, res.Read(d, mc))
				assert.Equal(t, c.displayName, d.Get("display_name"))
			})
		}
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()

			gomock.InOrder(
				mockSDK.EXPECT().Update(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id, gomock.Any()).Return(nil, nil),
				mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id).Return(c.structValue(t), nil),
			)

			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			d.SetId(c.id)
			require.NoError(t, res.Update(d, newGoMockProviderClient()))
		})

		t.Run("API error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()

			mockSDK.EXPECT().Update(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id, gomock.Any()).Return(nil, vapiErrors.InternalServerError{})

			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			d.SetId(c.id)
			require.Error(t, res.Update(d, newGoMockProviderClient()))
		})

		t.Run("empty ID", func(t *testing.T) {
			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			require.Error(t, res.Update(d, newGoMockProviderClient()))
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()

			mockSDK.EXPECT().Delete(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id, gomock.Any()).Return(nil)

			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			d.SetId(c.id)
			require.NoError(t, res.Delete(d, newGoMockProviderClient()))
		})

		t.Run("API error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbPersistenceMock(t, ctrl)
			defer restore()

			mockSDK.EXPECT().Delete(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, c.id, gomock.Any()).Return(vapiErrors.InternalServerError{})

			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			d.SetId(c.id)
			require.Error(t, res.Delete(d, newGoMockProviderClient()))
		})

		t.Run("empty ID", func(t *testing.T) {
			res := c.newRes()
			d := schema.TestResourceDataRaw(t, res.Schema, c.data())
			require.Error(t, res.Delete(d, newGoMockProviderClient()))
		})
	})

	t.Run("Schema and importer", func(t *testing.T) {
		res := c.newRes()
		ctxSchema, ok := res.Schema["context"]
		require.True(t, ok)
		assert.True(t, ctxSchema.Required)
		require.NotNil(t, res.Importer)

		runVpcLbImporterSuite(t, res, "vpc-lb-persistence-profiles")
	})
}

func TestMockResourceNsxtVpcLBCookiePersistenceProfile(t *testing.T) {
	runVpcLbPersistenceResourceSuite(t, vpcLbPersistCase{
		name:        "cookie",
		newRes:      resourceNsxtVpcLBCookiePersistenceProfile,
		id:          "vpc-lb-cookie-1",
		displayName: "Test VPC LB Cookie Persistence Profile",
		extraData: map[string]interface{}{
			"cookie_name": "NSXLB",
		},
		newModel: func(displayName string) (interface{}, bindings.BindingType) {
			return model.LBCookiePersistenceProfile{
				DisplayName:    &displayName,
				ResourceType:   model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE,
				CookieName:     strPtr("NSXLB"),
				CookieMode:     strPtr(model.LBCookiePersistenceProfile_COOKIE_MODE_INSERT),
				CookieFallback: boolPtr(true),
				CookieGarble:   boolPtr(true),
			}, model.LBCookiePersistenceProfileBindingType()
		},
		checkRead: func(t *testing.T, d *schema.ResourceData) {
			assert.Equal(t, "NSXLB", d.Get("cookie_name"))
		},
	})
}

// TestMockVpcLbPersistenceSslAPISelectors covers the VPC/infra client selection (including the
// unsupported-context nil result) for the persistence, client SSL and server SSL profile APIs.
func TestMockVpcLbPersistenceSslAPISelectors(t *testing.T) {
	vpcCtx := utl.SessionContext{ClientType: utl.VPC, ProjectID: vpcLbPrfProject, VPCID: vpcLbPrfVpc}
	localCtx := utl.SessionContext{ClientType: utl.Local}

	t.Run("persistence profile", func(t *testing.T) {
		assert.IsType(t, vpcLbPersistenceProfileAdapter{}, getLbPersistenceProfileAPI(vpcCtx, nil))
		assert.NotNil(t, getLbPersistenceProfileAPI(localCtx, nil))
		assert.NotEqual(t, reflect.TypeOf(vpcLbPersistenceProfileAdapter{}), reflect.TypeOf(getLbPersistenceProfileAPI(localCtx, nil)))

		origVpc, origInfra := cliVpcLbPersistenceProfilesClient, cliLbPersistenceProfilesClient
		defer func() { cliVpcLbPersistenceProfilesClient, cliLbPersistenceProfilesClient = origVpc, origInfra }()
		cliVpcLbPersistenceProfilesClient = func(utl.SessionContext, vapiProtocolClient.Connector) *apivpcs.VpcLbPersistenceProfileClientContext {
			return nil
		}
		cliLbPersistenceProfilesClient = func(utl.SessionContext, vapiProtocolClient.Connector) *infraapi.LBPersistenceProfileClientContext {
			return nil
		}
		assert.Nil(t, getLbPersistenceProfileAPI(vpcCtx, nil))
		assert.Nil(t, getLbPersistenceProfileAPI(localCtx, nil))
	})

	t.Run("client ssl profile", func(t *testing.T) {
		assert.IsType(t, vpcLbClientSslProfileAdapter{}, getLbClientSslProfileAPI(vpcCtx, nil))
		assert.NotNil(t, getLbClientSslProfileAPI(localCtx, nil))
		assert.NotEqual(t, reflect.TypeOf(vpcLbClientSslProfileAdapter{}), reflect.TypeOf(getLbClientSslProfileAPI(localCtx, nil)))

		origVpc, origInfra := cliVpcLbClientSslProfilesClient, cliLbClientSslProfilesClient
		defer func() { cliVpcLbClientSslProfilesClient, cliLbClientSslProfilesClient = origVpc, origInfra }()
		cliVpcLbClientSslProfilesClient = func(utl.SessionContext, vapiProtocolClient.Connector) *apivpcs.VpcLbClientSslProfileClientContext {
			return nil
		}
		cliLbClientSslProfilesClient = func(utl.SessionContext, vapiProtocolClient.Connector) *infraapi.LBClientSslProfileClientContext {
			return nil
		}
		assert.Nil(t, getLbClientSslProfileAPI(vpcCtx, nil))
		assert.Nil(t, getLbClientSslProfileAPI(localCtx, nil))
	})

	t.Run("server ssl profile", func(t *testing.T) {
		assert.IsType(t, vpcLbServerSslProfileAdapter{}, getLbServerSslProfileAPI(vpcCtx, nil))
		assert.NotNil(t, getLbServerSslProfileAPI(localCtx, nil))
		assert.NotEqual(t, reflect.TypeOf(vpcLbServerSslProfileAdapter{}), reflect.TypeOf(getLbServerSslProfileAPI(localCtx, nil)))

		origVpc, origInfra := cliVpcLbServerSslProfilesClient, cliLbServerSslProfilesClient
		defer func() { cliVpcLbServerSslProfilesClient, cliLbServerSslProfilesClient = origVpc, origInfra }()
		cliVpcLbServerSslProfilesClient = func(utl.SessionContext, vapiProtocolClient.Connector) *apivpcs.VpcLbServerSslProfileClientContext {
			return nil
		}
		cliLbServerSslProfilesClient = func(utl.SessionContext, vapiProtocolClient.Connector) *infraapi.LBServerSslProfileClientContext {
			return nil
		}
		assert.Nil(t, getLbServerSslProfileAPI(vpcCtx, nil))
		assert.Nil(t, getLbServerSslProfileAPI(localCtx, nil))
	})
}
