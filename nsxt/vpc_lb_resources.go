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

func resourceNsxtVpcLBHttpApplicationProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBHttpApplicationProfile(), "VPC LB HTTP Application Profile", "vpc-lb-app-profiles")
}

func resourceNsxtVpcLBFastTcpApplicationProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBFastTcpApplicationProfile(), "VPC LB Fast TCP Application Profile", "vpc-lb-app-profiles")
}

func resourceNsxtVpcLBFastUdpApplicationProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBFastUdpApplicationProfile(), "VPC LB Fast UDP Application Profile", "vpc-lb-app-profiles")
}

func resourceNsxtVpcLBHttpMonitorProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBHttpMonitorProfile(), "VPC LB HTTP Monitor Profile", "vpc-lb-monitor-profiles")
}

func resourceNsxtVpcLBHttpsMonitorProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBHttpsMonitorProfile(), "VPC LB HTTPS Monitor Profile", "vpc-lb-monitor-profiles")
}

func resourceNsxtVpcLBIcmpMonitorProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBIcmpMonitorProfile(), "VPC LB ICMP Monitor Profile", "vpc-lb-monitor-profiles")
}

func resourceNsxtVpcLBPassiveMonitorProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBPassiveMonitorProfile(), "VPC LB Passive Monitor Profile", "vpc-lb-monitor-profiles")
}

func resourceNsxtVpcLBTcpMonitorProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBTcpMonitorProfile(), "VPC LB TCP Monitor Profile", "vpc-lb-monitor-profiles")
}

func resourceNsxtVpcLBUdpMonitorProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBUdpMonitorProfile(), "VPC LB UDP Monitor Profile", "vpc-lb-monitor-profiles")
}

func resourceNsxtVpcLBCookiePersistenceProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBCookiePersistenceProfile(), "VPC LB Cookie Persistence Profile", "vpc-lb-persistence-profiles")
}

func resourceNsxtVpcLBGenericPersistenceProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBGenericPersistenceProfile(), "VPC LB Generic Persistence Profile", "vpc-lb-persistence-profiles")
}

func resourceNsxtVpcLBSourceIpPersistenceProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBSourceIpPersistenceProfile(), "VPC LB Source IP Persistence Profile", "vpc-lb-persistence-profiles")
}

func resourceNsxtVpcLBClientSslProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBClientSslProfile(), "VPC LB Client SSL Profile", "vpc-lb-client-ssl-profiles")
}

func resourceNsxtVpcLBServerSslProfile() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBServerSslProfile(), "VPC LB Server SSL Profile", "vpc-lb-server-ssl-profiles")
}

func resourceNsxtVpcLBPool() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBPool(), "VPC LB Pool", "vpc-lb-pools")
}

func resourceNsxtVpcLBVirtualServer() *schema.Resource {
	return vpcLbResource(resourceNsxtPolicyLBVirtualServer(), "VPC LB Virtual Server", "vpc-lb-virtual-servers")
}

func dataSourceNsxtVpcLBAppProfile() *schema.Resource {
	return vpcLbDataSource(dataSourceNsxtPolicyLBAppProfile(), "VPC LB Application Profile")
}

func dataSourceNsxtVpcLBClientSslProfile() *schema.Resource {
	return vpcLbDataSource(dataSourceNsxtPolicyLBClientSslProfile(), "VPC LB Client SSL Profile")
}

func dataSourceNsxtVpcLBMonitor() *schema.Resource {
	return vpcLbDataSource(dataSourceNsxtPolicyLBMonitor(), "VPC LB Monitor Profile")
}

func dataSourceNsxtVpcLBPersistenceProfile() *schema.Resource {
	return vpcLbDataSource(dataSourceNsxtPolicyLbPersistenceProfile(), "VPC LB Persistence Profile")
}

func dataSourceNsxtVpcLBServerSslProfile() *schema.Resource {
	return vpcLbDataSource(dataSourceNsxtPolicyLBServerSslProfile(), "VPC LB Server SSL Profile")
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
