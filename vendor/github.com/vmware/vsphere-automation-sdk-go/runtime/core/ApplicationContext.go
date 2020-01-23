/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import (
	"encoding/json"
)

type ApplicationContext struct {
	wireData map[string]*string
}

func NewApplicationContext(appContext map[string]*string) *ApplicationContext {
	if appContext == nil {
		appContext = map[string]*string{}
	}
	return &ApplicationContext{wireData: appContext}
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