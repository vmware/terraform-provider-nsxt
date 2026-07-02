//nolint:revive
package projects

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type DnsServiceClientContext utl.ClientContext

func NewDnsServicesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *DnsServiceClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Multitenancy, utl.VPC:
		client = client0.NewDnsServicesClient(connector)

	default:
		return nil
	}
	return &DnsServiceClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c DnsServiceClientContext) Delete(orgIdParam string, projectIdParam string, dnsServiceIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.DnsServicesClient)
		err = client.Delete(orgIdParam, projectIdParam, dnsServiceIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c DnsServiceClientContext) Get(orgIdParam string, projectIdParam string, dnsServiceIdParam string) (model0.DnsService, error) {
	var obj model0.DnsService
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
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

func (c DnsServiceClientContext) Patch(orgIdParam string, projectIdParam string, dnsServiceIdParam string, policyDnsServiceParam model0.DnsService) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.DnsServicesClient)
		err = client.Patch(orgIdParam, projectIdParam, dnsServiceIdParam, policyDnsServiceParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c DnsServiceClientContext) Update(orgIdParam string, projectIdParam string, dnsServiceIdParam string, policyDnsServiceParam model0.DnsService) (model0.DnsService, error) {
	var err error
	var obj model0.DnsService

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.DnsServicesClient)
		obj, err = client.Update(orgIdParam, projectIdParam, dnsServiceIdParam, policyDnsServiceParam)

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
