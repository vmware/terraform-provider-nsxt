/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security


import (
	"errors"
)

// Verifies subject confirmation data of token matches bearer token.
type BearerTokenSCVerifier struct {
}

func NewBearerTokenSCVerifier() *BearerTokenSCVerifier {
	return &BearerTokenSCVerifier{}
}

// verify if the json payload that is passed contains valid bearer token and corresponding scheme ID
func (j *BearerTokenSCVerifier) Process(jsonMessage *map[string]interface{}) error {

	securityContext, err := GetSecurityContext(jsonMessage)
	if err != nil {
		// does not have to propogated to higher layers.
		// it is okay for some requests to not include security context.
		return nil
	}

	if securityContext[AUTHENTICATION_SCHEME_ID] != SAML_BEARER_SCHEME_ID {
		return nil
	}

	if !isSamlBearerToken(securityContext[SAML_TOKEN].(string)) {
		return errors.New("SAML token passed is not of type Bearer")
	}

	return nil
}
