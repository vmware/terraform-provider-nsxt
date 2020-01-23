/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info

import "github.com/vmware/vsphere-automation-sdk-go/runtime/data"

type OperationInfo struct {
	privileges          []string
	serviceIdentifier   string
	operationIdentifier string
	parameterPrivileges map[string][]string

	errors      []*data.ErrorDefinition
	routingInfo *RoutingInfo

	// privilegeInfo holds info of propertypath & privileges
	// example:
	// {
	// 		.
	// 		.
	// 		"package.service" : {
	// 			operation.params#key : [p1, p2 , p3]
	// 			.
	// 			.
	// 		}
	//	}
	// PropertyPath: params#key
	// Privileges: [p1, p2 , p3]
	privilegeInfo []*PrivilegeInfo

	parameters map[string]*FieldInfo
	unions     map[string]*UnionInfo

	documentation        string
	returnValue          *FieldInfo
	authenticationScheme []*AuthenticationScheme
}

func NewOperationInfo() *OperationInfo {
	return &OperationInfo{}
}

func (op *OperationInfo) Identifier() string {
	return op.serviceIdentifier
}

func (op *OperationInfo) SetIdentifier(id string) {
	op.serviceIdentifier = id
}

func (op *OperationInfo) Privileges() []string {
	return op.privileges
}

func (op *OperationInfo) PrivilegeInfo() []*PrivilegeInfo {
	return op.privilegeInfo
}

func (op *OperationInfo) SetPrivileges(privileges []string) {
	op.privileges = privileges
}

func (op *OperationInfo) AddPrivilege(privilege string) {
	op.privileges = append(op.privileges, privilege)
}

func (op *OperationInfo) AddPrivilegeInfo(pInfo *PrivilegeInfo) {
	op.privilegeInfo = append(op.privilegeInfo, pInfo)
}

// parameters Fieldinfo
func (op *OperationInfo) FieldInfo(field string) *FieldInfo {
	return op.parameters[field]
}

func (op *OperationInfo) FieldInfoMap() map[string]*FieldInfo {
	return op.parameters
}

func (op *OperationInfo) SetFieldInfo(field string, fieldInfo *FieldInfo) {
	if op.parameters == nil {
		op.parameters = make(map[string]*FieldInfo)
	}
	op.parameters[field] = fieldInfo
}

func (op *OperationInfo) ListFieldInfo() []string {
	keys := []string{}
	for k, _ := range op.parameters {
		keys = append(keys, k)
	}
	return keys
}
