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

func stringSet(values ...string) *schema.Set {
	items := make([]interface{}, len(values))
	for i, v := range values {
		items[i] = v
	}
	return schema.NewSet(schema.HashString, items)
}

func TestUnitNsxt_validateGroupCriteriaSets(t *testing.T) {
	t.Run("empty block errors", func(t *testing.T) {
		_, err := validateGroupCriteriaSets([]interface{}{nil})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "empty criteria block")
	})

	t.Run("mixed expression types in a single block errors", func(t *testing.T) {
		block := map[string]interface{}{
			"condition":            []interface{}{map[string]interface{}{"key": "Tag"}},
			"ipaddress_expression": []interface{}{map[string]interface{}{"ip_addresses": stringSet("1.1.1.1")}},
		}
		_, err := validateGroupCriteriaSets([]interface{}{block})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "should be homogeneous")
	})

	t.Run("single expression is not nested, multiple are", func(t *testing.T) {
		singleBlock := map[string]interface{}{
			"condition": []interface{}{map[string]interface{}{"key": "Tag"}},
		}
		nestedBlock := map[string]interface{}{
			"condition": []interface{}{
				map[string]interface{}{"key": "Tag"},
				map[string]interface{}{"key": "Name"},
			},
		}
		meta, err := validateGroupCriteriaSets([]interface{}{singleBlock, nestedBlock})
		require.NoError(t, err)
		require.Len(t, meta, 2)
		assert.False(t, meta[0].IsNested)
		assert.True(t, meta[1].IsNested)
		assert.Equal(t, "condition", meta[0].ExpressionType)
	})
}

func TestUnitNsxt_validateGroupConjunctions(t *testing.T) {
	criteria := []criteriaMeta{
		{ExpressionType: "condition"},
		{ExpressionType: "ipaddress_expression"},
	}

	t.Run("AND requires matching expression types", func(t *testing.T) {
		conjunctions := []interface{}{
			map[string]interface{}{"operator": model.ConjunctionOperator_CONJUNCTION_OPERATOR_AND},
		}
		err := validateGroupConjunctions(conjunctions, criteria)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must use the same types")
	})

	t.Run("OR allows differing expression types", func(t *testing.T) {
		conjunctions := []interface{}{
			map[string]interface{}{"operator": model.ConjunctionOperator_CONJUNCTION_OPERATOR_OR},
		}
		err := validateGroupConjunctions(conjunctions, criteria)
		require.NoError(t, err)
	})

	t.Run("AND with matching expression types succeeds", func(t *testing.T) {
		sameCriteria := []criteriaMeta{
			{ExpressionType: "condition"},
			{ExpressionType: "condition"},
		}
		conjunctions := []interface{}{
			map[string]interface{}{"operator": model.ConjunctionOperator_CONJUNCTION_OPERATOR_AND},
		}
		err := validateGroupConjunctions(conjunctions, sameCriteria)
		require.NoError(t, err)
	})
}

