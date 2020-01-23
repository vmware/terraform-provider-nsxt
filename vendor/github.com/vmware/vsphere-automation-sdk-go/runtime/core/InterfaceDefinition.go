/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

type InterfaceDefinition struct {
	id        InterfaceIdentifier
	methodIds []MethodIdentifier
}

func NewInterfaceDefinition(id InterfaceIdentifier, methodIds []MethodIdentifier) InterfaceDefinition {
	return InterfaceDefinition{id: id, methodIds: methodIds}
}

func (interfaceDefinition InterfaceDefinition) Identifier() InterfaceIdentifier {
	return interfaceDefinition.id
}

func (interfaceDefinition InterfaceDefinition) MethodIdentifiers() []MethodIdentifier {
	return interfaceDefinition.methodIds
}

func (interfaceDefinition InterfaceDefinition) Equals(other InterfaceDefinition) bool {
	return interfaceDefinition.id.Equals(other.id)
	//TODO: sreeshas
	// write code to compare slice of MethodIdentifiers
}
