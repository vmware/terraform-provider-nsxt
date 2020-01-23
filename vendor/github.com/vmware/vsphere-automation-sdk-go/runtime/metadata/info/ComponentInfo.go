/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info

type ComponentInfo struct {
	packageInfo   map[string]*PackageInfo
	identifier    string
	documentation string
	checksum      string
}

func NewComponentInfo() *ComponentInfo {
	pkgInfo := make(map[string]*PackageInfo)
	compInfo := ComponentInfo{packageInfo: pkgInfo}
	return &compInfo
}

func (comp *ComponentInfo) ListPackageInfo() []string {
	keys := []string{}
	for k := range comp.packageInfo {
		keys = append(keys, k)
	}
	return keys
}

func (comp *ComponentInfo) PackageInfo(packageName string) *PackageInfo {
	return comp.packageInfo[packageName]
}

func (comp *ComponentInfo) PackageMap() map[string]*PackageInfo {
	return comp.packageInfo
}

func (comp *ComponentInfo) SetPackageInfo(packageName string, pkg *PackageInfo) {
	comp.packageInfo[packageName] = pkg
}

func (comp *ComponentInfo) Identifier() string {
	return comp.identifier
}

func (comp *ComponentInfo) SetIdentifier(identifier string) {
	comp.identifier = identifier
}

func (comp *ComponentInfo) Documentation() string {
	return comp.documentation
}

func (comp *ComponentInfo) SetDocumentation(documentation string) {
	comp.documentation = documentation
}

func (comp *ComponentInfo) Checksum() string {
	return comp.checksum
}

func (comp *ComponentInfo) SetChecksum(checksum string) {
	comp.checksum = checksum
}
