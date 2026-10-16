// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

// Adapters which let the LB resources and data sources share a single code path
// between the infra (/infra/lb-*) and the VPC (/orgs/.../vpcs/.../vpc-lb-*) APIs.
// The infra wrapper clients already satisfy the interfaces below; the VPC wrapper
// clients are adapted by binding the org/project/vpc parents.

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliVpcLbAppProfilesClient = vpcs.NewVpcLbAppProfilesClient
var cliVpcLbMonitorProfilesClient = vpcs.NewVpcLbMonitorProfilesClient
var cliVpcLbPersistenceProfilesClient = vpcs.NewVpcLbPersistenceProfilesClient
var cliVpcLbClientSslProfilesClient = vpcs.NewVpcLbClientSslProfilesClient
var cliVpcLbServerSslProfilesClient = vpcs.NewVpcLbServerSslProfilesClient
var cliVpcLbPoolsClient = vpcs.NewVpcLbPoolsClient
var cliVpcLbVirtualServersClient = vpcs.NewVpcLbVirtualServersClient
var cliVpcLbsClient = vpcs.NewVpcLbsClient

type lbAppProfileAPI interface {
	Get(id string) (*data.StructValue, error)
	Patch(id string, obj *data.StructValue) error
	Delete(id string, force *bool) error
	List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBAppProfileListResult, error)
}

type vpcLbAppProfileAdapter struct {
	client  *vpcs.VpcLbAppProfileClientContext
	parents []string
}

func (a vpcLbAppProfileAdapter) Get(id string) (*data.StructValue, error) {
	return a.client.Get(a.parents[0], a.parents[1], a.parents[2], id)
}

func (a vpcLbAppProfileAdapter) Patch(id string, obj *data.StructValue) error {
	return a.client.Patch(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbAppProfileAdapter) Delete(id string, force *bool) error {
	return a.client.Delete(a.parents[0], a.parents[1], a.parents[2], id, force)
}

func (a vpcLbAppProfileAdapter) List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBAppProfileListResult, error) {
	return a.client.List(a.parents[0], a.parents[1], a.parents[2], cursor, includeMarkForDelete, includedFields, pageSize, sortAscending, sortBy)
}

// getLbAppProfileAPI returns the VPC client for VPC contexts and the infra client otherwise.
// A nil interface is returned if the API is not supported for the context.
func getLbAppProfileAPI(sessionContext utl.SessionContext, connector client.Connector) lbAppProfileAPI {
	if sessionContext.ClientType == utl.VPC {
		c := cliVpcLbAppProfilesClient(sessionContext, connector)
		if c == nil {
			return nil
		}
		return vpcLbAppProfileAdapter{client: c, parents: getVpcParentsFromContext(sessionContext)}
	}
	c := cliLbAppProfilesClient(sessionContext, connector)
	if c == nil {
		return nil
	}
	return c
}

type lbMonitorProfileAPI interface {
	Get(id string) (*data.StructValue, error)
	Patch(id string, obj *data.StructValue) error
	Delete(id string, force *bool) error
	List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBMonitorProfileListResult, error)
}

type vpcLbMonitorProfileAdapter struct {
	client  *vpcs.VpcLbMonitorProfileClientContext
	parents []string
}

func (a vpcLbMonitorProfileAdapter) Get(id string) (*data.StructValue, error) {
	return a.client.Get(a.parents[0], a.parents[1], a.parents[2], id)
}

func (a vpcLbMonitorProfileAdapter) Patch(id string, obj *data.StructValue) error {
	return a.client.Patch(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbMonitorProfileAdapter) Delete(id string, force *bool) error {
	return a.client.Delete(a.parents[0], a.parents[1], a.parents[2], id, force)
}

func (a vpcLbMonitorProfileAdapter) List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBMonitorProfileListResult, error) {
	return a.client.List(a.parents[0], a.parents[1], a.parents[2], cursor, includeMarkForDelete, includedFields, pageSize, sortAscending, sortBy)
}

// getLbMonitorProfileAPI returns the VPC client for VPC contexts and the infra client otherwise.
// A nil interface is returned if the API is not supported for the context.
func getLbMonitorProfileAPI(sessionContext utl.SessionContext, connector client.Connector) lbMonitorProfileAPI {
	if sessionContext.ClientType == utl.VPC {
		c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
		if c == nil {
			return nil
		}
		return vpcLbMonitorProfileAdapter{client: c, parents: getVpcParentsFromContext(sessionContext)}
	}
	c := cliLbMonitorProfilesClient(sessionContext, connector)
	if c == nil {
		return nil
	}
	return c
}

type lbPersistenceProfileAPI interface {
	Get(id string) (*data.StructValue, error)
	Patch(id string, obj *data.StructValue) error
	Update(id string, obj *data.StructValue) (*data.StructValue, error)
	Delete(id string, force *bool) error
	List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBPersistenceProfileListResult, error)
}

