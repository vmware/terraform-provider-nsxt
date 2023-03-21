/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package http

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/http/internal"
	"io"
)

type VAPIFrameReader interface {
	ReadFrame() ([]byte, error)
}

// NewVapiFrameReader returns a new VAPIFrameReader that translates vAPI frames returned from http
// chunked format
func NewVapiFrameReader(r io.Reader) VAPIFrameReader {
	return internal.NewVapiFrameReader(r)
}
