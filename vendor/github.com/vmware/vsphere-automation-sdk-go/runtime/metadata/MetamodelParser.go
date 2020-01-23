/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package metadata

// Note:
// Parser is incomplete
// Only valid to find resource metadata for given params info

import (
	"encoding/json"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/metadata/info"
	"reflect"
	"strings"
)

type MetamodelParser struct {
	compinfo      *info.ComponentInfo
	structureInfo map[string]*info.StructureInfo
}

func NewMetamodelParser(compinfo *info.ComponentInfo) *MetamodelParser {
	strInfo := make(map[string]*info.StructureInfo)
	return &MetamodelParser{compinfo: compinfo, structureInfo: strInfo}
}

func (parser *MetamodelParser) ComponentInfo() *info.ComponentInfo {
	return parser.compinfo
}

func (parser *MetamodelParser) Parse(authzData []byte) error {

	var err error
	if parser.compinfo == nil {
		parser.compinfo = info.NewComponentInfo()
	}
	var jsonData interface{}

	err = json.Unmarshal(authzData, &jsonData)
	if err != nil {
		return err
	}

	// Metamodel Data
	if _, ok := jsonData.(map[string]interface{}); !ok {
		return getError("Json Metadata failed map[sting]interface assertion")
	} else if _, ok := jsonData.(map[string]interface{})[AUTH_METAMODEL]; !ok {
		return getError("Json Metadata lookup failed to get value of `metamodel`")
	}
	metamodelData := jsonData.(map[string]interface{})[AUTH_METAMODEL]

	// Component Data
	if _, ok := metamodelData.(map[string]interface{}); !ok {
		return getError("Metamodel data failed map[sting]interface assertion")
	} else if _, ok := metamodelData.(map[string]interface{})[AUTH_COMPONENT]; !ok {
		return getError("Metadata lookup failed to get value of `component`")
	}
	componentData := metamodelData.(map[string]interface{})[AUTH_COMPONENT]

	// Packages
	if _, ok := componentData.(map[string]interface{}); !ok {
		return getError("Component data failed map[sting]interface assertion")
	} else if _, ok := componentData.(map[string]interface{})[AUTH_PACKAGES]; !ok {
		return getError("Component data lookup failed to get value of `packages`")
	}
	packagesArray := componentData.(map[string]interface{})[AUTH_PACKAGES]

	if _, ok := packagesArray.([]interface{}); !ok {
		return getError("Packages array failed to assert to type []interface")
	}
	packages := packagesArray.([]interface{})

	// Process Component's Name
	if _, ok := componentData.(map[string]interface{})[NAME]; !ok {
		return getError("Component data lookup failed to get value of `name`")
	} else if _, ok := componentData.(map[string]interface{})[NAME].(string); !ok {
		return getError("value of Component data lookup of `name` failed to assert to type string")
	}
	name := componentData.(map[string]interface{})[NAME].(string)
	parser.compinfo.SetIdentifier(name)

	// process Component's Packages
	err = parser.parsePackages(packages)
	if err != nil {
		return err
	}

	return nil
}

