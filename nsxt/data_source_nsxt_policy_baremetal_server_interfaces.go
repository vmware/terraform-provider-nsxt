// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyBareMetalServerInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyBareMetalServerInterfacesRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"display_name": {
				Type:        schema.TypeString,
				Description: "Regex to filter bare metal server interfaces by display name",
				Optional:    true,
			},
			"bms_external_id": {
				Type:        schema.TypeString,
				Description: "Filter by parent bare metal server external ID",
				Optional:    true,
			},
			"source_id": {
				Type:        schema.TypeString,
				Description: "Filter by source (bare metal controller) id",
				Optional:    true,
			},
			"tag_scope": {
				Type:        schema.TypeString,
				Description: "Filter by tag scope (matches Tag.scope)",
				Optional:    true,
			},
			"tag": {
				Type:        schema.TypeString,
				Description: "Filter by tag value (matches Tag.tag)",
				Optional:    true,
			},
			"results": {
				Type:        schema.TypeList,
				Description: "Bare metal server interfaces matching the filters",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_id": {
							Type:        schema.TypeString,
							Description: "External ID of the bare metal server interface",
							Computed:    true,
						},
						"display_name": {
							Type:        schema.TypeString,
							Description: "Display name of the bare metal server interface",
							Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceNsxtPolicyBareMetalServerInterfacesRead(d *schema.ResourceData, m interface{}) error {
	if err := validateBMSVersionRequirement(); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}

	// Use enhanced search with server-side filtering
	results, err := listInventoryResourcesByType(connector, ctx, "BareMetalServerInterface", buildBareMetalTagsFilter(d))
	if err != nil {
		return fmt.Errorf("error listing bare metal server interfaces: %v", err)
	}

	all, err := convertSearchResultToBareMetalServerInterfaceList(results)
	if err != nil {
		return err
	}

	// Apply additional client-side filters
	var filtered []model.BareMetalServerInterface

	var re *regexp.Regexp
	if v, ok := d.GetOk("display_name"); ok {
		re, err = regexp.Compile(v.(string))
		if err != nil {
			return fmt.Errorf("invalid regex for display_name: %v", err)
		}
	}

	bmsExternalID := d.Get("bms_external_id").(string)
	sourceID := d.Get("source_id").(string)

	for _, iface := range all {
		match := true

		// Display name regex filter
		if re != nil && iface.DisplayName != nil {
			if !re.MatchString(*iface.DisplayName) {
				match = false
			}
		}

		// Parent server filter
		if bmsExternalID != "" && (iface.BmsExternalId == nil || *iface.BmsExternalId != bmsExternalID) {
			match = false
		}

		// Source ID filter
		if sourceID != "" && (iface.SourceId == nil || *iface.SourceId != sourceID) {
			match = false
		}

		if match {
			filtered = append(filtered, iface)
		}
	}

	// Convert to schema format
	var items []map[string]interface{}
	for _, iface := range filtered {
		item := map[string]interface{}{
			"external_id":     getStringValue(iface.ExternalId),
			"display_name":    getStringValue(iface.DisplayName),
			"bms_external_id": getStringValue(iface.BmsExternalId),
			"source_id":       getStringValue(iface.SourceId),
			"state":           getStringValue(iface.State),
			"resource_type":   getStringValue(iface.ResourceType),
			"last_sync_time":  getInt64Value(iface.LastSyncTime),
		}

		if iface.IsMgmtInterface != nil {
			item["is_mgmt_interface"] = *iface.IsMgmtInterface
		}

		// IP addresses
		if iface.IpAddresses != nil {
			item["ip_addresses"] = iface.IpAddresses
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
			item["tags"] = tags
		}

		items = append(items, item)
	}

	d.SetId(newUUID())
	d.Set("results", items)
	return nil
}

// convertSearchResultToBareMetalServerInterfaceList converts search results to BareMetalServerInterface list
func convertSearchResultToBareMetalServerInterfaceList(searchResults []*data.StructValue) ([]model.BareMetalServerInterface, error) {
	var interfaces []model.BareMetalServerInterface
	converter := bindings.NewTypeConverter()

	for _, item := range searchResults {
		dv, errs := converter.ConvertToGolang(item, model.BareMetalServerInterfaceBindingType())
		if len(errs) > 0 {
			return interfaces, errs[0]
		}
		interfaces = append(interfaces, dv.(model.BareMetalServerInterface))
	}
	return interfaces, nil
}
