/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import "github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"

type OpaqueDefinition struct{}

func NewOpaqueDefinition() OpaqueDefinition {
	return OpaqueDefinition{}
}

func (opaqueDefinition OpaqueDefinition) Type() DataType {
	return OPAQUE
}

func (opaqueDefinition OpaqueDefinition) CompleteValue(value DataValue) {

}
func (opaqueDefinition OpaqueDefinition) String() string {
	return OPAQUE.String()
}

/**
 * Validates that the specified DataValue is an instance of this data
 * definition.
 *
 * <p>Only validates that supplied value is not <code>nil</code>.
 */
func (opaqueDefinition OpaqueDefinition) Validate(value DataValue) []error {
	if value == nil {
		var msg = l10n.NewRuntimeError("vapi.data.opaque.definition.null.value", map[string]string{})
		return []error{msg}
	}
	return nil
}

func (opaqueDefinition OpaqueDefinition) ValidInstanceOf(value DataValue) bool {
	return true
}
