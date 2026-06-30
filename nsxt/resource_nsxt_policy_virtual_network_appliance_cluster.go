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

var cliVNAClustersClient = enforcement_points.NewVirtualNetworkApplianceClustersClient

var vnaClusterFormFactorValues = []string{
	model.VirtualNetworkApplianceCluster_APPLIANCE_FORM_FACTOR_SMALL,
	model.VirtualNetworkApplianceCluster_APPLIANCE_FORM_FACTOR_MEDIUM,
	model.VirtualNetworkApplianceCluster_APPLIANCE_FORM_FACTOR_LARGE,
	model.VirtualNetworkApplianceCluster_APPLIANCE_FORM_FACTOR_XLARGE,
}

var vnaClusterCoreAllocationProfileValues = []string{
	model.VirtualNetworkApplianceClusterAdvancedConfiguration_CORE_ALLOCATION_PROFILE_L4FORWARDING,
	model.VirtualNetworkApplianceClusterAdvancedConfiguration_CORE_ALLOCATION_PROFILE_L7LBSERVICE,
	model.VirtualNetworkApplianceClusterAdvancedConfiguration_CORE_ALLOCATION_PROFILE_L4LBSERVICE,
}

var vnaClusterServiceTypeValues = []string{
	model.VirtualNetworkApplianceCluster_SERVICE_TYPE_VPC_SERVICES,
	model.VirtualNetworkApplianceCluster_SERVICE_TYPE_ROUTE_CONTROLLER,
}

