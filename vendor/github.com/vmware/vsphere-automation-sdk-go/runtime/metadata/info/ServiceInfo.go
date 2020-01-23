/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


type ServiceInfo struct {
	identifier           string
	documentation        string
	privileges           []string
	authenticationScheme []*AuthenticationScheme

	enumerationInfo map[string]*EnumerationInfo
	operationInfo   map[string]*OperationInfo
	routingInfo     map[string]*RoutingInfo
	structureInfo   map[string]*StructureInfo
}

func NewServiceInfo() *ServiceInfo {
	enumInfo := make(map[string]*EnumerationInfo)
	operInfo := make(map[string]*OperationInfo)
	routInfo := make(map[string]*RoutingInfo)
	strucInfo := make(map[string]*StructureInfo)
	return &ServiceInfo{
		operationInfo:   operInfo,
		enumerationInfo: enumInfo,
		routingInfo:     routInfo,
		structureInfo:   strucInfo,
	}
}

// OperationInfo
func (ser *ServiceInfo) OperationInfo(operation string) *OperationInfo {
	return ser.operationInfo[operation]
}

func (ser *ServiceInfo) SetOperationInfo(operation string, operInfo *OperationInfo) {
	ser.operationInfo[operation] = operInfo
}

func (ser *ServiceInfo) OperationInfoMap() map[string]*OperationInfo {
	return ser.operationInfo
}

func (comp *ServiceInfo) ListOperationInfo() []string {
	keys := []string{}
	for k, _ := range comp.operationInfo {
		keys = append(keys, k)
	}
	return keys
}

// StructureInfo
func (ser *ServiceInfo) StructureInfo(structure string) *StructureInfo {
	return ser.structureInfo[structure]
}

func (ser *ServiceInfo) SetStructureInfo(structure string, structureInfo *StructureInfo) {
	ser.structureInfo[structure] = structureInfo
}

func (ser *ServiceInfo) StructureInfoMap() map[string]*StructureInfo {
	return ser.structureInfo
}

func (ser *ServiceInfo) ListStructureInfo() []string {
	keys := []string{}
	for k, _ := range ser.structureInfo {
		keys = append(keys, k)
	}
	return keys
}

// Identifier
func (ser *ServiceInfo) Identifier() string {
	return ser.identifier
}

func (ser *ServiceInfo) SetIdentifier(identifier string) {
	ser.identifier = identifier
}
