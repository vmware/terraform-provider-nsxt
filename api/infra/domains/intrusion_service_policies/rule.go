//nolint:revive
package intrusionservicepolicies

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains/intrusion_service_policies"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra/domains/intrusion_service_policies"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type IntrusionServicePolicyRuleClientContext utl.ClientContext

func NewRulesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *IntrusionServicePolicyRuleClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewRulesClient(connector)

	case utl.Multitenancy:
		client = client1.NewRulesClient(connector)

	default:
		return nil
	}
	return &IntrusionServicePolicyRuleClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c IntrusionServicePolicyRuleClientContext) Get(domainIdParam string, policyIdParam string, ruleIdParam string) (model0.IdsRule, error) {
	var obj model0.IdsRule
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RulesClient)
		obj, err = client.Get(domainIdParam, policyIdParam, ruleIdParam)
		if err != nil {
			return obj, err
		}

	case utl.Multitenancy:
		client := c.Client.(client1.RulesClient)
		obj, err = client.Get(utl.DefaultOrgID, c.ProjectID, domainIdParam, policyIdParam, ruleIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c IntrusionServicePolicyRuleClientContext) Delete(domainIdParam string, policyIdParam string, ruleIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RulesClient)
		err = client.Delete(domainIdParam, policyIdParam, ruleIdParam)

	case utl.Multitenancy:
		client := c.Client.(client1.RulesClient)
		err = client.Delete(utl.DefaultOrgID, c.ProjectID, domainIdParam, policyIdParam, ruleIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c IntrusionServicePolicyRuleClientContext) Patch(domainIdParam string, policyIdParam string, ruleIdParam string, idsRuleParam model0.IdsRule) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RulesClient)
		err = client.Patch(domainIdParam, policyIdParam, ruleIdParam, idsRuleParam)

	case utl.Multitenancy:
		client := c.Client.(client1.RulesClient)
		err = client.Patch(utl.DefaultOrgID, c.ProjectID, domainIdParam, policyIdParam, ruleIdParam, idsRuleParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c IntrusionServicePolicyRuleClientContext) Update(domainIdParam string, policyIdParam string, ruleIdParam string, idsRuleParam model0.IdsRule) (model0.IdsRule, error) {
	var err error
	var obj model0.IdsRule

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RulesClient)
		obj, err = client.Update(domainIdParam, policyIdParam, ruleIdParam, idsRuleParam)

	case utl.Multitenancy:
		client := c.Client.(client1.RulesClient)
		obj, err = client.Update(utl.DefaultOrgID, c.ProjectID, domainIdParam, policyIdParam, ruleIdParam, idsRuleParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c IntrusionServicePolicyRuleClientContext) List(domainIdParam string, policyIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.IdsRuleListResult, error) {
	var err error
	var obj model0.IdsRuleListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RulesClient)
		obj, err = client.List(domainIdParam, policyIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	case utl.Multitenancy:
		client := c.Client.(client1.RulesClient)
		obj, err = client.List(utl.DefaultOrgID, c.ProjectID, domainIdParam, policyIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}
