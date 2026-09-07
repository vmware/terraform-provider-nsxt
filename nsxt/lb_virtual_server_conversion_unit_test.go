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

// These tests exercise the schema<->struct conversion helpers in
// resource_nsxt_policy_lb_virtual_server.go directly. They are pure logic (no SDK client
// calls), so no mocking is needed - only a *schema.ResourceData built from the resource's own
// schema via schema.TestResourceDataRaw.

func lbVsConverterStructValue(t *testing.T, obj interface{}, bt bindings.BindingType) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	sv, errs := converter.ConvertToVapi(obj, bt)
	require.Empty(t, errs)
	return sv.(*data.StructValue)
}

func TestUnitNsxt_getPolicyClientSSLBindingFromSchema(t *testing.T) {
	res := resourceNsxtPolicyLBVirtualServer()

	t.Run("empty when no client_ssl block is set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		assert.Nil(t, getPolicyClientSSLBindingFromSchema(d))
	})

	t.Run("parses a client_ssl block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"client_ssl": []interface{}{
				map[string]interface{}{
					"certificate_chain_depth":  3,
					"client_auth":              model.LBClientSslProfileBinding_CLIENT_AUTH_REQUIRED,
					"ca_paths":                 []interface{}{"/infra/certs/ca1"},
					"crl_paths":                []interface{}{"/infra/crls/crl1"},
					"default_certificate_path": "/infra/certs/default",
					"ssl_profile_path":         "/infra/ssl-profiles/p1",
					"sni_paths":                []interface{}{"/infra/certs/sni1"},
				},
			},
		})

		binding := getPolicyClientSSLBindingFromSchema(d)
		require.NotNil(t, binding)
		assert.Equal(t, int64(3), *binding.CertificateChainDepth)
		assert.Equal(t, model.LBClientSslProfileBinding_CLIENT_AUTH_REQUIRED, *binding.ClientAuth)
		assert.Equal(t, []string{"/infra/certs/ca1"}, binding.ClientAuthCaPaths)
		assert.Equal(t, "/infra/certs/default", *binding.DefaultCertificatePath)
		assert.Equal(t, "/infra/ssl-profiles/p1", *binding.SslProfilePath)
	})
}

func TestUnitNsxt_setPolicyClientSSLBindingInSchema(t *testing.T) {
	res := resourceNsxtPolicyLBVirtualServer()

	t.Run("clears the block when binding is nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		setPolicyClientSSLBindingInSchema(d, nil)
		assert.Empty(t, d.Get("client_ssl").([]interface{}))
	})

	t.Run("sets the block from a binding", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		depth := int64(2)
		auth := model.LBClientSslProfileBinding_CLIENT_AUTH_IGNORE
		cert := "/infra/certs/default"
		profile := "/infra/ssl-profiles/p1"
		setPolicyClientSSLBindingInSchema(d, &model.LBClientSslProfileBinding{
			CertificateChainDepth:  &depth,
			ClientAuth:             &auth,
			DefaultCertificatePath: &cert,
			SslProfilePath:         &profile,
		})

		list := d.Get("client_ssl").([]interface{})
		require.Len(t, list, 1)
		elem := list[0].(map[string]interface{})
		assert.Equal(t, 2, elem["certificate_chain_depth"])
		assert.Equal(t, auth, elem["client_auth"])
	})
}

func TestUnitNsxt_getPolicyServerSSLBindingFromSchema(t *testing.T) {
	res := resourceNsxtPolicyLBVirtualServer()

	t.Run("empty when no server_ssl block is set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		assert.Nil(t, getPolicyServerSSLBindingFromSchema(d))
	})

	t.Run("parses a server_ssl block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"server_ssl": []interface{}{
				map[string]interface{}{
					"certificate_chain_depth": 4,
					"server_auth":             model.LBServerSslProfileBinding_SERVER_AUTH_REQUIRED,
					"ca_paths":                []interface{}{"/infra/certs/ca2"},
					"crl_paths":               []interface{}{"/infra/crls/crl2"},
					"client_certificate_path": "/infra/certs/client",
					"ssl_profile_path":        "/infra/ssl-profiles/p2",
				},
			},
		})

		binding := getPolicyServerSSLBindingFromSchema(d)
		require.NotNil(t, binding)
		assert.Equal(t, int64(4), *binding.CertificateChainDepth)
		assert.Equal(t, model.LBServerSslProfileBinding_SERVER_AUTH_REQUIRED, *binding.ServerAuth)
		assert.Equal(t, "/infra/certs/client", *binding.ClientCertificatePath)
		assert.Equal(t, "/infra/ssl-profiles/p2", *binding.SslProfilePath)
	})
}

