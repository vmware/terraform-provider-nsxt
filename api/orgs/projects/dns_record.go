//nolint:revive
package projects

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type ProjectDnsRecordClientContext utl.ClientContext

func NewDnsRecordsClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *ProjectDnsRecordClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Multitenancy:
		client = client0.NewDnsRecordsClient(connector)

	default:
		return nil
	}
	return &ProjectDnsRecordClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c ProjectDnsRecordClientContext) Delete(orgIdParam string, projectIdParam string, dnsRecordIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsRecordsClient)
		err = client.Delete(orgIdParam, projectIdParam, dnsRecordIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c ProjectDnsRecordClientContext) Get(orgIdParam string, projectIdParam string, dnsRecordIdParam string) (model0.ProjectDnsRecord, error) {
	var obj model0.ProjectDnsRecord
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsRecordsClient)
		obj, err = client.Get(orgIdParam, projectIdParam, dnsRecordIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c ProjectDnsRecordClientContext) Patch(orgIdParam string, projectIdParam string, dnsRecordIdParam string, projectDnsRecordParam model0.ProjectDnsRecord) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsRecordsClient)
		err = client.Patch(orgIdParam, projectIdParam, dnsRecordIdParam, projectDnsRecordParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c ProjectDnsRecordClientContext) Update(orgIdParam string, projectIdParam string, dnsRecordIdParam string, projectDnsRecordParam model0.ProjectDnsRecord) (model0.ProjectDnsRecord, error) {
	var err error
	var obj model0.ProjectDnsRecord

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsRecordsClient)
		obj, err = client.Update(orgIdParam, projectIdParam, dnsRecordIdParam, projectDnsRecordParam)

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
