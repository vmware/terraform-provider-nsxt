/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import (
	"crypto"
	"encoding/base64"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

//Used for signing Json request messages.
type JSONSsoSigner struct {
}

func NewJSONSsoSigner() *JSONSsoSigner {
	return &JSONSsoSigner{}
}

//
// Sign the input JSON request message.
// The message is signed using user's private key. The digest and saml
// token is then added to the security context block of the execution
// context. A timestamp is also added to guard against replay attacks
// Sample input security context:
// {
//    'schemeId': 'SAML_TOKEN',
//    'privateKey': <PRIVATE_KEY>,
//    'samlToken': <SAML_TOKEN>,
//    'signatureAlgorithm': <ALGORITHM>,
// }
//
// Security context block before signing:
// {
//    'schemeId': 'SAML_TOKEN',
//    'signatureAlgorithm': <ALGORITHM>,
//    'timestamp': {
//        'created': '2012-10-26T12:24:18.941Z',
//        'expires': '2012-10-26T12:44:18.941Z',
// }
//
//
// Security context block after signing:
// {
//    'schemeId': 'SAML_TOKEN',
//    'signatureAlgorithm': <ALGORITHM>,
//    'signature': {
//        'samlToken': <SAML_TOKEN>,
//        'value': <DIGEST>
//    }
//    'timestamp': {
//        'created': '2012-10-26T12:24:18.941Z',
//        'expires': '2012-10-26T12:44:18.941Z',
//    }
// }
func (j *JSONSsoSigner) Process(jsonRequestBody *map[string]interface{}) error {
	securityContext, err := GetSecurityContext(jsonRequestBody)
	if err != nil {
		// does not have to be propogated to higher layers.
		// it is okay for some requests to not include security context.
		return nil
	}

	newSecurityContext := map[string]interface{}{}
	if schemeId, ok := securityContext[AUTHENTICATION_SCHEME_ID]; ok {
		if schemeId != SAML_HOK_SCHEME_ID {
			return nil
		}
		newSecurityContext[AUTHENTICATION_SCHEME_ID] = SAML_HOK_SCHEME_ID
	} else {
		log.Debugf("Security Context does not contain authentication scheme ID")
	}
	newSecurityContext[TIMESTAMP] = GenerateRequestTimeStamp()

	var jwsAlgorithm string
	if _, ok := securityContext[SIGNATURE_ALGORITHM]; ok {
		if temp, ok := securityContext[SIGNATURE_ALGORITHM].(string); ok {
			jwsAlgorithm = temp
		} else {
			log.Error("Value of signature algorithm extracted from Security Context failed assertion of type string")
			return l10n.NewRuntimeErrorNoParam("vapi.security.sso.signature.invalid")
		}
	} else {
		log.Debugf("Security Context does not contain signature algorithm")
		return l10n.NewRuntimeErrorNoParam("vapi.security.sso.signature.invalid")
	}

	var algorithm crypto.Hash
	if temp, ok := algorithmMap[jwsAlgorithm]; ok {
		algorithm = temp
	} else {
		log.Error("Invalid Hash algorithm was passed")
		return l10n.NewRuntimeErrorNoParam("vapi.security.sso.hash.invalid")
	}

	newSecurityContext[SIGNATURE_ALGORITHM] = jwsAlgorithm
	err = SetSecurityContext(jsonRequestBody, newSecurityContext)
	if err != nil {
		return err
	}

	var privateKey string
	if pvtK, ok := securityContext[PRIVATE_KEY]; ok {
		if temp, ok := pvtK.(string); ok {
			privateKey = temp
		} else {
			log.Error("Value of private key extracted from Security Context failed assertion of type string")
			return l10n.NewRuntimeErrorNoParam("vapi.security.sso.pvtkey.invalid")
		}
	} else {
		log.Debugf("Security Context does not contain private key")
		return l10n.NewRuntimeErrorNoParam("vapi.security.sso.pvtkey.invalid")
	}

	rsaPrivateKey, parseErr := ParsePrivateKey(preparePrivateKey(privateKey))
	if parseErr != nil {
		return parseErr
	}

	jce := JSONCanonicalEncoder{}
	toSign, canonErr := jce.Marshal(jsonRequestBody)
	if canonErr != nil {
		return canonErr
	}

	signedBytes, err := Sign(toSign, algorithm, rsaPrivateKey)
	if err != nil {
		return err
	}
	sig := base64.StdEncoding.EncodeToString(signedBytes)
	var samlToken interface{}
	if temp, ok := securityContext[SAML_TOKEN]; ok {
		samlToken = temp
	} else {
		log.Error("Security Context does not contain saml token")
		return l10n.NewRuntimeErrorNoParam("vapi.security.sso.samltoken.invalid")
	}
	newSecurityContext[SIGNATURE] = map[string]interface{}{SAML_TOKEN: samlToken, DIGEST: sig}
	err = SetSecurityContext(jsonRequestBody, newSecurityContext)
	return err
}
