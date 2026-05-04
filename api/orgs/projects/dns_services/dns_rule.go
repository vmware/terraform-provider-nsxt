//nolint:revive
package dns_services

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/dns_services"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type ProjectDnsRuleClientContext utl.ClientContext

func NewRulesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *ProjectDnsRuleClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Multitenancy:
		client = client0.NewRulesClient(connector)

	default:
		return nil
	}
	return &ProjectDnsRuleClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c ProjectDnsRuleClientContext) Delete(orgIdParam string, projectIdParam string, dnsServiceIdParam string, ruleIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.RulesClient)
		err = client.Delete(orgIdParam, projectIdParam, dnsServiceIdParam, ruleIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c ProjectDnsRuleClientContext) Get(orgIdParam string, projectIdParam string, dnsServiceIdParam string, ruleIdParam string) (model0.ProjectDnsRule, error) {
	var obj model0.ProjectDnsRule
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.RulesClient)
		obj, err = client.Get(orgIdParam, projectIdParam, dnsServiceIdParam, ruleIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c ProjectDnsRuleClientContext) Patch(orgIdParam string, projectIdParam string, dnsServiceIdParam string, ruleIdParam string, projectDnsRuleParam model0.ProjectDnsRule) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.RulesClient)
		err = client.Patch(orgIdParam, projectIdParam, dnsServiceIdParam, ruleIdParam, projectDnsRuleParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c ProjectDnsRuleClientContext) Update(orgIdParam string, projectIdParam string, dnsServiceIdParam string, ruleIdParam string, projectDnsRuleParam model0.ProjectDnsRule) (model0.ProjectDnsRule, error) {
	var err error
	var obj model0.ProjectDnsRule

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.RulesClient)
		obj, err = client.Update(orgIdParam, projectIdParam, dnsServiceIdParam, ruleIdParam, projectDnsRuleParam)

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
