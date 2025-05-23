// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyGatewayInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyGatewayInterfaceRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"gateway_path": {
				Type:        schema.TypeString,
				Description: "The name of the gateway to which interface is linked",
				Optional:    true,
			},
			"service_path": {
				Type:        schema.TypeString,
				Description: "The name of the locale service of the gateway to which interface is linked",
				Optional:    true,
			},
			"path": getPathSchema(),
			"edge_cluster_path": {
				Type:        schema.TypeString,
				Description: "The path of the edge cluster connected to the gateway linked to this interface. This is exported only for Tier0 gateways",
				Optional:    true,
				Computed:    false,
			},
			"segment_path": {
				Type:        schema.TypeString,
				Description: "Policy path for segment to be connected with the Gateway.",
				Optional:    true,
				Computed:    false,
			},
		},
	}
}

func dataSourceNsxtPolicyGatewayInterfaceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()
	var dataValue interface{}
	var path *string
	var edgePath *string
	var description *string
	var segmentPath *string
	var errors []error
	var id *string
	t0Gw := d.Get("gateway_path").(string)
	servicePath := d.Get("service_path").(string)
	if t0Gw == "" && servicePath == "" {
		return fmt.Errorf("One of gateway_path or service_path should be set")
	}
	query := make(map[string]string)
	if servicePath != "" {
		query["parent_path"] = servicePath
	} else {
		query["parent_path"] = t0Gw + "/locale-services/*"
	}
	var isT0 bool
	var err error
	if servicePath != "" {
		isT0, err = isT0Gw(servicePath)
	} else {
		isT0, err = isT0Gw(t0Gw)
	}
	if err != nil {
		return err
	}
	searchStr := "Tier1Interface"
	if isT0 {
		searchStr = "Tier0Interface"
	}
	obj, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), searchStr, query)
	if err != nil {
		return err
	}

	if isT0 {
		dataValue, errors = converter.ConvertToGolang(obj, model.Tier0InterfaceBindingType())
		currGwInterface := dataValue.(model.Tier0Interface)
		path = currGwInterface.Path
		edgePath = currGwInterface.EdgePath
		description = currGwInterface.Description
		segmentPath = currGwInterface.SegmentPath
		id = currGwInterface.Id
	} else {
		dataValue, errors = converter.ConvertToGolang(obj, model.Tier1InterfaceBindingType())
		currGwInterface := dataValue.(model.Tier1Interface)
		path = currGwInterface.Path
		description = currGwInterface.Description
		segmentPath = currGwInterface.SegmentPath
		id = currGwInterface.Id
	}
	if len(errors) > 0 {
		return errors[0]
	}
	d.Set("path", *path)
	if isT0 && edgePath != nil {
		err = d.Set("edge_cluster_path", *edgePath)
		if err != nil {
			return fmt.Errorf("Error while setting the interface edge cluster path : %v", err)
		}
	}
	d.Set("description", *description)
	d.Set("segment_path", *segmentPath)

	d.SetId(*id)
	return nil
}

func isT0Gw(path string) (bool, error) {
	segs := strings.Split(path, "/")
	if len(segs) < 3 {
		return false, fmt.Errorf("Not a valid gateway path")
	}
	if (len(segs) > 4 && segs[len(segs)-4] == "tier-0s") || segs[len(segs)-2] == "tier-0s" {
		return true, nil
	}
	return false, nil
}
