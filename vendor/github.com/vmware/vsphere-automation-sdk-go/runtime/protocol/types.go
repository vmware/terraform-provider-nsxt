/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package protocol

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
)

// RestClientOptions provides contract for managing rest connector options
// when vAPI connector is used in REST protocol. Set by client.UsingRest() client.ConnectorOption
type RestClientOptions interface {
	// SecurityContextSerializers list of serializers which provide mappings between core.SecurityContext object and http headers
	SecurityContextSerializers() map[string]SecurityContextSerializer
	// EnableDefaultContentType when true overrides Content-Type header to value 'application/json'
	EnableDefaultContentType() bool
}

// SecurityContextSerializer serializes core.SecurityContext object into http headers
// implemented by concrete context serializers, such as rest.UserPwdSecContextSerializer,
// rest.SessionSecContextSerializer and rest.OauthSecContextSerializer.
// Clients can also implement a custom serializer with special serialization requirements.
type SecurityContextSerializer interface {
	// Serialize provides http headers map to serialize core.SecurityContext object into authorization headers
	Serialize(core.SecurityContext) (map[string]interface{}, error)
}
