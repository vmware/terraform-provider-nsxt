// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

// VPC load balancer objects (/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-*) share
// their schema and CRUD logic with the matching infra load balancer objects
// (/infra/lb-*). The VPC flavor is derived from the infra resource/data source by adding
// the mandatory context (project_id + vpc_id), which makes the shared code select the
// VPC API clients (see lb_vpc_clients.go), and by gating on the minimal NSX version.

const vpcLbMinNsxVersion = "9.0.2"

var vpcLbPathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/%s/[id]"

// isLbCacheEnabledForRead disables read-through caching for VPC contexts, as the LB
// resource cache is populated from the infra API only.
func isLbCacheEnabledForRead(d *schema.ResourceData, m interface{}) bool {
	return isCacheEnabledForRead(d, m) && getSessionContext(d, m).ClientType != utl.VPC
}

func vpcLbSchema(src map[string]*schema.Schema) map[string]*schema.Schema {
	res := make(map[string]*schema.Schema, len(src)+1)
	for k, v := range src {
		res[k] = v
	}
	res["context"] = getContextSchemaExtended(true, false, true, true)
	return res
}

// vpcLbResource derives a VPC resource from its infra counterpart. apiSegment is the
// VPC API path segment (e.g. vpc-lb-pools), used for the import path hint.
func vpcLbResource(src *schema.Resource, description, apiSegment string) *schema.Resource {
	res := *src
	res.Schema = vpcLbSchema(src.Schema)
	res.Importer = &schema.ResourceImporter{
		State: nsxtVersionCheckImporter(vpcLbMinNsxVersion, description,
			getVpcPathResourceImporter(fmt.Sprintf(vpcLbPathExample, apiSegment))),
	}
	create := src.Create
	res.Create = func(d *schema.ResourceData, m interface{}) error {
		if !util.NsxVersionHigherOrEqual(vpcLbMinNsxVersion) {
			return fmt.Errorf("%s resource requires NSX version %s or higher", description, vpcLbMinNsxVersion)
		}
		return create(d, m)
	}
	return &res
}

// vpcLbDataSource derives a VPC data source from its infra counterpart.
func vpcLbDataSource(src *schema.Resource, description string) *schema.Resource {
	res := *src
	res.Schema = vpcLbSchema(src.Schema)
	read := src.Read
	res.Read = func(d *schema.ResourceData, m interface{}) error {
		if !util.NsxVersionHigherOrEqual(vpcLbMinNsxVersion) {
			return fmt.Errorf("%s data source requires NSX version %s or higher", description, vpcLbMinNsxVersion)
		}
		return read(d, m)
	}
	return &res
}

func resourceNsxtVpcLBService() *schema.Resource {
	res := vpcLbResource(resourceNsxtPolicyLBService(), "VPC LB Service", "vpc-lbs")
	// NSX connects a VPC LB service to its VPC, the path of which is assigned by the server
	res.Schema["connectivity_path"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Policy path of the VPC which the LB service is connected to",
		Computed:    true,
	}
	return res
}

func dataSourceNsxtVpcLBService() *schema.Resource {
	return vpcLbDataSource(dataSourceNsxtPolicyLbService(), "VPC LB Service")
}
