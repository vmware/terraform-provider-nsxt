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

func dataSourceNsxtPolicyBareMetalServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyBareMetalServerRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"external_id": {
				Type:        schema.TypeString,
				Description: "External ID of the bare metal server (exact match)",
				Optional:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name of the bare metal server (prefix match)",
				Optional:    true,
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
	}
}

func dataSourceNsxtPolicyBareMetalServerRead(d *schema.ResourceData, m interface{}) error {
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

	var server model.BareMetalServer
	var err error

	if externalID != "" {
		// Find by external ID (exact match)
		server, err = findBareMetalServerByExternalID(connector, ctx, externalID)
		if err != nil {
			return fmt.Errorf("Failed to find bare metal server with external_id '%s': %v", externalID, err)
		}
	} else {
		// Find by display name (prefix match returning exactly one result)
		results, err := listInventoryResourcesByType(connector, ctx, "BareMetalServer", nil)
		if err != nil {
			return fmt.Errorf("error listing bare metal servers: %v", err)
		}

		servers, err := convertSearchResultToBareMetalServerList(results)
		if err != nil {
			return err
		}

		var matches []model.BareMetalServer
		for _, srv := range servers {
			if srv.DisplayName != nil && *srv.DisplayName == displayName {
				matches = append(matches, srv)
			}
		}

		if len(matches) == 0 {
			return fmt.Errorf("No bare metal server found with display_name '%s'", displayName)
		}
		if len(matches) > 1 {
			return fmt.Errorf("Multiple bare metal servers found with display_name '%s'. Use external_id for exact match", displayName)
		}
		server = matches[0]
	}

	// Set attributes
	d.SetId(*server.ExternalId)
	d.Set("external_id", server.ExternalId)
	d.Set("display_name", server.DisplayName)
	d.Set("source_id", server.SourceId)
	d.Set("resource_type", server.ResourceType)
	d.Set("last_sync_time", server.LastSyncTime)

	if server.CpuCores != nil {
		d.Set("cpu_cores", int(*server.CpuCores))
	}

	// OS info
	if server.OsInfo != nil {
		d.Set("os_name", server.OsInfo.OsName)
		d.Set("os_version", server.OsInfo.OsVersion)
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
		d.Set("tags", tags)
	}

	return nil
}

// findBareMetalServerByExternalID finds a BMS by external ID using search
func findBareMetalServerByExternalID(connector client.Connector, ctx utl.SessionContext, externalID string) (model.BareMetalServer, error) {
	var found model.BareMetalServer

	resultValues, err := listInventoryResourcesByAnyFieldAndType(connector, ctx, externalID, "BareMetalServer", nil)
	if err != nil {
		return found, err
	}

	converter := bindings.NewTypeConverter()
	for _, item := range resultValues {
		dv, errs := converter.ConvertToGolang(item, model.BareMetalServerBindingType())
		if len(errs) > 0 {
			return found, errs[0]
		}
		bms := dv.(model.BareMetalServer)
		if bms.ExternalId != nil && *bms.ExternalId == externalID {
			return bms, nil
		}
	}
	return found, fmt.Errorf("could not find bare metal server with external_id: %s", externalID)
}
