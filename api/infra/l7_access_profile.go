//nolint:revive
package infra

// The following file has been autogenerated. Please avoid any changes!
import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type L7AccessProfileClientContext utl.ClientContext

func NewL7AccessProfilesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *L7AccessProfileClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewL7AccessProfilesClient(connector)

	case utl.Multitenancy:
		client = client1.NewL7AccessProfilesClient(connector)

	default:
		return nil
	}
	return &L7AccessProfileClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c L7AccessProfileClientContext) Get(l7AccessProfileIdParam string) (model0.L7AccessProfile, error) {
	var obj model0.L7AccessProfile
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.L7AccessProfilesClient)
		obj, err = client.Get(l7AccessProfileIdParam)
		if err != nil {
			return obj, err
		}

	case utl.Multitenancy:
		client := c.Client.(client1.L7AccessProfilesClient)
		obj, err = client.Get(utl.DefaultOrgID, c.ProjectID, l7AccessProfileIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c L7AccessProfileClientContext) Delete(l7AccessProfileIdParam string, overrideParam *bool) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.L7AccessProfilesClient)
		err = client.Delete(l7AccessProfileIdParam, overrideParam)

	case utl.Multitenancy:
		client := c.Client.(client1.L7AccessProfilesClient)
		err = client.Delete(utl.DefaultOrgID, c.ProjectID, l7AccessProfileIdParam, overrideParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c L7AccessProfileClientContext) Patch(l7AccessProfileIdParam string, l7AccessProfileParam model0.L7AccessProfile, overrideParam *bool) (model0.L7AccessProfile, error) {
	var err error
	var obj model0.L7AccessProfile

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.L7AccessProfilesClient)
		obj, err = client.Patch(l7AccessProfileIdParam, l7AccessProfileParam, overrideParam)

	case utl.Multitenancy:
		client := c.Client.(client1.L7AccessProfilesClient)
		obj, err = client.Patch(utl.DefaultOrgID, c.ProjectID, l7AccessProfileIdParam, l7AccessProfileParam, overrideParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c L7AccessProfileClientContext) Update(l7AccessProfileIdParam string, l7AccessProfileParam model0.L7AccessProfile, overrideParam *bool) (model0.L7AccessProfile, error) {
	var err error
	var obj model0.L7AccessProfile

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.L7AccessProfilesClient)
		obj, err = client.Update(l7AccessProfileIdParam, l7AccessProfileParam, overrideParam)

	case utl.Multitenancy:
		client := c.Client.(client1.L7AccessProfilesClient)
		obj, err = client.Update(utl.DefaultOrgID, c.ProjectID, l7AccessProfileIdParam, l7AccessProfileParam, overrideParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}
