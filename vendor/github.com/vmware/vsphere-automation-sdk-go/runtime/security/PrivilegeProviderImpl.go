/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/metadata"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/metadata/info"
	"reflect"
	"strings"
)

type PrivilegeProviderImpl struct {
	apiMetadata   *metadata.ApiMetadata
	components    []string
	operationInfo map[string]*info.OperationInfo
	packageInfo   map[string]*info.PackageInfo
	structureInfo map[string]*info.StructureInfo
}

func NewPrivilegeProviderImpl(apiMetadata *metadata.ApiMetadata) *PrivilegeProviderImpl {
	operationInfo := make(map[string]*info.OperationInfo)
	packageInfo := make(map[string]*info.PackageInfo)
	structureInfo := make(map[string]*info.StructureInfo)
	privilegeProviderImpl := PrivilegeProviderImpl{
		apiMetadata:   apiMetadata,
		operationInfo: operationInfo,
		structureInfo: structureInfo,
		packageInfo:   packageInfo,
	}
	privilegeProviderImpl.LoadMetadata()
	return &privilegeProviderImpl
}

func (privilegeProv *PrivilegeProviderImpl) LoadMetadata() {
	privilegeProv.processComponentMetadata(privilegeProv.apiMetadata.ComponentInfo())
}

func (privilegeProv *PrivilegeProviderImpl) processComponentMetadata(componentInfoResult map[string]*info.ComponentInfo) {
	// iterate through all the packages
	for _, compInfo := range componentInfoResult {
		// Pass throught packages and process structures and services
		for pkgName, pkgInfo := range compInfo.PackageMap() {
			privilegeProv.processStructureMetadata(pkgInfo.StructureInfoMap())
			privilegeProv.processServiceMetadata(pkgInfo.ServiceInfoMap())
			privilegeProv.packageInfo[pkgName] = pkgInfo
		}
	}
}

func (privilegeProv *PrivilegeProviderImpl) processServiceMetadata(serviceInfoResult map[string]*info.ServiceInfo) {
	for serName, serInfo := range serviceInfoResult {
		privilegeProv.processStructureMetadata(serInfo.StructureInfoMap())
		privilegeProv.processOperationMetadata(serInfo.OperationInfoMap(), serName)
	}
}

func (privilegeProv *PrivilegeProviderImpl) processStructureMetadata(structureInfoResult map[string]*info.StructureInfo) {
	for structName, structInfo := range structureInfoResult {
		privilegeProv.structureInfo[structName] = structInfo
	}
}

func (privilegeProv *PrivilegeProviderImpl) processOperationMetadata(operationInfoResult map[string]*info.OperationInfo, serviceName string) {
	for opName, opInfo := range operationInfoResult {
		fullyQualifiedOperName := serviceName + "." + opName
		privilegeProv.operationInfo[fullyQualifiedOperName] = opInfo
	}
}

func (privilegeProv *PrivilegeProviderImpl) GetPrivilegeInfo(fullyQualifiedOperName string, inputValue data.DataValue) (map[ResourceIdentifier][]string, error) {
	var operInfo *info.OperationInfo
	privilegesMap := make(map[ResourceIdentifier][]string)

	operInfo, ok := privilegeProv.operationInfo[fullyQualifiedOperName]
	if !ok {
		log.Error("Given fullyQualifiedOperName doesnot exists i.e. " + fullyQualifiedOperName)
		args := map[string]string{
			"msg": "given fullyQualifiedOperName doesnot exists " + fullyQualifiedOperName,
		}
		return nil, l10n.NewRuntimeError("vapi.security.authorization.exception", args)
	}

	if _, ok := inputValue.(*data.VoidValue); ok || inputValue == nil {
		log.Debugf("No Operation Input provided")
	} else {
		structInputValue, ok := inputValue.(*data.StructValue)
		if !ok {
			args := map[string]string{
				"msg": "Input Data value is not of type StructValue",
			}
			return nil, l10n.NewRuntimeError("vapi.security.authorization.exception", args)
		}

		OpMetaParser := metadata.NewOperationMetadataParser(structInputValue.Name(), operInfo, privilegeProv.structureInfo)

		// Add Privileges i.e we acquire from inputValue
		// paramPrivilegeInfo is of type map[string][]*info.ParamPrivilegeInfo
		paramPrivilegeInfo := OpMetaParser.GeneratePrivilegeInfo()

		visitDataValue(privilegesMap, inputValue, paramPrivilegeInfo)
	}

	// Add Privileges required by the operation itself
	operid := NewResourceIdentifier(true, fullyQualifiedOperName, "")
	privilegesMap[operid] = operInfo.Privileges()

	// Add Privileges required on package level
	err := privilegeProv.addPackagePrivileges(privilegesMap, fullyQualifiedOperName)
	if err != nil {
		return privilegesMap, err
	}

	return privilegesMap, nil
}