func TestUnitNsxt_setPolicyServerSSLBindingInSchema(t *testing.T) {
	res := resourceNsxtPolicyLBVirtualServer()

	t.Run("clears the block when binding is nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		setPolicyServerSSLBindingInSchema(d, nil)
		assert.Empty(t, d.Get("server_ssl").([]interface{}))
	})

	t.Run("sets the block from a binding", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		depth := int64(1)
		auth := model.LBServerSslProfileBinding_SERVER_AUTH_AUTO_APPLY
		profile := "/infra/ssl-profiles/p2"
		setPolicyServerSSLBindingInSchema(d, &model.LBServerSslProfileBinding{
			CertificateChainDepth: &depth,
			ServerAuth:            &auth,
			SslProfilePath:        &profile,
		})

		list := d.Get("server_ssl").([]interface{})
		require.Len(t, list, 1)
		elem := list[0].(map[string]interface{})
		assert.Equal(t, 1, elem["certificate_chain_depth"])
		assert.Equal(t, auth, elem["server_auth"])
	})
}

func TestUnitNsxt_getPolicyAccessListControlFromSchema(t *testing.T) {
	res := resourceNsxtPolicyLBVirtualServer()

	t.Run("empty when no access_list_control block is set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		assert.Nil(t, getPolicyAccessListControlFromSchema(d))
	})

	t.Run("parses an access_list_control block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"access_list_control": []interface{}{
				map[string]interface{}{
					"action":     model.LBAccessListControl_ACTION_DROP,
					"enabled":    true,
					"group_path": "/infra/domains/default/groups/g1",
				},
			},
		})

		control := getPolicyAccessListControlFromSchema(d)
		require.NotNil(t, control)
		assert.Equal(t, model.LBAccessListControl_ACTION_DROP, *control.Action)
		assert.True(t, *control.Enabled)
		assert.Equal(t, "/infra/domains/default/groups/g1", *control.GroupPath)
	})
}

func TestUnitNsxt_setPolicyAccessListControlInSchema(t *testing.T) {
	res := resourceNsxtPolicyLBVirtualServer()

	t.Run("clears the block when control is nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		setPolicyAccessListControlInSchema(d, nil)
		assert.Empty(t, d.Get("access_list_control").([]interface{}))
	})

	t.Run("sets the block from a control", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		action := model.LBAccessListControl_ACTION_ALLOW
		enabled := true
		groupPath := "/infra/domains/default/groups/g1"
		setPolicyAccessListControlInSchema(d, &model.LBAccessListControl{
			Action:    &action,
			Enabled:   &enabled,
			GroupPath: &groupPath,
		})

		list := d.Get("access_list_control").([]interface{})
		require.Len(t, list, 1)
		elem := list[0].(map[string]interface{})
		assert.Equal(t, action, elem["action"])
		assert.Equal(t, groupPath, elem["group_path"])
	})
}

func TestUnitNsxt_getRuleActionOrMethod(t *testing.T) {
	t.Run("empty list yields no results", func(t *testing.T) {
		ruleData := map[string]interface{}{"connection_drop": []interface{}{}}
		result := getRuleActionOrMethod(ruleData, "connection_drop", nil, nil, model.LBRuleAction_TYPE_LBCONNECTIONDROPACTION)
		assert.Empty(t, result)
	})

	t.Run("builds a StructValue with string and bool fields set", func(t *testing.T) {
		ruleData := map[string]interface{}{
			"http_request_method": []interface{}{
				map[string]interface{}{"method": "GET", "inverse": true},
			},
		}
		result := getRuleActionOrMethod(ruleData, "http_request_method", []string{"method"}, []string{"inverse"}, model.LBRuleCondition_TYPE_LBHTTPREQUESTMETHODCONDITION)
		require.Len(t, result, 1)

		converter := bindings.NewTypeConverter()
		golang, errs := converter.ConvertToGolang(result[0], model.LBHttpRequestMethodConditionBindingType())
		require.Empty(t, errs)
		cond := golang.(model.LBHttpRequestMethodCondition)
		assert.Equal(t, "GET", *cond.Method)
		assert.True(t, *cond.Inverse)
	})

	t.Run("skips empty string fields", func(t *testing.T) {
		ruleData := map[string]interface{}{
			"select_pool": []interface{}{
				map[string]interface{}{"pool_id": ""},
			},
		}
		result := getRuleActionOrMethod(ruleData, "select_pool", []string{"pool_id"}, nil, model.LBRuleAction_TYPE_LBSELECTPOOLACTION)
		require.Len(t, result, 1)

		converter := bindings.NewTypeConverter()
		golang, errs := converter.ConvertToGolang(result[0], model.LBSelectPoolActionBindingType())
		require.Empty(t, errs)
		action := golang.(model.LBSelectPoolAction)
		assert.Nil(t, action.PoolId)
	})
}

