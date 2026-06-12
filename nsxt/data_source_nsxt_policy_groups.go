// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	projects "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"
)

var cliSharedWithMeClient = func(connector vapiProtocolClient.Connector) projects.SharedWithMeClient {
	return projects.NewSharedWithMeClient(connector)
}

func dataSourceNsxtPolicyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyGroupsRead,

		Schema: map[string]*schema.Schema{
			"context": getContextSchema(false, false, false),
			"domain":  getDomainNameSchema(),
			"include_shared_groups": {
				Type:        schema.TypeBool,
				Description: "Include groups shared with the project",
				Optional:    true,
				Default:     false,
			},
			"items": {
				Type:        schema.TypeMap,
				Description: "Mapping of service policy path by display name",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNsxtPolicyGroupsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	domainName := d.Get("domain").(string)

	client := cliGroupsClient(getSessionContext(d, m), connector)

	boolFalse := false
	var cursor *string
	groupsMap := make(map[string]string)

	for {
		results, err := client.List(domainName, cursor, nil, nil, nil, nil, &boolFalse, nil)
		if err != nil {
			return err
		}
		for _, r := range results.Results {
			groupsMap[*r.DisplayName] = *r.Path
		}
		cursor = results.Cursor
		if cursor == nil {
			break
		}
	}

	includeSharedGroups := d.Get("include_shared_groups").(bool)
	sessionContext := getSessionContext(d, m)
	// If requested and within a project/VPC context, retrieve groups shared with the project
	// and merge them into the groups map.
	if includeSharedGroups && (sessionContext.ClientType == utl.Multitenancy || sessionContext.ClientType == utl.VPC) {
		sharedClient := cliSharedWithMeClient(connector)
		resourceType := "Group"
		sharedResult, err := sharedClient.List(utl.DefaultOrgID, sessionContext.ProjectID, &resourceType)
		if err != nil {
			return err
		}
		for _, sr := range sharedResult.Results {
			for _, ro := range sr.ResourceObjects {
				if ro.ResourcePath != nil && sr.DisplayName != nil {
					groupsMap[*sr.DisplayName] = *ro.ResourcePath
				}
			}
		}
	}

	d.Set("items", groupsMap)
	d.SetId(newUUID())
	return nil
}
