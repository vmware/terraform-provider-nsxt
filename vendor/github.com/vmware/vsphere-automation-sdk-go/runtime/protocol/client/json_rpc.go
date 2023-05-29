/* Copyright © 2019-2021, 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package client

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client/internal"
	"net/http"
)

// Deprecated: use NewConnector(address) instead
type JsonRpcConnector struct {
	*internal.JsonRpcHttpProtocol
	connector *connector
}

// NewJsonRpcConnector instantiates instance of JsonRpcConnector
// Deprecated: use NewConnector(url, WithHttpClient(client)) instead
func NewJsonRpcConnector(url string, client http.Client, options ...ConnectorOption) *JsonRpcConnector {
	options = append(options, WithHttpClient(&client))
	connector := NewConnector(url, options...)
	return &JsonRpcConnector{
		JsonRpcHttpProtocol: connector.protocol.(*internal.JsonRpcHttpProtocol),
		connector:           connector}
}

// SetApplicationContext specifies security context to be used by Connector instance
// Deprecated: security context should be specified when instantiating Connector instance
// Use WithApplicationContext ConnectorOption helper method
func (j *JsonRpcConnector) SetApplicationContext(context *core.ApplicationContext) {
	j.connector.appContext = context
}

// SetSecurityContext specifies application context to be used by Connector instance
// Deprecated: application context should be specified when instantiating Connector instance
// Use WithSecurityContext ConnectorOption helper method
func (j *JsonRpcConnector) SetSecurityContext(context core.SecurityContext) {
	j.connector.SetSecurityContext(context)
}
