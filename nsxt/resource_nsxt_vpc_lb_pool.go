// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var vpcLbPoolPathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-pools/[pool]"

func resourceNsxtVpcLBPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBPoolCreate,
		Read:   resourceNsxtVpcLBPoolRead,
		Update: resourceNsxtVpcLBPoolUpdate,
		Delete: resourceNsxtVpcLBPoolDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB Pool", getVpcPathResourceImporter(vpcLbPoolPathExample)),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
			"member":       getPoolMembersSchema(),
			"member_group": getPolicyPoolMemberGroupSchema(),
			"active_monitor_paths": {
				Type:          schema.TypeList,
				Description:   "Used by the load balancer to initiate new connections to the servers to check their health. Active healthchecks are deactivated by default and can be activated using this setting",
				Elem:          getPolicyPathSchemaSimple(),
				Optional:      true,
				ConflictsWith: []string{"active_monitor_path"},
			},
			"active_monitor_path": getPolicyPathSchemaExtended(false, false, "Active healthcheck is disabled by default and can be enabled using this setting", "This attribute is deprecated, please use active_monitor_paths", []string{"active_monitor_paths"}),
			"algorithm": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(lbPoolAlgorithmValues, false),
				Optional:     true,
				Default:      model.LBPool_ALGORITHM_ROUND_ROBIN,
			},
			"min_active_members": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Default:      1,
			},
			"passive_monitor_path": getPolicyPathSchema(false, false, "Policy path for passive health monitor"),
			"tcp_multiplexing_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tcp_multiplexing_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"snat": getPolicyPoolSnatSchema(),
		},
	}
}

func resourceNsxtVpcLBPoolExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPoolsClient(sessionContext, connector)
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcLBPoolCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Pool resource requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPoolsClient(sessionContext, connector)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBPoolExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	activeMonitorPaths := interfaceListToStringList(d.Get("active_monitor_paths").([]interface{}))
	if activeMonitorPaths == nil && d.Get("active_monitor_path") != "" {
		activeMonitorPath := d.Get("active_monitor_path").(string)
		activeMonitorPaths = []string{activeMonitorPath}
	}
	algorithm := d.Get("algorithm").(string)
	memberGroup := getPolicyPoolMemberGroupFromSchema(d)
	members := getPolicyPoolMembersFromSchema(d)
	minActiveMembers := int64(d.Get("min_active_members").(int))
	passiveMonitorPath := d.Get("passive_monitor_path").(string)
	snatTranslation, err := getPolicyPoolSnatFromSchema(d)
	if err != nil {
		return err
	}
	tcpMultiplexingEnabled := d.Get("tcp_multiplexing_enabled").(bool)
	tcpMultiplexingNumber := int64(d.Get("tcp_multiplexing_number").(int))

	obj := model.LBPool{
		DisplayName:            &displayName,
		Description:            &description,
		Tags:                   tags,
		ActiveMonitorPaths:     activeMonitorPaths,
		Algorithm:              &algorithm,
		MemberGroup:            memberGroup,
		Members:                members,
		PassiveMonitorPath:     &passiveMonitorPath,
		MinActiveMembers:       &minActiveMembers,
		SnatTranslation:        snatTranslation,
		TcpMultiplexingEnabled: &tcpMultiplexingEnabled,
		TcpMultiplexingNumber:  &tcpMultiplexingNumber,
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating VPC LBPool with ID %s", id)
	err = c.Patch(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleCreateError("VPC LBPool", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtVpcLBPoolRead(d, m)
}

func resourceNsxtVpcLBPoolRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPoolsClient(sessionContext, connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBPool ID")
	}

	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBPool", id, err)
	}

	if _, ok := d.GetOk("active_monitor_path"); ok {
		// deprecated value was set
		if len(obj.ActiveMonitorPaths) > 0 {
			d.Set("active_monitor_path", obj.ActiveMonitorPaths[0])
		} else if err := d.Set("active_monitor_path", ""); err != nil {
			return err
		}
	} else {
		// resource uses lists attribute (non-deprecated) or attribute
		// not specified. In this case set non-deprecated attr in state
		d.Set("active_monitor_paths", obj.ActiveMonitorPaths)
	}
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("algorithm", obj.Algorithm)
	err = setPolicyPoolMembersInSchema(d, obj.Members)
	if err != nil {
		return err
	}
	err = setPolicyPoolMemberGroupInSchema(d, obj.MemberGroup)
	if err != nil {
		return err
	}
	err = setPolicyPoolSnatInSchema(d, obj.SnatTranslation)
	if err != nil {
		return err
	}

	if obj.MinActiveMembers != nil {
		d.Set("min_active_members", obj.MinActiveMembers)
	} else {
		d.Set("min_active_members", 1)
	}
	d.Set("passive_monitor_path", obj.PassiveMonitorPath)
	d.Set("tcp_multiplexing_enabled", obj.TcpMultiplexingEnabled)
	d.Set("tcp_multiplexing_number", obj.TcpMultiplexingNumber)

	return nil
}

func resourceNsxtVpcLBPoolUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPoolsClient(sessionContext, connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBPool ID")
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	activeMonitorPaths := interfaceListToStringList(d.Get("active_monitor_paths").([]interface{}))
	if activeMonitorPaths == nil && d.Get("active_monitor_path") != "" {
		activeMonitorPath := d.Get("active_monitor_path").(string)
		activeMonitorPaths = []string{activeMonitorPath}
	}
	algorithm := d.Get("algorithm").(string)
	memberGroup := getPolicyPoolMemberGroupFromSchema(d)
	members := getPolicyPoolMembersFromSchema(d)
	passiveMonitorPath := d.Get("passive_monitor_path").(string)
	snatTranslation, err := getPolicyPoolSnatFromSchema(d)
	if err != nil {
		return err
	}
	tcpMultiplexingEnabled := d.Get("tcp_multiplexing_enabled").(bool)
	tcpMultiplexingNumber := int64(d.Get("tcp_multiplexing_number").(int))
	minActiveMembers := int64(d.Get("min_active_members").(int))
	revision := int64(d.Get("revision").(int))

	obj := model.LBPool{
		DisplayName:            &displayName,
		Description:            &description,
		Tags:                   tags,
		ActiveMonitorPaths:     activeMonitorPaths,
		Algorithm:              &algorithm,
		MemberGroup:            memberGroup,
		Members:                members,
		MinActiveMembers:       &minActiveMembers,
		PassiveMonitorPath:     &passiveMonitorPath,
		SnatTranslation:        snatTranslation,
		TcpMultiplexingEnabled: &tcpMultiplexingEnabled,
		TcpMultiplexingNumber:  &tcpMultiplexingNumber,
		Revision:               &revision,
	}

	// Update the resource using PUT
	_, err = c.Update(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleUpdateError("VPC LBPool", id, err)
	}
	return resourceNsxtVpcLBPoolRead(d, m)
}

func resourceNsxtVpcLBPoolDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBPool ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPoolsClient(sessionContext, connector)

	force := true
	err := c.Delete(parents[0], parents[1], parents[2], id, &force)
	if err != nil {
		return handleDeleteError("VPC LBPool", id, err)
	}

	return nil
}
