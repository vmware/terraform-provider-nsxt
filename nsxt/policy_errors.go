/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
)

func logVapiErrorData(message string, apiErrorDataValue *data.StructValue) error {
	var typeConverter = bindings.NewTypeConverter()
	typeConverter.SetMode(bindings.REST)
	data, err := typeConverter.ConvertToGolang(apiErrorDataValue, model.ApiErrorBindingType())

	if err != nil {
		log.Printf("[ERROR] Failed to extract error details: %s", err)
		return fmt.Errorf("%s (failed to extract additional details due to %s)", message, err)
	}
	apiError := data.(model.ApiError)
	details := fmt.Sprintf(" %s: %s (code %v)", message, *apiError.ErrorMessage, *apiError.ErrorCode)

	if len(apiError.RelatedErrors) > 0 {
		details += "\nRelated errors:\n"
		for _, relatedErr := range apiError.RelatedErrors {
			details += fmt.Sprintf("%s (code %v)", *relatedErr.ErrorMessage, relatedErr.ErrorCode)
		}
	}
	log.Printf("[ERROR]: %s", details)
	return fmt.Errorf(details)
}

func logAPIError(message string, err error) error {
	if vapiError, ok := err.(errors.InvalidRequest); ok {
		return logVapiErrorData(message, vapiError.Data)
	}
	if vapiError, ok := err.(errors.NotFound); ok {
		return logVapiErrorData(message, vapiError.Data)
	}
	if vapiError, ok := err.(errors.Unauthorized); ok {
		return logVapiErrorData(message, vapiError.Data)
	}
	if vapiError, ok := err.(errors.Unauthenticated); ok {
		return logVapiErrorData(message, vapiError.Data)
	}
	if vapiError, ok := err.(errors.InternalServerError); ok {
		return logVapiErrorData(message, vapiError.Data)
	}
	if vapiError, ok := err.(errors.ServiceUnavailable); ok {
		return logVapiErrorData(message, vapiError.Data)
	}

	return err
}

func isNotFoundError(err error) bool {
	if _, ok := err.(errors.NotFound); ok {
		return true
	}

	return false
}

func handleCreateError(resourceType string, resourceID string, err error) error {
	msg := fmt.Sprintf("Failed to create %s %s", resourceType, resourceID)
	return logAPIError(msg, err)
}

func handleUpdateError(resourceType string, resourceID string, err error) error {
	msg := fmt.Sprintf("Failed to update %s %s", resourceType, resourceID)
	return logAPIError(msg, err)
}

func handleReadError(d *schema.ResourceData, resourceType string, resourceID string, err error) error {
	msg := fmt.Sprintf("Failed to read %s %s", resourceType, resourceID)
	if isNotFoundError(err) {
		d.SetId("")
		log.Printf(msg)
		return nil
	}
	return logAPIError(msg, err)
}

func handleDeleteError(resourceType string, resourceID string, err error) error {
	if isNotFoundError(err) {
		log.Printf("[WARNING] %s %s not found on backend", resourceType, resourceID)
		// We don't want to fail apply on this
		return nil
	}
	msg := fmt.Sprintf("Failed to delete %s %s", resourceType, resourceID)
	return logAPIError(msg, err)
}
