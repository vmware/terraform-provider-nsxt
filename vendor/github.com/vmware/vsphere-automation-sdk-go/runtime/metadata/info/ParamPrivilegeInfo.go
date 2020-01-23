/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


type ParamPrivilegeInfo struct {
	paramName          string
	resourceType       string
	resourceTypeHolder string
	privileges         []string
}

func NewParamPrivilegeInfo(paramName string, resourceType string, resourceTypeHolder string, privileges []string) *ParamPrivilegeInfo {
	return &ParamPrivilegeInfo{paramName: paramName, resourceType: resourceType, resourceTypeHolder: resourceTypeHolder, privileges: privileges}
}

func (pp *ParamPrivilegeInfo) Privileges() []string {
	return pp.privileges
}

func (pp *ParamPrivilegeInfo) ParamName() string {
	return pp.paramName
}

func (pp *ParamPrivilegeInfo) ResourceType() string {
	return pp.resourceType
}

func (pp *ParamPrivilegeInfo) ResourceTypeHolder() string {
	return pp.resourceTypeHolder
}
