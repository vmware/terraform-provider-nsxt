/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

var algorithmMap = map[string]crypto.Hash{RS256: crypto.SHA256, RS384: crypto.SHA384, RS512: crypto.SHA512}

// Used to verify the authenticity of the request
// message by verifying the digest present in the security context block.
type JSONSsoVerifier struct {
}

func NewJSONSsoVerifier() *JSONSsoVerifier {
	return &JSONSsoVerifier{}
}

//Verify the input JSON message.
//
//       For verification, we need 4 things:
//
//       1. algorithm: extracted from security context
//       2. certificate: public key of the principal embedded in the
//       SAML token is used
//       3. digest: value field from signature block
//       4. canonical msg: signature block is removed from the request
//       and the remaining part is canonicalized
//
//       Sample input security context:
//       {
//           'schemeId': 'SAML_TOKEN',
//           'signatureAlgorithm': <ALGORITHM>,
//           'signature': {
//               'samlToken': <SAML_TOKEN>,
//               'value': <DIGEST>
//           }
//           'timestamp': {
//               'created': '2012-10-26T12:24:18.941Z',
//               'expires': '2012-10-26T12:44:18.941Z',
//           }
//       }

func (j *JSONSsoVerifier) Process(jsonMessage *map[string]interface{}) error {
	securityContext, err := GetSecurityContext(jsonMessage)
	if err != nil {
		// does not have to propogated to higher layers.
		// it is okay for some requests to not include security context.
		return nil
	}
	if schemeId, ok := securityContext[AUTHENTICATION_SCHEME_ID]; ok {
		if schemeId != SAML_HOK_SCHEME_ID {
			return nil
		}
	} else {
		log.Debugf("Security Context does not contain authentication scheme ID")
	}

	var signature map[string]interface{}
	if sig, ok := securityContext[SIGNATURE]; ok {
		if temp, ok := sig.(map[string]interface{}); ok {
			signature = temp
		} else {
			log.Error("Value of signature extracted from Security Context failed assertion of type map[string]interface")
			return l10n.NewRuntimeErrorNoParam("vapi.security.sso.signature.invalid")
		}
	} else {
		log.Error("Security Context does not contain signature")
		return l10n.NewRuntimeErrorNoParam("vapi.security.sso.signature.invalid")
	}
	delete(securityContext, SIGNATURE)

	var digest string
	if dig, ok := signature[DIGEST]; ok {
		if temp, ok := dig.(string); ok {
			digest = temp
		} else {
			log.Error("Value of digest extracted from signature failed assertion of type string")
			return l10n.NewRuntimeErrorNoParam("vapi.security.sso.digest.invalid")
		}
	} else {
		log.Error("Signature (Map data structure) does not contain digest")
		return l10n.NewRuntimeErrorNoParam("vapi.security.sso.digest.invalid")
	}
	base64DecodedBytes, decodeErr := base64.StdEncoding.DecodeString(digest)
	if decodeErr != nil {
		return decodeErr
	}

	var jwsAlgorithm string
	if _, ok := securityContext[SIGNATURE_ALGORITHM]; ok {
		if temp, ok := securityContext[SIGNATURE_ALGORITHM].(string); ok {
			jwsAlgorithm = temp
		} else {
			log.Error("Value of signature algorithm extracted from Security Context failed assertion of type string")
			return l10n.NewRuntimeErrorNoParam("vapi.security.sso.signature.algorithm.invalid")
		}
	} else {
		log.Debugf("Security Context does not contain signature algorithm")
		return l10n.NewRuntimeErrorNoParam("vapi.security.sso.signature.algorithm.invalid")
	}

	var algo crypto.Hash
	if temp, ok := algorithmMap[jwsAlgorithm]; ok {
		algo = temp
	} else {
		log.Error("Invalid Hash algorithm")
		return l10n.NewRuntimeErrorNoParam("vapi.security.sso.hash.invalid")
	}

	var samlToken string
	if _, ok := signature[SAML_TOKEN]; ok {
		if temp, ok := signature[SAML_TOKEN].(string); ok {
			samlToken = temp
		} else {
			log.Error("Value of saml token extracted from signature failed assertion of type string")
			return l10n.NewRuntimeErrorNoParam("vapi.security.sso.samltoken.invalid")
		}
	} else {
		log.Error("Signature (Map data structure) does not contain saml token")
		return l10n.NewRuntimeErrorNoParam("vapi.security.sso.samltoken.invalid")
	}
	certificate, exErr := ExtractCertificate(samlToken)
	if exErr != nil {
		return exErr
	}
	jce := NewJSONCanonicalEncoder()
	canonicalizedMsg, canError := jce.Marshal(*jsonMessage)
	if canError != nil {
		return canError
	}

	var pubKey *rsa.PublicKey
	if temp, ok := certificate.PublicKey.(*rsa.PublicKey); ok {
		pubKey = temp
	} else {
		log.Error("Public key in certificate failed type assertion")
		return l10n.NewRuntimeErrorNoParam("vapi.security.sso.pubkey.invalid")
	}
	err = VerifySignature(pubKey, algo, canonicalizedMsg, base64DecodedBytes)
	if err == nil {
		securityContext[AUTHENTICATED] = true
	}
	securityContext[SAML_TOKEN] = samlToken
	securityContext[SIGNATURE_ALGORITHM] = jwsAlgorithm
	return err

}
