// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var rtepAssignments = []string{
	"static_ipv4_list",
	"static_ipv4_pool",
}

func resourceNsxtPolicyEdgeTransportNodeRTEP() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyEdgeTransportNodeRTEPCreate,
		Read:   resourceNsxtPolicyEdgeTransportNodeRTEPRead,
		Update: resourceNsxtPolicyEdgeTransportNodeRTEPUpdate,
		Delete: resourceNsxtPolicyEdgeTransportNodeRTEPDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyEdgeTransportNodeRTEPImporter,
		},
		Schema: map[string]*schema.Schema{
			"edge_transport_node_path": {
				Type:        schema.TypeString,
				Description: "Policy path of Edge transport node to associate with remote tunnel endpoint.",
				Required:    true,
				ForceNew:    true,
			},
			"host_switch_name": {
				Type:        schema.TypeString,
				Description: "The host switch name to be used for the remote tunnel endpoint",
				Required:    true,
				ForceNew:    true,
			},
			"ip_assignment": getPolicyIPAssignmentSchema(false, 1, 1, rtepAssignments),
			"named_teaming_policy": {
				Type:        schema.TypeString,
				Description: "The named teaming policy to be used by the remote tunnel endpoint",
				Required:    true,
			},
			"vlan": {
				Type:         schema.TypeInt,
				Description:  "VLAN id for remote tunnel endpoint",
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 4094),
			},
		},
	}
}

func setNsxtPolicyEdgeTransportNodeRTEP(d *schema.ResourceData, m interface{}, op string) error {
	connector := getPolicyConnector(m)

	tnPath := d.Get("edge_transport_node_path").(string)
	client := enforcement_points.NewEdgeTransportNodesClient(connector)

	siteID, epID, edgeID := getEdgeTransportNodeKeysFromPath(tnPath)
	obj, err := client.Get(siteID, epID, edgeID)
	if err != nil {
		return err
	}

	hswName := d.Get("host_switch_name").(string)
	found := false

	for i, hsw := range obj.SwitchSpec.Switches {
		if *hsw.SwitchName == hswName {
			{
				found = true
				if len(hsw.RemoteTunnelEndpoint) > 0 && op == "create" {
					return fmt.Errorf("remote tunnel endpoint for Edge transport node %s already exists", tnPath)
				}

				var rteps []model.PolicyEdgeTransportNodeRtepConfig

				if op != "delete" {
					ipAssignments, err := getPolicyIPAssignmentsFromSchema(d.Get("ip_assignment"))
					if err != nil {
						return err
					}
					namedTeamingPolicy := d.Get("named_teaming_policy").(string)
					vlan := int64(d.Get("vlan").(int))
					rteps = []model.PolicyEdgeTransportNodeRtepConfig{
						{
							IpAssignmentSpecs:  ipAssignments,
							NamedTeamingPolicy: &namedTeamingPolicy,
							Vlan:               &vlan,
						},
					}
				}

				obj.SwitchSpec.Switches[i].RemoteTunnelEndpoint = rteps

				if op == "create" {
					err = client.Patch(siteID, epID, edgeID, obj)
				} else {
					_, err = client.Update(siteID, epID, edgeID, obj)
				}
				if err != nil {
					return err
				}

				if op == "create" {
					d.SetId(fmt.Sprintf("%s:%s", tnPath, hswName))
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("switch %s not found for Edge transport node %s", hswName, tnPath)
	}

	if op == "delete" {
		return nil
	}
	return resourceNsxtPolicyEdgeTransportNodeRTEPRead(d, m)
}

func resourceNsxtPolicyEdgeTransportNodeRTEPCreate(d *schema.ResourceData, m interface{}) error {
	return setNsxtPolicyEdgeTransportNodeRTEP(d, m, "create")
}

func resourceNsxtPolicyEdgeTransportNodeRTEPRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	tnPath := d.Get("edge_transport_node_path").(string)
	client := enforcement_points.NewEdgeTransportNodesClient(connector)

	siteID, epID, edgeID := getEdgeTransportNodeKeysFromPath(tnPath)
	obj, err := client.Get(siteID, epID, edgeID)
	if err != nil {
		return err
	}

	hswName := d.Get("host_switch_name").(string)

	for _, hsw := range obj.SwitchSpec.Switches {
		if *hsw.SwitchName == hswName {
			{
				if len(hsw.RemoteTunnelEndpoint) == 0 {
					return errors.NotFound{}
				}

				ipAssignment, err := setPolicyIPAssignmentsInSchema(hsw.RemoteTunnelEndpoint[0].IpAssignmentSpecs)
				if err != nil {
					return err
				}

				// Only one RTEP is supported
				d.Set("ip_assignment", ipAssignment)
				d.Set("named_teaming_policy", hsw.RemoteTunnelEndpoint[0].NamedTeamingPolicy)
				d.Set("vlan", hsw.RemoteTunnelEndpoint[0].Vlan)

				return nil
			}
		}
	}

	return errors.NotFound{}
}

func resourceNsxtPolicyEdgeTransportNodeRTEPUpdate(d *schema.ResourceData, m interface{}) error {
	return setNsxtPolicyEdgeTransportNodeRTEP(d, m, "update")
}

func resourceNsxtPolicyEdgeTransportNodeRTEPDelete(d *schema.ResourceData, m interface{}) error {
	return setNsxtPolicyEdgeTransportNodeRTEP(d, m, "delete")
}

func resourceNsxtPolicyEdgeTransportNodeRTEPImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	segs := strings.Split(importID, ":")
	if len(segs) != 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("import parameter %s should have a policy path to PolicyEdgeTransportNode and a switch_name, separated by ':'", importID)
	}
	tnPath := segs[0]

	err := validateImportPolicyPath(tnPath)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	hswName := segs[1]

	d.Set("edge_transport_node_path", tnPath)
	d.Set("host_switch_name", hswName)

	d.SetId(fmt.Sprintf("%s:%s", tnPath, hswName))

	return []*schema.ResourceData{d}, nil
}
