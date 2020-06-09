/* Copyright © 2019-2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

type PermissionValidator interface {
	Validate(username string, groups []string, requiredPrivileges map[ResourceIdentifier][]string) (bool, error)
}
