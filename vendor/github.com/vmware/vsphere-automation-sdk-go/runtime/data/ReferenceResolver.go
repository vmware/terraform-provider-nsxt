/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
)

/**
 * Context which keeps track of structure definitions and structure references
 * in a graph of definition objects. The context can resolve the unresolved
 * references after it learns about all definitions and references in the
 * definition graph.
 */
type ReferenceResolver struct {
	definitions map[string]StructDefinition
	references  []*StructRefDefinition
}

func NewReferenceResolver() *ReferenceResolver {
	return &ReferenceResolver{definitions: map[string]StructDefinition{}, references: []*StructRefDefinition{}}
}

func (r *ReferenceResolver) AddDefinition(structDef StructDefinition) {
	r.definitions[structDef.Name()] = structDef
}

func (r *ReferenceResolver) AddReference(structRefDef *StructRefDefinition) {
	r.references = append(r.references, structRefDef)
}

func (r *ReferenceResolver) IsDefined(structName string) bool {
	if _, ok := r.definitions[structName]; ok {
		return true
	}
	return false
}

/**
 * Traverses all references and resolves the unresolved ones.
 * throws error if a reference cannot be resolved.
 */
func (r *ReferenceResolver) ResolveReferences() []error {
	for _, ref := range r.references {
		if ref.Target() == nil {
			if structDef, ok := r.definitions[ref.Name()]; ok {
				ref.SetTarget(&structDef)
			} else {
				return []error{l10n.NewRuntimeError("vapi.data.structref.structure.not.defined",
					map[string]string{"referenceType": ref.Name()})}
			}
		}
	}
	return nil
}

func (r *ReferenceResolver) Definition(structName string) StructDefinition {
	return r.definitions[structName]
}
