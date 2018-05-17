/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/loadbalancer"
	"log"
	"net/http"
)

var poolAlgTypeValues = []string{"ROUND_ROBIN", "WEIGHTED_ROUND_ROBIN", "LEAST_CONNECTION", "WEIGHTED_LEAST_CONNECTION", "IP_HASH"}
var poolSnatTranslationTypeValues = []string{"LbSnatAutoMap", "LbSnatIpPool", "Transparent"}
var memberAdminStateTypeValues = []string{"ENABLED", "DISABLED", "GRACEFUL_DISABLED"}

func resourceNsxtLbPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLbPoolCreate,
		Read:   resourceNsxtLbPoolRead,
		Update: resourceNsxtLbPoolUpdate,
		Delete: resourceNsxtLbPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"algorithm": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Load balancing algorithm controls how the incoming connections are distributed among the members",
				ValidateFunc: validation.StringInSlice(poolAlgTypeValues, false),
				Optional:     true,
				Default:      "ROUND_ROBIN",
			},
			"min_active_members": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The minimum number of members for the pool to be considered active",
				Optional:    true,
				Default:     1,
			},
			"tcp_multiplexing_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "TCP multiplexing allows the same TCP connection between load balancer and the backend server to be used for sending multiple client requests from different client TCP connections",
				Optional:    true,
				Default:     false,
			},
			"tcp_multiplexing_number": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The maximum number of TCP connections per pool that are idly kept alive for sending future client requests",
				Optional:    true,
				Default:     6,
			},
			"active_monitor_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Active health monitor Id. If one is not set, the active healthchecks will be disabled",
				Optional:    true,
			},
			"passive_monitor_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Passive health monitor Id. If one is not set, the passive healthchecks will be disabled",
				Optional:    true,
			},
			"snat_translation_type": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Type of SNAT performed to ensure reverse traffic from the server can be received and processed by the loadbalancer",
				ValidateFunc: validation.StringInSlice(poolSnatTranslationTypeValues, false),
				Optional:     true,
				Default:      "Transparent",
			},
			"snat_translation_ip": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Ip address or Ip range for SNAT of type LbSnatIpPool",
				ValidateFunc: validateIPOrRange(),
				Optional:     true,
			},
			"member": getPoolMembersSchema(),
		},
	}
}

func getPoolMembersSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of server pool members. Each pool member is identified, typically, by an IP address and a port",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": &schema.Schema{
					Type:        schema.TypeString,
					Description: "Pool member name",
					Optional:    true,
					Computed:    true,
				},
				"admin_state": &schema.Schema{
					Type:         schema.TypeString,
					Description:  "Member admin state",
					Optional:     true,
					ValidateFunc: validation.StringInSlice(memberAdminStateTypeValues, false),
				},
				"backup_member": &schema.Schema{
					Type:        schema.TypeBool,
					Description: "A boolean flag which reflects whether this is a backup pool member",
					Optional:    true,
					Default:     false,
				},
				"ip_address": &schema.Schema{
					Type:         schema.TypeString,
					Description:  "Pool member IP address",
					Required:     true,
					ValidateFunc: validateSingleIP(),
				},
				"max_concurrent_connections": &schema.Schema{
					Type:        schema.TypeInt,
					Description: "To ensure members are not overloaded, connections to a member can be capped by the load balancer. When a member reaches this limit, it is skipped during server selection. If it is not specified, it means that connections are unlimited",
					Optional:    true,
				},
				"port": &schema.Schema{
					Type:         schema.TypeString,
					Description:  "If port is specified, all connections will be sent to this port. Only single port is supported. If unset, the same port the client connected to will be used, it could be overrode by default_pool_member_port setting in virtual server. The port should not specified for port range case",
					Optional:     true,
					ValidateFunc: validateSinglePort(),
				},
				"weight": &schema.Schema{
					Type:        schema.TypeInt,
					Description: "Pool member weight is used for WEIGHTED_ROUND_ROBIN balancing algorithm. The weight value would be ignored in other algorithms",
					Optional:    true,
				},
			},
		},
	}
}

func getActiveMonitorIdsFromSchema(d *schema.ResourceData) []string {
	activeMonitorId := d.Get("active_monitor_id").(string)
	var activeMonitorIds []string
	if activeMonitorId != "" {
		activeMonitorIds = append(activeMonitorIds, activeMonitorId)
	}
	return activeMonitorIds
}

func getSnatTranslationFromSchema(d *schema.ResourceData) *loadbalancer.LbSnatTranslation {
	// Transparent type should not create an object
	trType := d.Get("snat_translation_type").(string)
	if trType == "Transparent" {
		return nil
	}
	if trType == "LbSnatIpPool" {
		ipAddresses := make([]loadbalancer.LbSnatIpElement, 0, 1)
		ipAddress := d.Get("snat_translation_ip").(string)
		elem := loadbalancer.LbSnatIpElement{IpAddress: ipAddress}
		ipAddresses = append(ipAddresses, elem)

		return &loadbalancer.LbSnatTranslation{
			Type_:       trType,
			IpAddresses: ipAddresses,
		}
	}
	// For LbSnatAutoMap type
	return &loadbalancer.LbSnatTranslation{Type_: trType}
}

func setPoolMembersInSchema(d *schema.ResourceData, members []loadbalancer.PoolMember) error {
	var membersList []map[string]interface{}
	for _, member := range members {
		elem := make(map[string]interface{})
		elem["display_name"] = member.DisplayName
		elem["admin_state"] = member.AdminState
		elem["backup_member"] = member.BackupMember
		elem["ip_address"] = member.IpAddress
		elem["max_concurrent_connections"] = member.MaxConcurrentConnections
		elem["port"] = member.Port
		elem["weight"] = member.Weight

		membersList = append(membersList, elem)
	}
	err := d.Set("member", membersList)
	return err
}

