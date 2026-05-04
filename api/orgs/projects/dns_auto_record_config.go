//nolint:revive
package projects

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type ProjectDnsAutoRecordConfigClientContext utl.ClientContext

func NewDnsAutoRecordConfigsClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *ProjectDnsAutoRecordConfigClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Multitenancy:
		client = client0.NewDnsAutoRecordConfigsClient(connector)

	default:
		return nil
	}
	return &ProjectDnsAutoRecordConfigClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c ProjectDnsAutoRecordConfigClientContext) Delete(orgIdParam string, projectIdParam string, configIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsAutoRecordConfigsClient)
		err = client.Delete(orgIdParam, projectIdParam, configIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c ProjectDnsAutoRecordConfigClientContext) Get(orgIdParam string, projectIdParam string, configIdParam string) (model0.ProjectDnsAutoRecordConfig, error) {
	var obj model0.ProjectDnsAutoRecordConfig
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsAutoRecordConfigsClient)
		obj, err = client.Get(orgIdParam, projectIdParam, configIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c ProjectDnsAutoRecordConfigClientContext) Patch(orgIdParam string, projectIdParam string, configIdParam string, projectDnsAutoRecordConfigParam model0.ProjectDnsAutoRecordConfig) error {
	var err error

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsAutoRecordConfigsClient)
		err = client.Patch(orgIdParam, projectIdParam, configIdParam, projectDnsAutoRecordConfigParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c ProjectDnsAutoRecordConfigClientContext) Update(orgIdParam string, projectIdParam string, configIdParam string, projectDnsAutoRecordConfigParam model0.ProjectDnsAutoRecordConfig) (model0.ProjectDnsAutoRecordConfig, error) {
	var err error
	var obj model0.ProjectDnsAutoRecordConfig

	switch c.ClientType {

	case utl.Multitenancy:
		client := c.Client.(client0.DnsAutoRecordConfigsClient)
		obj, err = client.Update(orgIdParam, projectIdParam, configIdParam, projectDnsAutoRecordConfigParam)

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}
