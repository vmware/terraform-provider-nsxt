/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

// User Identity class represents result for Authentication
// Handler authenticate method.
type UserIdentity struct {
	userName string
	domain   *string
	groups   []string
}

func NewUserIdentity(userName string) *UserIdentity {
	return &UserIdentity{userName: userName}
}

func (u *UserIdentity) Groups() []string {
	return u.groups
}

func (u *UserIdentity) SetGroups(groups []string) {
	u.groups = groups
}

func (u *UserIdentity) SetDomain(domain *string) {
	u.domain = domain
}

func (u *UserIdentity) UserName() string {
	return u.userName
}

func (u *UserIdentity) Domain() *string {
	return u.domain
}
