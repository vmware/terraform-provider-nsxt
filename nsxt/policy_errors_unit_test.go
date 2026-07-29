//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiStd "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func apiErrorToStructValue(t *testing.T, apiError model.ApiError) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(apiError, model.ApiErrorBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestUnitNsxt_printAPIError(t *testing.T) {
	msg := "something failed"
	details := "extra details"
	code := int64(42)

	assert.Equal(t, "", printAPIError(model.ApiError{}))
	assert.Equal(t, msg, printAPIError(model.ApiError{ErrorMessage: &msg}))
	assert.Equal(t, msg+": "+details, printAPIError(model.ApiError{ErrorMessage: &msg, Details: &details}))
	assert.Equal(t, msg+": "+details+" (code 42)", printAPIError(model.ApiError{ErrorMessage: &msg, Details: &details, ErrorCode: &code}))
}

func TestUnitNsxt_printRelatedAPIError(t *testing.T) {
	msg := "related failure"
	code := int64(7)

	assert.Equal(t, "", printRelatedAPIError(model.RelatedApiError{}))
	assert.Equal(t, msg, printRelatedAPIError(model.RelatedApiError{ErrorMessage: &msg}))
	assert.Equal(t, "(code 7)", printRelatedAPIError(model.RelatedApiError{ErrorCode: &code}))
	assert.Equal(t, msg+" (code 7)", printRelatedAPIError(model.RelatedApiError{ErrorMessage: &msg, ErrorCode: &code}))
}

func TestUnitNsxt_logRawVapiErrorData(t *testing.T) {
	t.Run("encodable data produces JSON-encoded message", func(t *testing.T) {
		msg := "underlying failure"
		sv := apiErrorToStructValue(t, model.ApiError{ErrorMessage: &msg})
		err := logRawVapiErrorData("op failed", nil, sv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "op failed")
		assert.Contains(t, err.Error(), msg)
	})
}

func TestUnitNsxt_logVapiErrorData(t *testing.T) {
	t.Run("nil data with messages uses first message", func(t *testing.T) {
		messages := []vapiStd.LocalizableMessage{{DefaultMessage: "root cause"}}
		err := logVapiErrorData("failed op", messages, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "root cause")
	})

	t.Run("nil data with type but no messages", func(t *testing.T) {
		vt := vapiErrors.ErrorType_ERROR
		err := logVapiErrorData("failed op", nil, &vt, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), string(vt))
	})

	t.Run("nil data no messages no type", func(t *testing.T) {
		err := logVapiErrorData("failed op", nil, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no additional details provided")
	})

	t.Run("valid ApiError data with related errors", func(t *testing.T) {
		msg := "top level failure"
		relatedMsg := "related failure"
		sv := apiErrorToStructValue(t, model.ApiError{
			ErrorMessage:  &msg,
			RelatedErrors: []model.RelatedApiError{{ErrorMessage: &relatedMsg}},
		})
		err := logVapiErrorData("op failed", nil, nil, sv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), msg)
		assert.Contains(t, err.Error(), "Related errors")
		assert.Contains(t, err.Error(), relatedMsg)
	})

	t.Run("data of an unrelated type falls back to raw output", func(t *testing.T) {
		// A StructValue built from a type other than ApiError either fails
		// conversion or converts into an empty ApiError; either way an error
		// describing the failure must still be returned.
		converter := bindings.NewTypeConverter()
		val, errs := converter.ConvertToVapi(model.Tag{Tag: str("t"), Scope: str("s")}, model.TagBindingType())
		require.Empty(t, errs)
		err := logVapiErrorData("op failed", nil, nil, val.(*data.StructValue))
		require.Error(t, err)
	})
}

func TestUnitNsxt_getInvalidRequestErrorCode(t *testing.T) {
	t.Run("non-InvalidRequest error returns NotFound", func(t *testing.T) {
		_, err := getInvalidRequestErrorCode(vapiErrors.NotFound{})
		require.Error(t, err)
	})

	t.Run("InvalidRequest with nil Data returns NotFound", func(t *testing.T) {
		_, err := getInvalidRequestErrorCode(vapiErrors.InvalidRequest{})
		require.Error(t, err)
		_, isNotFound := err.(vapiErrors.NotFound)
		assert.True(t, isNotFound)
	})

	t.Run("InvalidRequest with valid ApiError data returns code", func(t *testing.T) {
		code := int64(505)
		sv := apiErrorToStructValue(t, model.ApiError{ErrorCode: &code})
		got, err := getInvalidRequestErrorCode(vapiErrors.InvalidRequest{Data: sv})
		require.NoError(t, err)
		assert.Equal(t, code, got)
	})
}

