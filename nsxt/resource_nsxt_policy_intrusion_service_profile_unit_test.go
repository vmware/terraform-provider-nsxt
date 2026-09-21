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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestUnitNsxt_buildIdsProfileCriteriaFilter(t *testing.T) {
	sv, err := buildIdsProfileCriteriaFilter(model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, []string{"a", "b"})
	require.NoError(t, err)
	require.NotNil(t, sv)

	converter := bindings.NewTypeConverter()
	golang, errs := converter.ConvertToGolang(sv, model.IdsProfileFilterCriteriaBindingType())
	require.Empty(t, errs)
	criteria := golang.(model.IdsProfileFilterCriteria)
	assert.Equal(t, model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, *criteria.FilterName)
	assert.Equal(t, []string{"a", "b"}, criteria.FilterValue)
}

func TestUnitNsxt_buildIdsProfileCriteriaOperator(t *testing.T) {
	sv, err := buildIdsProfileCriteriaOperator()
	require.NoError(t, err)
	require.NotNil(t, sv)

	converter := bindings.NewTypeConverter()
	golang, errs := converter.ConvertToGolang(sv, model.IdsProfileConjunctionOperatorBindingType())
	require.Empty(t, errs)
	operator := golang.(model.IdsProfileConjunctionOperator)
	assert.Equal(t, "AND", *operator.Operator)
}

func idsProfileTestResource() *schema.Resource {
	return resourceNsxtPolicyIntrusionServiceProfile()
}

func TestUnitNsxt_getIdsProfileCriteriaFromSchema(t *testing.T) {
	res := idsProfileTestResource()
	converter := bindings.NewTypeConverter()

	t.Run("no criteria returns empty list", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		result, err := getIdsProfileCriteriaFromSchema(d)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("single criteria field produces one filter, no trailing operator", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"criteria": []interface{}{
				map[string]interface{}{
					"attack_types": []interface{}{"malware"},
				},
			},
		})
		result, err := getIdsProfileCriteriaFromSchema(d)
		require.NoError(t, err)
		require.Len(t, result, 1)

		golang, errs := converter.ConvertToGolang(result[0], model.IdsProfileFilterCriteriaBindingType())
		require.Empty(t, errs)
		criteria := golang.(model.IdsProfileFilterCriteria)
		assert.Equal(t, model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, *criteria.FilterName)
	})

	t.Run("multiple criteria fields are joined with AND operators", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"criteria": []interface{}{
				map[string]interface{}{
					"attack_types":   []interface{}{"malware"},
					"attack_targets": []interface{}{"server"},
				},
			},
		})
		result, err := getIdsProfileCriteriaFromSchema(d)
		require.NoError(t, err)
		// filter, AND, filter -> the trailing operator after the last filter is stripped
		require.Len(t, result, 3)

		golang0, errs := converter.ConvertToGolang(result[0], model.IdsProfileFilterCriteriaBindingType())
		require.Empty(t, errs)
		assert.Equal(t, model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, *golang0.(model.IdsProfileFilterCriteria).FilterName)

		golangOp, errs := converter.ConvertToGolang(result[1], model.IdsProfileConjunctionOperatorBindingType())
		require.Empty(t, errs)
		assert.Equal(t, "AND", *golangOp.(model.IdsProfileConjunctionOperator).Operator)

		golang2, errs := converter.ConvertToGolang(result[2], model.IdsProfileFilterCriteriaBindingType())
		require.Empty(t, errs)
		assert.Equal(t, model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TARGET, *golang2.(model.IdsProfileFilterCriteria).FilterName)
	})
}

func TestUnitNsxt_setIdsProfileCriteriaInSchema(t *testing.T) {
	res := idsProfileTestResource()

	t.Run("filter criteria without operators populate the schema", func(t *testing.T) {
		d := res.TestResourceData()
		filter, err := buildIdsProfileCriteriaFilter(model.IdsProfileFilterCriteria_FILTER_NAME_CVSS, []string{"HIGH"})
		require.NoError(t, err)

		require.NoError(t, setIdsProfileCriteriaInSchema([]*data.StructValue{filter}, d))

		criteria := d.Get("criteria").([]interface{})
		require.Len(t, criteria, 1)
		elem := criteria[0].(map[string]interface{})
		assert.ElementsMatch(t, []string{"HIGH"}, elem["cvss"].(*schema.Set).List())
	})

	t.Run("AND operators at odd indices are skipped", func(t *testing.T) {
		d := res.TestResourceData()
		filter1, err := buildIdsProfileCriteriaFilter(model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, []string{"malware"})
		require.NoError(t, err)
		op, err := buildIdsProfileCriteriaOperator()
		require.NoError(t, err)
		filter2, err := buildIdsProfileCriteriaFilter(model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TARGET, []string{"server"})
		require.NoError(t, err)

		require.NoError(t, setIdsProfileCriteriaInSchema([]*data.StructValue{filter1, op, filter2}, d))

		criteria := d.Get("criteria").([]interface{})
		require.Len(t, criteria, 1)
		elem := criteria[0].(map[string]interface{})
		assert.ElementsMatch(t, []string{"malware"}, elem["attack_types"].(*schema.Set).List())
		assert.ElementsMatch(t, []string{"server"}, elem["attack_targets"].(*schema.Set).List())
	})

	t.Run("empty list clears the schema", func(t *testing.T) {
		d := res.TestResourceData()
		require.NoError(t, setIdsProfileCriteriaInSchema(nil, d))
		assert.Empty(t, d.Get("criteria").([]interface{}))
	})
}

func TestUnitNsxt_getIdsProfileSignaturesFromSchema(t *testing.T) {
	res := idsProfileTestResource()

	t.Run("no signatures returns empty list", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		result := getIdsProfileSignaturesFromSchema(d)
		assert.Empty(t, result)
	})

	t.Run("converts overridden signatures", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"overridden_signature": []interface{}{
				map[string]interface{}{
					"signature_id": "sig-1",
					"enabled":      true,
					"action":       model.IdsProfileLocalSignature_ACTION_DROP,
				},
			},
		})
		result := getIdsProfileSignaturesFromSchema(d)
		require.Len(t, result, 1)
		assert.Equal(t, "sig-1", *result[0].SignatureId)
		assert.True(t, *result[0].Enable)
		assert.Equal(t, model.IdsProfileLocalSignature_ACTION_DROP, *result[0].Action)
	})
}

func TestUnitNsxt_setIdsProfileSignaturesInSchema(t *testing.T) {
	res := idsProfileTestResource()

	t.Run("populates schema from signature models", func(t *testing.T) {
		d := res.TestResourceData()
		require.NoError(t, setIdsProfileSignaturesInSchema([]model.IdsProfileLocalSignature{
			{SignatureId: strPtr("sig-1"), Enable: boolPtr(true), Action: strPtr(model.IdsProfileLocalSignature_ACTION_REJECT)},
		}, d))

		signatures := d.Get("overridden_signature").(*schema.Set).List()
		require.Len(t, signatures, 1)
		elem := signatures[0].(map[string]interface{})
		assert.Equal(t, "sig-1", elem["signature_id"])
		assert.Equal(t, true, elem["enabled"])
		assert.Equal(t, model.IdsProfileLocalSignature_ACTION_REJECT, elem["action"])
	})

	t.Run("empty list clears the schema", func(t *testing.T) {
		d := res.TestResourceData()
		require.NoError(t, setIdsProfileSignaturesInSchema(nil, d))
		assert.Empty(t, d.Get("overridden_signature").(*schema.Set).List())
	})
}