type vpcLbPersistenceProfileAdapter struct {
	client  *vpcs.VpcLbPersistenceProfileClientContext
	parents []string
}

func (a vpcLbPersistenceProfileAdapter) Get(id string) (*data.StructValue, error) {
	return a.client.Get(a.parents[0], a.parents[1], a.parents[2], id)
}

func (a vpcLbPersistenceProfileAdapter) Patch(id string, obj *data.StructValue) error {
	return a.client.Patch(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbPersistenceProfileAdapter) Update(id string, obj *data.StructValue) (*data.StructValue, error) {
	return a.client.Update(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbPersistenceProfileAdapter) Delete(id string, force *bool) error {
	return a.client.Delete(a.parents[0], a.parents[1], a.parents[2], id, force)
}

func (a vpcLbPersistenceProfileAdapter) List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBPersistenceProfileListResult, error) {
	return a.client.List(a.parents[0], a.parents[1], a.parents[2], cursor, includeMarkForDelete, includedFields, pageSize, sortAscending, sortBy)
}

// getLbPersistenceProfileAPI returns the VPC client for VPC contexts and the infra client otherwise.
// A nil interface is returned if the API is not supported for the context.
func getLbPersistenceProfileAPI(sessionContext utl.SessionContext, connector client.Connector) lbPersistenceProfileAPI {
	if sessionContext.ClientType == utl.VPC {
		c := cliVpcLbPersistenceProfilesClient(sessionContext, connector)
		if c == nil {
			return nil
		}
		return vpcLbPersistenceProfileAdapter{client: c, parents: getVpcParentsFromContext(sessionContext)}
	}
	c := cliLbPersistenceProfilesClient(sessionContext, connector)
	if c == nil {
		return nil
	}
	return c
}

type lbClientSslProfileAPI interface {
	Get(id string) (model.LBClientSslProfile, error)
	Patch(id string, obj model.LBClientSslProfile) error
	Delete(id string, force *bool) error
	List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBClientSslProfileListResult, error)
}

type vpcLbClientSslProfileAdapter struct {
	client  *vpcs.VpcLbClientSslProfileClientContext
	parents []string
}

func (a vpcLbClientSslProfileAdapter) Get(id string) (model.LBClientSslProfile, error) {
	return a.client.Get(a.parents[0], a.parents[1], a.parents[2], id)
}

func (a vpcLbClientSslProfileAdapter) Patch(id string, obj model.LBClientSslProfile) error {
	return a.client.Patch(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbClientSslProfileAdapter) Delete(id string, force *bool) error {
	return a.client.Delete(a.parents[0], a.parents[1], a.parents[2], id, force)
}

func (a vpcLbClientSslProfileAdapter) List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBClientSslProfileListResult, error) {
	return a.client.List(a.parents[0], a.parents[1], a.parents[2], cursor, includeMarkForDelete, includedFields, pageSize, sortAscending, sortBy)
}

// getLbClientSslProfileAPI returns the VPC client for VPC contexts and the infra client otherwise.
// A nil interface is returned if the API is not supported for the context.
func getLbClientSslProfileAPI(sessionContext utl.SessionContext, connector client.Connector) lbClientSslProfileAPI {
	if sessionContext.ClientType == utl.VPC {
		c := cliVpcLbClientSslProfilesClient(sessionContext, connector)
		if c == nil {
			return nil
		}
		return vpcLbClientSslProfileAdapter{client: c, parents: getVpcParentsFromContext(sessionContext)}
	}
	c := cliLbClientSslProfilesClient(sessionContext, connector)
	if c == nil {
		return nil
	}
	return c
}

type lbServerSslProfileAPI interface {
	Get(id string) (model.LBServerSslProfile, error)
	Patch(id string, obj model.LBServerSslProfile) error
	Update(id string, obj model.LBServerSslProfile) (model.LBServerSslProfile, error)
	Delete(id string, force *bool) error
	List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBServerSslProfileListResult, error)
}

type vpcLbServerSslProfileAdapter struct {
	client  *vpcs.VpcLbServerSslProfileClientContext
	parents []string
}

func (a vpcLbServerSslProfileAdapter) Get(id string) (model.LBServerSslProfile, error) {
	return a.client.Get(a.parents[0], a.parents[1], a.parents[2], id)
}

