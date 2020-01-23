/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import "github.com/vmware/vsphere-automation-sdk-go/runtime/data"

type PrivilegeProvider interface {
	GetPrivilegeInfo(fullyQualifiedOperName string, inputValue data.DataValue) (map[ResourceIdentifier][]string, error)
}
