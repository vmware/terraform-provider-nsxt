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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliVNACRUDClient = enforcement_points.NewVirtualNetworkApplianceCRUDClient

func resourceNsxtPolicyVirtualNetworkAppliance() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyVirtualNetworkApplianceCreate,
		Read:   resourceNsxtPolicyVirtualNetworkApplianceRead,
		Update: resourceNsxtPolicyVirtualNetworkApplianceUpdate,
		Delete: resourceNsxtPolicyVirtualNetworkApplianceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyVirtualNetworkApplianceImporter,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"cluster_path": {
				Type:         schema.TypeString,
				Description:  "Policy path of the parent Virtual Network Appliance Cluster",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "Desired hostname or FQDN for the Virtual Network Appliance VM",
				Optional:    true,
				Computed:    true,
			},
			"failure_domain_path": {
				Type:         schema.TypeString,
				Description:  "Policy path of the failure domain",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"credentials": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Credentials for the Virtual Network Appliance node",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Node audit user password",
						},
						"audit_username": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Audit username (computed)",
						},
						"cli_password": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							DiffSuppressFunc: suppressIfEmptyPriorState,
							Description:      "Node CLI admin password",
						},
						"cli_username": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLI admin username (computed)",
						},
						"root_password": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							DiffSuppressFunc: suppressIfEmptyPriorState,
							Description:      "Node root user password",
						},
					},
				},
			},
			"management_interface": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Management interface configuration for the Virtual Network Appliance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Distributed portgroup identifier for the management vNIC",
						},
						"ip_assignment": getPolicyIPAssignmentSchema(true, 1, 2, managementAssignments),
					},
				},
			},
			"vm_deployment_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "vSphere deployment configuration for the Virtual Network Appliance VM",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_manager_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "NSX compute manager (vCenter) identifier",
						},
						"cluster_or_resource_pool_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "vSphere cluster or resource pool identifier",
						},
						"datastore_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "vSphere datastore identifier",
						},
						"reservation_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Resource reservation settings",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu_reservation_in_mhz": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "CPU reservation in MHz",
									},
									"cpu_reservation_in_shares": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      model.CPUReservation_RESERVATION_IN_SHARES_HIGH_PRIORITY,
										Description:  "CPU reservation in shares",
										ValidateFunc: validation.StringInSlice(policyCpuReservationValues, false),
									},
									"memory_reservation_percentage": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      100,
										Description:  "Memory reservation percentage",
										ValidateFunc: validation.IntBetween(0, 100),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func getVNAClusterPathComponents(clusterPath string) (siteID, epID, clusterID string, err error) {
	siteID = getResourceIDFromResourcePath(clusterPath, "sites")
	if siteID == "" {
		return "", "", "", fmt.Errorf("error obtaining site ID from cluster path %s", clusterPath)
	}
	epID = getResourceIDFromResourcePath(clusterPath, "enforcement-points")
	if epID == "" {
		return "", "", "", fmt.Errorf("error obtaining enforcement-point ID from cluster path %s", clusterPath)
	}
	clusterID = getResourceIDFromResourcePath(clusterPath, "virtual-network-appliance-clusters")
	if clusterID == "" {
		return "", "", "", fmt.Errorf("error obtaining cluster ID from cluster path %s", clusterPath)
	}
	return siteID, epID, clusterID, nil
}

func getVNACredentialsFromSchema(creds interface{}) *model.VirtualNetworkApplianceCredential {
	if creds == nil {
		return nil
	}
	list, ok := creds.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	cred := list[0].(map[string]interface{})
	obj := &model.VirtualNetworkApplianceCredential{}
	if v, ok := cred["cli_password"].(string); ok && v != "" {
		obj.CliPassword = &v
	}
	if v, ok := cred["root_password"].(string); ok && v != "" {
		obj.RootPassword = &v
	}
	if v, ok := cred["audit_password"].(string); ok && v != "" {
		obj.AuditPassword = &v
	}
	return obj
}

