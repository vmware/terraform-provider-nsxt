/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import (
	"errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

/**
 * Reference to a {@link StructDefinition}. If the reference is resolved, it is
 * bound to a specific {@link StructDefinition} target. If the reference is
 * unresolved, its target is <code>nil</code>.
 */
/**
 * StructRef is different when compared to other Definition structs.
 * It returns a pointer instead of value because of the resolving task involved.
 * When StructRef gets resolved, all its references will also be resolved.
 */
type StructRefDefinition struct {
	name             string
	structDefinition *StructDefinition
}

func NewStructRefDefinition(name string, structDefinition *StructDefinition) *StructRefDefinition {
	if name == "" {
		log.Error("StructRef name missing")
	}
	return &StructRefDefinition{name: name, structDefinition: structDefinition}
}

func (structRefDefinition StructRefDefinition) Name() string {
	return structRefDefinition.name
}

func (structRefDefinition StructRefDefinition) Target() *StructDefinition {
	return structRefDefinition.structDefinition
}

func (structRefDefinition *StructRefDefinition) SetTarget(structDefinition *StructDefinition) error {
	if structDefinition == nil {
		return errors.New("StructDefinition may not be nil")
	}
	if structRefDefinition.structDefinition != nil {
		return errors.New("Structref already resolved")
	}
	if structRefDefinition.name != structDefinition.name {
		return errors.New("type mismatch")
	}

	structRefDefinition.structDefinition = structDefinition
	return nil
}

func (structRefDefinition StructRefDefinition) Type() DataType {
	return STRUCTURE_REF
}

func (structRefDefinition StructRefDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = structRefDefinition.Validate(value)
	if len(errors) != 0 {
		return false
	}
	return true
}
func (structRefDefinition StructRefDefinition) Validate(value DataValue) []error {
	var errorList = structRefDefinition.CheckResolved()
	if errorList != nil {
		return errorList
	}
	return structRefDefinition.structDefinition.Validate(value)

}
func (structRefDefinition StructRefDefinition) CompleteValue(value DataValue) {
	var error = structRefDefinition.CheckResolved()
	if error == nil {
		structRefDefinition.structDefinition.CompleteValue(value)
		return
	} else {
		log.Error(error)
	}
}
func (structRefDefinition StructRefDefinition) String() string {
	return STRUCTURE_REF.String()
}

func (structRefDefinition StructRefDefinition) CheckResolved() []error {
	if structRefDefinition.structDefinition == nil {
		return []error{l10n.NewRuntimeError("vapi.data.structref.not.resolved",
			map[string]string{"referenceType": structRefDefinition.name})}
	}
	return nil
}
