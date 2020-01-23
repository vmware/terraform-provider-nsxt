/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


const ROUTING_PATH_SEPARATOR = ","

type RoutingInfo struct {
	routingStrategy string
	routingPath     string
	routingPaths    []string
	operationHints  []string
	idTypes         map[string]string
}

func NewRoutingInfo() *RoutingInfo {
	routingPaths := []string{}
	operationHints := []string{}
	idTypes := make(map[string]string)

	routingInfo := RoutingInfo{routingPaths: routingPaths, operationHints: operationHints, idTypes: idTypes}
	return &routingInfo
}

// Routing Strategy
func (r *RoutingInfo) RoutingStrategy() string {
	return r.routingStrategy
}

func (r *RoutingInfo) SetRoutingStrategy(routingStrategy string) {
	r.routingStrategy = routingStrategy
}

// Routing Path
func (r *RoutingInfo) RoutingPath() string {
	return r.routingPath
}

func (r *RoutingInfo) SetRoutingPath(routingPath string) {
	r.routingPath = routingPath
}

// Routing Paths
func (r *RoutingInfo) RoutingPaths() []string {
	return r.routingPaths
}

func (r *RoutingInfo) SetRoutingPaths(routingPaths []string) {
	r.routingPaths = routingPaths
}

func (r *RoutingInfo) AddToRoutingPaths(routingPath string) {
	if r.routingPaths == nil {
		r.routingPaths = []string{routingPath}
	} else {
		r.routingPaths = append(r.routingPaths, routingPath)
	}
}

// Operation Hints
func (r *RoutingInfo) OperationHints() []string {
	return r.operationHints
}

func (r *RoutingInfo) SetOperationHints(operationHints []string) {
	r.operationHints = operationHints
}

func (r *RoutingInfo) AddToOperationHints(operationHint string) {
	if r.operationHints == nil {
		r.operationHints = []string{operationHint}
	} else {
		r.operationHints = append(r.operationHints, operationHint)
	}
}

// Id Types
func (r *RoutingInfo) IdTypes(id string) string {
	return r.idTypes[id]
}

func (r *RoutingInfo) SetIdTypes(id string, idTypes string) {
	r.idTypes[id] = idTypes
}
