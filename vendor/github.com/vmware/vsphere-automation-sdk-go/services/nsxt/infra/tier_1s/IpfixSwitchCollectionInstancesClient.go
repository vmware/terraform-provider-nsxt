// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: IpfixSwitchCollectionInstances
// Used by client-side stubs.

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = core.SupportedByRuntimeVersion1

type IpfixSwitchCollectionInstancesClient interface {

	// API deletes IPFIX Switch Collection Instance.Flow forwarding to selected collector will be stopped. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-profiles
	//
	// @param tier1IdParam Tier-1 ID (required)
	// @param ipfixSwitchCollectionInstanceIdParam IPFIX Switch Collection Instance ID (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string) error

	// API will return details of IPFIX switch collection. If instance does not exist, it will return 404. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-profiles
	//
	// @param tier1IdParam Tier-1 ID (required)
	// @param ipfixSwitchCollectionInstanceIdParam IPFIX switch collection id (required)
	// @return com.vmware.nsx_policy.model.IPFIXSwitchCollectionInstance
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string) (model.IPFIXSwitchCollectionInstance, error)

	// API provides list IPFIX Switch collection instances available on selected logical switch. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-profiles
	//
	// @param tier1IdParam Tier-1 ID (required)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_policy.model.IPFIXSwitchCollectionInstanceListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(tier1IdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPFIXSwitchCollectionInstanceListResult, error)

	// Create a new IPFIX switch collection instance if the IPFIX switch collection instance with given id does not already exist. If the IPFIX switch collection instance with the given id already exists, patch with the existing IPFIX switch collection instance. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-profiles
	//
	// @param tier1IdParam Tier-1 ID (required)
	// @param ipfixSwitchCollectionInstanceIdParam IPFIX Switch Collection Instance ID (required)
	// @param iPFIXSwitchCollectionInstanceParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string, iPFIXSwitchCollectionInstanceParam model.IPFIXSwitchCollectionInstance) error

	// Create or replace IPFIX switch collection instance. Instance will start forwarding data to provided IPFIX collector. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-profiles
	//
	// @param tier1IdParam Tier-1 ID (required)
	// @param ipfixSwitchCollectionInstanceIdParam IPFIX Switch Collection Instance ID (required)
	// @param iPFIXSwitchCollectionInstanceParam (required)
	// @return com.vmware.nsx_policy.model.IPFIXSwitchCollectionInstance
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string, iPFIXSwitchCollectionInstanceParam model.IPFIXSwitchCollectionInstance) (model.IPFIXSwitchCollectionInstance, error)
}

type ipfixSwitchCollectionInstancesClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewIpfixSwitchCollectionInstancesClient(connector client.Connector) *ipfixSwitchCollectionInstancesClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.tier_1s.ipfix_switch_collection_instances")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	iIface := ipfixSwitchCollectionInstancesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &iIface
}

func (iIface *ipfixSwitchCollectionInstancesClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := iIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (iIface *ipfixSwitchCollectionInstancesClient) Delete(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string) error {
	typeConverter := iIface.connector.TypeConverter()
	executionContext := iIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(ipfixSwitchCollectionInstancesDeleteInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("IpfixSwitchCollectionInstanceId", ipfixSwitchCollectionInstanceIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := ipfixSwitchCollectionInstancesDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	iIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := iIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.ipfix_switch_collection_instances", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (iIface *ipfixSwitchCollectionInstancesClient) Get(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string) (model.IPFIXSwitchCollectionInstance, error) {
	typeConverter := iIface.connector.TypeConverter()
	executionContext := iIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(ipfixSwitchCollectionInstancesGetInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("IpfixSwitchCollectionInstanceId", ipfixSwitchCollectionInstanceIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.IPFIXSwitchCollectionInstance
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := ipfixSwitchCollectionInstancesGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	iIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := iIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.ipfix_switch_collection_instances", "get", inputDataValue, executionContext)
	var emptyOutput model.IPFIXSwitchCollectionInstance
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ipfixSwitchCollectionInstancesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.IPFIXSwitchCollectionInstance), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (iIface *ipfixSwitchCollectionInstancesClient) List(tier1IdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPFIXSwitchCollectionInstanceListResult, error) {
	typeConverter := iIface.connector.TypeConverter()
	executionContext := iIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(ipfixSwitchCollectionInstancesListInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.IPFIXSwitchCollectionInstanceListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := ipfixSwitchCollectionInstancesListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	iIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := iIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.ipfix_switch_collection_instances", "list", inputDataValue, executionContext)
	var emptyOutput model.IPFIXSwitchCollectionInstanceListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ipfixSwitchCollectionInstancesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.IPFIXSwitchCollectionInstanceListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (iIface *ipfixSwitchCollectionInstancesClient) Patch(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string, iPFIXSwitchCollectionInstanceParam model.IPFIXSwitchCollectionInstance) error {
	typeConverter := iIface.connector.TypeConverter()
	executionContext := iIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(ipfixSwitchCollectionInstancesPatchInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("IpfixSwitchCollectionInstanceId", ipfixSwitchCollectionInstanceIdParam)
	sv.AddStructField("IPFIXSwitchCollectionInstance", iPFIXSwitchCollectionInstanceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := ipfixSwitchCollectionInstancesPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	iIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := iIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.ipfix_switch_collection_instances", "patch", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (iIface *ipfixSwitchCollectionInstancesClient) Update(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string, iPFIXSwitchCollectionInstanceParam model.IPFIXSwitchCollectionInstance) (model.IPFIXSwitchCollectionInstance, error) {
	typeConverter := iIface.connector.TypeConverter()
	executionContext := iIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(ipfixSwitchCollectionInstancesUpdateInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("IpfixSwitchCollectionInstanceId", ipfixSwitchCollectionInstanceIdParam)
	sv.AddStructField("IPFIXSwitchCollectionInstance", iPFIXSwitchCollectionInstanceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.IPFIXSwitchCollectionInstance
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := ipfixSwitchCollectionInstancesUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	iIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := iIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.ipfix_switch_collection_instances", "update", inputDataValue, executionContext)
	var emptyOutput model.IPFIXSwitchCollectionInstance
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ipfixSwitchCollectionInstancesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.IPFIXSwitchCollectionInstance), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