func getPoolMembersFromSchema(d *schema.ResourceData) []loadbalancer.PoolMember {
	members := d.Get("member").([]interface{})
	var memberList []loadbalancer.PoolMember
	for _, member := range members {
		data := member.(map[string]interface{})
		elem := loadbalancer.PoolMember{
			AdminState:               data["admin_state"].(string),
			BackupMember:             data["backup_member"].(bool),
			DisplayName:              data["display_name"].(string),
			IpAddress:                data["ip_address"].(string),
			MaxConcurrentConnections: int64(data["max_concurrent_connections"].(int)),
			Port:   data["port"].(string),
			Weight: int64(data["weight"].(int)),
		}

		memberList = append(memberList, elem)
	}
	return memberList
}

func resourceNsxtLbPoolCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	activeMonitorIds := getActiveMonitorIdsFromSchema(d)
	passiveMonitorId := d.Get("passive_monitor_id").(string)
	algorithm := d.Get("algorithm").(string)
	members := getPoolMembersFromSchema(d)
	minActiveMembers := int64(d.Get("min_active_members").(int))
	tcpMultiplexingEnabled := d.Get("tcp_multiplexing_enabled").(bool)
	tcpMultiplexingNumber := int64(d.Get("tcp_multiplexing_number").(int))
	snatTranslation := getSnatTranslationFromSchema(d)
	lbPool := loadbalancer.LbPool{
		Description:            description,
		DisplayName:            displayName,
		Tags:                   tags,
		ActiveMonitorIds:       activeMonitorIds,
		Algorithm:              algorithm,
		MinActiveMembers:       minActiveMembers,
		PassiveMonitorId:       passiveMonitorId,
		SnatTranslation:        snatTranslation,
		TcpMultiplexingEnabled: tcpMultiplexingEnabled,
		TcpMultiplexingNumber:  tcpMultiplexingNumber,
		Members:                members,
	}

	lbPool, resp, err := nsxClient.ServicesApi.CreateLoadBalancerPool(nsxClient.Context, lbPool)

	if err != nil {
		return fmt.Errorf("Error during LbPool create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during LbPool create: %v", resp.StatusCode)
	}
	d.SetId(lbPool.Id)

	return resourceNsxtLbPoolRead(d, m)
}

func resourceNsxtLbPoolRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	lbPool, resp, err := nsxClient.ServicesApi.ReadLoadBalancerPool(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbPool %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LbPool read: %v", err)
	}

	d.Set("revision", lbPool.Revision)
	d.Set("description", lbPool.Description)
	d.Set("display_name", lbPool.DisplayName)
	setTagsInSchema(d, lbPool.Tags)
	if lbPool.ActiveMonitorIds != nil && len(lbPool.ActiveMonitorIds) > 0 {
		d.Set("active_monitor_id", lbPool.ActiveMonitorIds[0])
	} else {
		d.Set("active_monitor_id", "")
	}
	d.Set("passive_monitor_id", lbPool.PassiveMonitorId)
	d.Set("algorithm", lbPool.Algorithm)
	err = setPoolMembersInSchema(d, lbPool.Members)
	if err != nil {
		return fmt.Errorf("Error during LB Pool members set in schema: %v", err)
	}
	d.Set("min_active_members", lbPool.MinActiveMembers)
	if lbPool.SnatTranslation != nil {
		d.Set("snat_translation_type", lbPool.SnatTranslation.Type_)
		if lbPool.SnatTranslation.IpAddresses != nil && len(lbPool.SnatTranslation.IpAddresses) > 0 {
			d.Set("snat_translation_ip", lbPool.SnatTranslation.IpAddresses[0].IpAddress)
		}
	} else {
		d.Set("snat_translation_type", "Transparent")
	}
	d.Set("tcp_multiplexing_enabled", lbPool.TcpMultiplexingEnabled)
	d.Set("tcp_multiplexing_number", lbPool.TcpMultiplexingNumber)

	return nil
}

func resourceNsxtLbPoolUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int32(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	activeMonitorIds := getActiveMonitorIdsFromSchema(d)
	passiveMonitorId := d.Get("passive_monitor_id").(string)
	algorithm := d.Get("algorithm").(string)
	members := getPoolMembersFromSchema(d)
	minActiveMembers := int64(d.Get("min_active_members").(int))
	snatTranslation := getSnatTranslationFromSchema(d)
	tcpMultiplexingEnabled := d.Get("tcp_multiplexing_enabled").(bool)
	tcpMultiplexingNumber := int64(d.Get("tcp_multiplexing_number").(int))
	lbPool := loadbalancer.LbPool{
		Revision:               revision,
		Description:            description,
		DisplayName:            displayName,
		Tags:                   tags,
		ActiveMonitorIds:       activeMonitorIds,
		Algorithm:              algorithm,
		MinActiveMembers:       minActiveMembers,
		PassiveMonitorId:       passiveMonitorId,
		SnatTranslation:        snatTranslation,
		TcpMultiplexingEnabled: tcpMultiplexingEnabled,
		TcpMultiplexingNumber:  tcpMultiplexingNumber,
		Members:                members,
	}

	lbPool, resp, err := nsxClient.ServicesApi.UpdateLoadBalancerPool(nsxClient.Context, id, lbPool)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LbPool update: %v", err)
	}

	return resourceNsxtLbPoolRead(d, m)
}

func resourceNsxtLbPoolDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.ServicesApi.DeleteLoadBalancerPool(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbPool delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbPool %s not found", id)
		d.SetId("")
	}
	return nil
}
