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

	t.Run("converts a rule with every action and condition type populated", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		name := "rule-full"
		strategy := model.LBRule_MATCH_STRATEGY_ALL
		phase := model.LBRule_PHASE_HTTP_ACCESS
		rules := []model.LBRule{
			{
				DisplayName:     &name,
				MatchStrategy:   &strategy,
				Phase:           &phase,
				Actions:         fullLbRuleActionStructValues(t),
				MatchConditions: fullLbRuleConditionStructValues(t),
			},
		}

		setPolicyLbRulesInSchema(d, rules)

		list := d.Get("rule").([]interface{})
		require.Len(t, list, 1)
		ruleElem := list[0].(map[string]interface{})
		actionElem := ruleElem["action"].([]interface{})[0].(map[string]interface{})
		conditionElem := ruleElem["condition"].([]interface{})[0].(map[string]interface{})

		for _, key := range []string{
			"connection_drop", "select_pool", "http_redirect", "http_request_uri_rewrite",
			"http_request_header_rewrite", "http_reject", "http_response_header_rewrite",
			"http_request_header_delete", "http_response_header_delete", "variable_assignment",
			"variable_persistence_on", "variable_persistence_learn", "jwt_auth", "ssl_mode_selection",
		} {
			assert.Lenf(t, actionElem[key], 1, "action block %q", key)
		}
		for _, key := range []string{
			"http_request_body", "http_request_uri", "http_request_header", "http_request_method",
			"http_request_uri_arguments", "http_request_version", "http_request_cookie",
			"http_response_header", "tcp_header", "ip_header", "variable", "ssl_sni", "http_ssl",
		} {
			assert.Lenf(t, conditionElem[key], 1, "condition block %q", key)
		}

		httpSslElem := conditionElem["http_ssl"].([]interface{})[0].(map[string]interface{})
		assert.Equal(t, 1, httpSslElem["client_certificate_issuer_dn"].(*schema.Set).Len())
		assert.Equal(t, 1, httpSslElem["client_certificate_subject_dn"].(*schema.Set).Len())
	})
}

