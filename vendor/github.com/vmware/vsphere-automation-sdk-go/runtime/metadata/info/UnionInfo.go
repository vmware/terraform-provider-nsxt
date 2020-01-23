/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


type UnionInfo struct {
	tag   string
	cases []*UnionCase
}

func NewUnionInfo() *UnionInfo {
	cases := []*UnionCase{}
	unionInfo := UnionInfo{cases: cases}
	return &unionInfo
}

func (u *UnionInfo) Tag() string {
	return u.tag
}

func (u *UnionInfo) SetTag(tag string) {
	u.tag = tag
}

func (u *UnionInfo) Cases() []*UnionCase {
	return u.cases
}

func (u *UnionInfo) SetCases(cases []*UnionCase) {
	u.cases = cases
}

func (u *UnionInfo) AddToCase(ucase *UnionCase) {
	u.cases = append(u.cases, ucase)
}
