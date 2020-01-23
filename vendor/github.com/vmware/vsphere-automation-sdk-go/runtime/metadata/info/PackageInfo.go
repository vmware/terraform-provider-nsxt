/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info

type PackageInfo struct {
	identifier           string
	documentation        string
	privileges           []string
	authenticationScheme []*AuthenticationScheme

	enumerationInfo map[string]*EnumerationInfo
	routingInfo     map[string]*RoutingInfo
	serviceInfo     map[string]*ServiceInfo
	structureInfo   map[string]*StructureInfo
}

func NewPackageInfo() *PackageInfo {
	return &PackageInfo{}
}

// Enumeration Info
func (pkg *PackageInfo) ListEnumerationInfo() []string {
	keys := []string{}
	for k := range pkg.enumerationInfo {
		keys = append(keys, k)
	}
	return keys
}

func (pkg *PackageInfo) EnumerationInfo(enum string) *EnumerationInfo {
	return pkg.enumerationInfo[enum]
}

func (pkg *PackageInfo) SetEnumerationInfo(enum string, emumInfo *EnumerationInfo) {
	if pkg.enumerationInfo == nil {
		pkg.enumerationInfo = make(map[string]*EnumerationInfo)
	}
	pkg.enumerationInfo[enum] = emumInfo
}

// Routing Info
func (pkg *PackageInfo) ListRoutingInfo() []string {
	keys := []string{}
	for k := range pkg.routingInfo {
		keys = append(keys, k)
	}
	return keys
}

func (pkg *PackageInfo) RoutingInfo(id string) *RoutingInfo {
	return pkg.routingInfo[id]
}

func (pkg *PackageInfo) SetRoutingInfo(id string, rtgInfo *RoutingInfo) {
	pkg.routingInfo[id] = rtgInfo
}

// Service Info
func (pkg *PackageInfo) ListServiceInfo() []string {
	keys := []string{}
	for k := range pkg.serviceInfo {
		keys = append(keys, k)
	}
	return keys
}

func (pkg *PackageInfo) ServiceInfoMap() map[string]*ServiceInfo {
	return pkg.serviceInfo
}

func (pkg *PackageInfo) ServiceInfo(service string) *ServiceInfo {
	return pkg.serviceInfo[service]
}

func (pkg *PackageInfo) SetServiceInfo(service string, ser *ServiceInfo) {
	if pkg.serviceInfo == nil {
		pkg.serviceInfo = make(map[string]*ServiceInfo)
	}
	pkg.serviceInfo[service] = ser
}

// StructureInfo
func (pkg *PackageInfo) ListStructureInfo() []string {
	keys := []string{}
	for k := range pkg.structureInfo {
		keys = append(keys, k)
	}
	return keys
}

func (pkg *PackageInfo) StructureInfoMap() map[string]*StructureInfo {
	return pkg.structureInfo
}

func (pkg *PackageInfo) StructureInfo(str string) *StructureInfo {
	return pkg.structureInfo[str]
}

func (pkg *PackageInfo) SetStructureInfo(str string, strInfo *StructureInfo) {
	if pkg.structureInfo == nil {
		pkg.structureInfo = make(map[string]*StructureInfo)
	}
	pkg.structureInfo[str] = strInfo
}

// Privileges
func (pkg *PackageInfo) Privileges() []string {
	return pkg.privileges
}

func (pkg *PackageInfo) SetPrivileges(privileges []string) {
	pkg.privileges = privileges
}

func (pkg *PackageInfo) AddPrivilege(privilege string) {
	if pkg.privileges != nil {
		pkg.privileges = append(pkg.privileges, privilege)
	} else {
		pkg.privileges = []string{privilege}
	}
}

// Documentation
func (pkg *PackageInfo) Documentation() string {
	return pkg.documentation
}

func (pkg *PackageInfo) SetDocumentation(documentation string) {
	pkg.documentation = documentation
}

// Identifier
func (pkg *PackageInfo) Identifier() string {
	return pkg.identifier
}

func (pkg *PackageInfo) SetIdentifier(identifier string) {
	pkg.identifier = identifier
}

// Authentication Scheme
func (pkg *PackageInfo) AuthenticationScheme() []*AuthenticationScheme {
	return pkg.authenticationScheme
}

func (pkg *PackageInfo) SetAuthenticationScheme(authScheme []*AuthenticationScheme) {
	pkg.authenticationScheme = authScheme
}

func (pkg *PackageInfo) AddToAuthenticationScheme(elem *AuthenticationScheme) {
	pkg.authenticationScheme = append(pkg.authenticationScheme, elem)
}
