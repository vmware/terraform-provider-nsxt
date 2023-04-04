/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// JwtVapiClaims is an implementation of the interface jwt.Claims, which
// provides a NOOP jwt.Claims#Valid. The type serves only as a container for
// acquiring and preserving the claims parsed by jwt.Parser. Claims validation
// is carried out by JwtVapiClaimsValidator.
//
// Note - a JSON string cannot be deserialized into JwtVapiClaimsValidator
// Hence, the intermediate container - JwtVapiClaims map[string]interface{},
// is introduced.
type JwtVapiClaims map[string]interface{}

// Valid provides NOOP implementation of jwt.Claims interface
func (m JwtVapiClaims) Valid() error {
	return nil
}

var missingClaimErr = errors.New("Missing claim in JWT")

func getClaim(claims *JwtVapiClaims, claim string) (interface{}, error) {
	if val, ok := (*claims)[claim]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("%w - '%v'", missingClaimErr, claim)
}

func getClaimAsString(claims *JwtVapiClaims, claim string) (string, error) {
	claimVal, err := getClaim(claims, claim)
	if err != nil {
		return "", err
	}
	switch claimVal.(type) {
	case string:
		return claimVal.(string), nil
	}
	return "", fmt.Errorf("Invalid JWT '%v' claim with value '%v' of type '%T'."+
		" Expected claim of type string", claim, claimVal, claimVal)
}

func getClaimAsStringSlice(claims *JwtVapiClaims, claim string) ([]string, error) {
	claimVal, err := getClaim(claims, claim)
	if err != nil {
		return nil, err
	}
	switch claimVal.(type) {
	case []interface{}:
		strs, err := assertInterfaceSliceIsStringSlice(claimVal.([]interface{}))
		if err != nil {
			return nil, fmt.Errorf("Invalid JWT '%v' claim with value '%v' of type '%T'."+
				" Expected claim of type []string. %w", claim, claimVal, claimVal, err)
		}
		return strs, nil
	}
	return nil, fmt.Errorf("Invalid JWT '%v' claim with value '%v' of type '%T'."+
		" Expected claim of type []string", claim, claimVal, claimVal)
}

func getStringOrSliceClaimAsSlice(claims *JwtVapiClaims, claim string) ([]string, error) {
	claimVal, err := getClaim(claims, claim)
	if err != nil {
		return nil, err
	}
	str, strErr := getClaimAsString(claims, claim)
	strs, strsErr := getClaimAsStringSlice(claims, claim)
	if strErr != nil && strsErr != nil {
		return nil, fmt.Errorf("Invalid JWT '%v' claim with value '%v' of type '%T'."+
			" Expected claim of type []string or string", claim, claimVal, claimVal)
	}
	if strs != nil {
		return strs, nil
	}
	return []string{str}, nil
}

func getClaimAsInt64(claims *JwtVapiClaims, claim string) (int64, error) {
	claimVal, err := getClaim(claims, claim)
	if err != nil {
		return -1, err
	}
	switch claimValTyped := claimVal.(type) {
	case float64:
		return int64(claimValTyped), nil
	case json.Number:
		claimInt64, castErr := claimValTyped.Int64()
		if castErr != nil {
			break
		}
		return claimInt64, nil
	}
	return -1, fmt.Errorf("Invalid JWT '%v' claim with value '%v' of type '%T'."+
		" Expected claim of type int64", claim, claimVal, claimVal)
}

// JwtVapiClaimsValidator validator describing the vAPI specific JWT validation procedure.
// JwtVapiClaimsValidator contains a pointer to the JwtVapiClaims being validated.
// A separate validator object should be created per request to avoid race conditions
// during different JwtVapiClaims validation.
//
// JwtVapiClaimsValidator takes into consideration clockSkew while validating
// unix time associated claims and compares the 'aud' claim against a
// predefined set of audience.
type JwtVapiClaimsValidator struct {
	claims              *JwtVapiClaims
	maxClockSkew        int64
	acceptableAudiences []string
}

// NewJwtVapiClaimsValidator factory method creating JwtVapiClaimsValidator
func NewJwtVapiClaimsValidator(claims *JwtVapiClaims, maxClockSkew int64, acceptableAudiences []string) *JwtVapiClaimsValidator {
	return &JwtVapiClaimsValidator{
		claims:              claims,
		maxClockSkew:        maxClockSkew,
		acceptableAudiences: acceptableAudiences,
	}
}

func (v *JwtVapiClaimsValidator) verifyStringClaim(claim string) error {
	_, err := getClaimAsString(v.claims, claim)
	return err
}

func (v *JwtVapiClaimsValidator) verifyIssuedAt(now int64) error {
	iat, err := getClaimAsInt64(v.claims, CLAIM_ISSUED_AT)
	if err != nil {
		return err
	}
	if now+v.maxClockSkew > iat {
		return nil
	}
	return fmt.Errorf("JWT issued in the future. '%v' claim: %v. Current time: %v", CLAIM_ISSUED_AT, iat, now)
}

func (v *JwtVapiClaimsValidator) verifyExpiresAt(now int64) error {
	exp, err := getClaimAsInt64(v.claims, CLAIM_EXPIRES_AT)
	if err != nil {
		return err
	}
	if now < exp+v.maxClockSkew {
		return nil
	}
	return fmt.Errorf("Expired JWT. '%v' claim: %v. Current time: %v", CLAIM_EXPIRES_AT, exp, now)
}

func (v *JwtVapiClaimsValidator) verifyNotBefore(now int64) error {
	nbf, err := getClaimAsInt64(v.claims, CLAIM_NOT_BEFORE)
	if err != nil {
		if errors.Is(err, missingClaimErr) {
			return nil
		}
		return err
	}
	if now+v.maxClockSkew >= nbf {
		return nil
	}
	return fmt.Errorf("JWT before use time. '%v' claim: %v. Current time: %v", CLAIM_NOT_BEFORE, nbf, now)
}

func (v *JwtVapiClaimsValidator) verifyAudience() error {
	auds, err := getStringOrSliceClaimAsSlice(v.claims, CLAIM_AUDIENCE)
	if err != nil {
		return err
	}

	for _, audience := range auds {
		if contains(v.acceptableAudiences, audience) {
			return nil
		}
	}
	return fmt.Errorf("No acceptable audience in JWT with '%v' claim: %v", CLAIM_AUDIENCE, auds)
}

func (v *JwtVapiClaimsValidator) verifyGroups() error {
	_, err := getClaimAsStringSlice(v.claims, CLAIM_GROUP_NAMES)
	if errors.Is(err, missingClaimErr) {
		return nil
	}
	return err
}

func (v *JwtVapiClaimsValidator) Valid() error {
	now := time.Now().Unix()
	var err error
	if err = v.verifyStringClaim(CLAIM_ISSUER); err != nil {
		return err
	}
	if err = v.verifyStringClaim(CLAIM_SUBJECT); err != nil {
		return err
	}
	if err = v.verifyIssuedAt(now); err != nil {
		return err
	}
	if err = v.verifyExpiresAt(now); err != nil {
		return err
	}
	if err = v.verifyNotBefore(now); err != nil {
		return err
	}
	if err = v.verifyAudience(); err != nil {
		return err
	}
	if err = v.verifyGroups(); err != nil {
		return err
	}
	return err
}