func (parser *MetamodelParser) parsePackages(packages []interface{}) error {

	for _, pkg := range packages {

		if _, ok := pkg.(map[string]interface{}); !ok {
			return getError("Packages failed to assert for type map[string]interface{}")
		}
		packge := pkg.(map[string]interface{})
		var pkginfo *info.PackageInfo

		// process Package's name
		if _, ok := packge[NAME]; !ok {
			return getError("packages attribute lookup for attribute `name` failed")
		} else if _, ok := packge[NAME].(string); !ok {
			return getError("packages attribute `name` assertion interface to string failed")
		} else {
			name := packge[NAME].(string)
			pkginfo = parser.compinfo.PackageInfo(name)
			if pkginfo == nil {
				pkginfo = info.NewPackageInfo()
			}
			pkginfo.SetIdentifier(name)
			parser.compinfo.SetPackageInfo(name, pkginfo)
		}

		// process Package's Services
		if _, ok := packge[SERVICES]; !ok {
			return getError("packages attribute lookup for attribute `services` failed")
		} else if _, ok := packge[SERVICES].([]interface{}); !ok {
			return getError("packages attribute `services` assertion to []interface{} failed")
		} else {
			services := packge[SERVICES].([]interface{})
			err := parser.parseServices(services, pkginfo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (parser *MetamodelParser) parseServices(services []interface{}, pkg *info.PackageInfo) error {

	for _, service := range services {

		if _, ok := service.(map[string]interface{}); !ok {
			return getError("services failed to assert for type map[string]interface")
		}
		ser := service.(map[string]interface{})
		var serinfo *info.ServiceInfo

		// process Service's name
		if _, ok := ser[NAME]; !ok {
			return getError("Service attribute lookup for attribute `name` failed")
		} else if _, ok := ser["name"].(string); !ok {
			return getError("Service attribute `name` assertion interface to string failed")
		} else {
			name := ser["name"].(string)
			serinfo = pkg.ServiceInfo(name)
			if serinfo == nil {
				serinfo = info.NewServiceInfo()
			}
			serinfo.SetIdentifier(name)
			pkg.SetServiceInfo(name, serinfo)
		}

		// process Service's Structures
		if _, ok := ser[STRUCTURES]; !ok {
			return getError("Service attribute lookup for attribute `structures` failed")
		} else if _, ok := ser[STRUCTURES].([]interface{}); !ok {
			return getError("Service attribute `structures` assertion to []interface failed")
		} else {
			structures := ser[STRUCTURES].([]interface{})
			err := parser.parseStructures(structures, serinfo)
			if err != nil {
				return err
			}
		}

		// process Service's operations
		if _, ok := ser[OPERATIONS]; !ok {
			return getError("Service attribute lookup for attribute `operations` failed")
		} else if _, ok := ser[OPERATIONS].([]interface{}); !ok {
			return getError("Service attribute `operations` assertion to []interface failed")
		} else {
			operations := ser[OPERATIONS].([]interface{})
			err := parser.parseOperations(operations, serinfo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (parser *MetamodelParser) parseStructures(structures []interface{}, ser *info.ServiceInfo) error {

	for _, structure := range structures {

		if _, ok := structure.(map[string]interface{}); !ok {
			return getError("Structure failed to assert for type map[string]interface")
		}
		structMap := structure.(map[string]interface{})
		structinfo := info.NewStructureInfo()

		// process Structure's name
		if _, ok := structMap[NAME]; !ok {
			return getError("Structures attribute lookup for attribute `name` failed")
		} else if _, ok := structMap[NAME].(string); !ok {
			return getError("Structures attribute `name` assertion to string failed")
		} else {
			name := structMap[NAME].(string)
			structinfo.SetIdentifier(name)
			ser.SetStructureInfo(name, structinfo)
			parser.structureInfo[name] = structinfo
		}

		// process Structure's Fields
		if _, ok := structMap[FIELDS]; !ok {
			return getError("Structures attribute lookup for attribute `fields` failed")
		} else if _, ok := structMap[FIELDS].([]interface{}); !ok {
			return getError("Structures attribute `fields` assertion to []interface failed")
		} else {
			fields := structMap[FIELDS].([]interface{})
			err := parser.parseFields(fields, structinfo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (parser *MetamodelParser) parseFields(fields []interface{}, str info.Info) error {
	for _, f := range fields {

		if _, ok := f.(map[string]interface{}); !ok {
			return getError("fields failed to assert for type map[string]interface")
		}
		field := f.(map[string]interface{})
		var fieldinfo *info.FieldInfo

		// process Field's name
		if _, ok := field[NAME]; !ok {
			return getError("fields attribute lookup for attribute `name` failed")
		} else if _, ok := field[NAME].(string); !ok {
			return getError("fields attribute `name` assertion to string failed")
		} else {
			name := field[NAME].(string)
			fieldinfo = info.NewFieldInfo(name)
			str.SetFieldInfo(name, fieldinfo)
		}

		// process Field's Metadata Resource/HasFieldsOf/IsOneOf
		if meta, ok := field["metadata"].(map[string]interface{}); ok {
			for metakey, val := range meta {
				if metavalue, ok := val.(map[string]interface{}); ok {

					// Processing HasFieldsOf metadata
					if metakey == HASFIELDSOF {
						var interSlice []interface{}
						switch reflect.TypeOf(metavalue[VALUE]) {

						// Process String
						case reflect.TypeOf(""):
							fieldinfo.AddHasFieldsOfTypeName(metavalue[VALUE].(string))

						// Process list of String
						case reflect.TypeOf(interSlice):
							temp := []string{}
							for _, tempVal := range metavalue[VALUE].([]interface{}) {
								temp = append(temp, tempVal.(string))
							}
							fieldinfo.SetHasFieldsOfTypeNames(temp)

						}

						// Processing Resource metadata
					} else if metakey == RESOURCE {
						for resourceKey, resourceValue := range metavalue {
							if resourceKey == TYPEHOLDER {
								fieldinfo.SetIsIdentifierPolymorphic(true)
								fieldinfo.SetIdentifierTypeHolder(resourceValue.(string))
							} else if resourceKey == VALUE {
								if resourceStr, ok := resourceValue.(string); ok {
									fieldinfo.SetIdentifierType(resourceStr)
								} else if resourceList, ok := resourceValue.([]string); ok {
									fieldinfo.SetIdentifierTypes(resourceList)
								}
							}
						}
					}
					// TODO: Handle IsOneOf cases
				}
			}
		} else {
			args := map[string]string{
				"msg":        "fields attribute metadata doesnot exists or its assertion to map[string]interface failed",
				"field name": field["name"].(string),
			}
			return l10n.NewRuntimeError("vapi.metadata.parser.failure", args)
		}

		// process Field's type Resource
		if _, ok := field[TYPE]; !ok {
			return getError("fields attribute lookup for attribute `type` failed")
		} else if _, ok := field[TYPE].(map[string]interface{}); !ok {
			return getError("fields attribute `type` assertion to map[string]interface failed")
		} else {
			typ := field[TYPE].(map[string]interface{})
			var typeinfo info.TypeInfo
			typeinfo, err := parser.parseTypeInfo(typ)
			if err != nil {
				return err
			}
			fieldinfo.SetTypeInfo(typeinfo)
		}
	}
	return nil
}

func (parser *MetamodelParser) parseTypeInfo(data map[string]interface{}) (info.TypeInfo, error) {
	// check category: if it's primitive, generic or user_defined

	if _, ok := data[CATEGORY]; !ok {
		return nil, getError("Parse TypeInfo Data attribute lookup for attribute `category` failed")
	} else if _, ok := data[CATEGORY].(string); !ok {
		return nil, getError("Parse TypeInfo Data attribute `category` assertion to string failed")
	}

	category := data[CATEGORY].(string)
	if category == PRIMITIVE {
		return parser.parsePrimitive(data)
	} else if category == GENERIC {
		return parser.parseGeneric(data)
	} else if category == USER_DEFINED {
		name, ok := data[USER_DEFINED_TYPE_NAME].(string)
		if !ok {
			return nil, getError("TypeInfo.user_defined attribute user_defined_type_name doesnot exists or assertion to string failed")
		}
		return info.NewUserDefinedTypeInfo(name, parser.structureInfo[name]), nil
	}
	return nil, getError("Parse TypeInfo detected unknown category")
}

func (parser *MetamodelParser) parseGeneric(data map[string]interface{}) (info.TypeInfo, error) {

	genericInfo := info.NewGenericTypeInfo("")

	if _, ok := data[GENERIC_TYPE]; !ok {
		return nil, getError("Generic TypeInfo attribute generic_type lookup failed")
	} else if _, ok := data[GENERIC_TYPE].(string); !ok {
		return nil, getError("Generic TypeInfo attribute generic_type assertion to string failed")
	}

	switch data[GENERIC_TYPE].(string) {
	case "List":
		// TODO: check what attributes are generated and use them
	case "Optional":

		if _, ok := data[Element_TYPE]; !ok {
			return nil, getError("Generic.Optional TypeInfo attribute element_type lookup failed")
		} else if _, ok := data[Element_TYPE].(map[string]interface{}); !ok {
			return nil, getError("Generic.Optional TypeInfo attribute element_type assertion to map[string]interface failed")
		}

		var elementType info.TypeInfo
		elementType, err := parser.parseTypeInfo(data[Element_TYPE].(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		genericInfo.SetIsOptional(true)
		genericInfo.SetElementType(elementType)

	case "Set":
		// TODO: check what attributes are generated and use them
	case "Map":
		var keyTypeInfo info.TypeInfo
		var valTypeInfo info.TypeInfo
		genericInfo.SetIsMap(true)

		keyTypeInfo, err := parser.parseTypeInfo(data[MAP_KEY].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		valTypeInfo, err = parser.parseTypeInfo(data[MAP_VALUE].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		genericInfo.SetMapKeyType(keyTypeInfo)
		genericInfo.SetMapValueType(valTypeInfo)
	}
	// assign to type info interface to generic info
	return genericInfo, nil
}

func (parser *MetamodelParser) parsePrimitive(data map[string]interface{}) (info.TypeInfo, error) {

	var datatype string

	if _, ok := data[PRIMITIVE_TYPE]; !ok {
		return nil, getError("TypeInfo.Primitive attribute primitive_typelookup failed")
	} else if _, ok := data[PRIMITIVE_TYPE].(string); !ok {
		return nil, getError("TypeInfo.Primitive attribute primitive_type assertion to string failed")
	}

	datatype = strings.ToUpper(data[PRIMITIVE_TYPE].(string))

	return info.NewPrimitiveTypeInfo("", info.PrimitiveTypeFromString(datatype)), nil
}

func (parser *MetamodelParser) parseOperations(operations []interface{}, ser *info.ServiceInfo) error {
	for _, operation := range operations {

		if _, ok := operation.(map[string]interface{}); !ok {
			return getError("Operation failed to assert for type map[string]interface")
		}
		op := operation.(map[string]interface{})
		opInfo := info.NewOperationInfo()

		// process Operation's name
		if _, ok := op[NAME]; !ok {
			return getError("Operations attribute lookup for attribute `name` failed")
		} else if _, ok := op[NAME].(string); !ok {
			return getError("Operations attribute `name` assertion to string failed")
		} else {
			name := op[NAME].(string)
			opInfo.SetIdentifier(name)
			ser.SetOperationInfo(name, opInfo)
		}

		// process Operation's params
		if _, ok := op[PARAMS]; !ok {
			return getError("Operations attribute lookup for attribute `params` failed")
		} else if _, ok := op[PARAMS].([]interface{}); !ok {
			return getError("Operations attribute `params` assertion to []interface failed")
		} else {
			err := parser.parseFields(op[PARAMS].([]interface{}), opInfo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
