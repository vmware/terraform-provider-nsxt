/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbProtocolParsers
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbProtocolParsersClient interface {

    // Delete the ALBProtocolParser along with all the entities contained by this ALBProtocolParser.
    //
    // @param albProtocolparserIdParam ALBProtocolParser ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albProtocolparserIdParam string, forceParam *bool) error

    // Read a ALBProtocolParser.
    //
    // @param albProtocolparserIdParam ALBProtocolParser ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBProtocolParser
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albProtocolparserIdParam string) (model.ALBProtocolParser, error)

    // Paginated list of all ALBProtocolParser for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBProtocolParserApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBProtocolParserApiResponse, error)

    // If a ALBprotocolparser with the alb-protocolparser-id is not already present, create a new ALBprotocolparser. If it already exists, update the ALBprotocolparser. This is a full replace.
    //
    // @param albProtocolparserIdParam ALBprotocolparser ID (required)
    // @param aLBProtocolParserParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albProtocolparserIdParam string, aLBProtocolParserParam model.ALBProtocolParser) error

    // If a ALBProtocolParser with the alb-ProtocolParser-id is not already present, create a new ALBProtocolParser. If it already exists, update the ALBProtocolParser. This is a full replace.
    //
    // @param albProtocolparserIdParam ALBProtocolParser ID (required)
    // @param aLBProtocolParserParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBProtocolParser
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albProtocolparserIdParam string, aLBProtocolParserParam model.ALBProtocolParser) (model.ALBProtocolParser, error)
}
