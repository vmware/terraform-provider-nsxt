/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"strings"
)

type MethodDefinition struct {
	id                  MethodIdentifier
	input               data.DataDefinition
	output              data.DataDefinition
	errorDefinitionList []data.ErrorDefinition
	errorDefinitionMap  map[string]data.ErrorDefinition
}

func NewMethodDefinition(id MethodIdentifier, input data.DataDefinition, output data.DataDefinition, errorDefs []data.ErrorDefinition) MethodDefinition {

	if input == nil {
		log.Errorf("Input data definition of the MethodIdentifier: %s is missing.", id)
		return MethodDefinition{}
	}
	if output == nil {
		log.Errorf("Output data definition of the MethodIdentifier: %s is missing.", id)
		return MethodDefinition{}
	}

	if errorDefs == nil {
		var errorDefinitionList = make([]data.ErrorDefinition, 0)
		var errorDefinitionMap = make(map[string]data.ErrorDefinition)
		return MethodDefinition{id: id, input: input, output: output, errorDefinitionList: errorDefinitionList, errorDefinitionMap: errorDefinitionMap}
	} else {
		var errorDefMap = make(map[string]data.ErrorDefinition)
		for _, errorDef := range errorDefs {
			errorDefMap[errorDef.Name()] = errorDef
		}
		return MethodDefinition{id: id, input: input, output: output, errorDefinitionList: errorDefs, errorDefinitionMap: errorDefMap}
	}

}

func (md MethodDefinition) Identifier() MethodIdentifier {
	return md.id
}

func (md MethodDefinition) InputDefinition() data.DataDefinition {
	return md.input
}

func (md MethodDefinition) OutputDefinition() data.DataDefinition {
	return md.output
}

func (md MethodDefinition) ErrorDefinitions() []data.ErrorDefinition {
	return md.errorDefinitionList
}

func (md MethodDefinition) ErrorDefinition(errorName string) data.ErrorDefinition {
	return md.errorDefinitionMap[errorName]
}

func (md MethodDefinition) HasErrorDefinition(errorName string) bool {
	if _, ok := md.errorDefinitionMap[errorName]; ok {
		return true
	}
	return false
}

func (md MethodDefinition) String() string {
	var sb strings.Builder
	if (md.id != MethodIdentifier{}) {
		sb.WriteString("Name: " + md.id.FullyQualifiedName())
	}
	if md.input != nil {
		sb.WriteString("Input: " + md.input.String())
	}
	if md.output != nil {
		sb.WriteString("Output: " + md.output.String())
	}
	for i := 0; i < len(md.errorDefinitionList); i++ {
		sb.WriteString("Error: " + md.errorDefinitionList[i].Name())
	}
	return sb.String()
}
