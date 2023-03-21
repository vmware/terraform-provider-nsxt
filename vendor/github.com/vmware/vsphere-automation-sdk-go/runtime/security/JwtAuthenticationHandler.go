/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"strings"
	"time"
)

type JwtAuthenticationHandler struct {
	keyProvider         VerificationKeyProvider
	maxClockSkew        int64
	acceptableAudiences []string
}

// NewJwtAuthenticationHandler factory method creating JwtAuthenticationHandler.
// If the supplied VerificationKeyProvider is also a VerificationKeyCache,
// the created JwtAuthenticationHandler is able to handle change of verification
// keys - if the signature of a JWT cannot be verified with the keys returned by
// VerificationKeyCache#Get, will refresh the key provider via
// VerificationKeyCache#Refresh and will retry signature verification. Retry is done only
// once.
// parameters
//  keyProvider  VerificationKeyProvider used for retrieving signing keys during
//               JWT validation. Must not be nil.
//  configOptions  a set of optional JwtHandlerConfigOption
// returns
//  By default the handler is created with a:
//    - maximum clock skew value of 10 minutes
//    - acceptable audiences containing a single element - "vmware-tes:vapi"
//  values above are configurable via the configOptions parameter.
func NewJwtAuthenticationHandler(keyProvider VerificationKeyProvider, configOptions ...JwtHandlerConfigOption) (*JwtAuthenticationHandler, error) {
	if keyProvider == nil {
		return nil, fmt.Errorf("key provider must not be nil")
	}

	jwtHandler := &JwtAuthenticationHandler{
		keyProvider:         keyProvider,
		maxClockSkew:        10 * 60,
		acceptableAudiences: []string{"vmware-tes:vapi"},
	}

	for _, fn := range configOptions {
		err := fn(jwtHandler)
		if err != nil {
			return nil, err
		}
	}

	return jwtHandler, nil
}

func (j *JwtAuthenticationHandler) Authenticate(ctx core.SecurityContext) (*UserIdentity, error) {
	if ctx.Property(AUTHENTICATION_SCHEME_ID) != j.SupportedScheme() {
		//Returning nil as error so the AuthenticationFilter can continue and try
		//with the next auth handler if one exists
		return nil, nil
	}

	tokenString, err := acquireTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	token, parts, err := parseUnverified(tokenString)
	if err != nil {
		return nil, err
	}

	// Assertion is safe as we are underlining *JwtVapiClaims in our facade parsing
	claims := token.Claims.(*JwtVapiClaims)
	err = j.validateClaims(claims)
	if err != nil {
		log.Debugf("Claims validation failed - %v", err)
		return nil, err
	}

	// Already validated
	iss, _ := getClaimAsString(claims, CLAIM_ISSUER)

	err = j.validateSignature(token, parts, iss)

	// if retry is configured and signature validation failed
	// attempt a single retry post VerificationKeyCache#Refresh
	if keyCache, ok := j.keyProvider.(VerificationKeyCache); ok && err != nil {
		refreshErr := keyCache.Refresh(iss)
		if refreshErr != nil {
			log.Debugf("Error during refresh of key cache - %v", refreshErr)
			return nil, fmt.Errorf("Failed to refresh key cache - %v", refreshErr)
		}
		err = j.validateSignature(token, parts, iss)
	}

	if err == nil {
		user := createIdentityFromClaims(claims)
		log.Debugf("Authenticated user with username: %v and domain: %v", (*user).userName, *(*user).domain)
		return user, nil
	}

	log.Debugf("Authentication failed - %v", err)
	return nil, err
}

// Validates the passed claims via a created JwtVapiClaimsValidator
// returns
//  nil for successful validation
//  error otherwise
func (j *JwtAuthenticationHandler) validateClaims(claims *JwtVapiClaims) error {
	claimsValidator := *NewJwtVapiClaimsValidator(claims, j.maxClockSkew, j.acceptableAudiences)
	return claimsValidator.Valid()
}