// suppressIfEmptyPriorState suppresses a diff for a write-only sensitive field
// when the prior state holds an empty string AND the resource already exists.
// This covers the import scenario: the API never returns passwords, so after
// `terraform import` the provider stores "" in state. Suppressing the
// config→state diff prevents a spurious "credentials needs to be changed" plan
// after import while still allowing real password changes once a value has been
// persisted in state.
// The d.Id() != "" guard ensures the suppression does not fire during Create
// (where the resource has no ID yet) or inside schema.TestResourceDataRaw
// (which builds a ResourceData from nil prior state with an empty ID).
func suppressIfEmptyPriorState(k, old, new string, d *schema.ResourceData) bool {
	return old == "" && d.Id() != ""
}

func setVNACredentialsInSchema(d *schema.ResourceData, obj *model.VirtualNetworkApplianceCredential) error {
	if obj == nil {
		return nil
	}
	// Passwords are write-only and not returned by GET. We only update state
	// when a credentials block is already present (user configured credentials
	// OR the importer seeded an empty block). When no block is in state the
	// early return is intentional: writing an empty block here would cause
	// spurious drift for VNAs that have no credentials in the Terraform config.
	// The import scenario is handled by the Importer, which seeds an empty
	// block before Read is called so the suppressIfEmptyPriorState
	// DiffSuppressFunc can suppress the password diff (bug 3715433).
	c := d.Get("credentials").([]interface{})
	if len(c) == 0 {
		return nil
	}
	// The Terraform Plugin SDK v2 normalises a TypeList element that was set
	// as an empty map (e.g. by the importer seeding an empty credentials block)
	// to nil when the element is read back via d.Get. Guard against that here
	// so we don't panic on the type assertion.
	creds, ok := c[0].(map[string]interface{})
	if !ok || creds == nil {
		creds = make(map[string]interface{})
	}
	if obj.AuditUsername != nil {
		creds["audit_username"] = *obj.AuditUsername
	}
	if obj.CliUsername != nil {
		creds["cli_username"] = *obj.CliUsername
	}
	return d.Set("credentials", []interface{}{creds})
}

func getVNAManagementInterfaceFromSchema(mgtIntface interface{}) (*model.VirtualNetworkApplianceManagementInterface, error) {
	if mgtIntface == nil {
		return nil, nil
	}
	list, ok := mgtIntface.([]interface{})
	if !ok || len(list) == 0 {
		return nil, nil
	}
	mgt := list[0].(map[string]interface{})
	ipAssignments, err := getPolicyIPAssignmentsFromSchema(mgt["ip_assignment"])
	if err != nil {
		return nil, err
	}
	networkID := mgt["network_id"].(string)
	return &model.VirtualNetworkApplianceManagementInterface{
		IpAssignmentSpecs: ipAssignments,
		NetworkId:         &networkID,
	}, nil
}

func setVNAManagementInterfaceInSchema(d *schema.ResourceData, obj *model.VirtualNetworkApplianceManagementInterface) error {
	if obj == nil {
		return nil
	}
	mgt := make(map[string]interface{})
	mgt["network_id"] = obj.NetworkId
	var err error
	mgt["ip_assignment"], err = setPolicyIPAssignmentsInSchema(obj.IpAssignmentSpecs)
	if err != nil {
		return err
	}
	return d.Set("management_interface", []interface{}{mgt})
}

func getVNADeploymentConfigFromSchema(cfgRaw interface{}) *model.VirtualNetworkApplianceDeploymentConfig {
	if cfgRaw == nil {
		return nil
	}
	list, ok := cfgRaw.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	cfg := list[0].(map[string]interface{})
	obj := &model.VirtualNetworkApplianceDeploymentConfig{}
	if v, ok := cfg["compute_manager_id"].(string); ok && v != "" {
		obj.ComputeManagerId = &v
	}
	if v, ok := cfg["cluster_or_resource_pool_id"].(string); ok && v != "" {
		obj.ClusterOrResourcePoolId = &v
	}
	if v, ok := cfg["datastore_id"].(string); ok && v != "" {
		obj.DatastoreId = &v
	}
	if riList, ok := cfg["reservation_info"].([]interface{}); ok && len(riList) > 0 {
		ri := riList[0].(map[string]interface{})
		reservationInfo := &model.ReservationInfo{}
		cpuMhz := int64(ri["cpu_reservation_in_mhz"].(int))
		cpuShares := ri["cpu_reservation_in_shares"].(string)
		memPct := int64(ri["memory_reservation_percentage"].(int))
		reservationInfo.CpuReservation = &model.CPUReservation{
			ReservationInMhz:    &cpuMhz,
			ReservationInShares: &cpuShares,
		}
		reservationInfo.MemoryReservation = &model.MemoryReservation{
			ReservationPercentage: &memPct,
		}
		obj.ReservationInfo = reservationInfo
	}
	return obj
}

