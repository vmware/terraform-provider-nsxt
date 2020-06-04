/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
)

var lbClientSSLModeValues = []string{
	model.LBClientSslProfileBinding_CLIENT_AUTH_REQUIRED,
	model.LBClientSslProfileBinding_CLIENT_AUTH_IGNORE,
}

var lbServerSSLModeValues = []string{
	model.LBServerSslProfileBinding_SERVER_AUTH_REQUIRED,
	model.LBServerSslProfileBinding_SERVER_AUTH_IGNORE,
	model.LBServerSslProfileBinding_SERVER_AUTH_AUTO_APPLY,
}

var lbAccessListControlActionValues = []string{
	model.LBAccessListControl_ACTION_ALLOW,
	model.LBAccessListControl_ACTION_DROP,
}

func resourceNsxtPolicyLBVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyLBVirtualServerCreate,
		Read:   resourceNsxtPolicyLBVirtualServerRead,
		Update: resourceNsxtPolicyLBVirtualServerUpdate,
		Delete: resourceNsxtPolicyLBVirtualServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":                   getNsxIDSchema(),
			"path":                     getPathSchema(),
			"display_name":             getDisplayNameSchema(),
			"description":              getDescriptionSchema(),
			"revision":                 getRevisionSchema(),
			"tag":                      getTagsSchema(),
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
				Description: "Flag to log significant events in access log, if access log is enabed",
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
			// TODO: add rules
		},
	}
}

func getPolicyLbClientSSLBindingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_auth": {
				Type:         schema.TypeString,
				Description:  "Client authentication mode",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(lbClientSSLModeValues, false),
				Default:      model.LBClientSslProfileBinding_CLIENT_AUTH_IGNORE,
			},
			"certificate_chain_depth": {
				Type:         schema.TypeInt,
				Description:  "Certificate chain depth",
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"ca_paths": {
				Type:        schema.TypeList,
				Description: "If client auth type is REQUIRED, client certificate must be signed by one Certificate Authorities",
				Optional:    true,
				Elem:        getElemPolicyPathSchema(),
			},
			"crl_paths": {
				Type:        schema.TypeList,
				Description: "Certificate Revocation Lists can be specified to disallow compromised certificates",
				Optional:    true,
				Elem:        getElemPolicyPathSchema(),
			},
			"default_certificate_path": getPolicyPathSchema(true, false, "Default Certificate Path"),
			"sni_paths": {
				Type:        schema.TypeList,
				Description: "This setting allows multiple certificates, for different hostnames, to be bound to the same virtual server",
				Optional:    true,
				Elem:        getElemPolicyPathSchema(),
			},
			"ssl_profile_path": getPolicyPathSchema(false, false, "Client SSL Profile Path"),
		},
	}
}

func getPolicyLbServerSSLBindingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_auth": {
				Type:         schema.TypeString,
				Description:  "Server authentication mode",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(lbServerSSLModeValues, false),
				Default:      model.LBServerSslProfileBinding_SERVER_AUTH_AUTO_APPLY,
			},
			"client_certificate_path": getPolicyPathSchema(false, false, "Client certificat path for client authentication"),
			"certificate_chain_depth": {
				Type:         schema.TypeInt,
				Description:  "Certificate chain depth",
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"ca_paths": {
				Type:        schema.TypeList,
				Description: "If server auth type is REQUIRED, server certificate must be signed by one Certificate Authorities",
				Optional:    true,
				Elem:        getElemPolicyPathSchema(),
			},
			"crl_paths": {
				Type:        schema.TypeList,
				Description: "Certificate Revocation Lists can be specified disallow compromised certificates",
				Optional:    true,
				Elem:        getElemPolicyPathSchema(),
			},
			"ssl_profile_path": getPolicyPathSchema(false, false, "Server SSL Profile Path"),
		},
	}
}

func getPolicyLbAccessListControlSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:         schema.TypeString,
				Description:  "Action to apply to connections matching the group",
				Required:     true,
				ValidateFunc: validation.StringInSlice(lbAccessListControlActionValues, false),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable access list control",
				Optional:    true,
				Default:     true,
			},
			"group_path": getPolicyPathSchema(true, false, "The path of grouping object which defines the IP addresses or ranges to match the client IP"),
		},
	}
}

func getPolicyClientSSLBindingFromSchema(d *schema.ResourceData) *model.LBClientSslProfileBinding {
	bindings := d.Get("client_ssl").([]interface{})
	for _, binding := range bindings {
		// maximum count is one
		data := binding.(map[string]interface{})
		chainDepth := int64(data["certificate_chain_depth"].(int))
		clientAuth := data["client_auth"].(string)
		caList := interface2StringList(data["ca_paths"].([]interface{}))
		crlList := interface2StringList(data["crl_paths"].([]interface{}))
		certPath := data["default_certificate_path"].(string)
		profilePath := data["ssl_profile_path"].(string)
		profileBinding := model.LBClientSslProfileBinding{
			CertificateChainDepth:  &chainDepth,
			ClientAuth:             &clientAuth,
			ClientAuthCaPaths:      caList,
			ClientAuthCrlPaths:     crlList,
			DefaultCertificatePath: &certPath,
			SslProfilePath:         &profilePath,
		}

		return &profileBinding
	}

	return nil
}

func setPolicyClientSSLBindingInSchema(d *schema.ResourceData, binding *model.LBClientSslProfileBinding) error {
	var bindingList []map[string]interface{}
	if binding != nil {
		elem := make(map[string]interface{})
		if binding.CertificateChainDepth != nil {
			elem["certificate_chain_depth"] = *binding.CertificateChainDepth
		}
		if binding.ClientAuth != nil {
			elem["client_auth"] = *binding.ClientAuth
		}
		elem["ca_paths"] = binding.ClientAuthCaPaths
		elem["crl_paths"] = binding.ClientAuthCrlPaths
		elem["default_certificate_path"] = binding.DefaultCertificatePath
		if binding.SslProfilePath != nil {
			elem["ssl_profile_path"] = *binding.SslProfilePath
		}

		bindingList = append(bindingList, elem)
	}

	err := d.Set("client_ssl", bindingList)
	return err
}

func getPolicyServerSSLBindingFromSchema(d *schema.ResourceData) *model.LBServerSslProfileBinding {
	bindings := d.Get("server_ssl").([]interface{})
	for _, binding := range bindings {
		// maximum count is one
		data := binding.(map[string]interface{})
		chainDepth := int64(data["certificate_chain_depth"].(int))
		serverAuth := data["server_auth"].(string)
		caList := interface2StringList(data["ca_paths"].([]interface{}))
		crlList := interface2StringList(data["crl_paths"].([]interface{}))
		certPath := data["client_certificate_path"].(string)
		profilePath := data["ssl_profile_path"].(string)
		profileBinding := model.LBServerSslProfileBinding{
			CertificateChainDepth: &chainDepth,
			ServerAuth:            &serverAuth,
			ServerAuthCaPaths:     caList,
			ServerAuthCrlPaths:    crlList,
			ClientCertificatePath: &certPath,
			SslProfilePath:        &profilePath,
		}

		return &profileBinding
	}

	return nil
}

func setPolicyServerSSLBindingInSchema(d *schema.ResourceData, binding *model.LBServerSslProfileBinding) error {
	var bindingList []map[string]interface{}
	if binding != nil {
		elem := make(map[string]interface{})
		if binding.CertificateChainDepth != nil {
			elem["certificate_chain_depth"] = *binding.CertificateChainDepth
		}
		if binding.ServerAuth != nil {
			elem["server_auth"] = *binding.ServerAuth
		}
		elem["ca_paths"] = binding.ServerAuthCaPaths
		elem["crl_paths"] = binding.ServerAuthCrlPaths
		if binding.ClientCertificatePath != nil {
			elem["client_certificate_path"] = *binding.ClientCertificatePath
		}
		if binding.SslProfilePath != nil {
			elem["ssl_profile_path"] = *binding.SslProfilePath
		}

		bindingList = append(bindingList, elem)
	}

	err := d.Set("server_ssl", bindingList)
	return err
}