func populatePrivilegesForId(accprivileges map[ResourceIdentifier][]string, input data.DataValue, resourceType string, privileges []string) {
	switch reflect.TypeOf(input) {

	case data.StringValuePtr:
		strVal := input.(*data.StringValue)

		id := NewResourceIdentifier(false, strVal.Value(), resourceType)
		accprivileges[id] = privileges

	case data.ListValuePtr:
		listVal := input.(*data.ListValue)
		for _, element := range listVal.List() {
			populatePrivilegesForId(accprivileges, element, resourceType, privileges)
		}

	case data.OptionalValuePtr:
		optVal := input.(*data.OptionalValue)
		if optVal.IsSet() {
			populatePrivilegesForId(accprivileges, optVal.Value(), resourceType, privileges)
		}
	}

}

func visitElementValue(privileges map[ResourceIdentifier][]string, input data.DataValue, paramPrivilegeInfo *info.ParamPrivilegeInfo, containerStructure *data.StructValue) {
	resourceType := paramPrivilegeInfo.ResourceType()
	populatePrivilegesForId(privileges, input, resourceType, paramPrivilegeInfo.Privileges())
}

func visitDataValue(privileges map[ResourceIdentifier][]string, input data.DataValue, paramPrivilegeInfo map[string][]*info.ParamPrivilegeInfo) {
	switch reflect.TypeOf(input) {

	case data.StructValuePtr:
		structInput := input.(*data.StructValue)
		lookUpKey := structInput.Name()

		var idFieldsMap map[string]*info.ParamPrivilegeInfo = make(map[string]*info.ParamPrivilegeInfo)
		if _, ok := paramPrivilegeInfo[lookUpKey]; ok {
			paramInfos := paramPrivilegeInfo[lookUpKey]
			for _, paramInfo := range paramInfos {
				idFieldsMap[paramInfo.ParamName()] = paramInfo
			}
		}

		for key, dv := range structInput.Fields() {
			var fieldValue data.DataValue = dv
			if _, ok := idFieldsMap[key]; ok {
				visitElementValue(privileges, fieldValue, idFieldsMap[key], structInput)
			} else {
				visitDataValue(privileges, fieldValue, paramPrivilegeInfo)
			}
		}

	case data.ListValuePtr:
		listVal := input.(*data.ListValue)
		for _, element := range listVal.List() {
			visitDataValue(privileges, element, paramPrivilegeInfo)
		}

	case data.OptionalValuePtr:
		optVal := input.(*data.OptionalValue)
		if optVal.IsSet() {
			visitDataValue(privileges, optVal.Value(), paramPrivilegeInfo)
		}
	}
}

func (privilegeProv *PrivilegeProviderImpl) addPackagePrivileges(privilegesMap map[ResourceIdentifier][]string, fullyQualifiedOperName string) error {
	// Package Name
	tokens := strings.Split(fullyQualifiedOperName, ".")
	if !(len(tokens) > 2) {
		args := map[string]string{
			"msg": "Package level privilege extraction failed due to invalid fully qualified operation name",
		}
		return l10n.NewRuntimeError("vapi.security.authorization.exception", args)
	}
	pkgName := strings.Join(tokens[:len(tokens)-2], ".")

	if pkgInfo, ok := privilegeProv.packageInfo[pkgName]; ok {
		pkgid := NewResourceIdentifier(false, pkgName, "")
		privilegesMap[pkgid] = pkgInfo.Privileges()
	}
	return nil
}