func setVNADeploymentConfigInSchema(d *schema.ResourceData, obj *model.VirtualNetworkApplianceDeploymentConfig) error {
	if obj == nil {
		return nil
	}
	cfg := make(map[string]interface{})
	if obj.ComputeManagerId != nil {
		cfg["compute_manager_id"] = *obj.ComputeManagerId
	}
	if obj.ClusterOrResourcePoolId != nil {
		cfg["cluster_or_resource_pool_id"] = *obj.ClusterOrResourcePoolId
	}
	if obj.DatastoreId != nil {
		cfg["datastore_id"] = *obj.DatastoreId
	}
	if obj.ReservationInfo != nil {
		ri := make(map[string]interface{})
		if obj.ReservationInfo.CpuReservation != nil {
			if obj.ReservationInfo.CpuReservation.ReservationInMhz != nil {
				ri["cpu_reservation_in_mhz"] = int(*obj.ReservationInfo.CpuReservation.ReservationInMhz)
			}
			if obj.ReservationInfo.CpuReservation.ReservationInShares != nil {
				ri["cpu_reservation_in_shares"] = *obj.ReservationInfo.CpuReservation.ReservationInShares
			}
		}
		if obj.ReservationInfo.MemoryReservation != nil && obj.ReservationInfo.MemoryReservation.ReservationPercentage != nil {
			ri["memory_reservation_percentage"] = int(*obj.ReservationInfo.MemoryReservation.ReservationPercentage)
		}
		cfg["reservation_info"] = []interface{}{ri}
	}
	return d.Set("vm_deployment_config", []interface{}{cfg})
}

func policyVNABuildObject(d *schema.ResourceData) (model.VirtualNetworkAppliance, error) {
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.VirtualNetworkAppliance{
		Description:  &description,
		DisplayName:  &displayName,
		Tags:         tags,
		ResourceType: model.VirtualNetworkAppliance__TYPE_IDENTIFIER,
	}

	if v, ok := d.GetOk("hostname"); ok {
		s := v.(string)
		obj.Hostname = &s
	}
	if v, ok := d.GetOk("failure_domain_path"); ok {
		s := v.(string)
		obj.FailureDomainPath = &s
	}

	obj.Credentials = getVNACredentialsFromSchema(d.Get("credentials"))

	mgtIntf, err := getVNAManagementInterfaceFromSchema(d.Get("management_interface"))
	if err != nil {
		return obj, err
	}
	obj.ManagementInterface = mgtIntf

	obj.VmDeploymentConfig = getVNADeploymentConfigFromSchema(d.Get("vm_deployment_config"))

	return obj, nil
}

