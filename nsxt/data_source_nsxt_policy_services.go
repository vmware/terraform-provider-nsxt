// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	orgs "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs"
	projects "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"
)

var cliOrgSharedWithMeClient = func(connector vapiProtocolClient.Connector) orgs.SharedWithMeClient {
	return orgs.NewSharedWithMeClient(connector)
}

var cliProjectSharedWithMeClient = func(connector vapiProtocolClient.Connector) projects.SharedWithMeClient {
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

type defaultServiceInfo struct {
	displayName string
	isDefault   bool
}

func getDefaultServicesInfo(connector vapiProtocolClient.Connector, isGlobal bool) (map[string]defaultServiceInfo, error) {
	defaultSessionContext := utl.SessionContext{
		ProjectID:  "",
		VPCID:      "",
		FromGlobal: false,
	}
	if isGlobal {
		defaultSessionContext.ClientType = utl.Global
	} else {
		defaultSessionContext.ClientType = utl.Local
	}
	defaultClient := cliServicesClient(defaultSessionContext, connector)
	servicesInfo := make(map[string]defaultServiceInfo)
	var cursor *string
	for {
		defaultResults, err := defaultClient.List(cursor, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		if defaultResults.Results != nil {
			for _, r := range defaultResults.Results {
				if r.Path != nil && r.DisplayName != nil {
					isDefault := r.IsDefault != nil && *r.IsDefault
					servicesInfo[*r.Path] = defaultServiceInfo{
						displayName: *r.DisplayName,
						isDefault:   isDefault,
					}
				}
			}
		}
		cursor = defaultResults.Cursor
		if cursor == nil {
			break
		}
	}
	return servicesInfo, nil
}

func dataSourceNsxtPolicyServicesRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	builtInOnly := d.Get("built_in_only").(bool)
	includeSharedServices := d.Get("include_shared_services").(bool)

	sessionContext := getSessionContext(d, m)
	isProjectContext := sessionContext.ProjectID != "" && sessionContext.VPCID == ""

	servicesMap := make(map[string]string)

	if !isProjectContext {
		// Default Space behavior:
		// - We do not support built_in_only and include_shared_services. They are ignored.
		// - We just return all services under the given context.
		client := cliServicesClient(sessionContext, connector)
		var cursor *string
		for {
			results, err := client.List(cursor, nil, nil, nil, nil, nil, nil)
			if err != nil {
				return err
			}
			if results.Results != nil {
				for _, r := range results.Results {
					if r.DisplayName != nil && r.Path != nil {
						servicesMap[*r.DisplayName] = *r.Path
					}
				}
			}
			cursor = results.Cursor
			if cursor == nil {
				break
			}
		}
	} else {
		// User Created Project behavior (strictly Project Context):
		// We always fetch the project's own services (since they are always custom services or project-built-in).
		// - If built_in_only is true, we only keep those with IsDefault == true.
		// - If built_in_only is false, we keep all of them.
		client := cliServicesClient(sessionContext, connector)
		var cursor *string
		for {
			results, err := client.List(cursor, nil, nil, nil, nil, nil, nil)
			if err != nil {
				return err
			}
			if results.Results != nil {
				for _, r := range results.Results {
					if r.DisplayName != nil && r.Path != nil {
						isDefault := r.IsDefault != nil && *r.IsDefault
						if builtInOnly && !isDefault {
							continue
						}
						servicesMap[*r.DisplayName] = *r.Path
					}
				}
			}
			cursor = results.Cursor
			if cursor == nil {
				break
			}
		}

		// Then, if includeSharedServices is true, we also fetch and append the shared services.
		if includeSharedServices {
			defaultServices, err := getDefaultServicesInfo(connector, isPolicyGlobalManager(m))
			if err != nil {
				return err
			}

			resourceType := "Service"

			// Helper to process shared resources and filter
			addSharedResource := func(sr nsxModel.SharedResource) {
				for _, ro := range sr.ResourceObjects {
					if ro.ResourcePath != nil {
						info, exists := defaultServices[*ro.ResourcePath]
						if !exists {
							continue
						}
						if builtInOnly && !info.isDefault {
							continue
						}
						servicesMap[info.displayName] = *ro.ResourcePath
					}
				}
			}

			// 1. Fetch Org-level shared resources
			orgSharedClient := cliOrgSharedWithMeClient(connector)
			orgResult, err := orgSharedClient.List(utl.DefaultOrgID, &resourceType)
			if err != nil {
				if !isNotFoundError(err) {
					return err
				}
			} else if orgResult.Results != nil {
				for _, sr := range orgResult.Results {
					addSharedResource(sr)
				}
			}

			// 2. Fetch Project-level shared resources
			if sessionContext.ProjectID != "" {
				projectSharedClient := cliProjectSharedWithMeClient(connector)
				projectResult, err := projectSharedClient.List(utl.DefaultOrgID, sessionContext.ProjectID, &resourceType)
				if err != nil {
					if !isNotFoundError(err) {
						return err
					}
				} else if projectResult.Results != nil {
					for _, sr := range projectResult.Results {
						addSharedResource(sr)
					}
				}
			}
		}
	}

	d.Set("items", servicesMap)
	d.SetId(newUUID())
	return nil
}
