/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package customtypes

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable = (*NSXLicenseType)(nil)
	_ xattr.TypeWithValidate  = (*NSXLicenseType)(nil)
)

// NSXLicenseType is an attribute type that represents a valid NSX license string. No semantic equality
type NSXLicenseType struct {
	basetypes.StringType
}

// String returns a human readable string of the type name.
func (t NSXLicenseType) String() string {
	return "customtypes.NSXLicenseType"
}

// ValueType returns the Value type.
func (t NSXLicenseType) ValueType(ctx context.Context) attr.Value {
	return NSXLicense{}
}

// Equal returns true if the given type is equivalent.
func (t NSXLicenseType) Equal(o attr.Type) bool {
	other, ok := o.(NSXLicenseType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// Validate implements type validation. This type requires the value provided to be a String value that is a valid NSX license.
// This utilizes the Go `net/netip` library for parsing so leading zeroes will be rejected as invalid.
func (t NSXLicenseType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.String) {
		err := fmt.Errorf("expected String value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"MSX License Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var valueString string

	if err := in.As(&valueString); err != nil {
		diags.AddAttributeError(
			path,
			"NSX License Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)

		return diags
	}

	r := regexp.MustCompile("^[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}$")
	if !r.MatchString(valueString) {
		diags.AddAttributeError(
			path,
			"Invalid NSX License",
			fmt.Sprintf("Value %s must be a valid nsx license key matching: ^[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}$", valueString),
		)

		return diags
	}

	return diags
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t NSXLicenseType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return NSXLicense{
		StringValue: in,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t NSXLicenseType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}