// fullLbRuleActionStructValues returns one *data.StructValue per action type (mirroring
// fullLbRuleActionData's schema-side fixture), so a single setPolicyLbRulesInSchema call
// exercises every action-type case in its switch statement.
func fullLbRuleActionStructValues(t *testing.T) []*data.StructValue {
	t.Helper()
	poolID := "/infra/lb-pools/pool-1"
	redirectStatus, redirectURL := "301", "http://example.com"
	uri, uriArgs := "/new", "a=b"
	headerName, headerValue := "X-Foo", "bar"
	replyMsg, replyStatus := "denied", "403"
	respHeaderName, respHeaderValue := "X-Bar", "baz"
	respDeleteHeaderName := "X-Bar"
	varName, varValue := "v1", "val1"
	persistPath, persistVar := "/infra/lb-persistence-profiles/p1", "v2"
	persistHash := true
	sslMode := model.LBSslModeSelectionAction_SSL_MODE_OFFLOAD
	realm := "myrealm"
	passToPool := true
	certPath := "/infra/certs/jwt"
	keySV := lbVsConverterStructValue(t, model.LBJwtCertificateKey{
		Type_:           model.LBJwtKey_TYPE_LBJWTCERTIFICATEKEY,
		CertificatePath: &certPath,
	}, model.LBJwtCertificateKeyBindingType())

	return []*data.StructValue{
		lbVsConverterStructValue(t, model.LBConnectionDropAction{Type_: model.LBRuleAction_TYPE_LBCONNECTIONDROPACTION}, model.LBConnectionDropActionBindingType()),
		lbVsConverterStructValue(t, model.LBSelectPoolAction{Type_: model.LBRuleAction_TYPE_LBSELECTPOOLACTION, PoolId: &poolID}, model.LBSelectPoolActionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRedirectAction{Type_: model.LBRuleAction_TYPE_LBHTTPREDIRECTACTION, RedirectStatus: &redirectStatus, RedirectUrl: &redirectURL}, model.LBHttpRedirectActionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRequestUriRewriteAction{Type_: model.LBRuleAction_TYPE_LBHTTPREQUESTURIREWRITEACTION, Uri: &uri, UriArguments: &uriArgs}, model.LBHttpRequestUriRewriteActionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRequestHeaderRewriteAction{Type_: model.LBRuleAction_TYPE_LBHTTPREQUESTHEADERREWRITEACTION, HeaderName: &headerName, HeaderValue: &headerValue}, model.LBHttpRequestHeaderRewriteActionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRejectAction{Type_: model.LBRuleAction_TYPE_LBHTTPREJECTACTION, ReplyMessage: &replyMsg, ReplyStatus: &replyStatus}, model.LBHttpRejectActionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpResponseHeaderRewriteAction{Type_: model.LBRuleAction_TYPE_LBHTTPRESPONSEHEADERREWRITEACTION, HeaderName: &respHeaderName, HeaderValue: &respHeaderValue}, model.LBHttpResponseHeaderRewriteActionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRequestHeaderDeleteAction{Type_: model.LBRuleAction_TYPE_LBHTTPREQUESTHEADERDELETEACTION, HeaderName: &headerName}, model.LBHttpRequestHeaderDeleteActionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpResponseHeaderDeleteAction{Type_: model.LBRuleAction_TYPE_LBHTTPRESPONSEHEADERDELETEACTION, HeaderName: &respDeleteHeaderName}, model.LBHttpResponseHeaderDeleteActionBindingType()),
		lbVsConverterStructValue(t, model.LBVariableAssignmentAction{Type_: model.LBRuleAction_TYPE_LBVARIABLEASSIGNMENTACTION, VariableName: &varName, VariableValue: &varValue}, model.LBVariableAssignmentActionBindingType()),
		lbVsConverterStructValue(t, model.LBVariablePersistenceOnAction{Type_: model.LBRuleAction_TYPE_LBVARIABLEPERSISTENCEONACTION, PersistenceProfilePath: &persistPath, VariableName: &persistVar, VariableHashEnabled: &persistHash}, model.LBVariablePersistenceOnActionBindingType()),
		lbVsConverterStructValue(t, model.LBVariablePersistenceLearnAction{Type_: model.LBRuleAction_TYPE_LBVARIABLEPERSISTENCELEARNACTION, PersistenceProfilePath: &persistPath, VariableName: &persistVar, VariableHashEnabled: &persistHash}, model.LBVariablePersistenceLearnActionBindingType()),
		lbVsConverterStructValue(t, model.LBSslModeSelectionAction{Type_: model.LBRuleAction_TYPE_LBSSLMODESELECTIONACTION, SslMode: &sslMode}, model.LBSslModeSelectionActionBindingType()),
		lbVsConverterStructValue(t, model.LBJwtAuthAction{
			Type_:         model.LBRuleAction_TYPE_LBJWTAUTHACTION,
			Realm:         &realm,
			PassJwtToPool: &passToPool,
			Key:           keySV,
			Tokens:        []string{"Authorization"},
		}, model.LBJwtAuthActionBindingType()),
	}
}

