// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	enforcement_points "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var cliVNSClustersClient = enforcement_points.NewVirtualNetworkServiceClustersClient

const vnsClusterMinNsxVersion = "9.2.0"

var vnsClusterNodeFormFactorValues = []string{
	model.VirtualNetworkServiceCluster_NODE_FORM_FACTOR_SMALL,
	model.VirtualNetworkServiceCluster_NODE_FORM_FACTOR_MEDIUM,
	model.VirtualNetworkServiceCluster_NODE_FORM_FACTOR_LARGE,
}

func resourceNsxtPolicyVirtualNetworkServiceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyVirtualNetworkServiceClusterCreate,
		Read:   resourceNsxtPolicyVirtualNetworkServiceClusterRead,
		Update: resourceNsxtPolicyVirtualNetworkServiceClusterUpdate,
		Delete: resourceNsxtPolicyVirtualNetworkServiceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyVirtualNetworkServiceClusterImporter,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"site_path": {
				Type:         schema.TypeString,
				Description:  "Path to the site this Virtual Network Service Cluster belongs to",
				Optional:     true,
				ForceNew:     true,
				Default:      defaultInfraSitePath,
				ValidateFunc: validatePolicyPath(),
			},
			"enforcement_point": {
				Type:        schema.TypeString,
				Description: "ID of the enforcement point this Virtual Network Service Cluster belongs to",
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
			},
			"node_form_factor": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Form factor of the Virtual Network Service nodes in the cluster",
				ValidateFunc: validation.StringInSlice(vnsClusterNodeFormFactorValues, false),
			},
			"number_of_nodes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Number of Virtual Network Service nodes to deploy in the cluster",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"bfd_profile_path": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Path of the BFD profile applied to all Service Hosts in this cluster",
				ValidateFunc: validatePolicyPath(),
			},
			"default_project_service_allocation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "When true, this cluster is available for default project service allocation",
			},
			"vsphere_deployment_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "vSphere deployment configurations for the Virtual Network Service cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_manager_id": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "NSX identifier of the registered vCenter server",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"cluster_ids": {
							Type:        schema.TypeSet,
							Required:    true,
							MinItems:    1,
							Description: "vCenter clusters on which the Virtual Network Service nodes are deployed",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"fault_domain_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Infrastructure fault domains or fault domain groups restricting the deployment within the clusters",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func vnsClusterCheckVersion() error {
	if !util.NsxVersionHigherOrEqual(vnsClusterMinNsxVersion) {
		return fmt.Errorf("this resource requires NSX version %s or higher", vnsClusterMinNsxVersion)
	}
	return nil
}

func setVNSClusterDeploymentConfigInSchema(d *schema.ResourceData, cfgs []model.VirtualNetworkServiceDeploymentConfig) error {
	var list []interface{}
	for _, cfg := range cfgs {
		entry := make(map[string]interface{})
		if cfg.ComputeManagerId != nil {
			entry["compute_manager_id"] = *cfg.ComputeManagerId
		}
		entry["cluster_ids"] = stringList2Interface(cfg.ClusterIds)
		entry["fault_domain_ids"] = stringList2Interface(cfg.FaultDomainIds)
		list = append(list, entry)
	}
	return d.Set("vsphere_deployment_config", list)
}

func getVNSClusterDeploymentConfigFromSchema(d *schema.ResourceData) []model.VirtualNetworkServiceDeploymentConfig {
	var cfgs []model.VirtualNetworkServiceDeploymentConfig
	for _, raw := range d.Get("vsphere_deployment_config").([]interface{}) {
		data := raw.(map[string]interface{})
		cmID := data["compute_manager_id"].(string)
		cfg := model.VirtualNetworkServiceDeploymentConfig{
			ComputeManagerId: &cmID,
			ClusterIds:       interface2StringList(data["cluster_ids"].(*schema.Set).List()),
		}
		if fd, ok := data["fault_domain_ids"].(*schema.Set); ok && fd.Len() > 0 {
			cfg.FaultDomainIds = interface2StringList(fd.List())
		}
		cfgs = append(cfgs, cfg)
	}
	return cfgs
}

