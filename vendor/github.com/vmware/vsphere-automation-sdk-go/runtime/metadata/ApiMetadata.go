/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package metadata


import (
	"encoding/json"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/metadata/info"
	"io/ioutil"
)

type ApiMetadata struct {
	compInfo           map[string]*info.ComponentInfo
	metamodelFilePaths []string
	privilegeFilePaths []string
}

func (apimeta *ApiMetadata) ComponentInfo() map[string]*info.ComponentInfo {
	return apimeta.compInfo
}

func NewApiMetadata(metamodelFilePaths []string, privilegeFilePaths []string) (*ApiMetadata, error) {
	comp := make(map[string]*info.ComponentInfo)
	apimeta := ApiMetadata{
		compInfo:           comp,
		metamodelFilePaths: metamodelFilePaths,
		privilegeFilePaths: privilegeFilePaths,
	}
	err := apimeta.LoadAuthzMetadata()
	if err != nil {
		return nil, err
	}
	return &apimeta, nil
}

func getMetadataName(metadata string, component string, data []byte) (string, error) {

	var jsonData interface{}

	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Error(err)
	}

	if _, ok := jsonData.(map[string]interface{}); !ok {
		return "", getError("Json Metadata failed map[sting]interface assertion")
	} else if _, ok := jsonData.(map[string]interface{})[AUTH_METAMODEL]; !ok {
		return "", getError("Json Metadata lookup failed to get value of `metamodel`")
	}
	metamodelData := jsonData.(map[string]interface{})[AUTH_METAMODEL]

	// Component Data
	if _, ok := metamodelData.(map[string]interface{}); !ok {
		return "", getError("Metamodel data failed map[sting]interface assertion")
	} else if _, ok := metamodelData.(map[string]interface{})[AUTH_COMPONENT]; !ok {
		return "", getError("Metadata lookup failed to get value of `component`")
	}
	componentData := metamodelData.(map[string]interface{})[AUTH_COMPONENT]

	// name
	if _, ok := componentData.(map[string]interface{})[NAME]; !ok {
		return "", getError("Component data lookup failed to get value of `name`")
	} else if _, ok := componentData.(map[string]interface{})[NAME].(string); !ok {
		return "", getError("value of Component data lookup of `name` failed to assert to type string")
	}
	name := componentData.(map[string]interface{})[NAME].(string)

	return name, nil
}

func (apimeta *ApiMetadata) LoadAuthzMetadata() error {

	var compInfo *info.ComponentInfo

	// Parse through metamodel
	for _, filepath := range apimeta.metamodelFilePaths {
		metamodeldata, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Error("Error reading authorization metadata file %s", filepath)
			return err
		}

		name, err := getMetadataName(AUTH_METAMODEL, AUTH_COMPONENT, metamodeldata)
		if err != nil {
			return err
		}

		compInfo = apimeta.compInfo[name]
		if compInfo == nil {
			compInfo = info.NewComponentInfo()
		}

		meta := NewMetamodelParser(compInfo)
		err = meta.Parse(metamodeldata)
		if err != nil {
			return err
		}

		apimeta.compInfo[name] = compInfo
	}

	// Parse through Privilege paths
	for _, filepath := range apimeta.privilegeFilePaths {
		privilegeData, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Error("Error reading authorization metadata file %s", filepath)
			return err
		}

		name, err := getMetadataName(AUTH_PRIVILEGE, AUTH_COMPONENT, privilegeData)
		if err != nil {
			return err
		}

		compInfo = apimeta.compInfo[name]
		if compInfo == nil {
			compInfo = info.NewComponentInfo()
		}

		priv := NewPrivilegeParser(compInfo)
		err = priv.Parse(privilegeData)
		if err != nil {
			return err
		}
	}

	return nil
}