// fullLbRuleConditionStructValues returns one *data.StructValue per condition type (mirroring
// fullLbRuleConditionData's schema-side fixture), so a single setPolicyLbRulesInSchema call
// exercises every condition-type case in its switch statement.
func fullLbRuleConditionStructValues(t *testing.T) []*data.StructValue {
	t.Helper()
	trueVal, falseVal := true, false
	matchEquals := model.LBHttpRequestBodyCondition_MATCH_TYPE_EQUALS
	bodyValue := "v"
	cookieName, cookieValue := "c1", "v1"
	headerName, headerValue := "X-Foo", "bar"
	method := "GET"
	uriArgs := "a=b"
	uri := "/foo"
	version := "HTTP_VERSION_1_1"
	respHeaderName, respHeaderValue := "X-Bar", "baz"
	groupPath, sourceAddr := "/infra/domains/default/groups/g1", "1.2.3.4"
	sni := "example.com"
	sourcePort := "8080"
	varName, varValue := "v1", "val1"
	issuerDn := "CN=issuer"
	subjectDn := "CN=subject"
	sessionReused, usedProtocol, usedCipher := "IGNORE", "TLS_V1_2", "TLS_RSA_WITH_AES_128_CBC_SHA"

	return []*data.StructValue{
		lbVsConverterStructValue(t, model.LBHttpRequestBodyCondition{Type_: model.LBRuleCondition_TYPE_LBHTTPREQUESTBODYCONDITION, BodyValue: &bodyValue, MatchType: &matchEquals, CaseSensitive: &trueVal, Inverse: &falseVal}, model.LBHttpRequestBodyConditionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRequestUriCondition{Type_: model.LBRuleCondition_TYPE_LBHTTPREQUESTURICONDITION, Uri: &uri, MatchType: &matchEquals, CaseSensitive: &trueVal, Inverse: &falseVal}, model.LBHttpRequestUriConditionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRequestHeaderCondition{Type_: model.LBRuleCondition_TYPE_LBHTTPREQUESTHEADERCONDITION, HeaderName: &headerName, HeaderValue: &headerValue, MatchType: &matchEquals, CaseSensitive: &trueVal, Inverse: &falseVal}, model.LBHttpRequestHeaderConditionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRequestMethodCondition{Type_: model.LBRuleCondition_TYPE_LBHTTPREQUESTMETHODCONDITION, Method: &method, Inverse: &falseVal}, model.LBHttpRequestMethodConditionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRequestUriArgumentsCondition{Type_: model.LBRuleCondition_TYPE_LBHTTPREQUESTURIARGUMENTSCONDITION, UriArguments: &uriArgs, MatchType: &matchEquals, CaseSensitive: &trueVal, Inverse: &falseVal}, model.LBHttpRequestUriArgumentsConditionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRequestVersionCondition{Type_: model.LBRuleCondition_TYPE_LBHTTPREQUESTVERSIONCONDITION, Version: &version, Inverse: &falseVal}, model.LBHttpRequestVersionConditionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpRequestCookieCondition{Type_: model.LBRuleCondition_TYPE_LBHTTPREQUESTCOOKIECONDITION, CookieName: &cookieName, CookieValue: &cookieValue, MatchType: &matchEquals, CaseSensitive: &trueVal, Inverse: &falseVal}, model.LBHttpRequestCookieConditionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpResponseHeaderCondition{Type_: model.LBRuleCondition_TYPE_LBHTTPRESPONSEHEADERCONDITION, HeaderName: &respHeaderName, HeaderValue: &respHeaderValue, MatchType: &matchEquals, CaseSensitive: &trueVal, Inverse: &falseVal}, model.LBHttpResponseHeaderConditionBindingType()),
		lbVsConverterStructValue(t, model.LBTcpHeaderCondition{Type_: model.LBRuleCondition_TYPE_LBTCPHEADERCONDITION, SourcePort: &sourcePort, Inverse: &falseVal}, model.LBTcpHeaderConditionBindingType()),
		lbVsConverterStructValue(t, model.LBIpHeaderCondition{Type_: model.LBRuleCondition_TYPE_LBIPHEADERCONDITION, GroupPath: &groupPath, SourceAddress: &sourceAddr, Inverse: &falseVal}, model.LBIpHeaderConditionBindingType()),
		lbVsConverterStructValue(t, model.LBVariableCondition{Type_: model.LBRuleCondition_TYPE_LBVARIABLECONDITION, VariableName: &varName, VariableValue: &varValue, MatchType: &matchEquals, CaseSensitive: &trueVal, Inverse: &falseVal}, model.LBVariableConditionBindingType()),
		lbVsConverterStructValue(t, model.LBSslSniCondition{Type_: model.LBRuleCondition_TYPE_LBSSLSNICONDITION, Sni: &sni, MatchType: &matchEquals, CaseSensitive: &trueVal, Inverse: &falseVal}, model.LBSslSniConditionBindingType()),
		lbVsConverterStructValue(t, model.LBHttpSslCondition{
			Type_:   model.LBRuleCondition_TYPE_LBHTTPSSLCONDITION,
			Inverse: &falseVal,
			ClientCertificateIssuerDn: &model.LBClientCertificateIssuerDnCondition{
				CaseSensitive: &trueVal, IssuerDn: &issuerDn, MatchType: &matchEquals,
			},
			ClientCertificateSubjectDn: &model.LBClientCertificateSubjectDnCondition{
				CaseSensitive: &trueVal, SubjectDn: &subjectDn, MatchType: &matchEquals,
			},
			ClientSupportedSslCiphers: []string{"TLS_RSA_WITH_AES_128_CBC_SHA"},
			SessionReused:             &sessionReused,
			UsedProtocol:              &usedProtocol,
			UsedSslCipher:             &usedCipher,
		}, model.LBHttpSslConditionBindingType()),
	}
}

