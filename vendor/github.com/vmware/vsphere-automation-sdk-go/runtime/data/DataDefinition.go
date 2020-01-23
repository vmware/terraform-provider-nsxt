/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

type DataDefinition interface {
	/**
	* Returns the {@link DataType} for this type.
	*
	* @return  the {@link DataType} for this type
	 */
	Type() DataType

	/**
	* Validates that the specified {@link DataValue} is an instance of this
	* data-definition.
	*
	* @param value  the datavalue to validate
	* @return       true if the value is an instance of this data definition,
	*               false otherwise
	 */
	ValidInstanceOf(value DataValue) bool

	/**
	 * Validates that the specified {@link DataValue} is an instance of this
	 * data definition.
	 *
	 * <p>Validates that supplied <code>value</code> is not <code>nil</code>
	 * and it's type matches the type of this definition.
	 *
	 * @param value  the <code>DataValue</code> to validate
	 * @return       a list of messages describing validation problems; empty
	 *               list indicates that validation was successful
	 */
	Validate(value DataValue) []error

	/**
	* Check the value to see if the data-definition can fill in any missing
	* data before validation
	*
	* @param value  the value to check
	 */
	CompleteValue(value DataValue)

	String() string
}
