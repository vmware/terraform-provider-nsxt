//nolint:revive
package dns_services

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/dns_services"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type DnsRuleClientContext utl.ClientContext

func NewRulesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *DnsRuleClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Multitenancy, utl.VPC:
		client = client0.NewRulesClient(connector)

	default:
		return nil
	}
	return &DnsRuleClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c DnsRuleClientContext) Delete(orgIdParam string, projectIdParam string, dnsServiceIdParam string, ruleIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.RulesClient)
		err = client.Delete(orgIdParam, projectIdParam, dnsServiceIdParam, ruleIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c DnsRuleClientContext) Get(orgIdParam string, projectIdParam string, dnsServiceIdParam string, ruleIdParam string) (model0.DnsRule, error) {
	var obj model0.DnsRule
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
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

func (c DnsRuleClientContext) Patch(orgIdParam string, projectIdParam string, dnsServiceIdParam string, ruleIdParam string, projectDnsRuleParam model0.DnsRule) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.RulesClient)
		err = client.Patch(orgIdParam, projectIdParam, dnsServiceIdParam, ruleIdParam, projectDnsRuleParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c DnsRuleClientContext) Update(orgIdParam string, projectIdParam string, dnsServiceIdParam string, ruleIdParam string, projectDnsRuleParam model0.DnsRule) (model0.DnsRule, error) {
	var err error
	var obj model0.DnsRule

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.RulesClient)
		obj, err = client.Update(orgIdParam, projectIdParam, dnsServiceIdParam, ruleIdParam, projectDnsRuleParam)

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
