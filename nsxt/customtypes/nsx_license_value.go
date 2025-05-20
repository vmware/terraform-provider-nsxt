/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package customtypes

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable = (*NSXLicense)(nil)
)

// NSXLicense represents a valid NSX license string.
type NSXLicense struct {
	basetypes.StringValue
}

// Type returns an NSXLicenseType.
func (v NSXLicense) Type(_ context.Context) attr.Type {
	return NSXLicenseType{}
}

// Equal returns true if the given value is equivalent.
func (v NSXLicense) Equal(o attr.Value) bool {
	other, ok := o.(NSXLicense)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// NewNSXLicenseNull creates an NSXLicense with a null value. Determine whether the value is null via IsNull method.
func NewNSXLicenseNull() NSXLicense {
	return NSXLicense{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewNSXLicenseUnknown creates an NSXLicense with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewNSXLicenseUnknown() NSXLicense {
	return NSXLicense{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewNSXLicenseValue creates an NSXLicense with a known value. Access the value via ValueString method.
func NewNSXLicenseValue(value string) NSXLicense {
	return NSXLicense{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewNSXLicensePointerValue creates an NSXLicense with a null value if nil or a known value. Access the value via ValueStringPointer method.
func NewNSXLicensePointerValue(value *string) NSXLicense {
	return NSXLicense{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