func getPolicyAccessListControlFromSchema(d *schema.ResourceData) *model.LBAccessListControl {
	controls := d.Get("access_list_control").([]interface{})
	for _, control := range controls {
		// maximum count is one
		data := control.(map[string]interface{})
		action := data["action"].(string)
		enabled := data["enabled"].(bool)
		groupPath := data["group_path"].(string)
		result := model.LBAccessListControl{
			Action:    &action,
			Enabled:   &enabled,
			GroupPath: &groupPath,
		}

		return &result
	}

	return nil
}

func setPolicyAccessListControlInSchema(d *schema.ResourceData, control *model.LBAccessListControl) error {
	var controlList []map[string]interface{}
	if control != nil {
		elem := make(map[string]interface{})
		elem["action"] = control.Action
		elem["enabled"] = control.Enabled
		elem["group_path"] = control.GroupPath

		controlList = append(controlList, elem)
	}

	err := d.Set("access_list_control", controlList)
	return err
}

func policyLBVirtualServerVersionDepenantSet(d *schema.ResourceData, obj *model.LBVirtualServer) {
	if nsxVersionHigherOrEqual("3.0.0") {
		logSignificantOnly := d.Get("log_significant_event_only").(bool)
		obj.LogSignificantEventOnly = &logSignificantOnly
		obj.AccessListControl = getPolicyAccessListControlFromSchema(d)
	}
}

func resourceNsxtPolicyLBVirtualServerExists(id string, connector *client.RestConnector) bool {
	client := infra.NewDefaultLbVirtualServersClient(connector)

	_, err := client.Get(id)
	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving resource", err)

	return false
}

func resourceNsxtPolicyLBVirtualServerCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultLbVirtualServersClient(connector)

	if client == nil {
		return policyResourceNotSupportedError()
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, connector, resourceNsxtPolicyLBVirtualServerExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
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
	}

	policyLBVirtualServerVersionDepenantSet(d, &obj)

	if maxNewConnectionRate > 0 {
		obj.MaxNewConnectionRate = &maxNewConnectionRate
	}

	if maxConcurrentConnections > 0 {
		obj.MaxConcurrentConnections = &maxConcurrentConnections
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating LBVirtualServer with ID %s", id)
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("LBVirtualServer", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyLBVirtualServerRead(d, m)
}

func resourceNsxtPolicyLBVirtualServerRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultLbVirtualServersClient(connector)

	if client == nil {
		return policyResourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBVirtualServer ID")
	}

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "LBVirtualServer", id, err)
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

	return nil
}

func resourceNsxtPolicyLBVirtualServerUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultLbVirtualServersClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBVirtualServer ID")
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

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
	}

	policyLBVirtualServerVersionDepenantSet(d, &obj)

	if maxNewConnectionRate > 0 {
		obj.MaxNewConnectionRate = &maxNewConnectionRate
	}

	if maxConcurrentConnections > 0 {
		obj.MaxConcurrentConnections = &maxConcurrentConnections
	}

	// Update the resource using PATCH
	err := client.Patch(id, obj)
	if err != nil {
		return handleUpdateError("LBVirtualServer", id, err)
	}

	return resourceNsxtPolicyLBVirtualServerRead(d, m)
}

func resourceNsxtPolicyLBVirtualServerDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBVirtualServer ID")
	}

	connector := getPolicyConnector(m)
	client := infra.NewDefaultLbVirtualServersClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	force := true
	err := client.Delete(id, &force)
	if err != nil {
		return handleDeleteError("LBVirtualServer", id, err)
	}

	return nil
}
