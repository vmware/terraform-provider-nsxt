/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package http

import (
	"context"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"net/http"
)

// ClientResponseHandler provides contract for handling http responses
// Connector can have many response handlers but only one can handle http response.
// To skip handling of response return nil core.MethodResult and nil for error
type ClientResponseHandler interface {
	// HandleResponse either handles http response or skips it by returing nil, nil
	HandleResponse(ctx context.Context, response *http.Response) (core.MethodResult, error)
}

// ClientFramesResponseHandler provides contract for streamed response handlers
type ClientFramesResponseHandler interface {
	ClientResponseHandler
	// SetClientFrameDeserializer is used to override default ClientFrameDeserializer on a frames response handler
	SetClientFrameDeserializer(deserializer ClientFrameDeserializer)
}

// ClientFrameDeserializer takes care of deserialization of raw frame data into runtime
// specific core.MethodResult objects
type ClientFrameDeserializer interface {
	// DeserializeFrames method translates channel of []byte data into []core.MethodResult.
	// Implementations of this method should take care of creating, returning, and closing core.MethodResult channel.
	// Usually this happens by initiating a go routine and do the actual deserialization in it so that bindings calls
	// do not get blocked.
	// If returned channel data is not needed for some reason return closed empty channel instead of nil.
	DeserializeFrames(ctx context.Context, frames chan []byte) (chan core.MonoResult, error)
}
