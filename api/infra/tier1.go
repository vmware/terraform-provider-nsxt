//nolint:revive
package infra

// The following file has been autogenerated. Please avoid any changes!
import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	model1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client2 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type Tier1ClientContext utl.ClientContext

func NewTier1sClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *Tier1ClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewTier1sClient(connector)

	case utl.Global:
		client = client1.NewTier1sClient(connector)

	case utl.Multitenancy:
		client = client2.NewTier1sClient(connector)

	default:
		return nil
	}
	return &Tier1ClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c Tier1ClientContext) Get(tier1IdParam string) (model0.Tier1, error) {
	var obj model0.Tier1
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.Tier1sClient)
		obj, err = client.Get(tier1IdParam)
		if err != nil {
			return obj, err
		}

	case utl.Global:
		client := c.Client.(client1.Tier1sClient)
		gmObj, err1 := client.Get(tier1IdParam)
		if err1 != nil {
			return obj, err1
		}
		var rawObj interface{}
		rawObj, err = utl.ConvertModelBindingType(gmObj, model1.Tier1BindingType(), model0.Tier1BindingType())
		obj = rawObj.(model0.Tier1)

	case utl.Multitenancy:
		client := c.Client.(client2.Tier1sClient)
		obj, err = client.Get(utl.DefaultOrgID, c.ProjectID, tier1IdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c Tier1ClientContext) Patch(tier1IdParam string, tier1Param model0.Tier1) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.Tier1sClient)
		err = client.Patch(tier1IdParam, tier1Param)

	case utl.Global:
		client := c.Client.(client1.Tier1sClient)
		gmObj, err1 := utl.ConvertModelBindingType(tier1Param, model0.Tier1BindingType(), model1.Tier1BindingType())
		if err1 != nil {
			return err1
		}
		err = client.Patch(tier1IdParam, gmObj.(model1.Tier1))

	case utl.Multitenancy:
		client := c.Client.(client2.Tier1sClient)
		err = client.Patch(utl.DefaultOrgID, c.ProjectID, tier1IdParam, tier1Param)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c Tier1ClientContext) Update(tier1IdParam string, tier1Param model0.Tier1) (model0.Tier1, error) {
	var err error
	var obj model0.Tier1

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.Tier1sClient)
		obj, err = client.Update(tier1IdParam, tier1Param)

	case utl.Global:
		client := c.Client.(client1.Tier1sClient)
		gmObj, err := utl.ConvertModelBindingType(tier1Param, model0.Tier1BindingType(), model1.Tier1BindingType())
		if err != nil {
			return obj, err
		}
		gmObj, err = client.Update(tier1IdParam, gmObj.(model1.Tier1))
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.Tier1BindingType(), model0.Tier1BindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.Tier1)

	case utl.Multitenancy:
		client := c.Client.(client2.Tier1sClient)
		obj, err = client.Update(utl.DefaultOrgID, c.ProjectID, tier1IdParam, tier1Param)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c Tier1ClientContext) Delete(tier1IdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.Tier1sClient)
		err = client.Delete(tier1IdParam)

	case utl.Global:
		client := c.Client.(client1.Tier1sClient)
		err = client.Delete(tier1IdParam)

	case utl.Multitenancy:
		client := c.Client.(client2.Tier1sClient)
		err = client.Delete(utl.DefaultOrgID, c.ProjectID, tier1IdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c Tier1ClientContext) List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.Tier1ListResult, error) {
	var err error
	var obj model0.Tier1ListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.Tier1sClient)
		obj, err = client.List(cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	case utl.Global:
		client := c.Client.(client1.Tier1sClient)
		gmObj, err := client.List(cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.Tier1ListResultBindingType(), model0.Tier1ListResultBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.Tier1ListResult)

	case utl.Multitenancy:
		client := c.Client.(client2.Tier1sClient)
		obj, err = client.List(utl.DefaultOrgID, c.ProjectID, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}