/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package metadata

import (
	"encoding/json"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/metadata/info"
	"strings"
)

type PrivilegeParser struct {
	compinfo *info.ComponentInfo
}

func NewPrivilegeParser(compinfo *info.ComponentInfo) *PrivilegeParser {
	temp := PrivilegeParser{compinfo: compinfo}
	return &temp
}

func (parser *PrivilegeParser) ComponentInfo() *info.ComponentInfo {
	return parser.compinfo
}

func (parser *PrivilegeParser) Parse(authzData []byte) error {

	var err error
	if parser.compinfo == nil {
		parser.compinfo = info.NewComponentInfo()
	}
	var jsonData interface{}
	var packageInfo map[string]*info.PackageInfo = parser.compinfo.PackageMap()

	err = json.Unmarshal(authzData, &jsonData)
	if err != nil {
		return err
	}

	privilegeData, ok := jsonData.(map[string]interface{})[AUTH_PRIVILEGE]
	if !ok {
		return getError("Json Data attribute privilege doesnot exists or assertion to map[string]interface{} failed")
	}
	componentData := privilegeData.(map[string]interface{})[AUTH_COMPONENT]
	if !ok {
		return getError("Json Data attribute component doesnot exists or assertion to map[string]interface{} failed")
	}

	for key, val := range componentData.(map[string]interface{}) {
		var privilegeMap map[string][]string
		if key == AUTH_PRIVILEGE_COMPONENT_NAME {
			continue
		} else {
			// marshal component data in form of interface{} to byte array
			byteData, err := json.Marshal(val)
			if err != nil {
				return err
			}

			// Unmarshal byte array to read component data in form of map[string][]string
			err = json.Unmarshal(byteData, &privilegeMap)
			if err != nil {
				return err
			}

			if key == AUTH_PRIVILEGE_DEFAUlT {
				handleDefaultData(privilegeMap, packageInfo)
			} else {
				handleServiceData(key, privilegeMap, packageInfo)
			}
		}
	}

	return nil
}

func handleDefaultData(data map[string][]string, pkgInfo map[string]*info.PackageInfo) {
	// get Default key value pair i.e <package>: [list of privileges]
	for pkg, privilegeArray := range data {
		// Handling Package Info
		if pkgInfo[pkg] == nil {
			pkgInfo[pkg] = info.NewPackageInfo()
		}
		pkgInfo[pkg].SetPrivileges(privilegeArray)
	}
}

func handleServiceData(serviceName string, data map[string][]string, pkgInfo map[string]*info.PackageInfo) {

	packageName := GetPackageName(serviceName)

	if pkgInfo[packageName] == nil {
		pkgInfo[packageName] = info.NewPackageInfo()
	}

	var serviceInfo *info.ServiceInfo = pkgInfo[packageName].ServiceInfo(serviceName)

	if serviceInfo == nil {
		serviceInfo = info.NewServiceInfo()
		pkgInfo[packageName].SetServiceInfo(serviceName, serviceInfo)
	}

	var operInfo *info.OperationInfo
	for operation, privilegeArray := range data {
		tokens := strings.Split(operation, ".")
		// Tokens > 1, privileges provided are for a parameters otherwise it is for the named operation
		if len(tokens) > 1 {
			// Value at 0th index of tokens represent functions name while rest of it Parameter details
			// Handling params PrivilegeInfo
			// Check if operationInfo already exists for operation tokens[0]
			operInfo = serviceInfo.OperationInfo(tokens[0])
			if operInfo == nil {
				operInfo = info.NewOperationInfo()
				serviceInfo.SetOperationInfo(operation, operInfo)
			}
			// Creating params PrivilegeInfo
			privilegeInfo := info.NewPrivilegeInfo()
			privilegeInfo.SetPropertyPath(strings.Join(tokens[1:], "."))
			privilegeInfo.SetPrivileges(privilegeArray)

			operInfo.AddPrivilegeInfo(privilegeInfo)
		} else {
			operInfo = serviceInfo.OperationInfo(operation)
			// Handling Operation/Service Info
			if operInfo == nil {
				operInfo = info.NewOperationInfo()
			}
			operInfo.SetPrivileges(privilegeArray)
			serviceInfo.SetOperationInfo(operation, operInfo)
		}
	}
}

func GetPackageName(key string) string {
	// serviceName is always in the format: package_name.service_name i.e.
	// there is always at least a single '.' separator.
	// Example:
	// "privilege": {
	//     "component": {
	//         "name": "...."
	//         "default": { ... },
	//         "package_name.service_name": { ... }
	//         ...
	//         ...
	//     }
	// }
	// input: package_name.service_name -> output: package_name
	return key[0:strings.LastIndex(key, ".")]
}
