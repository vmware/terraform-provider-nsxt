// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyBareMetalServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyBareMetalServersRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"display_name": {
				Type:        schema.TypeString,
				Description: "Regex to filter bare metal servers by display name",
				Optional:    true,
			},
			"source_id": {
				Type:        schema.TypeString,
				Description: "Filter by source (bare metal controller) id",
				Optional:    true,
			},
			"os_name": {
				Type:        schema.TypeString,
				Description: "Filter by OS name (case-insensitive substring)",
				Optional:    true,
			},
			"os_version": {
				Type:        schema.TypeString,
				Description: "Filter by OS version (case-insensitive substring)",
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
				Description: "Bare metal servers matching the filters",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_id": {
							Type:        schema.TypeString,
							Description: "External ID of the bare metal server",
							Computed:    true,
						},
						"display_name": {
							Type:        schema.TypeString,
							Description: "Display name of the bare metal server",
							Computed:    true,
						},
						"source_id": {
							Type:        schema.TypeString,
							Description: "Source ID of the bare metal server",
							Computed:    true,
						},
						"cpu_cores": {
							Type:        schema.TypeInt,
							Description: "Number of CPU cores",
							Computed:    true,
						},
						"os_name": {
							Type:        schema.TypeString,
							Description: "Operating system name",
							Computed:    true,
						},
						"os_version": {
							Type:        schema.TypeString,
							Description: "Operating system version",
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
							Description: "List of tags applied to the server",
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

func dataSourceNsxtPolicyBareMetalServersRead(d *schema.ResourceData, m interface{}) error {
	if err := validateBMSVersionRequirement(); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}

	// Use enhanced search with server-side filtering
	results, err := listInventoryResourcesByType(connector, ctx, "BareMetalServer", buildBareMetalTagsFilter(d))
	if err != nil {
		return fmt.Errorf("error listing bare metal servers: %v", err)
	}

	all, err := convertSearchResultToBareMetalServerList(results)
	if err != nil {
		return err
	}

	// Apply additional client-side filters
	var filtered []model.BareMetalServer

	var re *regexp.Regexp
	if v, ok := d.GetOk("display_name"); ok {
		re, err = regexp.Compile(v.(string))
		if err != nil {
			return fmt.Errorf("invalid regex for display_name: %v", err)
		}
	}

	sourceID := d.Get("source_id").(string)
	osName := strings.ToLower(d.Get("os_name").(string))
	osVersion := strings.ToLower(d.Get("os_version").(string))

	for _, server := range all {
		match := true

		// Display name regex filter
		if re != nil && server.DisplayName != nil {
			if !re.MatchString(*server.DisplayName) {
				match = false
			}
		}

		// Source ID filter
		if sourceID != "" && (server.SourceId == nil || *server.SourceId != sourceID) {
			match = false
		}

		// OS name filter (case-insensitive substring)
		if osName != "" {
			if server.OsInfo == nil || server.OsInfo.OsName == nil {
				match = false
			} else if !strings.Contains(strings.ToLower(*server.OsInfo.OsName), osName) {
				match = false
			}
		}

		// OS version filter (case-insensitive substring)
		if osVersion != "" {
			if server.OsInfo == nil || server.OsInfo.OsVersion == nil {
				match = false
			} else if !strings.Contains(strings.ToLower(*server.OsInfo.OsVersion), osVersion) {
				match = false
			}
		}

		if match {
			filtered = append(filtered, server)
		}
	}

	// Convert to schema format
	var items []map[string]interface{}
	for _, server := range filtered {
		item := map[string]interface{}{
			"external_id":    getStringValue(server.ExternalId),
			"display_name":   getStringValue(server.DisplayName),
			"source_id":      getStringValue(server.SourceId),
			"resource_type":  getStringValue(server.ResourceType),
			"last_sync_time": getInt64Value(server.LastSyncTime),
		}

		if server.CpuCores != nil {
			item["cpu_cores"] = int(*server.CpuCores)
		}

		// OS info
		if server.OsInfo != nil {
			item["os_name"] = getStringValue(server.OsInfo.OsName)
			item["os_version"] = getStringValue(server.OsInfo.OsVersion)
		}

		// Tags
		if server.Tags != nil {
			var tags []map[string]interface{}
			for _, tag := range server.Tags {
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

// convertSearchResultToBareMetalServerList converts search results to BareMetalServer list
func convertSearchResultToBareMetalServerList(searchResults []*data.StructValue) ([]model.BareMetalServer, error) {
	var servers []model.BareMetalServer
	converter := bindings.NewTypeConverter()

	for _, item := range searchResults {
		dv, errs := converter.ConvertToGolang(item, model.BareMetalServerBindingType())
		if len(errs) > 0 {
			return servers, errs[0]
		}
		servers = append(servers, dv.(model.BareMetalServer))
	}
	return servers, nil
}
