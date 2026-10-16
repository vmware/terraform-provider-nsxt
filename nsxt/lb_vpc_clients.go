// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

// Adapters which let the LB resources and data sources share a single code path
// between the infra (/infra/lb-*) and the VPC (/orgs/.../vpcs/.../vpc-lb-*) APIs.
// The infra wrapper clients already satisfy the interfaces below; the VPC wrapper
// clients are adapted by binding the org/project/vpc parents.

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliVpcLbsClient = vpcs.NewVpcLbsClient

type lbServiceAPI interface {
	Get(id string) (model.LBService, error)
	Patch(id string, obj model.LBService) error
	Delete(id string, force *bool) error
}

type vpcLbServiceAdapter struct {
	client  *vpcs.VpcLbClientContext
	parents []string
}

func (a vpcLbServiceAdapter) Get(id string) (model.LBService, error) {
	return a.client.Get(a.parents[0], a.parents[1], a.parents[2], id)
}

func (a vpcLbServiceAdapter) Patch(id string, obj model.LBService) error {
	return a.client.Patch(a.parents[0], a.parents[1], a.parents[2], id, obj, nil)
}

func (a vpcLbServiceAdapter) Delete(id string, force *bool) error {
	return a.client.Delete(a.parents[0], a.parents[1], a.parents[2], id, force)
}

// getLbServiceAPI returns the VPC client for VPC contexts and the infra client otherwise.
// A nil interface is returned if the API is not supported for the context.
func getLbServiceAPI(sessionContext utl.SessionContext, connector client.Connector) lbServiceAPI {
	if sessionContext.ClientType == utl.VPC {
		c := cliVpcLbsClient(sessionContext, connector)
		if c == nil {
			return nil
		}
		return vpcLbServiceAdapter{client: c, parents: getVpcParentsFromContext(sessionContext)}
	}
	c := cliLbServicesClient(sessionContext, connector)
	if c == nil {
		return nil
	}
	return c
}
