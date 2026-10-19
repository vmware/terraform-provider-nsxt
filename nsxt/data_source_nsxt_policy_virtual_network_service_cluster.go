// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyVirtualNetworkServiceCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyVirtualNetworkServiceClusterRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"site_path": {
				Type:         schema.TypeString,
				Description:  "Path to the site this Virtual Network Service Cluster belongs to",
				Optional:     true,
				Default:      defaultInfraSitePath,
				ValidateFunc: validatePolicyPath(),
			},
			"enforcement_point": {
				Type:        schema.TypeString,
				Description: "ID of the enforcement point this Virtual Network Service Cluster belongs to",
				Optional:    true,
				Default:     "default",
			},
			"node_form_factor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Form factor of the Virtual Network Service nodes in the cluster",
			},
			"number_of_nodes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of Virtual Network Service nodes in the cluster",
			},
			"bfd_profile_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Path of the BFD profile applied to all Service Hosts in this cluster",
			},
			"default_project_service_allocation": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this cluster is available for default project service allocation",
			},
			"vsphere_deployment_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "vSphere deployment configurations for the cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_manager_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"fault_domain_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func listVirtualNetworkServiceClusters(siteID, epID string, connector client.Connector, sessionContext utl.SessionContext) ([]model.VirtualNetworkServiceCluster, error) {
	client := cliVNSClustersClient(sessionContext, connector)
	if client == nil {
		return nil, policyResourceNotSupportedError()
	}

	var results []model.VirtualNetworkServiceCluster
	boolFalse := false
	var cursor *string
	total := 0

	for {
		listResult, err := client.List(siteID, epID, cursor, &boolFalse, nil, nil, nil, nil)
		if err != nil {
			return results, err
		}
		results = append(results, listResult.Results...)
		if total == 0 && listResult.ResultCount != nil {
			total = int(*listResult.ResultCount)
		}
		cursor = listResult.Cursor
		if len(results) >= total {
			return results, nil
		}
	}
}

func dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d *schema.ResourceData, m interface{}) error {
	if err := vnsClusterCheckVersion(); err != nil {
		return err
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	sitePath := d.Get("site_path").(string)
	siteID := getResourceIDFromResourcePath(sitePath, "sites")
	if siteID == "" {
		return fmt.Errorf("error obtaining Site ID from site path %s", sitePath)
	}
	epID := d.Get("enforcement_point").(string)
	if epID == "" {
		epID = getPolicyEnforcementPoint(m)
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	var obj model.VirtualNetworkServiceCluster

	client := cliVNSClustersClient(sessionContext, connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	if objID != "" {
		objGet, err := client.Get(siteID, epID, objID)
		if err != nil {
			return handleDataSourceReadError(d, "VirtualNetworkServiceCluster", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("error obtaining VirtualNetworkServiceCluster ID or name during read")
	} else {
		objList, err := listVirtualNetworkServiceClusters(siteID, epID, connector, sessionContext)
		if err != nil {
			return handleListError("VirtualNetworkServiceCluster", err)
		}

		var perfectMatch []model.VirtualNetworkServiceCluster
		var prefixMatch []model.VirtualNetworkServiceCluster
		for _, objInList := range objList {
			if objInList.DisplayName == nil {
				continue
			}
			if strings.HasPrefix(*objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if *objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("found multiple VirtualNetworkServiceClusters with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("found multiple VirtualNetworkServiceClusters with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("VirtualNetworkServiceCluster with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("node_form_factor", obj.NodeFormFactor)
	d.Set("number_of_nodes", obj.NumberOfNodes)
	d.Set("bfd_profile_path", obj.BfdProfilePath)
	if obj.AllocationConfig != nil {
		d.Set("default_project_service_allocation", obj.AllocationConfig.DefaultProjectServiceAllocation)
	}

	return setVNSClusterDeploymentConfigInSchema(d, obj.VsphereDeploymentConfig)
}
