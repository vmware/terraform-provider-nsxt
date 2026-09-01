//nolint:revive
package dns_services

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/dns_services"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type DnsZoneClientContext utl.ClientContext

func NewZonesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *DnsZoneClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Multitenancy, utl.VPC:
		client = client0.NewZonesClient(connector)

	default:
		return nil
	}
	return &DnsZoneClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c DnsZoneClientContext) Delete(orgIdParam string, projectIdParam string, dnsServiceIdParam string, zoneIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.ZonesClient)
		err = client.Delete(orgIdParam, projectIdParam, dnsServiceIdParam, zoneIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c DnsZoneClientContext) Get(orgIdParam string, projectIdParam string, dnsServiceIdParam string, zoneIdParam string) (model0.DnsZone, error) {
	var obj model0.DnsZone
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.ZonesClient)
		obj, err = client.Get(orgIdParam, projectIdParam, dnsServiceIdParam, zoneIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c DnsZoneClientContext) Patch(orgIdParam string, projectIdParam string, dnsServiceIdParam string, zoneIdParam string, projectDnsZoneParam model0.DnsZone) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.ZonesClient)
		err = client.Patch(orgIdParam, projectIdParam, dnsServiceIdParam, zoneIdParam, projectDnsZoneParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c DnsZoneClientContext) Update(orgIdParam string, projectIdParam string, dnsServiceIdParam string, zoneIdParam string, projectDnsZoneParam model0.DnsZone) (model0.DnsZone, error) {
	var err error
	var obj model0.DnsZone

	switch c.ClientType {

	case utl.Multitenancy, utl.VPC:
		client := c.Client.(client0.ZonesClient)
		obj, err = client.Update(orgIdParam, projectIdParam, dnsServiceIdParam, zoneIdParam, projectDnsZoneParam)

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
