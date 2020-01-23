/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


type PrivilegeInfo struct {
	propertyPath string
	privileges   []string
}

func NewPrivilegeInfo() *PrivilegeInfo {
	return &PrivilegeInfo{}
}

func (pr *PrivilegeInfo) PropertyPath() string {
	return pr.propertyPath
}

func (pr *PrivilegeInfo) SetPropertyPath(pp string) {
	pr.propertyPath = pp
}

func (pr *PrivilegeInfo) Privileges() []string {
	return pr.privileges
}

func (pr *PrivilegeInfo) SetPrivileges(privileges []string) {
	pr.privileges = privileges
}

func (pr *PrivilegeInfo) AddPrivilege(privilege string) {
	if pr.privileges == nil {
		pr.privileges = []string{privilege}
	} else {
		pr.privileges = append(pr.privileges, privilege)
	}
}
