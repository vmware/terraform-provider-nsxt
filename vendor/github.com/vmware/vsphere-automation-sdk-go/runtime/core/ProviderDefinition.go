/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

type ProviderDefinition struct {
	id       string
	checkSum string
}

func (pd *ProviderDefinition) Identifier() string {
	return pd.id
}

func (pd *ProviderDefinition) CheckSum() string {
	return pd.checkSum
}

func NewProviderDefinition(id string, checkSum string) *ProviderDefinition {
	return &ProviderDefinition{id: id, checkSum: checkSum}
}
