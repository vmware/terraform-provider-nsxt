/* Copyright © 2019, 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import (
	"encoding/json"
)

// ApplicationContext provides additional execution context information.
// Additional information is in key-value string pairs and is set by the
// client. It is provided to the server implementation as is.
// This class is not safe for concurrent use.
// Use Copy method for such scenarios.
type ApplicationContext struct {
	wireData map[string]*string
}

// NewApplicationContext instantiates ApplicationContext object using provided
// map.
func NewApplicationContext(wireData map[string]*string) *ApplicationContext {
	if wireData == nil {
		wireData = map[string]*string{}
	}
	return &ApplicationContext{wireData: wireData}
}

// Copy makes a copy of the current object
func (a *ApplicationContext) Copy() *ApplicationContext {
	wireDataCopy := make(map[string]*string)
	for k, v := range a.wireData {
		wireDataCopy[k] = v
	}
	return &ApplicationContext{wireData: wireDataCopy}
}

func (a *ApplicationContext) GetProperty(key string) *string {
	return a.wireData[key]
}

func (a *ApplicationContext) HasProperty(key string) bool {
	if _, ok := a.wireData[key]; ok {
		return true
	}
	return false
}

func (a *ApplicationContext) GetAllProperties() map[string]*string {
	return a.wireData
}

func (a *ApplicationContext) SetProperty(key string, value *string) {
	a.wireData[key] = value
}

func (a *ApplicationContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.wireData)
}
