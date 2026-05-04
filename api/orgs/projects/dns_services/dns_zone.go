//nolint:revive
package dns_services

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/dns_services"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type ProjectDnsZoneClientContext utl.ClientContext

func NewZonesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *ProjectDnsZoneClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Multitenancy:
		client = client0.NewZonesClient(connector)

	default:
		return nil
	}
	return &ProjectDnsZoneClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c ProjectDnsZoneClientContext) Delete(orgIdParam string, projectIdParam string, dnsServiceIdParam string, zoneIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.ZonesClient)
		err = client.Delete(orgIdParam, projectIdParam, dnsServiceIdParam, zoneIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c ProjectDnsZoneClientContext) Get(orgIdParam string, projectIdParam string, dnsServiceIdParam string, zoneIdParam string) (model0.ProjectDnsZone, error) {
	var obj model0.ProjectDnsZone
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
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

func (c ProjectDnsZoneClientContext) Patch(orgIdParam string, projectIdParam string, dnsServiceIdParam string, zoneIdParam string, projectDnsZoneParam model0.ProjectDnsZone) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.ZonesClient)
		err = client.Patch(orgIdParam, projectIdParam, dnsServiceIdParam, zoneIdParam, projectDnsZoneParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c ProjectDnsZoneClientContext) Update(orgIdParam string, projectIdParam string, dnsServiceIdParam string, zoneIdParam string, projectDnsZoneParam model0.ProjectDnsZone) (model0.ProjectDnsZone, error) {
	var err error
	var obj model0.ProjectDnsZone

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.ZonesClient)
		obj, err = client.Update(orgIdParam, projectIdParam, dnsServiceIdParam, zoneIdParam, projectDnsZoneParam)

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
