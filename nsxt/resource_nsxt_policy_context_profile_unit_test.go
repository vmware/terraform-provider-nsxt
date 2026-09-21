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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestUnitNsxt_validateSubAttributes(t *testing.T) {
	res := resourceNsxtPolicyContextProfile()

	t.Run("single value with sub-attributes is valid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"app_id": []interface{}{
				map[string]interface{}{
					"value": []interface{}{"HTTP"},
					"sub_attribute": []interface{}{
						map[string]interface{}{"tls_version": []interface{}{"TLS_V12"}},
					},
				},
			},
		})
		attrs := d.Get("app_id").(*schema.Set).List()
		require.NoError(t, validateSubAttributes(attrs))
	})

	t.Run("multiple values with sub-attributes is invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"app_id": []interface{}{
				map[string]interface{}{
					"value": []interface{}{"HTTP", "HTTPS"},
					"sub_attribute": []interface{}{
						map[string]interface{}{"tls_version": []interface{}{"TLS_V12"}},
					},
				},
			},
		})
		attrs := d.Get("app_id").(*schema.Set).List()
		err := validateSubAttributes(attrs)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Sub-attributes are only applicable")
	})

	t.Run("multiple values without sub-attributes is valid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"app_id": []interface{}{
				map[string]interface{}{
					"value": []interface{}{"HTTP", "HTTPS"},
				},
			},
		})
		attrs := d.Get("app_id").(*schema.Set).List()
		require.NoError(t, validateSubAttributes(attrs))
	})
}

func TestUnitNsxt_constructSubAttributeModelList(t *testing.T) {
	res := resourceNsxtPolicyContextProfile()

	t.Run("converts populated sub-attribute keys", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"app_id": []interface{}{
				map[string]interface{}{
					"value": []interface{}{"HTTP"},
					"sub_attribute": []interface{}{
						map[string]interface{}{
							"tls_version":      []interface{}{"TLS_V12"},
							"cifs_smb_version": []interface{}{"SMB_V1"},
						},
					},
				},
			},
		})
		attrs := d.Get("app_id").(*schema.Set).List()
		attrMap := attrs[0].(map[string]interface{})
		subAttrs := attrMap["sub_attribute"].(*schema.Set).List()

		result, err := constructSubAttributeModelList(subAttrs)
		require.NoError(t, err)
		require.Len(t, result, 2)

		keys := make(map[string][]string)
		for _, sa := range result {
			keys[*sa.Key] = sa.Value
		}
		assert.ElementsMatch(t, []string{"TLS_V12"}, keys[model.PolicySubAttributes_KEY_TLS_VERSION])
		assert.ElementsMatch(t, []string{"SMB_V1"}, keys[model.PolicySubAttributes_KEY_CIFS_SMB_VERSION])
	})

	t.Run("empty sub-attributes returns empty list", func(t *testing.T) {
		result, err := constructSubAttributeModelList(nil)
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestUnitNsxt_constructAttributesModelList(t *testing.T) {
	res := resourceNsxtPolicyContextProfile()

	t.Run("app_id attribute includes sub-attributes", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"app_id": []interface{}{
				map[string]interface{}{
					"description": "my app",
					"value":       []interface{}{"HTTP"},
					"sub_attribute": []interface{}{
						map[string]interface{}{"tls_version": []interface{}{"TLS_V12"}},
					},
				},
			},
		})
		attrs := d.Get("app_id").(*schema.Set).List()

		result, err := constructAttributesModelList(attrs, "app_id")
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, model.PolicyAttributes_KEY_APP_ID, *result[0].Key)
		assert.ElementsMatch(t, []string{"HTTP"}, result[0].Value)
		require.Len(t, result[0].SubAttributes, 1)
		assert.Nil(t, result[0].CustomUrlPartialMatch)
	})

	t.Run("custom_url attribute sets partial match flag", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"custom_url": []interface{}{
				map[string]interface{}{
					"description":              "my url",
					"value":                    []interface{}{"example.com"},
					"custom_url_partial_match": true,
				},
			},
		})
		attrs := d.Get("custom_url").(*schema.Set).List()

		result, err := constructAttributesModelList(attrs, "custom_url")
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, model.PolicyAttributes_KEY_CUSTOM_URL, *result[0].Key)
		require.NotNil(t, result[0].CustomUrlPartialMatch)
		assert.True(t, *result[0].CustomUrlPartialMatch)
		assert.Empty(t, result[0].SubAttributes)
	})
}

func TestUnitNsxt_fillAttributesInSchema(t *testing.T) {
	res := resourceNsxtPolicyContextProfile()

	t.Run("app_id attribute with sub-attributes populates schema", func(t *testing.T) {
		d := res.TestResourceData()
		appIDKey := model.PolicyAttributes_KEY_APP_ID
		description := "my app"
		tlsVersionKey := model.PolicySubAttributes_KEY_TLS_VERSION
		isAlgType := false

		fillAttributesInSchema(d, []model.PolicyAttributes{
			{
				Key:         &appIDKey,
				Description: &description,
				Value:       []string{"HTTP"},
				IsALGType:   &isAlgType,
				SubAttributes: []model.PolicySubAttributes{
					{Key: &tlsVersionKey, Value: []string{"TLS_V12"}},
				},
			},
		})

		appIDs := d.Get("app_id").(*schema.Set).List()
		require.Len(t, appIDs, 1)
		elem := appIDs[0].(map[string]interface{})
		assert.Equal(t, description, elem["description"])
		subAttrs := elem["sub_attribute"].(*schema.Set).List()
		require.Len(t, subAttrs, 1)
		subAttrMap := subAttrs[0].(map[string]interface{})
		assert.ElementsMatch(t, []string{"TLS_V12"}, subAttrMap["tls_version"].(*schema.Set).List())
	})

	t.Run("custom_url attribute populates partial match flag", func(t *testing.T) {
		d := res.TestResourceData()
		customURLKey := model.PolicyAttributes_KEY_CUSTOM_URL
		description := "my url"
		partialMatch := true

		fillAttributesInSchema(d, []model.PolicyAttributes{
			{
				Key:                   &customURLKey,
				Description:           &description,
				Value:                 []string{"example.com"},
				CustomUrlPartialMatch: &partialMatch,
			},
		})

		urls := d.Get("custom_url").(*schema.Set).List()
		require.Len(t, urls, 1)
		elem := urls[0].(map[string]interface{})
		assert.Equal(t, true, elem["custom_url_partial_match"])
	})
}

func TestUnitNsxt_fillSubAttributesInSchema(t *testing.T) {
	t.Run("converts sub-attributes into a single schema element", func(t *testing.T) {
		tlsVersionKey := model.PolicySubAttributes_KEY_TLS_VERSION
		cifsKey := model.PolicySubAttributes_KEY_CIFS_SMB_VERSION

		result := fillSubAttributesInSchema([]model.PolicySubAttributes{
			{Key: &tlsVersionKey, Value: []string{"TLS_V12"}},
			{Key: &cifsKey, Value: []string{"SMB_V1"}},
		})

		require.Len(t, result, 1)
		elem := result[0].(map[string]interface{})
		assert.ElementsMatch(t, []string{"TLS_V12"}, elem["tls_version"])
		assert.ElementsMatch(t, []string{"SMB_V1"}, elem["cifs_smb_version"])
	})
}
