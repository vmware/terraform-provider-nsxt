/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security


type ResourceIdentifier struct {
	isOperation  bool
	id           string
	resourceType string
}

func NewResourceIdentifier(isOperation bool, id string, resourceType string) ResourceIdentifier {
	return ResourceIdentifier{isOperation: isOperation, id: id, resourceType: resourceType}
}

func (r ResourceIdentifier) IsOperation() bool {
	return r.isOperation
}

func (r ResourceIdentifier) Id() string {
	return r.id
}

func (r ResourceIdentifier) ResourceType() string {
	return r.resourceType
}