func TestUnitNsxt_logAPIError(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"InvalidRequest", vapiErrors.InvalidRequest{}},
		{"NotFound", vapiErrors.NotFound{}},
		{"Unauthorized", vapiErrors.Unauthorized{}},
		{"Unauthenticated", vapiErrors.Unauthenticated{}},
		{"InternalServerError", vapiErrors.InternalServerError{}},
		{"ServiceUnavailable", vapiErrors.ServiceUnavailable{}},
		{"ConcurrentChange", vapiErrors.ConcurrentChange{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := logAPIError("op", c.err)
			require.Error(t, err)
		})
	}

	t.Run("unrecognized error type is returned as-is", func(t *testing.T) {
		orig := vapiErrors.TimedOut{}
		err := logAPIError("op", orig)
		require.Error(t, err)
		assert.Equal(t, orig, err)
	})
}

func TestUnitNsxt_isErrorTypeCheckers(t *testing.T) {
	assert.True(t, isUnauthorizedError(vapiErrors.Unauthorized{}))
	assert.False(t, isUnauthorizedError(vapiErrors.NotFound{}))

	assert.True(t, isNotFoundError(vapiErrors.NotFound{}))
	assert.False(t, isNotFoundError(vapiErrors.Unauthorized{}))

	assert.True(t, isServiceUnavailableError(vapiErrors.ServiceUnavailable{}))
	assert.False(t, isServiceUnavailableError(vapiErrors.NotFound{}))

	assert.True(t, isTimeoutError(vapiErrors.TimedOut{}))
	assert.False(t, isTimeoutError(vapiErrors.NotFound{}))

	assert.True(t, isInternalServerError(vapiErrors.InternalServerError{}))
	assert.False(t, isInternalServerError(vapiErrors.NotFound{}))
}

func TestUnitNsxt_handleCreateError(t *testing.T) {
	err := handleCreateError("Segment", "seg-1", vapiErrors.InternalServerError{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to create Segment seg-1")
}

func TestUnitNsxt_handleUpdateError(t *testing.T) {
	err := handleUpdateError("Segment", "seg-1", vapiErrors.InternalServerError{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to update Segment seg-1")
}

func TestUnitNsxt_handleListError(t *testing.T) {
	err := handleListError("Segment", vapiErrors.InternalServerError{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to read Segments")
}

func TestUnitNsxt_handleReadError(t *testing.T) {
	t.Run("not found clears ID and returns nil", func(t *testing.T) {
		res := &schema.Resource{Schema: map[string]*schema.Schema{}}
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("seg-1")

		err := handleReadError(d, "Segment", "seg-1", vapiErrors.NotFound{})
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("other errors are propagated", func(t *testing.T) {
		res := &schema.Resource{Schema: map[string]*schema.Schema{}}
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("seg-1")

		err := handleReadError(d, "Segment", "seg-1", vapiErrors.InternalServerError{})
		require.Error(t, err)
		assert.Equal(t, "seg-1", d.Id())
	})
}

func TestUnitNsxt_handleDataSourceReadError(t *testing.T) {
	err := handleDataSourceReadError(nil, "Segment", "seg-1", vapiErrors.InternalServerError{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to read Segment seg-1")
}

func TestUnitNsxt_handleDeleteError(t *testing.T) {
	t.Run("not found is swallowed", func(t *testing.T) {
		err := handleDeleteError("Segment", "seg-1", vapiErrors.NotFound{})
		require.NoError(t, err)
	})

	t.Run("other errors are propagated", func(t *testing.T) {
		err := handleDeleteError("Segment", "seg-1", vapiErrors.InternalServerError{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to delete Segment seg-1")
	})
}

func TestUnitNsxt_handleMultitenancyTier0Error(t *testing.T) {
	err := handleMultitenancyTier0Error()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context use not supported with Tier0 gateways")
}
