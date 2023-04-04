/* Copyright © 2019, 2021-2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import (
	"encoding/json"
)

// SessionSecurityContext represents the security context needed for authentication using session ID.
type SessionSecurityContext struct {
	properties map[string]interface{}
}

func NewSessionSecurityContext(sessionID string) *SessionSecurityContext {
	properties := map[string]interface{}{}
	properties[SESSION_ID] = sessionID
	properties[AUTHENTICATION_SCHEME_ID] = SESSION_SCHEME_ID
	return &SessionSecurityContext{properties: properties}
}

func (s *SessionSecurityContext) Property(key string) interface{} {
	return s.properties[key]
}

func (s *SessionSecurityContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.properties)
}

func (s *SessionSecurityContext) GetAllProperties() map[string]interface{} {
	return s.properties
}

func (s *SessionSecurityContext) SetProperty(property string, value interface{}) {
	s.properties[property] = value
}

//UserPasswordSecurityContext represents a security context suitable for user/password authentication.
type UserPasswordSecurityContext struct {
	properties map[string]interface{}
}

func NewUserPasswordSecurityContext(username string, password string) *UserPasswordSecurityContext {
	properties := map[string]interface{}{}
	properties[USER_KEY] = username
	properties[PASSWORD_KEY] = password
	properties[AUTHENTICATION_SCHEME_ID] = USER_PASSWORD_SCHEME_ID
	return &UserPasswordSecurityContext{properties: properties}
}

func (u *UserPasswordSecurityContext) Property(key string) interface{} {
	return u.properties[key]
}

func (u *UserPasswordSecurityContext) User() string {
	return u.properties[USER_KEY].(string)
}

func (u *UserPasswordSecurityContext) Password() string {
	return u.properties[PASSWORD_KEY].(string)
}

func (u *UserPasswordSecurityContext) GetAllProperties() map[string]interface{} {
	return u.properties
}

func (u *UserPasswordSecurityContext) SetProperty(key string, value interface{}) {
	u.properties[key] = value
}

func (u *UserPasswordSecurityContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.properties)
}

// OauthSecurityContext represents a security context suitable for oauth authentication.
type OauthSecurityContext struct {
	properties map[string]interface{}
}

func NewOauthSecurityContext(accessToken string) *OauthSecurityContext {
	properties := map[string]interface{}{}
	properties[AUTHENTICATION_SCHEME_ID] = OAUTH_SCHEME_ID
	properties[ACCESS_TOKEN] = accessToken
	return &OauthSecurityContext{properties: properties}
}

func (o *OauthSecurityContext) Property(key string) interface{} {
	return o.properties[key]
}

func (o *OauthSecurityContext) Token() string {
	return o.properties[ACCESS_TOKEN].(string)
}

func (o *OauthSecurityContext) GetAllProperties() map[string]interface{} {
	return o.properties
}

func (o *OauthSecurityContext) SetProperty(key string, value interface{}) {
	o.properties[key] = value
}

func (o *OauthSecurityContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.properties)
}

// SAMLSecurityContext represents a security context for SAML tokens.
type SAMLSecurityContext struct {
	properties map[string]interface{}
}

// NewSAMLSecurityContext creates SAML security context to sign request with provided token and private key
// 	example:
//		connector := client.NewConnector(
//		"",
//		client.WithSecurityContext(NewSAMLSecurityContext(
//			"myToken",
//			"myKey",
//			security.RS256)))
//		client := NewSampleClient(connector)
//		client.MyOperation()
func NewSAMLSecurityContext(token, privateKey, signAlgorithm string) *SAMLSecurityContext {
	properties := map[string]interface{}{}
	properties[AUTHENTICATION_SCHEME_ID] = SAML_HOK_SCHEME_ID
	properties[SAML_TOKEN] = token
	properties[PRIVATE_KEY] = privateKey
	properties[SIGNATURE_ALGORITHM] = signAlgorithm
	return &SAMLSecurityContext{properties: properties}
}

func (o *SAMLSecurityContext) Property(key string) interface{} {
	return o.properties[key]
}

func (o *SAMLSecurityContext) Token() string {
	return o.properties[ACCESS_TOKEN].(string)
}

func (o *SAMLSecurityContext) GetAllProperties() map[string]interface{} {
	return o.properties
}

func (o *SAMLSecurityContext) SetProperty(key string, value interface{}) {
	o.properties[key] = value
}

func (o *SAMLSecurityContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.properties)
}
