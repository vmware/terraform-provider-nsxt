/* Copyright © 2019-2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

const AUTHENTICATION_SCHEME_ID = "schemeId"
const SESSION_SCHEME_ID = "com.vmware.vapi.std.security.session_id"

const SESSION_ID = "sessionId"
const NO_AUTH = "com.vmware.vapi.std.security.no_authentication"
const AUTHN_IDENTITY = "authnIdentity"
const SAML_BEARER_SCHEME_ID = "com.vmware.vapi.std.security.saml_bearer_token"

const SAML_BEARER_TYPE = "urn:oasis:names:tc:SAML:2.0:cm:bearer"
const SAML_HOK_TYPE = "urn:oasis:names:tc:SAML:2.0:cm:holder-of-key"
const SAML_HOK_TOKEN = "SAML_HOK_TOKEN"
const SAML_BEARER_TOKEN = "SAML_BEARER_TOKEN"

const USER_PASSWORD_SCHEME_ID = "com.vmware.vapi.std.security.user_pass"
const USER_KEY = "userName"
const PASSWORD_KEY = "password"
const SAML_HOK_SCHEME_ID = "com.vmware.vapi.std.security.saml_hok_token"
const SESSION_ID_KEY = "vmware-api-session-id"
const PRIVATE_KEY = "privateKey"
const SAML_TOKEN = "samlToken"
const SIGNATURE_ALGORITHM = "signatureAlgorithm"
const SIGNATURE = "signature"
const REQUEST_VALIDITY = 20 //in minutes
const TIMESTAMP = "timestamp"
const DIGEST = "value"
const AUTHENTICATED = "requestAuthenticated"

const OAUTH_SCHEME_ID = "com.vmware.vapi.std.security.oauth"
const CSP_AUTH_TOKEN_KEY = "csp-auth-token"
const ACCESS_TOKEN = "accessToken"

const RS256 = "RS256"
const RS384 = "RS384"
const RS512 = "RS512"

const TS_EXPIRES_KEY = "expires"
const TS_CREATED_KEY = "created"