func TestUnitNsxt_validateGroupCriteriaAndConjunctions(t *testing.T) {
	t.Run("no criteria and no conjunctions is valid", func(t *testing.T) {
		meta, err := validateGroupCriteriaAndConjunctions(nil, nil)
		require.NoError(t, err)
		assert.Nil(t, meta)
	})

	t.Run("even total with too few conjunctions reports missing conjunction", func(t *testing.T) {
		criteriaSets := []interface{}{
			map[string]interface{}{"condition": []interface{}{map[string]interface{}{"key": "Tag"}}},
			map[string]interface{}{"condition": []interface{}{map[string]interface{}{"key": "Name"}}},
		}
		_, err := validateGroupCriteriaAndConjunctions(criteriaSets, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Missing conjunction")
	})

	t.Run("even total with enough conjunctions reports missing trailing criteria", func(t *testing.T) {
		criteriaSets := []interface{}{
			map[string]interface{}{"condition": []interface{}{map[string]interface{}{"key": "Tag"}}},
		}
		conjunctions := []interface{}{
			map[string]interface{}{"operator": model.ConjunctionOperator_CONJUNCTION_OPERATOR_OR},
		}
		_, err := validateGroupCriteriaAndConjunctions(criteriaSets, conjunctions)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Missing criteria")
	})

	t.Run("odd total with valid conjunctions succeeds", func(t *testing.T) {
		criteriaSets := []interface{}{
			map[string]interface{}{"condition": []interface{}{map[string]interface{}{"key": "Tag"}}},
			map[string]interface{}{"condition": []interface{}{map[string]interface{}{"key": "Name"}}},
		}
		conjunctions := []interface{}{
			map[string]interface{}{"operator": model.ConjunctionOperator_CONJUNCTION_OPERATOR_OR},
		}
		meta, err := validateGroupCriteriaAndConjunctions(criteriaSets, conjunctions)
		require.NoError(t, err)
		require.Len(t, meta, 2)
	})
}

func TestUnitNsxt_buildGroupConditionData(t *testing.T) {
	condition := map[string]interface{}{
		"key":         model.Condition_KEY_TAG,
		"member_type": model.Condition_MEMBER_TYPE_VIRTUALMACHINE,
		"operator":    model.Condition_OPERATOR_EQUALS,
		"value":       "prod",
	}
	sv, err := buildGroupConditionData(condition)
	require.NoError(t, err)

	condMap, err := groupConditionDataToMap(sv)
	require.NoError(t, err)
	assert.Equal(t, model.Condition_KEY_TAG, *(condMap["key"].(*string)))
	assert.Equal(t, "prod", *(condMap["value"].(*string)))
}

func TestUnitNsxt_buildGroupConjunctionData(t *testing.T) {
	sv, err := buildGroupConjunctionData(model.ConjunctionOperator_CONJUNCTION_OPERATOR_AND)
	require.NoError(t, err)
	require.NotNil(t, sv)

	converter := bindings.NewTypeConverter()
	golang, errs := converter.ConvertToGolang(sv, model.ConjunctionOperatorBindingType())
	require.Empty(t, errs)
	conj := golang.(model.ConjunctionOperator)
	assert.Equal(t, model.ConjunctionOperator_CONJUNCTION_OPERATOR_AND, *conj.ConjunctionOperator)
}

func TestUnitNsxt_buildGroupIPAddressData(t *testing.T) {
	ipaddr := map[string]interface{}{"ip_addresses": stringSet("1.1.1.1", "2.2.2.2")}
	sv, err := buildGroupIPAddressData(ipaddr)
	require.NoError(t, err)

	criteria, _, err := fromGroupExpressionData([]*data.StructValue{sv})
	require.NoError(t, err)
	require.Len(t, criteria, 1)
	ipExpr := criteria[0]["ipaddress_expression"].([]map[string]interface{})
	assert.ElementsMatch(t, []string{"1.1.1.1", "2.2.2.2"}, ipExpr[0]["ip_addresses"])
}

func TestUnitNsxt_buildGroupMacAddressData(t *testing.T) {
	addr := map[string]interface{}{"mac_addresses": stringSet("aa:bb:cc:dd:ee:ff")}
	sv, err := buildGroupMacAddressData(addr)
	require.NoError(t, err)

	criteria, _, err := fromGroupExpressionData([]*data.StructValue{sv})
	require.NoError(t, err)
	require.Len(t, criteria, 1)
	macExpr := criteria[0]["macaddress_expression"].([]map[string]interface{})
	assert.Equal(t, []string{"aa:bb:cc:dd:ee:ff"}, macExpr[0]["mac_addresses"])
}

func TestUnitNsxt_buildGroupExternalIDExpressionData(t *testing.T) {
	idMap := map[string]interface{}{
		"member_type":  model.ExternalIDExpression_MEMBER_TYPE_VIRTUALMACHINE,
		"external_ids": stringSet("id-1", "id-2"),
	}
	sv, err := buildGroupExternalIDExpressionData(idMap)
	require.NoError(t, err)

	criteria, _, err := fromGroupExpressionData([]*data.StructValue{sv})
	require.NoError(t, err)
	require.Len(t, criteria, 1)
	extExpr := criteria[0]["external_id_expression"].([]map[string]interface{})
	assert.Equal(t, model.ExternalIDExpression_MEMBER_TYPE_VIRTUALMACHINE, *(extExpr[0]["member_type"].(*string)))
	assert.ElementsMatch(t, []string{"id-1", "id-2"}, extExpr[0]["external_ids"])
}

func TestUnitNsxt_buildGroupMemberPathData(t *testing.T) {
	pathMap := map[string]interface{}{"member_paths": stringSet("/infra/domains/default/groups/g1")}
	sv, err := buildGroupMemberPathData(pathMap)
	require.NoError(t, err)

	criteria, _, err := fromGroupExpressionData([]*data.StructValue{sv})
	require.NoError(t, err)
	require.Len(t, criteria, 1)
	pathExpr := criteria[0]["path_expression"].([]map[string]interface{})
	assert.Equal(t, []string{"/infra/domains/default/groups/g1"}, pathExpr[0]["member_paths"])
}

func TestUnitNsxt_buildGroupExpressionDataFromType(t *testing.T) {
	t.Run("nil datum errors", func(t *testing.T) {
		_, err := buildGroupExpressionDataFromType("condition", nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Empty set is not supported")
	})

	t.Run("unknown expression type errors", func(t *testing.T) {
		_, err := buildGroupExpressionDataFromType("bogus_expression", map[string]interface{}{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unknown expression type")
	})

	t.Run("dispatches to the matching builder for each known type", func(t *testing.T) {
		cases := []struct {
			expressionType string
			datum          interface{}
		}{
			{"condition", map[string]interface{}{
				"key": model.Condition_KEY_TAG, "member_type": model.Condition_MEMBER_TYPE_VIRTUALMACHINE,
				"operator": model.Condition_OPERATOR_EQUALS, "value": "v",
			}},
			{"ipaddress_expression", map[string]interface{}{"ip_addresses": stringSet("1.1.1.1")}},
			{"path_expression", map[string]interface{}{"member_paths": stringSet("/infra/x")}},
			{"macaddress_expression", map[string]interface{}{"mac_addresses": stringSet("aa:bb:cc:dd:ee:ff")}},
			{"external_id_expression", map[string]interface{}{
				"member_type": model.ExternalIDExpression_MEMBER_TYPE_VIRTUALMACHINE, "external_ids": stringSet("id-1"),
			}},
		}
		for _, c := range cases {
			sv, err := buildGroupExpressionDataFromType(c.expressionType, c.datum)
			require.NoError(t, err, c.expressionType)
			assert.NotNil(t, sv, c.expressionType)
		}
	})
}

func TestUnitNsxt_buildNestedGroupExpressionData(t *testing.T) {
	cond1, err := buildGroupConditionData(map[string]interface{}{
		"key": model.Condition_KEY_TAG, "member_type": model.Condition_MEMBER_TYPE_VIRTUALMACHINE,
		"operator": model.Condition_OPERATOR_EQUALS, "value": "v1",
	})
	require.NoError(t, err)
	cond2, err := buildGroupConditionData(map[string]interface{}{
		"key": model.Condition_KEY_NAME, "member_type": model.Condition_MEMBER_TYPE_VIRTUALMACHINE,
		"operator": model.Condition_OPERATOR_EQUALS, "value": "v2",
	})
	require.NoError(t, err)

	sv, err := buildNestedGroupExpressionData([]*data.StructValue{cond1, cond2})
	require.NoError(t, err)

	converter := bindings.NewTypeConverter()
	golang, errs := converter.ConvertToGolang(sv, model.NestedExpressionBindingType())
	require.Empty(t, errs)
	nested := golang.(model.NestedExpression)
	// 2 conditions + 1 implicit AND conjunction between them
	assert.Len(t, nested.Expressions, 3)
}

func TestUnitNsxt_buildGroupExpressionData(t *testing.T) {
	criteria := []criteriaMeta{
		{
			ExpressionType: "condition",
			IsNested:       false,
			criteriaBlocks: []interface{}{map[string]interface{}{
				"key": model.Condition_KEY_TAG, "member_type": model.Condition_MEMBER_TYPE_VIRTUALMACHINE,
				"operator": model.Condition_OPERATOR_EQUALS, "value": "v1",
			}},
		},
		{
			ExpressionType: "ipaddress_expression",
			IsNested:       false,
			criteriaBlocks: []interface{}{map[string]interface{}{"ip_addresses": stringSet("1.1.1.1")}},
		},
	}
	conjunctions := []interface{}{
		map[string]interface{}{"operator": model.ConjunctionOperator_CONJUNCTION_OPERATOR_OR},
	}

	expressionData, err := buildGroupExpressionData(criteria, conjunctions)
	require.NoError(t, err)
	// condition, conjunction, ipaddress_expression
	require.Len(t, expressionData, 3)

	parsedCriteria, parsedConjunctions, err := fromGroupExpressionData(expressionData)
	require.NoError(t, err)
	assert.Len(t, parsedCriteria, 2)
	assert.Len(t, parsedConjunctions, 1)
}

func TestUnitNsxt_fromGroupExpressionData_unsupportedType(t *testing.T) {
	tag := model.Tag{Tag: str("t"), Scope: str("s")}
	converter := bindings.NewTypeConverter()
	sv, errs := converter.ConvertToVapi(tag, model.TagBindingType())
	require.Empty(t, errs)

	_, _, err := fromGroupExpressionData([]*data.StructValue{sv.(*data.StructValue)})
	require.Error(t, err)
}

func TestUnitNsxt_buildIdentityGroupExpressionListData(t *testing.T) {
	identityGroups := []interface{}{
		map[string]interface{}{
			"distinguished_name":             "CN=group1",
			"domain_base_distinguished_name": "DC=example,DC=com",
			"sid":                            "S-1-5-21",
		},
	}
	sv, err := buildIdentityGroupExpressionListData(identityGroups)
	require.NoError(t, err)

	parsed, err := getIdentityGroupsData([]*data.StructValue{sv})
	require.NoError(t, err)
	require.Len(t, parsed, 1)
	assert.Equal(t, "CN=group1", *(parsed[0]["distinguished_name"].(*string)))
	assert.Equal(t, "S-1-5-21", *(parsed[0]["sid"].(*string)))
}
