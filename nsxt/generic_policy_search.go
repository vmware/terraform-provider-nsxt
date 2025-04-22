// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type GenericResourceData struct {
	DisplayName string
	Id          string
}

func (g GenericResourceData) policyGenericDataSourceResourceFilterAndSet(resultValues []*data.StructValue, resourceType string) ([]*data.StructValue, error) {
	var perfectMatch, prefixMatch []policySearchDataValue
	var obj []policySearchDataValue
	objName := g.DisplayName
	objID := g.Id
	converter := bindings.NewTypeConverter()

	for _, result := range resultValues {
		dataValue, errors := converter.ConvertToGolang(result, model.PolicyResourceBindingType())
		if len(errors) > 0 {
			return nil, errors[0]
		}
		policyResource := dataValue.(model.PolicyResource)
		if resourceType != *policyResource.ResourceType {
			continue
		}

		if objID != "" {
			perfectMatch = append(perfectMatch, policySearchDataValue{StructValue: result, Resource: policyResource})
			break
		} else {
			if *policyResource.DisplayName == objName {
				perfectMatch = append(perfectMatch, policySearchDataValue{StructValue: result, Resource: policyResource})
			}
			if strings.HasPrefix(*policyResource.DisplayName, objName) {
				prefixMatch = append(prefixMatch, policySearchDataValue{StructValue: result, Resource: policyResource})
			}
		}
	}

	if len(perfectMatch) > 0 {
		// if len(perfectMatch) > 1 {
		// if objID != "" {
		// 	return nil, fmt.Errorf("Found multiple %s with ID '%s'", resourceType, objID)
		// }
		// 	return nil, fmt.Errorf("Found multiple %s with name '%s'", resourceType, objName)
		// }
		obj = perfectMatch
	} else if len(prefixMatch) > 0 {
		// if len(prefixMatch) > 1 {
		// 	return nil, fmt.Errorf("Found multiple %s with name starting with '%s'", resourceType, objName)
		// }
		obj = prefixMatch
	} else {
		if objID != "" {
			return nil, fmt.Errorf("%s with ID '%s' was not found", resourceType, objID)
		}
		return nil, fmt.Errorf("%s with name '%s' was not found", resourceType, objName)
	}
	objOp := []*data.StructValue{}
	for _, i := range obj {
		objOp = append(objOp, i.StructValue)
	}
	return objOp, nil
}

func (g GenericResourceData) policyGenericDataSourceResourceRead(connector client.Connector, context utl.SessionContext, resourceType string, additionalQuery map[string]string) ([]*data.StructValue, error) {
	return g.policyGenericDataSourceResourceReadWithValidation(connector, context, resourceType, additionalQuery, true)
}

func (g GenericResourceData) policyGenericDataSourceResourceReadWithValidation(connector client.Connector, context utl.SessionContext, resourceType string, additionalQuery map[string]string, paramsValidation bool) ([]*data.StructValue, error) {
	objName := g.DisplayName
	objID := g.Id
	var err error
	var resultValues []*data.StructValue
	additionalQueryString := buildQueryStringFromMap(additionalQuery)
	// if paramsValidation && objID == "" && objName == "" {
	// 	return nil, fmt.Errorf("No 'id' or 'display_name' specified for %s", resourceType)
	// }
	if objID != "" {
		if resourceType == "PolicyEdgeNode" {
			// Edge Node is a special case where id != nsx_id
			// TODO: consider switching all searches to nsx id
			resultValues, err = listPolicyResourcesByNsxID(connector, context, &objID, &additionalQueryString)
		} else {
			resultValues, err = listPolicyResourcesByID(connector, context, &objID, &additionalQueryString)
		}
	} else {
		resultValues, err = listPolicyResourcesByNameAndType(connector, context, objName, resourceType, &additionalQueryString)
	}
	if err != nil {
		return nil, err
	}

	return g.policyGenericDataSourceResourceFilterAndSet(resultValues, resourceType)
}
