/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"sort"
	"strings"
)

type MethodDefinition struct {
	id                 MethodIdentifier
	input              data.DataDefinition
	output             data.DataDefinition
	errorDefinitionMap map[string]data.ErrorDefinition
}

var allowedErrors = []data.ErrorDefinition{
	bindings.INTERNAL_SERVER_ERROR_DEF,
	bindings.INVALID_ARGUMENT_ERROR_DEF,
	bindings.UNEXPECTED_INPUT_ERROR_DEF,
	bindings.OP_NOT_FOUND_ERROR_DEF,
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

	errorDefinitionMap := make(map[string]data.ErrorDefinition)
	for _, errorDef := range allowedErrors {
		errorDefinitionMap[errorDef.Name()] = errorDef
	}

	for _, errorDef := range errorDefs {
		errorDefinitionMap[errorDef.Name()] = errorDef
	}

	return MethodDefinition{id: id, input: input, output: output, errorDefinitionMap: errorDefinitionMap}
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

// ErrorDefinitions Returns method's error definitions sorted by error name
func (md MethodDefinition) ErrorDefinitions() []data.ErrorDefinition {
	var errorDefinitions = make([]data.ErrorDefinition, 0, len(md.errorDefinitionMap))
	keys := make([]string, 0, len(md.errorDefinitionMap))

	for key, _ := range md.errorDefinitionMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		errorDefinitions = append(errorDefinitions, md.errorDefinitionMap[key])
	}

	return errorDefinitions
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
	for _, errorDef := range md.ErrorDefinitions() {
		sb.WriteString("Error: " + errorDef.Name())
	}
	return sb.String()
}
