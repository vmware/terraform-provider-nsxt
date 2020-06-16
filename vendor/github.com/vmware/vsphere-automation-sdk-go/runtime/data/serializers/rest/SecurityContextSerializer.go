/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package rest

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
)

// SecurityContextSerializer is implemented by concrete
// context serializers, such as UserPwdSecContextSerializer, SessionSecContextSerializer
// and OauthSecContextSerializer. Clients can also implement a custom serializer
// with special serialization requirements.
type SecurityContextSerializer interface {
	Serialize(core.SecurityContext) (map[string]interface{}, error)
}

type UserPwdSecContextSerializer struct {
}

func NewUserPwdSecContextSerializer() *UserPwdSecContextSerializer {
	return &UserPwdSecContextSerializer{}
}

// Serialize authorization headers for username security context
func (u *UserPwdSecContextSerializer) Serialize(ctx core.SecurityContext) (map[string]interface{}, error) {
	username, err := GetSecurityCtxStrValue(ctx, security.USER_KEY)
	if err != nil {
		return nil, err
	}

	password, err := GetSecurityCtxStrValue(ctx, security.PASSWORD_KEY)
	if err != nil {
		return nil, err
	}

	if username == nil || password == nil {
		err := errors.New("Username and password are required for UserPwdSecContextSerializer")
		return nil, err
	}

	credentialString := fmt.Sprintf("%s:%s", *username, *password)
	base64EncodedVal := base64.StdEncoding.EncodeToString([]byte(credentialString))
	return map[string]interface{}{"Authorization": "Basic " + base64EncodedVal}, nil
}

type SessionSecContextSerializer struct {
}

func NewSessionSecContextSerializer() *SessionSecContextSerializer {
	return &SessionSecContextSerializer{}
}

// Serialize authorization headers for session security context
func (s *SessionSecContextSerializer) Serialize(ctx core.SecurityContext) (map[string]interface{}, error) {
	sessionID, err := GetSecurityCtxStrValue(ctx, security.SESSION_ID)
	if err != nil {
		return nil, err
	}

	if sessionID == nil {
		err := errors.New("Session ID is required for SessionSecContextSerializer")
		return nil, err
	}

	return map[string]interface{}{security.SESSION_ID_KEY: *sessionID}, nil
}

type OauthSecContextSerializer struct {
}

func NewOauthSecContextSerializer() *OauthSecContextSerializer {
	return &OauthSecContextSerializer{}
}

// Serialize authorization headers for oauth security context.
func (o *OauthSecContextSerializer) Serialize(ctx core.SecurityContext) (map[string]interface{}, error) {
	oauthToken, err := GetSecurityCtxStrValue(ctx, security.ACCESS_TOKEN)
	if err != nil {
		return nil, err
	}

	if oauthToken == nil {
		err := errors.New("Oauth token is required for OauthSecContextSerializer")
		return nil, err
	}

	return map[string]interface{}{security.CSP_AUTH_TOKEN_KEY: *oauthToken}, nil
}
