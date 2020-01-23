/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package metadata


import "github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"

func getError(msg string) error {
	args := map[string]string{
		"msg": msg,
	}
	return l10n.NewRuntimeError("vapi.metadata.parser.failure", args)
}