// fullLbRuleActionData returns an "action" block populated with every simple
// (string/bool-field) action type plus a jwt_auth action using a certificate key, so a single
// getPolicyLbRulesFromSchema call exercises every action-type append call site.
func fullLbRuleActionData() map[string]interface{} {
	return map[string]interface{}{
		"connection_drop": []interface{}{map[string]interface{}{}},
		"http_redirect": []interface{}{map[string]interface{}{
			"redirect_status": "301",
			"redirect_url":    "http://example.com",
		}},
		"http_reject": []interface{}{map[string]interface{}{
			"reply_message": "denied",
			"reply_status":  "403",
		}},
		"http_request_header_delete": []interface{}{map[string]interface{}{
			"header_name": "X-Foo",
		}},
		"http_request_header_rewrite": []interface{}{map[string]interface{}{
			"header_name":  "X-Foo",
			"header_value": "bar",
		}},
		"http_request_uri_rewrite": []interface{}{map[string]interface{}{
			"uri":           "/new",
			"uri_arguments": "a=b",
		}},
		"http_response_header_delete": []interface{}{map[string]interface{}{
			"header_name": "X-Bar",
		}},
		"http_response_header_rewrite": []interface{}{map[string]interface{}{
			"header_name":  "X-Bar",
			"header_value": "baz",
		}},
		"select_pool": []interface{}{map[string]interface{}{
			"pool_id": "/infra/lb-pools/pool-1",
		}},
		"ssl_mode_selection": []interface{}{map[string]interface{}{
			"ssl_mode": "SSL",
		}},
		"variable_assignment": []interface{}{map[string]interface{}{
			"variable_name":  "v1",
			"variable_value": "val1",
		}},
		"variable_persistence_learn": []interface{}{map[string]interface{}{
			"persistence_profile_path": "/infra/lb-persistence-profiles/p1",
			"variable_name":            "v2",
			"variable_hash_enabled":    true,
		}},
		"variable_persistence_on": []interface{}{map[string]interface{}{
			"persistence_profile_path": "/infra/lb-persistence-profiles/p2",
			"variable_name":            "v3",
			"variable_hash_enabled":    false,
		}},
		"jwt_auth": []interface{}{map[string]interface{}{
			"realm":            "myrealm",
			"pass_jwt_to_pool": true,
			"key": []interface{}{map[string]interface{}{
				"certificate_path": "/infra/certs/jwt",
			}},
			"tokens": []interface{}{"Authorization"},
		}},
	}
}

