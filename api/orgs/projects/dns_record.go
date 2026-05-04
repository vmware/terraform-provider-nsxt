//nolint:revive
package projects

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type DnsRecordClientContext utl.ClientContext

func NewDnsRecordsClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *DnsRecordClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Multitenancy, utl.VPC:
		client = client0.NewDnsRecordsClient(connector)

	default:
		return nil
	}
	return &DnsRecordClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c DnsRecordClientContext) Delete(orgIdParam string, projectIdParam string, dnsRecordIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.DnsRecordsClient)
		err = client.Delete(orgIdParam, projectIdParam, dnsRecordIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c DnsRecordClientContext) Get(orgIdParam string, projectIdParam string, dnsRecordIdParam string) (model0.DnsRecord, error) {
	var obj model0.DnsRecord
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
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

func (c DnsRecordClientContext) Patch(orgIdParam string, projectIdParam string, dnsRecordIdParam string, projectDnsRecordParam model0.DnsRecord) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.DnsRecordsClient)
		err = client.Patch(orgIdParam, projectIdParam, dnsRecordIdParam, projectDnsRecordParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c DnsRecordClientContext) Update(orgIdParam string, projectIdParam string, dnsRecordIdParam string, projectDnsRecordParam model0.DnsRecord) (model0.DnsRecord, error) {
	var err error
	var obj model0.DnsRecord

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.DnsRecordsClient)
		obj, err = client.Update(orgIdParam, projectIdParam, dnsRecordIdParam, projectDnsRecordParam)

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
