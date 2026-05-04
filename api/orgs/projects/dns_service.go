//nolint:revive
package projects

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type PolicyDnsServiceClientContext utl.ClientContext

func NewDnsServicesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *PolicyDnsServiceClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Multitenancy:
		client = client0.NewDnsServicesClient(connector)

	default:
		return nil
	}
	return &PolicyDnsServiceClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c PolicyDnsServiceClientContext) Delete(orgIdParam string, projectIdParam string, dnsServiceIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsServicesClient)
		err = client.Delete(orgIdParam, projectIdParam, dnsServiceIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c PolicyDnsServiceClientContext) Get(orgIdParam string, projectIdParam string, dnsServiceIdParam string) (model0.PolicyDnsService, error) {
	var obj model0.PolicyDnsService
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsServicesClient)
		obj, err = client.Get(orgIdParam, projectIdParam, dnsServiceIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c PolicyDnsServiceClientContext) Patch(orgIdParam string, projectIdParam string, dnsServiceIdParam string, policyDnsServiceParam model0.PolicyDnsService) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsServicesClient)
		err = client.Patch(orgIdParam, projectIdParam, dnsServiceIdParam, policyDnsServiceParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c PolicyDnsServiceClientContext) Update(orgIdParam string, projectIdParam string, dnsServiceIdParam string, policyDnsServiceParam model0.PolicyDnsService) (model0.PolicyDnsService, error) {
	var err error
	var obj model0.PolicyDnsService

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsServicesClient)
		obj, err = client.Update(orgIdParam, projectIdParam, dnsServiceIdParam, policyDnsServiceParam)

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