// fullLbRuleConditionData returns a "condition" block populated with every simple
// (string/bool-field) condition type plus an http_ssl condition using both DN sets, so a
// single getPolicyLbRulesFromSchema call exercises every condition-type append call site.
func fullLbRuleConditionData() map[string]interface{} {
	return map[string]interface{}{
		"http_request_body": []interface{}{map[string]interface{}{
			"body_value": "v", "match_type": "EQUALS", "case_sensitive": true, "inverse": false,
		}},
		"http_request_cookie": []interface{}{map[string]interface{}{
			"cookie_name": "c1", "cookie_value": "v1", "match_type": "EQUALS", "case_sensitive": true, "inverse": false,
		}},
		"http_request_header": []interface{}{map[string]interface{}{
			"header_name": "X-Foo", "header_value": "bar", "match_type": "EQUALS", "case_sensitive": true, "inverse": false,
		}},
		"http_request_method": []interface{}{map[string]interface{}{
			"method": "GET", "inverse": false,
		}},
		"http_request_uri_arguments": []interface{}{map[string]interface{}{
			"uri_arguments": "a=b", "match_type": "EQUALS", "case_sensitive": true, "inverse": false,
		}},
		"http_request_uri": []interface{}{map[string]interface{}{
			"uri": "/foo", "match_type": "EQUALS", "case_sensitive": true, "inverse": false,
		}},
		"http_request_version": []interface{}{map[string]interface{}{
			"version": "HTTP_VERSION_1_1", "inverse": false,
		}},
		"http_response_header": []interface{}{map[string]interface{}{
			"header_name": "X-Bar", "header_value": "baz", "match_type": "EQUALS", "case_sensitive": true, "inverse": false,
		}},
		"ip_header": []interface{}{map[string]interface{}{
			"group_path": "/infra/domains/default/groups/g1", "source_address": "1.2.3.4", "inverse": false,
		}},
		"ssl_sni": []interface{}{map[string]interface{}{
			"sni": "example.com", "match_type": "EQUALS", "case_sensitive": true, "inverse": false,
		}},
		"tcp_header": []interface{}{map[string]interface{}{
			"source_port": "8080", "inverse": false,
		}},
		"variable": []interface{}{map[string]interface{}{
			"match_type": "EQUALS", "variable_name": "v1", "variable_value": "val1", "case_sensitive": true, "inverse": false,
		}},
		"http_ssl": []interface{}{map[string]interface{}{
			"inverse": false,
			"client_certificate_issuer_dn": []interface{}{map[string]interface{}{
				"issuer_dn": "CN=issuer", "case_sensitive": true, "match_type": "EQUALS",
			}},
			"client_certificate_subject_dn": []interface{}{map[string]interface{}{
				"subject_dn": "CN=subject", "case_sensitive": true, "match_type": "EQUALS",
			}},
			"client_supported_ssl_ciphers": []interface{}{"TLS_RSA_WITH_AES_128_CBC_SHA"},
			"session_reused":               "IGNORE",
			"used_protocol":                "TLS_V1_2",
			"used_ssl_cipher":              "TLS_RSA_WITH_AES_128_CBC_SHA",
		}},
	}
}

func TestUnitNsxt_getPolicyLbRulesFromSchema(t *testing.T) {
	res := resourceNsxtPolicyLBVirtualServer()

	t.Run("converts a rule with every action and condition type populated", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{
					"display_name":   "rule-1",
					"match_strategy": "ALL",
					"phase":          "HTTP_ACCESS",
					"action":         []interface{}{fullLbRuleActionData()},
					"condition":      []interface{}{fullLbRuleConditionData()},
				},
			},
		})

		rules := getPolicyLbRulesFromSchema(d)
		require.Len(t, rules, 1)
		rule := rules[0]
		assert.Equal(t, "rule-1", *rule.DisplayName)
		assert.Equal(t, "ALL", *rule.MatchStrategy)
		assert.Equal(t, "HTTP_ACCESS", *rule.Phase)
		// 13 simple actions + 1 jwt_auth = 14
		assert.Len(t, rule.Actions, 14)
		// 12 simple conditions + 1 http_ssl = 13
		assert.Len(t, rule.MatchConditions, 13)
	})

	t.Run("jwt_auth with a public key", func(t *testing.T) {
		action := fullLbRuleActionData()
		action["jwt_auth"] = []interface{}{map[string]interface{}{
			"realm": "myrealm",
			"key": []interface{}{map[string]interface{}{
				"public_key_content": "-----BEGIN PUBLIC KEY-----",
			}},
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{
					"display_name": "rule-2",
					"action":       []interface{}{action},
				},
			},
		})

		rules := getPolicyLbRulesFromSchema(d)
		require.Len(t, rules, 1)
		assert.Len(t, rules[0].Actions, 14)
	})

	t.Run("jwt_auth with a symmetric key", func(t *testing.T) {
		action := fullLbRuleActionData()
		action["jwt_auth"] = []interface{}{map[string]interface{}{
			"realm": "myrealm",
			"key": []interface{}{map[string]interface{}{
				"symmetric_key": "shared-secret",
			}},
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{
					"display_name": "rule-3",
					"action":       []interface{}{action},
				},
			},
		})

		rules := getPolicyLbRulesFromSchema(d)
		require.Len(t, rules, 1)
		assert.Len(t, rules[0].Actions, 14)
	})

	t.Run("no rules returns nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		assert.Nil(t, getPolicyLbRulesFromSchema(d))
	})
}
