/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import "github.com/vmware/vsphere-automation-sdk-go/runtime/core"

// The AuthenticationHandler interface is used to verify the authentication
// data provided in the security context against an identity source.
type AuthenticationHandler interface {

	// Verifies the provided authentication data against the relevant identity
	// source.
	// returns
	//  (UserIdentity, nil) on successful authentication
	//  (nil, error) for failed authentication
	//  (nil, nil) for invalid auth handler
	Authenticate(ctx core.SecurityContext) (*UserIdentity, error)

	// Get the scheme supported by this handler
	SupportedScheme() string
}
