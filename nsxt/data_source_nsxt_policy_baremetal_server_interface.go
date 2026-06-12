// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyBareMetalServerInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyBareMetalServerInterfaceRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"external_id": {
				Type:        schema.TypeString,
				Description: "External ID of the bare metal server interface (exact match)",
				Optional:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name of the bare metal server interface (prefix match)",
				Optional:    true,
			},
			"bms_external_id": {
				Type:        schema.TypeString,
				Description: "External ID of the parent bare metal server",
				Computed:    true,
			},
			"source_id": {
				Type:        schema.TypeString,
				Description: "Source ID of the bare metal server interface",
				Computed:    true,
			},
			"is_mgmt_interface": {
				Type:        schema.TypeBool,
				Description: "Whether this is a management interface",
				Computed:    true,
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Description: "IP addresses assigned to the interface",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"state": {
				Type:        schema.TypeString,
				Description: "State of the interface (UP, DOWN, etc.)",
				Computed:    true,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Description: "Resource type",
				Computed:    true,
			},
			"last_sync_time": {
				Type:        schema.TypeInt,
				Description: "Last synchronization time",
				Computed:    true,
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "List of tags applied to the interface",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope": {
							Type:        schema.TypeString,
							Description: "Tag scope",
							Computed:    true,
						},
						"tag": {
							Type:        schema.TypeString,
							Description: "Tag value",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyBareMetalServerInterfaceRead(d *schema.ResourceData, m interface{}) error {
	if err := validateBMSVersionRequirement(); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}

	externalID := d.Get("external_id").(string)
	displayName := d.Get("display_name").(string)

	if externalID == "" && displayName == "" {
		return fmt.Errorf("either external_id or display_name must be specified")
	}

	var iface model.BareMetalServerInterface
	var err error

	if externalID != "" {
		// Find by external ID (exact match)
		iface, err = findBareMetalServerInterfaceByExternalID(connector, ctx, externalID)
		if err != nil {
			return fmt.Errorf("Failed to find bare metal server interface with external_id '%s': %v", externalID, err)
		}
	} else {
		// Find by display name (prefix match returning exactly one result)
		results, err := listInventoryResourcesByType(connector, ctx, "BareMetalServerInterface", nil)
		if err != nil {
			return fmt.Errorf("error listing bare metal server interfaces: %v", err)
		}

		interfaces, err := convertSearchResultToBareMetalServerInterfaceList(results)
		if err != nil {
			return err
		}

		var matches []model.BareMetalServerInterface
		for _, intf := range interfaces {
			if intf.DisplayName != nil && *intf.DisplayName == displayName {
				matches = append(matches, intf)
			}
		}

		if len(matches) == 0 {
			return fmt.Errorf("No bare metal server interface found with display_name '%s'", displayName)
		}
		if len(matches) > 1 {
			return fmt.Errorf("Multiple bare metal server interfaces found with display_name '%s'. Use external_id for exact match", displayName)
		}
		iface = matches[0]
	}

	// Set attributes
	d.SetId(*iface.ExternalId)
	d.Set("external_id", iface.ExternalId)
	d.Set("display_name", iface.DisplayName)
	d.Set("bms_external_id", iface.BmsExternalId)
	d.Set("source_id", iface.SourceId)
	d.Set("state", iface.State)
	d.Set("resource_type", iface.ResourceType)
	d.Set("last_sync_time", iface.LastSyncTime)

	if iface.IsMgmtInterface != nil {
		d.Set("is_mgmt_interface", *iface.IsMgmtInterface)
	}

	// IP addresses
	if iface.IpAddresses != nil {
		d.Set("ip_addresses", iface.IpAddresses)
	}

	// Tags
	if iface.Tags != nil {
		var tags []map[string]interface{}
		for _, tag := range iface.Tags {
			tagItem := map[string]interface{}{
				"scope": getStringValue(tag.Scope),
				"tag":   getStringValue(tag.Tag),
			}
			tags = append(tags, tagItem)
		}
		d.Set("tags", tags)
	}

	return nil
}

// findBareMetalServerInterfaceByExternalID finds a BMS interface by external ID using search
func findBareMetalServerInterfaceByExternalID(connector client.Connector, ctx utl.SessionContext, externalID string) (model.BareMetalServerInterface, error) {
	var found model.BareMetalServerInterface

	resultValues, err := listInventoryResourcesByAnyFieldAndType(connector, ctx, externalID, "BareMetalServerInterface", nil)
	if err != nil {
		return found, err
	}

	converter := bindings.NewTypeConverter()
	for _, item := range resultValues {
		dv, errs := converter.ConvertToGolang(item, model.BareMetalServerInterfaceBindingType())
		if len(errs) > 0 {
			return found, errs[0]
		}
		bmsi := dv.(model.BareMetalServerInterface)
		if bmsi.ExternalId != nil && *bmsi.ExternalId == externalID {
			return bmsi, nil
		}
	}
	return found, fmt.Errorf("could not find bare metal server interface with external_id: %s", externalID)
}
