/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package errors

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
)

var ERROR_BINDINGS_MAP = map[string]bindings.BindingType{
	AlreadyExists{}.Error():               AlreadyExistsBindingType(),
	AlreadyInDesiredState{}.Error():       AlreadyInDesiredStateBindingType(),
	Canceled{}.Error():                    CanceledBindingType(),
	ConcurrentChange{}.Error():            ConcurrentChangeBindingType(),
	Error{}.Error():                       ErrorBindingType(),
	FeatureInUse{}.Error():                FeatureInUseBindingType(),
	InternalServerError{}.Error():         InternalServerErrorBindingType(),
	InvalidArgument{}.Error():             InvalidArgumentBindingType(),
	InvalidElementConfiguration{}.Error(): InvalidElementConfigurationBindingType(),
	InvalidElementType{}.Error():          InvalidElementTypeBindingType(),
	InvalidRequest{}.Error():              InvalidRequestBindingType(),
	NotFound{}.Error():                    NotFoundBindingType(),
	NotAllowedInCurrentState{}.Error():    NotAllowedInCurrentStateBindingType(),
	OperationNotFound{}.Error():           OperationNotFoundBindingType(),
	ResourceBusy{}.Error():                ResourceBusyBindingType(),
	ResourceInUse{}.Error():               ResourceInUseBindingType(),
	ResourceInaccessible{}.Error():        ResourceInaccessibleBindingType(),
	ServiceUnavailable{}.Error():          ServiceUnavailableBindingType(),
	TimedOut{}.Error():                    TimedOutBindingType(),
	UnableToAllocateResource{}.Error():    UnableToAllocateResourceBindingType(),
	Unauthenticated{}.Error():             UnauthenticatedBindingType(),
	Unauthorized{}.Error():                UnauthorizedBindingType(),
	UnexpectedInput{}.Error():             UnexpectedInputBindingType(),
	Unsupported{}.Error():                 UnsupportedBindingType(),
	UnverifiedPeer{}.Error():              UnverifiedPeerBindingType(),
}
