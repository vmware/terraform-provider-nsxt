// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyEdgeNode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyEdgeNodeRead,

		Schema: map[string]*schema.Schema{
			"edge_cluster_path": getPolicyPathSchema(true, false, "Edge cluster Path"),
			"member_index": {
				Type:         schema.TypeInt,
				Description:  "Index of this node within edge cluster",
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the edge node to retrieve. Note: after the edge node is located, this attribute is overwritten with the value NSX reports as the node's id, which is actually its member index within the edge cluster rather than a stable unique identifier. Use path, unique_id or realization_id to reference this resource unambiguously.",
			},
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"unique_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique identifier assigned by the system",
			},
			"realization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID used to realize the entity",
			},
		},
	}
}

func dataSourceNsxtPolicyEdgeNodeRead(d *schema.ResourceData, m interface{}) error {
	// Read an edge node by name or id
	edgeClusterPath := d.Get("edge_cluster_path").(string)
	// Note - according to the documentation GetOkExists should be used
	// for bool types, but in this case it works and GetOk doesn't
	memberIndex, memberIndexSet := d.GetOkExists("member_index")

	query := make(map[string]string)
	query["parent_path"] = escapeSpecialCharacters(edgeClusterPath)
	if memberIndexSet {
		query["member_index"] = strconv.Itoa(memberIndex.(int))
	}
	obj, err := policyDataSourceResourceReadWithValidation(d, getPolicyConnector(m), getSessionContext(d, m), "PolicyEdgeNode", query, false)
	if err != nil {
		return err
	}
	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(obj, model.PolicyEdgeNodeBindingType())
	if len(errors) > 0 {
		return errors[0]
	}
	policyEdgeNode := dataValue.(model.PolicyEdgeNode)
	d.Set("member_index", policyEdgeNode.MemberIndex)
	d.Set("unique_id", policyEdgeNode.UniqueId)
	d.Set("realization_id", policyEdgeNode.RealizationId)
	return nil
}
