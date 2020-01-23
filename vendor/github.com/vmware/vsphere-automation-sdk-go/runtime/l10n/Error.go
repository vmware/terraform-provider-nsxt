/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package l10n

import "github.com/vmware/vsphere-automation-sdk-go/runtime/log"

type Error struct {
	id             string
	defaultMessage string
	args           map[string]string
}

func NewError(id string, defaultMessage string, args map[string]string) *Error {
	return &Error{id: id, defaultMessage: defaultMessage, args: args}
}

//Creates a Error by looking for id in the runtime message bundle
func NewRuntimeError(id string, args map[string]string) *Error {
	if args == nil {
		args = map[string]string{}
	}
	runtimeMessages := NewDefaultRuntimeMessageFormatter()
	msg := runtimeMessages.GetLocalizedMessage(id, args)
	log.Error(msg)
	return &Error{id: id, defaultMessage: msg, args: args}
}

func NewRuntimeErrorNoParam(id string) *Error {
	runtimeMessages := NewDefaultRuntimeMessageFormatter()
	args := map[string]string{}
	msg := runtimeMessages.GetLocalizedMessage(id, args)
	log.Error(msg)
	return &Error{id: id, defaultMessage: msg, args: args}
}

func (err *Error) ID() string {
	return err.id
}

func (err *Error) Error() string {
	return err.defaultMessage
}

func (err *Error) Args() map[string]string {
	return err.args
}
