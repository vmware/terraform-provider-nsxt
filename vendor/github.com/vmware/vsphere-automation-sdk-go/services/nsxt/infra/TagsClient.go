/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Tags
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type TagsClient interface {

    // Returns paginated list of all unique tags. Supports filtering by scope, tag and source from which tags are synched. Supports starts with, equals and contains operators on scope and tag values. To filter tags by starts with on scope or tag, use '\*' as prefix before the value. To filter tags by ends with on scope or tag, use '\*' as suffix after the value. To filter tags by contain on scope or tag, use '\*' as prefix and suffix on the value. Below special characters in the filter value needs to be escaped with hex values. - Character '&' needs to be escaped as '%26' - Character '[' needs to be escaped as '%5B' - Character ']' needs to be escaped as '%5D' - Character '+' needs to be escaped as '%2B' - Character '#' needs to be escaped as '%23'
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param scopeParam Tag scope (optional)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @param sourceParam Source from which tags are synced. (optional)
    // @param tagParam Tag value (optional)
    // @return com.vmware.nsx_policy.model.TagInfoListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, scopeParam *string, sortAscendingParam *bool, sortByParam *string, sourceParam *string, tagParam *string) (model.TagInfoListResult, error)
}
