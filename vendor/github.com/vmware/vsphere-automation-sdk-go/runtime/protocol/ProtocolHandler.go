/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package protocol

/**
 * ProtocolHanders are the classes that provide the endpoint for a given
 * protocol.
 */
type ProtocolHandler interface {
	Start()
	Stop()
}