func resourceNsxtPolicyVirtualNetworkApplianceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyVirtualNetworkApplianceClusterCreate,
		Read:   resourceNsxtPolicyVirtualNetworkApplianceClusterRead,
		Update: resourceNsxtPolicyVirtualNetworkApplianceClusterUpdate,
		Delete: resourceNsxtPolicyVirtualNetworkApplianceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyVirtualNetworkApplianceClusterImporter,
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
				Description:  "Path to the site this Virtual Network Appliance Cluster belongs to",
				Optional:     true,
				ForceNew:     true,
				Default:      defaultInfraSitePath,
				ValidateFunc: validatePolicyPath(),
			},
			"enforcement_point": {
				Type:        schema.TypeString,
				Description: "ID of the enforcement point this Virtual Network Appliance Cluster belongs to",
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
			},
			"appliance_form_factor": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Form factor for virtual network appliances in this cluster",
				ValidateFunc: validation.StringInSlice(vnaClusterFormFactorValues, false),
			},
			"appliance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Virtual network appliance type of the cluster",
			},
			"service_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Description:  "Service type of the cluster. Cannot be modified after creation",
				ValidateFunc: validation.StringInSlice(vnaClusterServiceTypeValues, false),
			},
			"password_managed_by_vcf": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable VCF password management for all virtual network appliances in the cluster",
			},
			"advanced_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				// Computed is required because NSX auto-populates this block with
				// server-side defaults (e.g. high_availability_profile,
				// overlay_transport_zone_path, core_allocation_profile) even when
				// the user omits it. Without Computed, Terraform plans perpetual
				// in-place updates to remove the server-populated values.
				Computed:    true,
				MaxItems:    1,
				Description: "Advanced configuration for virtual network appliances in the cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"core_allocation_profile": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "Core allocation profile for virtual network appliances in the cluster. Supported values: `L4FORWARDING`, `L7LBSERVICE`, `L4LBSERVICE`. For VPC_SERVICES clusters, `L4LBSERVICE` is the default when not set.",
							ValidateFunc: validation.StringInSlice(vnaClusterCoreAllocationProfileValues, false),
						},
						"high_availability_profile": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "Path to the high availability profile",
							ValidateFunc: validatePolicyPath(),
						},
						"overlay_transport_zone_path": {
							Type:     schema.TypeString,
							Optional: true,
							// Computed because NSX auto-populates this field server-side.
							Computed:     true,
							Description:  "Overlay transport zone path associated with the VNA host switch and TEP",
							ValidateFunc: validatePolicyPath(),
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyVirtualNetworkApplianceClusterRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	client := cliVNAClustersClient(sessionContext, connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	obj, err := client.Get(siteID, epID, id)
	if err != nil {
		return handleReadError(d, "VirtualNetworkApplianceCluster", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("appliance_form_factor", obj.ApplianceFormFactor)
	d.Set("appliance_type", obj.ApplianceType)
	d.Set("service_type", obj.ServiceType)
	d.Set("password_managed_by_vcf", obj.PasswordManagedByVcf)

	if err := setVNAClusterAdvancedConfigInSchema(d, obj.AdvancedConfiguration); err != nil {
		return err
	}

	return nil
}

func setVNAClusterAdvancedConfigInSchema(d *schema.ResourceData, cfg *model.VirtualNetworkApplianceClusterAdvancedConfiguration) error {
	if cfg == nil {
		return d.Set("advanced_configuration", nil)
	}
	entry := make(map[string]interface{})
	if cfg.CoreAllocationProfile != nil && util.NsxVersionHigherOrEqual("9.2.0") {
		entry["core_allocation_profile"] = *cfg.CoreAllocationProfile
	}
	if cfg.HighAvailabilityProfile != nil {
		entry["high_availability_profile"] = *cfg.HighAvailabilityProfile
	}
	if cfg.OverlayTransportZonePath != nil {
		entry["overlay_transport_zone_path"] = *cfg.OverlayTransportZonePath
	}
	return d.Set("advanced_configuration", []interface{}{entry})
}

func getVNAClusterAdvancedConfigFromSchema(d *schema.ResourceData) (*model.VirtualNetworkApplianceClusterAdvancedConfiguration, error) {
	rawList := d.Get("advanced_configuration").([]interface{})
	if len(rawList) == 0 {
		return nil, nil
	}
	data := rawList[0].(map[string]interface{})
	cfg := &model.VirtualNetworkApplianceClusterAdvancedConfiguration{}
	if v, ok := data["core_allocation_profile"].(string); ok && v != "" {
		if !util.NsxVersionHigherOrEqual("9.2.0") {
			return nil, fmt.Errorf("advanced_configuration.core_allocation_profile requires NSX version 9.2.0 or higher")
		}
		cfg.CoreAllocationProfile = &v
	}
	if v, ok := data["high_availability_profile"].(string); ok && v != "" {
		cfg.HighAvailabilityProfile = &v
	}
	if v, ok := data["overlay_transport_zone_path"].(string); ok && v != "" {
		cfg.OverlayTransportZonePath = &v
	}
	return cfg, nil
}

func resourceNsxtPolicyVirtualNetworkApplianceClusterExists(siteID, epID, id string, connector client.Connector) (bool, error) {
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
	vnaClient := cliVNAClustersClient(sessionContext, connector)
	if vnaClient == nil {
		return false, policyResourceNotSupportedError()
	}
	_, err = vnaClient.Get(siteID, epID, id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving VirtualNetworkApplianceCluster", err)
}

func policyVNAClusterBuildObject(d *schema.ResourceData) (model.VirtualNetworkApplianceCluster, error) {
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	advancedConfig, err := getVNAClusterAdvancedConfigFromSchema(d)
	if err != nil {
		return model.VirtualNetworkApplianceCluster{}, err
	}

	obj := model.VirtualNetworkApplianceCluster{
		Description:           &description,
		DisplayName:           &displayName,
		Tags:                  tags,
		AdvancedConfiguration: advancedConfig,
	}

	if v, ok := d.GetOk("appliance_form_factor"); ok {
		s := v.(string)
		obj.ApplianceFormFactor = &s
	}
	if v, ok := d.GetOk("appliance_type"); ok {
		s := v.(string)
		obj.ApplianceType = &s
	}
	if v, ok := d.GetOk("service_type"); ok {
		s := v.(string)
		obj.ServiceType = &s
	}
	if v, ok := d.GetOkExists("password_managed_by_vcf"); ok { //nolint:staticcheck
		b := v.(bool)
		obj.PasswordManagedByVcf = &b
	}

	return obj, nil
}

func resourceNsxtPolicyVirtualNetworkApplianceClusterCreate(d *schema.ResourceData, m interface{}) error {
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

	exists, err := resourceNsxtPolicyVirtualNetworkApplianceClusterExists(siteID, epID, id, connector)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("resource with ID %s already exists", id)
	}

	log.Printf("[INFO] Creating VirtualNetworkApplianceCluster with ID %s under site %s enforcement point %s", id, siteID, epID)

	sessionContext := getSessionContext(d, m)
	client := cliVNAClustersClient(sessionContext, connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	obj, err := policyVNAClusterBuildObject(d)
	if err != nil {
		return err
	}

	err = client.Patch(siteID, epID, id, obj)
	if err != nil {
		return handleCreateError("VirtualNetworkApplianceCluster", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
}

func resourceNsxtPolicyVirtualNetworkApplianceClusterUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating VirtualNetworkApplianceCluster with ID %s", id)

	sessionContext := getSessionContext(d, m)
	client := cliVNAClustersClient(sessionContext, connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	obj, err := policyVNAClusterBuildObject(d)
	if err != nil {
		return err
	}
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	_, err = client.Update(siteID, epID, id, obj)
	if err != nil {
		return handleUpdateError("VirtualNetworkApplianceCluster", id, err)
	}

	return resourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
}

func resourceNsxtPolicyVirtualNetworkApplianceClusterDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	client := cliVNAClustersClient(sessionContext, connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting VirtualNetworkApplianceCluster with ID %s", id)
	err = client.Delete(siteID, epID, id)
	if err != nil {
		return handleDeleteError("VirtualNetworkApplianceCluster", id, err)
	}

	return nil
}

func resourceNsxtPolicyVirtualNetworkApplianceClusterImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return rd, err
	}

	epID, err := getParameterFromPolicyPath("/enforcement-points/", "/virtual-network-appliance-clusters/", importID)
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