func resourceNsxtPolicyVirtualNetworkApplianceExists(clusterPath, id string, connector client.Connector) (bool, error) {
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	siteID, epID, clusterID, err := getVNAClusterPathComponents(clusterPath)
	if err != nil {
		return false, err
	}
	vnaClient := cliVNACRUDClient(sessionContext, connector)
	if vnaClient == nil {
		return false, policyResourceNotSupportedError()
	}
	_, err = vnaClient.Get(siteID, epID, clusterID, id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving VirtualNetworkAppliance", err)
}

func resourceNsxtPolicyVirtualNetworkApplianceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining VirtualNetworkAppliance ID")
	}

	clusterPath := d.Get("cluster_path").(string)
	siteID, epID, clusterID, err := getVNAClusterPathComponents(clusterPath)
	if err != nil {
		return err
	}

	vnaClient := cliVNACRUDClient(sessionContext, connector)
	if vnaClient == nil {
		return policyResourceNotSupportedError()
	}

	obj, err := vnaClient.Get(siteID, epID, clusterID, id)
	if err != nil {
		return handleReadError(d, "VirtualNetworkAppliance", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("hostname", obj.Hostname)
	d.Set("failure_domain_path", obj.FailureDomainPath)

	if err := setVNACredentialsInSchema(d, obj.Credentials); err != nil {
		return err
	}
	if err := setVNAManagementInterfaceInSchema(d, obj.ManagementInterface); err != nil {
		return err
	}
	if err := setVNADeploymentConfigInSchema(d, obj.VmDeploymentConfig); err != nil {
		return err
	}

	return nil
}

func resourceNsxtPolicyVirtualNetworkApplianceCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Get("nsx_id").(string)
	if id == "" {
		id = newUUID()
	}

	clusterPath := d.Get("cluster_path").(string)
	siteID, epID, clusterID, err := getVNAClusterPathComponents(clusterPath)
	if err != nil {
		return err
	}

	exists, err := resourceNsxtPolicyVirtualNetworkApplianceExists(clusterPath, id, connector)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("resource with ID %s already exists", id)
	}

	log.Printf("[INFO] Creating VirtualNetworkAppliance with ID %s under cluster %s", id, clusterPath)

	sessionContext := getSessionContext(d, m)
	vnaClient := cliVNACRUDClient(sessionContext, connector)
	if vnaClient == nil {
		return policyResourceNotSupportedError()
	}

	obj, err := policyVNABuildObject(d)
	if err != nil {
		return err
	}
	obj.Id = &id

	if err := vnaClient.Patch(siteID, epID, clusterID, id, obj); err != nil {
		return handleCreateError("VirtualNetworkAppliance", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
}

func resourceNsxtPolicyVirtualNetworkApplianceUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining VirtualNetworkAppliance ID")
	}

	clusterPath := d.Get("cluster_path").(string)
	siteID, epID, clusterID, err := getVNAClusterPathComponents(clusterPath)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating VirtualNetworkAppliance with ID %s", id)

	sessionContext := getSessionContext(d, m)
	vnaClient := cliVNACRUDClient(sessionContext, connector)
	if vnaClient == nil {
		return policyResourceNotSupportedError()
	}

	obj, err := policyVNABuildObject(d)
	if err != nil {
		return err
	}
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	if _, err := vnaClient.Update(siteID, epID, clusterID, id, obj); err != nil {
		return handleUpdateError("VirtualNetworkAppliance", id, err)
	}

	return resourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
}

func resourceNsxtPolicyVirtualNetworkApplianceDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	vnaClient := cliVNACRUDClient(sessionContext, connector)
	if vnaClient == nil {
		return policyResourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining VirtualNetworkAppliance ID")
	}

	clusterPath := d.Get("cluster_path").(string)
	siteID, epID, clusterID, err := getVNAClusterPathComponents(clusterPath)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting VirtualNetworkAppliance with ID %s", id)
	if err := vnaClient.Delete(siteID, epID, clusterID, id); err != nil {
		return handleDeleteError("VirtualNetworkAppliance", id, err)
	}

	return nil
}

func resourceNsxtPolicyVirtualNetworkApplianceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return rd, err
	}

	epID, err := getParameterFromPolicyPath("/enforcement-points/", "/virtual-network-appliance-clusters/", importID)
	if err != nil {
		return nil, err
	}

	sitePath, err := getSitePathFromChildResourcePath(importID)
	if err != nil {
		return rd, err
	}

	clusterPath := fmt.Sprintf("%s/enforcement-points/%s/virtual-network-appliance-clusters/%s",
		sitePath, epID,
		getResourceIDFromResourcePath(importID, "virtual-network-appliance-clusters"),
	)
	d.Set("cluster_path", clusterPath)

	// Seed an empty credentials block so that setVNACredentialsInSchema writes
	// NSX-returned usernames to state during the Read that follows import.
	// Without this seed the function's early return (for no-credentials-in-state)
	// would leave credentials absent from state, causing the subsequent plan to
	// show +credentials drift. The suppressIfEmptyPriorState DiffSuppressFunc on
	// the password fields then ensures the plan shows no changes even though the
	// config passwords are not yet persisted in state (bug 3715433).
	d.Set("credentials", []interface{}{map[string]interface{}{}})

	return rd, nil
}
