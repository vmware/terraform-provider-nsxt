// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
)

var cliVpcLbAppProfilesClient = vpcs.NewVpcLbAppProfilesClient
var cliVpcLbMonitorProfilesClient = vpcs.NewVpcLbMonitorProfilesClient
var cliVpcLbPersistenceProfilesClient = vpcs.NewVpcLbPersistenceProfilesClient
var cliVpcLbClientSslProfilesClient = vpcs.NewVpcLbClientSslProfilesClient
var cliVpcLbServerSslProfilesClient = vpcs.NewVpcLbServerSslProfilesClient
var cliVpcLbPoolsClient = vpcs.NewVpcLbPoolsClient
var cliVpcLbVirtualServersClient = vpcs.NewVpcLbVirtualServersClient
