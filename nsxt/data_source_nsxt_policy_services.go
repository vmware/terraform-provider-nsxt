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

func dataSourceNsxtPolicyServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyServicesRead,

		Schema: map[string]*schema.Schema{
			"context": getContextSchema(false, false, false),
			"include_shared_services": {
				Type:        schema.TypeBool,
				Description: "Include services shared with the project",
				Optional:    true,
				Default:     false,
			},
			"built_in_only": {
				Type:        schema.TypeBool,
				Description: "Only include default built-in services",
				Optional:    true,
				Default:     false,
			},
			"items": {
				Type:        schema.TypeMap,
				Description: "Mapping of services policy path by display name",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNsxtPolicyServicesRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	builtInOnly := d.Get("built_in_only").(bool)

	servicesMap := make(map[string]string)

	if !builtInOnly {
		client := cliServicesClient(getSessionContext(d, m), connector)
		results, err := client.List(nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return err
		}
		for _, r := range results.Results {
			servicesMap[*r.DisplayName] = *r.Path
		}
	}

	includeSharedServices := d.Get("include_shared_services").(bool)
	sessionContext := getSessionContext(d, m)

	if builtInOnly || (includeSharedServices && (sessionContext.ClientType == utl.Multitenancy || sessionContext.ClientType == utl.VPC)) {
		defaultSessionContext := utl.SessionContext{
			ProjectID:  "",
			VPCID:      "",
			FromGlobal: false,
		}
		if isPolicyGlobalManager(m) {
			defaultSessionContext.ClientType = utl.Global
		} else {
			defaultSessionContext.ClientType = utl.Local
		}
		defaultClient := cliServicesClient(defaultSessionContext, connector)
		defaultResults, err := defaultClient.List(nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return err
		}
		for _, r := range defaultResults.Results {
			servicesMap[*r.DisplayName] = *r.Path
		}
	}

	if includeSharedServices && (sessionContext.ClientType == utl.Multitenancy || sessionContext.ClientType == utl.VPC) {
		sharedClient := cliSharedWithMeClient(connector)
		resourceType := "Service"
		sharedResult, err := sharedClient.List(utl.DefaultOrgID, sessionContext.ProjectID, &resourceType)
		if err != nil {
			return err
		}
		for _, sr := range sharedResult.Results {
			for _, ro := range sr.ResourceObjects {
				if ro.ResourcePath != nil && sr.DisplayName != nil {
					servicesMap[*sr.DisplayName] = *ro.ResourcePath
				}
			}
		}
	}

	d.Set("items", servicesMap)
	d.SetId(newUUID())
	return nil
}