func (a vpcLbServerSslProfileAdapter) Patch(id string, obj model.LBServerSslProfile) error {
	return a.client.Patch(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbServerSslProfileAdapter) Update(id string, obj model.LBServerSslProfile) (model.LBServerSslProfile, error) {
	return a.client.Update(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbServerSslProfileAdapter) Delete(id string, force *bool) error {
	return a.client.Delete(a.parents[0], a.parents[1], a.parents[2], id, force)
}

func (a vpcLbServerSslProfileAdapter) List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBServerSslProfileListResult, error) {
	return a.client.List(a.parents[0], a.parents[1], a.parents[2], cursor, includeMarkForDelete, includedFields, pageSize, sortAscending, sortBy)
}

// getLbServerSslProfileAPI returns the VPC client for VPC contexts and the infra client otherwise.
// A nil interface is returned if the API is not supported for the context.
func getLbServerSslProfileAPI(sessionContext utl.SessionContext, connector client.Connector) lbServerSslProfileAPI {
	if sessionContext.ClientType == utl.VPC {
		c := cliVpcLbServerSslProfilesClient(sessionContext, connector)
		if c == nil {
			return nil
		}
		return vpcLbServerSslProfileAdapter{client: c, parents: getVpcParentsFromContext(sessionContext)}
	}
	c := cliLbServerSslProfilesClient(sessionContext, connector)
	if c == nil {
		return nil
	}
	return c
}

type lbPoolAPI interface {
	Get(id string) (model.LBPool, error)
	Patch(id string, obj model.LBPool) error
	Update(id string, obj model.LBPool) (model.LBPool, error)
	Delete(id string, force *bool) error
	List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBPoolListResult, error)
}

type vpcLbPoolAdapter struct {
	client  *vpcs.VpcLbPoolClientContext
	parents []string
}

func (a vpcLbPoolAdapter) Get(id string) (model.LBPool, error) {
	return a.client.Get(a.parents[0], a.parents[1], a.parents[2], id)
}

func (a vpcLbPoolAdapter) Patch(id string, obj model.LBPool) error {
	return a.client.Patch(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbPoolAdapter) Update(id string, obj model.LBPool) (model.LBPool, error) {
	return a.client.Update(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbPoolAdapter) Delete(id string, force *bool) error {
	return a.client.Delete(a.parents[0], a.parents[1], a.parents[2], id, force)
}

func (a vpcLbPoolAdapter) List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBPoolListResult, error) {
	return a.client.List(a.parents[0], a.parents[1], a.parents[2], cursor, includeMarkForDelete, includedFields, pageSize, sortAscending, sortBy)
}

// getLbPoolAPI returns the VPC client for VPC contexts and the infra client otherwise.
// A nil interface is returned if the API is not supported for the context.
func getLbPoolAPI(sessionContext utl.SessionContext, connector client.Connector) lbPoolAPI {
	if sessionContext.ClientType == utl.VPC {
		c := cliVpcLbPoolsClient(sessionContext, connector)
		if c == nil {
			return nil
		}
		return vpcLbPoolAdapter{client: c, parents: getVpcParentsFromContext(sessionContext)}
	}
	c := cliLbPoolsClient(sessionContext, connector)
	if c == nil {
		return nil
	}
	return c
}

type lbVirtualServerAPI interface {
	Get(id string) (model.LBVirtualServer, error)
	Patch(id string, obj model.LBVirtualServer) error
	Update(id string, obj model.LBVirtualServer) (model.LBVirtualServer, error)
	Delete(id string, force *bool) error
	List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBVirtualServerListResult, error)
}

type vpcLbVirtualServerAdapter struct {
	client  *vpcs.VpcLbVirtualServerClientContext
	parents []string
}

func (a vpcLbVirtualServerAdapter) Get(id string) (model.LBVirtualServer, error) {
	return a.client.Get(a.parents[0], a.parents[1], a.parents[2], id)
}

func (a vpcLbVirtualServerAdapter) Patch(id string, obj model.LBVirtualServer) error {
	return a.client.Patch(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbVirtualServerAdapter) Update(id string, obj model.LBVirtualServer) (model.LBVirtualServer, error) {
	return a.client.Update(a.parents[0], a.parents[1], a.parents[2], id, obj)
}

func (a vpcLbVirtualServerAdapter) Delete(id string, force *bool) error {
	return a.client.Delete(a.parents[0], a.parents[1], a.parents[2], id, force)
}

func (a vpcLbVirtualServerAdapter) List(cursor *string, includeMarkForDelete *bool, includedFields *string, pageSize *int64, sortAscending *bool, sortBy *string) (model.LBVirtualServerListResult, error) {
	return a.client.List(a.parents[0], a.parents[1], a.parents[2], cursor, includeMarkForDelete, includedFields, pageSize, sortAscending, sortBy)
}

// getLbVirtualServerAPI returns the VPC client for VPC contexts and the infra client otherwise.
// A nil interface is returned if the API is not supported for the context.
func getLbVirtualServerAPI(sessionContext utl.SessionContext, connector client.Connector) lbVirtualServerAPI {
	if sessionContext.ClientType == utl.VPC {
		c := cliVpcLbVirtualServersClient(sessionContext, connector)
		if c == nil {
			return nil
		}
		return vpcLbVirtualServerAdapter{client: c, parents: getVpcParentsFromContext(sessionContext)}
	}
	c := cliLbVirtualServersClient(sessionContext, connector)
	if c == nil {
		return nil
	}
	return c
}

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