func TestUnitNsxt_setPolicyLbRulesInSchema(t *testing.T) {
	res := resourceNsxtPolicyLBVirtualServer()

	t.Run("empty rule list clears the block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		setPolicyLbRulesInSchema(d, nil)
		assert.Empty(t, d.Get("rule").([]interface{}))
	})

	t.Run("converts a rule with simple actions and conditions", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		poolID := "/infra/lb-pools/pool-1"
		selectPoolSV := lbVsConverterStructValue(t, model.LBSelectPoolAction{
			Type_:  model.LBRuleAction_TYPE_LBSELECTPOOLACTION,
			PoolId: &poolID,
		}, model.LBSelectPoolActionBindingType())

		method := "GET"
		methodSV := lbVsConverterStructValue(t, model.LBHttpRequestMethodCondition{
			Type_:  model.LBRuleCondition_TYPE_LBHTTPREQUESTMETHODCONDITION,
			Method: &method,
		}, model.LBHttpRequestMethodConditionBindingType())

		name := "rule-1"
		strategy := model.LBRule_MATCH_STRATEGY_ALL
		phase := model.LBRule_PHASE_HTTP_ACCESS
		rules := []model.LBRule{
			{
				DisplayName:     &name,
				MatchStrategy:   &strategy,
				Phase:           &phase,
				Actions:         []*data.StructValue{selectPoolSV},
				MatchConditions: []*data.StructValue{methodSV},
			},
		}

		setPolicyLbRulesInSchema(d, rules)

		list := d.Get("rule").([]interface{})
		require.Len(t, list, 1)
		ruleElem := list[0].(map[string]interface{})
		assert.Equal(t, name, ruleElem["display_name"])

		actionList := ruleElem["action"].([]interface{})
		require.Len(t, actionList, 1)
		actionElem := actionList[0].(map[string]interface{})
		selectPoolList := actionElem["select_pool"].([]interface{})
		require.Len(t, selectPoolList, 1)
		assert.Equal(t, poolID, selectPoolList[0].(map[string]interface{})["pool_id"])

		conditionList := ruleElem["condition"].([]interface{})
		require.Len(t, conditionList, 1)
		conditionElem := conditionList[0].(map[string]interface{})
		methodList := conditionElem["http_request_method"].([]interface{})
		require.Len(t, methodList, 1)
		assert.Equal(t, method, methodList[0].(map[string]interface{})["method"])
	})

	t.Run("converts a jwt_auth action with a certificate key", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		realm := "myrealm"
		passToPool := true
		certPath := "/infra/certs/jwt"
		keySV := lbVsConverterStructValue(t, model.LBJwtCertificateKey{
			Type_:           model.LBJwtKey_TYPE_LBJWTCERTIFICATEKEY,
			CertificatePath: &certPath,
		}, model.LBJwtCertificateKeyBindingType())
		jwtAction := model.LBJwtAuthAction{
			Type_:         model.LBRuleAction_TYPE_LBJWTAUTHACTION,
			Realm:         &realm,
			PassJwtToPool: &passToPool,
			Key:           keySV,
			Tokens:        []string{"Authorization"},
		}
		jwtSV := lbVsConverterStructValue(t, jwtAction, model.LBJwtAuthActionBindingType())

		name := "rule-jwt"
		strategy := model.LBRule_MATCH_STRATEGY_ANY
		phase := model.LBRule_PHASE_HTTP_ACCESS
		rules := []model.LBRule{
			{
				DisplayName:   &name,
				MatchStrategy: &strategy,
				Phase:         &phase,
				Actions:       []*data.StructValue{jwtSV},
			},
		}

		setPolicyLbRulesInSchema(d, rules)

		list := d.Get("rule").([]interface{})
		require.Len(t, list, 1)
		actionElem := list[0].(map[string]interface{})["action"].([]interface{})[0].(map[string]interface{})
		jwtList := actionElem["jwt_auth"].([]interface{})
		require.Len(t, jwtList, 1)
		jwtElem := jwtList[0].(map[string]interface{})
		assert.Equal(t, realm, jwtElem["realm"])
		keySet := jwtElem["key"].(*schema.Set)
		require.Equal(t, 1, keySet.Len())
	})
}
