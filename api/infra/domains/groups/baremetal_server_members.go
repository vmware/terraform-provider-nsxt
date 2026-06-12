//nolint:revive
package groups

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains/groups/members"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type BmsMembersClientContext utl.ClientContext

func NewBmsMembersClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *BmsMembersClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewBmsClient(connector)

	default:
		return nil
	}
	return &BmsMembersClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c BmsMembersClientContext) List(domainIdParam string, groupIdParam string, cursorParam *string, enforcementPointPathParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.BareMetalServerListResult, error) {
	var obj model0.BareMetalServerListResult
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.BmsClient)
		obj, err = client.List(domainIdParam, groupIdParam, cursorParam, enforcementPointPathParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}
