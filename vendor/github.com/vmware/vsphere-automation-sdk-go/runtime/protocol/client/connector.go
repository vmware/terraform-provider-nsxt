/* Copyright © 2019, 2021, 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package client

import (
	"sync"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/common"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client/internal"
)

// Connector provides and keeps common information for executing vapi bindings requests
// todo: embedding internal interfaces is not a good option
// as it might mislead developer go make backwards incompatible changes.
// Move interface to protocol package in order to remove need for referencing internal interface
// this would require bindings update as currently they use client.Connector
// internal code, currently, can not reference this interface as cyclic imports are not allowed in go
type Connector interface {
	internal.Connector
}

type connector struct {
	address            string
	protocol           core.APIProvider
	provider           core.APIProvider
	decorators         []core.APIProviderDecorator
	appContext         *core.ApplicationContext
	connectionMetadata map[string]interface{}
	typeConverter      *bindings.TypeConverter
	// mu guards the SecurityContext
	mu              sync.Mutex
	securityContext core.SecurityContext
}

// NewConnector instantiates connector object, used by vAPI bindings for client-server communication.
// returned Connector instance is safe to be used across go routines (once VAPI-4899 gets resolved).
// Default vAPI protocol used for communication is JSON-RPC. To use a different one such as REST,
// appropriate ConnectorOption need to be provided.
// connectorOptions are used also to adjust the connector instance appropriately. Such as setting
// specific application and security contexts, specific http.Client instance and so on.
func NewConnector(address string, connectorOptions ...ConnectorOption) *connector {
	defaultConnector := getDefaultConnector(address)

	for _, fn := range connectorOptions {
		fn(defaultConnector)
	}

	chainDecorators(defaultConnector)

	return defaultConnector
}

func getDefaultConnector(address string) *connector {
	c := &connector{
		address:       address,
		typeConverter: bindings.NewTypeConverter(),
	}
	c.protocol = internal.NewJsonRpcHttpProtocol(c)

	return c
}

func chainDecorators(connector *connector) {
	connector.provider = connector
	for _, decorator := range connector.decorators {
		connector.provider = decorator(connector.provider)
	}
}

// <-- Connector interface implementation start

// Address gets connector's address
func (c *connector) Address() string {
	return c.address
}

// ApplicationContext gets connector's application context
func (c *connector) ApplicationContext() *core.ApplicationContext {
	return c.appContext
}

// SecurityContext gets connector's security context
func (c *connector) SecurityContext() core.SecurityContext {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.securityContext
}

// SetSecurityContext replaces the connector's security context.
// Useful when credentials change, e.g. session expires or password is changed.
func (c *connector) SetSecurityContext(context core.SecurityContext) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.securityContext = context
}

// GetApiProvider gets current connector's provider. By default this is the connector
// instance itself. This provider could be decorated by WithDecorators ConnectorOption
// to extend connector's Invoke functionality by calling a different core.APIProvider beforehand.
func (c *connector) GetApiProvider() core.APIProvider {
	return c.provider
}

// TypeConverter gets connector's type converter
func (c *connector) TypeConverter() *bindings.TypeConverter {
	return c.typeConverter
}

// NewExecutionContext creates vAPI execution context from connector's security and application contexts
func (c *connector) NewExecutionContext() *core.ExecutionContext {
	if c.appContext == nil {
		c.appContext = core.NewApplicationContext(nil)
	}
	// use application context copy for thread safety
	appContextCopy := c.appContext.Copy()
	common.InsertOperationId(appContextCopy)
	executionCtx := core.NewExecutionContext(appContextCopy, c.SecurityContext())
	// Set default accepted response type.
	executionCtx.SetConnectionMetadata(core.ResponseTypeKey, core.OnlyMonoResponse)
	return executionCtx
}

// verify we implement Connector interface
var _ Connector = &connector{}

// <-- Connector interface implementation end

// <-- APIProvider interface implementation start

// Invoke does the actual call to the server based on the specific connector's protocol.
// By default this is the JSON-RPC protocol.
func (c *connector) Invoke(
	serviceId string,
	operationId string,
	inputValue data.DataValue,
	ctx *core.ExecutionContext) core.MethodResult {

	return c.protocol.Invoke(
		serviceId,
		operationId,
		inputValue,
		ctx)
}

// verify we implement core.APIProvider interface
var _ core.APIProvider = &connector{}

// <-- APIProvider interface implementation end
