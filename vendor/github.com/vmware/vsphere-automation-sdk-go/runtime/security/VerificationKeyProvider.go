/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

// VerificationKeyProvider provides keys used to validate the authenticity of the JWT token.
type VerificationKeyProvider interface {

	// Retrieves keys from issuer.
	// returns
	//  ([]interface{}, nil) keys utilized for JWT verification
	//  (nil, error) for invalid retrieval
	Get(issuer string) ([]interface{}, error)
}
