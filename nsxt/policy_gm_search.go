/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/search"
)

type policySearchDataValue struct {
	StructValue *data.StructValue
	Resource    model.PolicyResource
}

func policyDataSourceResourceFilterAndSet(d *schema.ResourceData, resultValues []*data.StructValue, resourceType string) (*data.StructValue, error) {
	var perfectMatch, prefixMatch []policySearchDataValue
	var obj policySearchDataValue
	objName := d.Get("display_name").(string)
	objID := d.Get("id").(string)
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	for _, result := range resultValues {
		dataValue, errors := converter.ConvertToGolang(result, model.PolicyResourceBindingType())
		if len(errors) > 0 {
			return nil, errors[0]
		}
		policyResource := dataValue.(model.PolicyResource)

		if objID != "" && resourceType == *policyResource.ResourceType {
			perfectMatch = append(perfectMatch, policySearchDataValue{StructValue: result, Resource: policyResource})
			continue
		}
		if *policyResource.DisplayName == objName {
			perfectMatch = append(perfectMatch, policySearchDataValue{StructValue: result, Resource: policyResource})
		}
		if strings.HasPrefix(*policyResource.DisplayName, objName) {
			prefixMatch = append(prefixMatch, policySearchDataValue{StructValue: result, Resource: policyResource})
		}
	}

	if len(perfectMatch) > 0 {
		if len(perfectMatch) > 1 {
			if objID != "" {
				return nil, fmt.Errorf("Found multiple %s with ID '%s'", resourceType, objID)
			}
			return nil, fmt.Errorf("Found multiple %s with name '%s'", resourceType, objName)
		}
		obj = perfectMatch[0]
	} else if len(prefixMatch) > 0 {
		if len(prefixMatch) > 1 {
			return nil, fmt.Errorf("Found multiple %s with name starting with '%s'", resourceType, objName)
		}
		obj = prefixMatch[0]
	} else {
		if objID != "" {
			return nil, fmt.Errorf("%s with ID '%s' was not found", resourceType, objID)
		}
		return nil, fmt.Errorf("%s with name '%s' was not found", resourceType, objName)
	}

	d.SetId(*obj.Resource.Id)
	d.Set("display_name", obj.Resource.DisplayName)
	d.Set("description", obj.Resource.Description)
	d.Set("path", obj.Resource.Path)

	return obj.StructValue, nil
}

func policyDataSourceResourceRead(d *schema.ResourceData, connector *client.RestConnector, resourceType string, additionalQuery map[string]string) (*data.StructValue, error) {
	return policyDataSourceResourceReadWithValidation(d, connector, resourceType, additionalQuery, true)
}

func policyDataSourceResourceReadWithValidation(d *schema.ResourceData, connector *client.RestConnector, resourceType string, additionalQuery map[string]string, paramsValidation bool) (*data.StructValue, error) {
	objName := d.Get("display_name").(string)
	objID := d.Get("id").(string)
	var err error
	var resultValues []*data.StructValue
	additionalQueryString := buildQueryStringFromMap(additionalQuery)
	if paramsValidation && objID == "" && objName == "" {
		return nil, fmt.Errorf("No 'id' or 'display_name' specified for %s", resourceType)
	}
	if objID != "" {
		resultValues, err = listPolicyResourcesByID(connector, &objID, &additionalQueryString)
	} else {
		resultValues, err = listPolicyResourcesByType(connector, &resourceType, &additionalQueryString)
	}
	if err != nil {
		return nil, err
	}

	return policyDataSourceResourceFilterAndSet(d, resultValues, resourceType)
}

func listPolicyResourcesByType(connector *client.RestConnector, resourceType *string, additionalQuery *string) ([]*data.StructValue, error) {
	query := fmt.Sprintf("resource_type:%s", *resourceType)
	return searchPolicyResources(connector, *buildPolicyResourcesQuery(&query, additionalQuery))
}

func listPolicyResourcesByID(connector *client.RestConnector, resourceID *string, additionalQuery *string) ([]*data.StructValue, error) {
	query := fmt.Sprintf("id:%s", *resourceID)
	return searchPolicyResources(connector, *buildPolicyResourcesQuery(&query, additionalQuery))
}

func searchPolicyResources(connector *client.RestConnector, query string) ([]*data.StructValue, error) {
	client := search.NewDefaultQueryClient(connector)
	var results []*data.StructValue
	var cursor *string
	total := 0

	for {
		searchResponse, err := client.List(query, cursor, nil, nil, nil, nil)
		if err != nil {
			return results, err
		}
		results = append(results, searchResponse.Results...)
		if total == 0 {
			// first response
			total = int(*searchResponse.ResultCount)
		}
		cursor = searchResponse.Cursor
		if len(results) >= total {
			return results, nil
		}
	}
}

func buildPolicyResourcesQuery(query *string, additionalQuery *string) *string {
	if additionalQuery != nil && *additionalQuery != "" {
		*query = *query + " AND " + *additionalQuery
	}
	return query
}