func policyVNSClusterBuildObject(d *schema.ResourceData) model.VirtualNetworkServiceCluster {
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.VirtualNetworkServiceCluster{
		Description:             &description,
		DisplayName:             &displayName,
		Tags:                    tags,
		VsphereDeploymentConfig: getVNSClusterDeploymentConfigFromSchema(d),
	}

	if v, ok := d.GetOk("node_form_factor"); ok {
		s := v.(string)
		obj.NodeFormFactor = &s
	}
	if v, ok := d.GetOk("number_of_nodes"); ok {
		n := int64(v.(int))
		obj.NumberOfNodes = &n
	}
	if v, ok := d.GetOk("bfd_profile_path"); ok {
		s := v.(string)
		obj.BfdProfilePath = &s
	}
	if v, ok := d.GetOkExists("default_project_service_allocation"); ok { //nolint:staticcheck
		b := v.(bool)
		obj.AllocationConfig = &model.VirtualNetworkServiceAllocationConfig{DefaultProjectServiceAllocation: &b}
	}

	return obj
}

func resourceNsxtPolicyVirtualNetworkServiceClusterExists(siteID, epID, id string, connector client.Connector) (bool, error) {
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	siteClient := cliSitesClient(sessionContext, connector)
	if siteClient == nil {
		return false, policyResourceNotSupportedError()
	}
	_, err := siteClient.Get(siteID)
	if err != nil {
		msg := fmt.Sprintf("failed to read site %s", siteID)
		return false, logAPIError(msg, err)
	}
	vnsClient := cliVNSClustersClient(sessionContext, connector)
	if vnsClient == nil {
		return false, policyResourceNotSupportedError()
	}
	_, err = vnsClient.Get(siteID, epID, id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving VirtualNetworkServiceCluster", err)
}

func resourceNsxtPolicyVirtualNetworkServiceClusterCreate(d *schema.ResourceData, m interface{}) error {
	if err := vnsClusterCheckVersion(); err != nil {
		return err
	}
	connector := getPolicyConnector(m)

	id := d.Get("nsx_id").(string)
	if id == "" {
		id = newUUID()
	}
	sitePath := d.Get("site_path").(string)
	siteID := getResourceIDFromResourcePath(sitePath, "sites")
	if siteID == "" {
		return fmt.Errorf("error obtaining Site ID from site path %s", sitePath)
	}
	epID := d.Get("enforcement_point").(string)
	if epID == "" {
		epID = getPolicyEnforcementPoint(m)
	}

	exists, err := resourceNsxtPolicyVirtualNetworkServiceClusterExists(siteID, epID, id, connector)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("resource with ID %s already exists", id)
	}

	log.Printf("[INFO] Creating VirtualNetworkServiceCluster with ID %s under site %s enforcement point %s", id, siteID, epID)

	client := cliVNSClustersClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	err = client.Patch(siteID, epID, id, policyVNSClusterBuildObject(d))
	if err != nil {
		return handleCreateError("VirtualNetworkServiceCluster", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m)
}

func resourceNsxtPolicyVirtualNetworkServiceClusterRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := cliVNSClustersClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	obj, err := client.Get(siteID, epID, id)
	if err != nil {
		return handleReadError(d, "VirtualNetworkServiceCluster", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("node_form_factor", obj.NodeFormFactor)
	d.Set("number_of_nodes", obj.NumberOfNodes)
	d.Set("bfd_profile_path", obj.BfdProfilePath)
	if obj.AllocationConfig != nil {
		d.Set("default_project_service_allocation", obj.AllocationConfig.DefaultProjectServiceAllocation)
	}

	return setVNSClusterDeploymentConfigInSchema(d, obj.VsphereDeploymentConfig)
}

func resourceNsxtPolicyVirtualNetworkServiceClusterUpdate(d *schema.ResourceData, m interface{}) error {
	if err := vnsClusterCheckVersion(); err != nil {
		return err
	}
	connector := getPolicyConnector(m)

	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating VirtualNetworkServiceCluster with ID %s", id)

	client := cliVNSClustersClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	obj := policyVNSClusterBuildObject(d)
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	_, err = client.Update(siteID, epID, id, obj)
	if err != nil {
		return handleUpdateError("VirtualNetworkServiceCluster", id, err)
	}

	return resourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m)
}

func resourceNsxtPolicyVirtualNetworkServiceClusterDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := cliVNSClustersClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting VirtualNetworkServiceCluster with ID %s", id)
	err = client.Delete(siteID, epID, id)
	if err != nil {
		return handleDeleteError("VirtualNetworkServiceCluster", id, err)
	}

	return nil
}

func resourceNsxtPolicyVirtualNetworkServiceClusterImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return rd, err
	}

	epID, err := getParameterFromPolicyPath("/enforcement-points/", "/virtual-network-service-clusters/", importID)
	if err != nil {
		return nil, err
	}
	d.Set("enforcement_point", epID)

	sitePath, err := getSitePathFromChildResourcePath(importID)
	if err != nil {
		return rd, err
	}
	d.Set("site_path", sitePath)

	return rd, nil
}
