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

var vpcLbVirtualServerPathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-virtual-servers/[virtual server]"

func resourceNsxtVpcLBVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBVirtualServerCreate,
		Read:   resourceNsxtVpcLBVirtualServerRead,
		Update: resourceNsxtVpcLBVirtualServerUpdate,
		Delete: resourceNsxtVpcLBVirtualServerDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB Virtual Server", resourceNsxtVpcLBVirtualServerImport),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":                   getNsxIDSchema(),
			"path":                     getPathSchema(),
			"display_name":             getDisplayNameSchema(),
			"description":              getDescriptionSchema(),
			"revision":                 getRevisionSchema(),
			"tag":                      getTagsSchema(),
			"context":                  getContextSchemaExtended(true, false, true, true),
			"application_profile_path": getPolicyPathSchema(true, false, "Application profile for this virtual server"),
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable Virtual Server",
				Optional:    true,
				Default:     true,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Description:  "Virtual Server IP address",
				Required:     true,
				ValidateFunc: validateSingleIP(),
			},
			"ports": {
				Type:        schema.TypeList,
				Description: "Virtual Server ports",
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePortRange(),
				},
			},
			"persistence_profile_path": getPolicyPathSchema(false, false, "Path to persistence profile allowing related client connections to be sent to the same backend server."),
			"service_path":             getPolicyPathSchema(false, false, "Virtual Server can be associated with Load Balancer Service"),
			"default_pool_member_ports": {
				Type:        schema.TypeList,
				Description: "Default pool member ports when member port is not defined",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePortRange(),
				},
				MaxItems: 14,
			},
			"access_log_enabled": {
				Type:        schema.TypeBool,
				Description: "If enabled, all connections/requests sent to virtual server are logged to the access log file",
				Optional:    true,
				Default:     false,
			},
			"max_concurrent_connections": {
				Type:         schema.TypeInt,
				Description:  "To ensure one virtual server does not over consume resources, connections to a virtual server can be capped.",
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"max_new_connection_rate": {
				Type:         schema.TypeInt,
				Description:  "To ensure one virtual server does not over consume resources, connections to a member can be rate limited.",
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"pool_path":       getPolicyPathSchema(false, false, "Path for Load Balancer Pool"),
			"sorry_pool_path": getPolicyPathSchema(false, false, "When load balancer can not select server in default pool or pool in rules, the request would be served by sorry pool"),
			"client_ssl": {
				Type:        schema.TypeList,
				Description: "This setting is used when load balancer terminates client SSL connection",
				Elem:        getPolicyLbClientSSLBindingSchema(),
				Optional:    true,
				MaxItems:    1,
			},
			"server_ssl": {
				Type:        schema.TypeList,
				Description: "This setting is used when load balancer establishes connection to the backend server",
				Elem:        getPolicyLbServerSSLBindingSchema(),
				Optional:    true,
				MaxItems:    1,
			},
			"log_significant_event_only": {
				Type:        schema.TypeBool,
				Description: "Flag to log significant events in access log, if access log is enabled",
				Optional:    true,
				Default:     false,
			},
			"access_list_control": {
				Type:        schema.TypeList,
				Description: "IP access list control for filtering the connections from clients",
				Elem:        getPolicyLbAccessListControlSchema(),
				Optional:    true,
				MaxItems:    1,
			},
			"rule": {
				Type:        schema.TypeList,
				Description: "List of load balancer rules for layer 7 virtual servers with LBHttpProfile",
				Optional:    true,
				MaxItems:    4000,
				Elem:        getPolicyLbRuleBindingSchema(),
			},
		},
	}
}

func resourceNsxtVpcLBVirtualServerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	resources, err := getVpcPathResourceImporter(vpcLbVirtualServerPathExample)(d, m)
	if err != nil {
		return nil, err
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbVirtualServersClient(sessionContext, getPolicyConnector(m))
	obj, err := c.Get(parents[0], parents[1], parents[2], d.Id())
	if err != nil {
		return nil, logAPIError("Error importing VPC LBVirtualServer", err)
	}

	// Only an explicit import adopts rules that Terraform does not yet manage.
	setPolicyLbRulesInSchema(d, obj.Rules)
	return resources, nil
}

func resourceNsxtVpcLBVirtualServerExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbVirtualServersClient(sessionContext, connector)
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcLBVirtualServerCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Virtual Server resource requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbVirtualServersClient(sessionContext, connector)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBVirtualServerExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	accessLogEnabled := d.Get("access_log_enabled").(bool)
	applicationProfilePath := d.Get("application_profile_path").(string)
	clientSSLProfileBinding := getPolicyClientSSLBindingFromSchema(d)
	defaultPoolMemberPorts := getStringListFromSchemaList(d, "default_pool_member_ports")
	enabled := d.Get("enabled").(bool)
	ipAddress := d.Get("ip_address").(string)
	lbPersistenceProfilePath := d.Get("persistence_profile_path").(string)
	lbServicePath := d.Get("service_path").(string)
	maxConcurrentConnections := int64(d.Get("max_concurrent_connections").(int))
	maxNewConnectionRate := int64(d.Get("max_new_connection_rate").(int))
	poolPath := d.Get("pool_path").(string)
	ports := getStringListFromSchemaList(d, "ports")
	serverSSLProfileBinding := getPolicyServerSSLBindingFromSchema(d)
	sorryPoolPath := d.Get("sorry_pool_path").(string)
	rules := getPolicyLbRulesFromSchema(d)

	obj := model.LBVirtualServer{
		DisplayName:              &displayName,
		Description:              &description,
		Tags:                     tags,
		AccessLogEnabled:         &accessLogEnabled,
		ApplicationProfilePath:   &applicationProfilePath,
		ClientSslProfileBinding:  clientSSLProfileBinding,
		DefaultPoolMemberPorts:   defaultPoolMemberPorts,
		Enabled:                  &enabled,
		IpAddress:                &ipAddress,
		LbPersistenceProfilePath: &lbPersistenceProfilePath,
		LbServicePath:            &lbServicePath,
		PoolPath:                 &poolPath,
		Ports:                    ports,
		ServerSslProfileBinding:  serverSSLProfileBinding,
		SorryPoolPath:            &sorryPoolPath,
		Rules:                    rules,
	}

	policyLBVirtualServerVersionDependantSet(d, &obj)

	if maxNewConnectionRate > 0 {
		obj.MaxNewConnectionRate = &maxNewConnectionRate
	}

	if maxConcurrentConnections > 0 {
		obj.MaxConcurrentConnections = &maxConcurrentConnections
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating VPC LBVirtualServer with ID %s", id)
	err = c.Patch(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleCreateError("VPC LBVirtualServer", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcLBVirtualServerRead(d, m)
}

func resourceNsxtVpcLBVirtualServerRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbVirtualServersClient(sessionContext, connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBVirtualServer ID")
	}

	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBVirtualServer", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("access_log_enabled", obj.AccessLogEnabled)
	d.Set("application_profile_path", obj.ApplicationProfilePath)
	setPolicyClientSSLBindingInSchema(d, obj.ClientSslProfileBinding)
	d.Set("default_pool_member_ports", obj.DefaultPoolMemberPorts)
	d.Set("enabled", obj.Enabled)
	d.Set("ip_address", obj.IpAddress)
	d.Set("persistence_profile_path", obj.LbPersistenceProfilePath)
	d.Set("service_path", obj.LbServicePath)
	d.Set("max_concurrent_connections", obj.MaxConcurrentConnections)
	d.Set("max_new_connection_rate", obj.MaxNewConnectionRate)
	d.Set("pool_path", obj.PoolPath)
	d.Set("ports", obj.Ports)
	setPolicyServerSSLBindingInSchema(d, obj.ServerSslProfileBinding)
	d.Set("sorry_pool_path", obj.SorryPoolPath)
	d.Set("log_significant_event_only", obj.LogSignificantEventOnly)
	setPolicyAccessListControlInSchema(d, obj.AccessListControl)

	// Refresh only rules already managed by Terraform. The importer separately
	// adopts existing rules; an ordinary refresh must leave unmanaged rules alone.
	if _, ok := d.GetOk("rule"); ok {
		setPolicyLbRulesInSchema(d, obj.Rules)
	} else {
		log.Printf("[INFO] Ignoring rules since the user did not specify them in configuration")
	}

	return nil
}

func resourceNsxtVpcLBVirtualServerUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbVirtualServersClient(sessionContext, connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBVirtualServer ID")
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)

	accessLogEnabled := d.Get("access_log_enabled").(bool)
	clientSSLProfileBinding := getPolicyClientSSLBindingFromSchema(d)
	applicationProfilePath := d.Get("application_profile_path").(string)
	defaultPoolMemberPorts := getStringListFromSchemaList(d, "default_pool_member_ports")
	enabled := d.Get("enabled").(bool)
	ipAddress := d.Get("ip_address").(string)
	lbPersistenceProfilePath := d.Get("persistence_profile_path").(string)
	lbServicePath := d.Get("service_path").(string)
	maxConcurrentConnections := int64(d.Get("max_concurrent_connections").(int))
	maxNewConnectionRate := int64(d.Get("max_new_connection_rate").(int))
	poolPath := d.Get("pool_path").(string)
	ports := getStringListFromSchemaList(d, "ports")
	serverSSLProfileBinding := getPolicyServerSSLBindingFromSchema(d)
	sorryPoolPath := d.Get("sorry_pool_path").(string)
	rules := getPolicyLbRulesFromSchema(d)

	obj := model.LBVirtualServer{
		DisplayName:              &displayName,
		Description:              &description,
		Tags:                     tags,
		AccessLogEnabled:         &accessLogEnabled,
		ApplicationProfilePath:   &applicationProfilePath,
		ClientSslProfileBinding:  clientSSLProfileBinding,
		DefaultPoolMemberPorts:   defaultPoolMemberPorts,
		Enabled:                  &enabled,
		IpAddress:                &ipAddress,
		LbPersistenceProfilePath: &lbPersistenceProfilePath,
		LbServicePath:            &lbServicePath,
		PoolPath:                 &poolPath,
		Ports:                    ports,
		ServerSslProfileBinding:  serverSSLProfileBinding,
		SorryPoolPath:            &sorryPoolPath,
		Rules:                    rules,
	}

	policyLBVirtualServerVersionDependantSet(d, &obj)

	// Rules are only managed by the provider when the user configured them, otherwise rules which
	// were defined outside of Terraform are preserved:
	//   no rules in resource / no rules in state / no rules in NSX => no change
	//   no rules in resource / no rules in state /    rules in NSX => no change (legacy behaviour)
	//      rules in resource / no rules in state / no rules in NSX => rules will be created
	//      rules in resource /    rules in state /    rules in NSX => rules will be changed
	//   no rules in resource /    rules in state /    rules in NSX => rules will be deleted
	if !d.HasChange("rule") {
		log.Printf("[INFO] No changes detected in rule section")
		existingObj, err := c.Get(parents[0], parents[1], parents[2], id)
		if err != nil {
			return handleUpdateError("VPC LBVirtualServer", id, err)
		}

		if len(existingObj.Rules) > 0 {
			obj.Rules = existingObj.Rules
		}
	} else {
		log.Printf("[INFO] Changes detected in rule section")
	}

	if maxNewConnectionRate > 0 {
		obj.MaxNewConnectionRate = &maxNewConnectionRate
	}

	if maxConcurrentConnections > 0 {
		obj.MaxConcurrentConnections = &maxConcurrentConnections
	}

	// Update the resource using PATCH
	err := c.Patch(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleUpdateError("VPC LBVirtualServer", id, err)
	}

	return resourceNsxtVpcLBVirtualServerRead(d, m)
}

func resourceNsxtVpcLBVirtualServerDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBVirtualServer ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbVirtualServersClient(sessionContext, connector)

	force := true
	err := c.Delete(parents[0], parents[1], parents[2], id, &force)
	if err != nil {
		return handleDeleteError("VPC LBVirtualServer", id, err)
	}

	return nil
}
