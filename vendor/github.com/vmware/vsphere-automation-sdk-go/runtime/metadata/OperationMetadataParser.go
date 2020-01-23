/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package metadata


// Note:
// Parser is incomplete
// Only valid to find resource/hasFieldOf metadata of structure info for given params info, Need to implement enumeration info

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/metadata/info"
	"reflect"
	"strings"
)

type OperationMetadataParser struct {
	operationInputName string
	operationInfo      *info.OperationInfo
	structureInfoMap   map[string]*info.StructureInfo
	paramPrivilegeInfo map[string][]*info.ParamPrivilegeInfo
}

func NewOperationMetadataParser(operationInputName string, operationInfo *info.OperationInfo, structureInfoMap map[string]*info.StructureInfo) *OperationMetadataParser {
	paramPrivilegeInfo := make(map[string][]*info.ParamPrivilegeInfo)
	OpMetaParser := OperationMetadataParser{
		operationInputName: operationInputName,
		operationInfo:      operationInfo,
		structureInfoMap:   structureInfoMap,
		paramPrivilegeInfo: paramPrivilegeInfo,
	}
	return &OpMetaParser
}

func (opMetaParser *OperationMetadataParser) GeneratePrivilegeInfo() map[string][]*info.ParamPrivilegeInfo {
	for _, priviInfo := range opMetaParser.operationInfo.PrivilegeInfo() {
		opMetaParser.generatePrivInfoForPath(priviInfo.PropertyPath(), priviInfo.Privileges())
	}
	return opMetaParser.paramPrivilegeInfo
}

func (opMetaParser *OperationMetadataParser) generatePrivInfoForPath(propertyPath string, privileges []string) {
	var propertyPathTokenizer []string = strings.Split(propertyPath, ".")
	var fieldInfos map[string]*info.FieldInfo = opMetaParser.operationInfo.FieldInfoMap()

	var currentTypeName string = opMetaParser.operationInputName
	for i, currentToken := range propertyPathTokenizer {
		if (i + 1) == len(propertyPathTokenizer) {
			// it is the last token - add ParamPrivilegeInfo to the overall result
			opMetaParser.processFieldInfo(currentTypeName, currentToken, fieldInfos, privileges)

		} else {

			// check if it is of type MAP
			if strings.Contains(currentToken, "#") {
				if strings.Contains(currentToken, "#key") {
					log.Error("'#key' is only supported at the end, detected: ", propertyPath)
				}
				currentToken = currentToken[:strings.Index(currentToken, "#")]
			}
			var fieldInfo *info.FieldInfo = fieldInfos[currentToken]
			var structinfo *info.StructureInfo = opMetaParser.processFieldType(fieldInfo, fieldInfo.TypeInfo())

			if structinfo != nil {
				fieldInfos = structinfo.FieldInfoMap()
				currentTypeName = structinfo.Identifier()
			}
		}
	}
}

func (opMetaParser *OperationMetadataParser) processFieldType(fieldInfo *info.FieldInfo, fieldType info.TypeInfo) *info.StructureInfo {

	switch reflect.TypeOf(fieldType) {
	case reflect.TypeOf(info.NewGenericTypeInfo("")):
		fieldTypeGeneric, _ := fieldType.(*info.GenericTypeInfo)
		if fieldTypeGeneric.IsMap() {
			return opMetaParser.processFieldType(fieldInfo, fieldTypeGeneric.MapValueType())
		} else {
			return opMetaParser.processFieldType(fieldInfo, fieldTypeGeneric.ElementType())
		}
	case reflect.TypeOf(info.NewUserDefinedTypeInfo("", nil)):
		return fieldType.(*info.UserDefinedTypeInfo).StructureInfo()

	case reflect.TypeOf(info.NewPrimitiveTypeInfo("", info.PrimitiveTypeFromString(""))):
		if fieldType.(*info.PrimitiveTypeInfo).PrimitiveType() == info.PrimitiveTypeFromString("DYNAMIC_STRUCTURE") || fieldType.(*info.PrimitiveTypeInfo).PrimitiveType() == info.PrimitiveTypeFromString("STRUCTURE") {
			hasFieldsOfTypeNames := fieldInfo.HasFieldsOfTypeNames()
			if len(hasFieldsOfTypeNames) == 0 {
				log.Error("No values for hasFieldsOfTypeNames detected")
				return nil
			} else if len(hasFieldsOfTypeNames) > 1 {
				log.Error("Mulitple values for hasFieldsOfTypeNames detected")
				return nil
			}
			// TODO: Add support to handle mulitple values of hasFieldsOf, for now send the value at index 0
			return opMetaParser.structureInfoMap[hasFieldsOfTypeNames[0]]
		}
	default:
		log.Error("processFieldType unknown fieldType info.TypeInfo provided : ", reflect.TypeOf(fieldType))
	}
	return nil
}

func (opMetaParser *OperationMetadataParser) processFieldInfo(typeName string, fieldName string, fieldInfos map[string]*info.FieldInfo, privileges []string) {
	var paramInfo *info.ParamPrivilegeInfo
	var fieldname string = fieldName
	if strings.Contains(fieldName, "#") {
		fieldname = fieldName[:strings.Index(fieldName, "#")]
	}
	var fieldInfo *info.FieldInfo = fieldInfos[fieldname]

	if strings.Contains(fieldName, "#") {
		typeName = typeName + "-" + fieldName
		indicator := fieldName[strings.Index(fieldName, "#")+1:]
		paramInfo = info.NewParamPrivilegeInfo(indicator, fieldInfo.IdentifierType(), fieldInfo.IdentifierTypeHolder(), privileges)
	} else {
		paramInfo = info.NewParamPrivilegeInfo(fieldName, fieldInfo.IdentifierType(), fieldInfo.IdentifierTypeHolder(), privileges)
	}
	log.Debugf("typeName: %s", typeName)

	if _, ok := opMetaParser.paramPrivilegeInfo[typeName]; !ok {
		opMetaParser.paramPrivilegeInfo[typeName] = make([]*info.ParamPrivilegeInfo, 0)
	}

	opMetaParser.paramPrivilegeInfo[typeName] = append(opMetaParser.paramPrivilegeInfo[typeName], paramInfo)
}
