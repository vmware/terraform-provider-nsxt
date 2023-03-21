/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

// VerificationKeyCache which caches keys and needs to be refreshed if the signing keys change.
type VerificationKeyCache interface {
	VerificationKeyProvider

	// Refresh of the provider (if the underlying implementation maintains a cache)
	// If the refresh is successful, all subsequent calls to Get must
	// return the refreshed keys.
	Refresh(issuer string) error
}