// Validates the signature of a parsed JWT token against keys fetched by the keyProvider.
// validateSignatureAgainstKey is utilized for validations against individual keys
// returns
//  nil for successful validation against any of the keys
//  error if no key matches the signature or the keyProvider fails to retrieve them
func (j *JwtAuthenticationHandler) validateSignature(token *jwt.Token, tokenParts []string, issuer string) error {
	keys, err := j.keyProvider.Get(issuer)
	if err != nil {
		log.Errorf("Failed to retrieve verification keys - %v", err)
		return err
	}

	for _, key := range keys {
		err := validateSignatureAgainstKey(token, tokenParts, key)
		if err == nil {
			log.Debug("Signature successfully validated")
			return nil
		}
	}

	log.Debug("JWT signature validation failed - no matching key")
	return fmt.Errorf("JWT signature validation failed")
}

func (*JwtAuthenticationHandler) SupportedScheme() string {
	return OAUTH_SCHEME_ID
}

func acquireTokenFromContext(ctx core.SecurityContext) (string, error) {
	token := ctx.Property(ACCESS_TOKEN)
	if token == nil {
		return "", fmt.Errorf("Missing JWT")
	}
	if tokenString, ok := token.(string); ok {
		return tokenString, nil
	}
	return "", fmt.Errorf("Malformed JWT")
}

// Facade around jwt.Parser.ParseUnverified used to create a jwt.Token
// object with parsed headers, JwtVapiClaims, signature and underlined jwt.SigningMethod
// via which the signature is to be verified. The jwt.SigningMethod is derived from
// the 'alg' header.
// Along with the parsed token object an array of strings representing the '.'
// separated token is returned.
// returns
//  (*jwt.Token, []string, nil) for successful parsing
//  (nil, nil, error) otherwise
func parseUnverified(tokenString string) (*jwt.Token, []string, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	claims := &JwtVapiClaims{}
	token, parts, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, nil, err
	}
	// populate the token signature as in jwt.Parser.ParseWithClaims; parts length is validated in ParseUnverified
	token.Signature = parts[2]
	return token, parts, nil
}

// Validates the signature of a parsed JWT token via jwt.SigningMethod.Verify,
// where the token's specific jwt.SigningMethod is derived from the 'alg' header
// returns
//  nil for successful validation
//  error otherwise
func validateSignatureAgainstKey(token *jwt.Token, tokenParts []string, key interface{}) error {
	// jwt.SigningMethod.Verify requires signing string and signature as separate inputs
	return token.Method.Verify(strings.Join(tokenParts[0:2], "."), token.Signature, key)
}

// Creates a UserIdentity object from JWT claims. Username and domain are
// derived form the `sub` claim - the `@` symbol acts as a separator between
// username and domain, example - 'sub': "Admin@vshpere" is converted to
// 'username': "Admin", 'domain': "vsphere". UserIdentity.groups are corresponding
// to the `group_names` claim.
func createIdentityFromClaims(claims *JwtVapiClaims) *UserIdentity {
	userDomainTuple := strings.SplitN((*claims)[CLAIM_SUBJECT].(string), "@", 2)
	var domain = ""
	if len(userDomainTuple) == 2 {
		domain = userDomainTuple[1]
	}
	userIdentity := UserIdentity{userName: userDomainTuple[0], domain: &domain}
	groups, _ := getClaimAsStringSlice(claims, CLAIM_GROUP_NAMES)
	userIdentity.groups = groups
	return &userIdentity
}

type JwtHandlerConfigOption func(*JwtAuthenticationHandler) error

// WithMaxClockSkew specifies the allowed time discrepancy between the client and the server.
// Fractions of a second are discarded. Must not be negative.
func WithMaxClockSkew(maxClockSkew time.Duration) JwtHandlerConfigOption {
	return func(cache *JwtAuthenticationHandler) error {
		if maxClockSkew < 0 {
			return fmt.Errorf("max clock skew must not be a negative duration")
		}
		cache.maxClockSkew = int64(maxClockSkew.Seconds())
		return nil
	}
}

// WithAcceptableAudiences specifies acceptable values for 'aud' claim in JWTs. A JWT must
// contain at least one of these values in its 'aud' claim. Must not be empty.
func WithAcceptableAudiences(acceptableAudiences []string) JwtHandlerConfigOption {
	return func(cache *JwtAuthenticationHandler) error {
		if len(acceptableAudiences) == 0 {
			return fmt.Errorf("acceptable audiences must contain at least one value")
		}
		cache.acceptableAudiences = acceptableAudiences
		return nil
	}
}
