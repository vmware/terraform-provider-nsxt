/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */
package security

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type OidcDiscoveryDoc struct {
	JwksUri string `json:"jwks_uri"`
	Issuer  string
}

type JwksDoc struct {
	Keys []Jwk
}

type Jwk struct {
	Kid string
	X5c []string
	Alg string
}

func ParseOidcConfig(input []byte) (*OidcDiscoveryDoc, error) {
	var doc OidcDiscoveryDoc
	err := json.Unmarshal(input, &doc)
	return &doc, err
}

func ParseJwks(input []byte) ([]interface{}, error) {
	var jwksDoc JwksDoc
	err := json.Unmarshal(input, &jwksDoc)
	if err != nil {
		return nil, err
	}
	keys := make([]interface{}, len(jwksDoc.Keys))
	for i, jwk := range jwksDoc.Keys {
		if len(jwk.X5c) == 0 {
			return nil, fmt.Errorf("unable to parse JWKS - missing certificate")
		}
		// The PKIX certificate containing the key value MUST be the first certificate
		// see https://datatracker.ietf.org/doc/html/rfc7517#section-4.7
		cert, err := decodeCertificate(jwk.X5c[0])
		if err != nil {
			return nil, err
		}
		keys[i] = (*cert).PublicKey
	}

	return keys, err
}

func decodeCertificate(encodedCert string) (*x509.Certificate, error) {
	decodedCert, err := base64.StdEncoding.DecodeString(encodedCert)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(decodedCert)
}
