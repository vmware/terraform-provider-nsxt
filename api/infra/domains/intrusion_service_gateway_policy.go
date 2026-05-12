//nolint:revive
package domains

import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type IntrusionServiceGatewayPolicyClientContext utl.ClientContext

func NewIntrusionServiceGatewayPoliciesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *IntrusionServiceGatewayPolicyClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewIntrusionServiceGatewayPoliciesClient(connector)

	default:
		return nil
	}
	return &IntrusionServiceGatewayPolicyClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c IntrusionServiceGatewayPolicyClientContext) Get(domainIdParam string, policyIdParam string) (model0.IdsGatewayPolicy, error) {
	var obj model0.IdsGatewayPolicy
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.IntrusionServiceGatewayPoliciesClient)
		obj, err = client.Get(domainIdParam, policyIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c IntrusionServiceGatewayPolicyClientContext) Delete(domainIdParam string, policyIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.IntrusionServiceGatewayPoliciesClient)
		err = client.Delete(domainIdParam, policyIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c IntrusionServiceGatewayPolicyClientContext) List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includeRuleCountParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.IdsGatewayPolicyListResult, error) {
	var err error
	var obj model0.IdsGatewayPolicyListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.IntrusionServiceGatewayPoliciesClient)
		obj, err = client.List(domainIdParam, cursorParam, includeMarkForDeleteObjectsParam, includeRuleCountParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c IntrusionServiceGatewayPolicyClientContext) Patch(domainIdParam string, policyIdParam string, idsGatewayPolicyParam model0.IdsGatewayPolicy) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.IntrusionServiceGatewayPoliciesClient)
		err = client.Patch(domainIdParam, policyIdParam, idsGatewayPolicyParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c IntrusionServiceGatewayPolicyClientContext) Update(domainIdParam string, policyIdParam string, idsGatewayPolicyParam model0.IdsGatewayPolicy) (model0.IdsGatewayPolicy, error) {
	var err error
	var obj model0.IdsGatewayPolicy

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.IntrusionServiceGatewayPoliciesClient)
		obj, err = client.Update(domainIdParam, policyIdParam, idsGatewayPolicyParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}
