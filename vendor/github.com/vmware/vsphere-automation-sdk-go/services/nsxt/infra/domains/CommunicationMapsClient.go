/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: CommunicationMaps
 * Used by client-side stubs.
 */

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type CommunicationMapsClient interface {

    // Deletes the communication map along with all the communication entries This API is deprecated. Please use the following API instead. DELETE /infra/domains/domain-id/security-policies/security-policy-id
    //
    // @param domainIdParam (required)
    // @param communicationMapIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(domainIdParam string, communicationMapIdParam string) error

    // Read communication-map for a domain. This API is deprecated. Please use the following API instead. GET /infra/domains/domain-id/security-policies/security-policy-id
    //
    // @param domainIdParam (required)
    // @param communicationMapIdParam (required)
    // @return com.vmware.nsx_policy.model.CommunicationMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(domainIdParam string, communicationMapIdParam string) (model.CommunicationMap, error)

    // List all communication maps for a domain. This API is deprecated. Please use the following API instead. GET /infra/domains/domain-id/security-policies
    //
    // @param domainIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.CommunicationMapListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.CommunicationMapListResult, error)

    // Patch the communication map for a domain. If a communication map for the given communication-map-id is not present, the object will get created and if it is present it will be updated. This is a full replace This API is deprecated. Please use the following API instead. PATCH /infra/domains/domain-id/security-policies/security-policy-id
    //
    // @param domainIdParam (required)
    // @param communicationMapIdParam (required)
    // @param communicationMapParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(domainIdParam string, communicationMapIdParam string, communicationMapParam model.CommunicationMap) error

    // This is used to set a precedence of a communication map w.r.t others. This API is deprecated. Please use the following API instead. POST /infra/domains/domain-id/security-policies/security-policy-id?action=revise
    //
    // @param domainIdParam (required)
    // @param communicationMapIdParam (required)
    // @param communicationMapParam (required)
    // @param anchorPathParam The communication map/communication entry path if operation is 'insert_after' or 'insert_before' (optional)
    // @param operationParam Operation (optional, default to insert_top)
    // @return com.vmware.nsx_policy.model.CommunicationMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Revise(domainIdParam string, communicationMapIdParam string, communicationMapParam model.CommunicationMap, anchorPathParam *string, operationParam *string) (model.CommunicationMap, error)

    // Create or Update the communication map for a domain. This is a full replace. All the CommunicationEntries are replaced. This API is deprecated. Please use the following API instead. PUT /infra/domains/domain-id/security-policies/security-policy-id
    //
    // @param domainIdParam (required)
    // @param communicationMapIdParam (required)
    // @param communicationMapParam (required)
    // @return com.vmware.nsx_policy.model.CommunicationMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(domainIdParam string, communicationMapIdParam string, communicationMapParam model.CommunicationMap) (model.CommunicationMap, error)
}
