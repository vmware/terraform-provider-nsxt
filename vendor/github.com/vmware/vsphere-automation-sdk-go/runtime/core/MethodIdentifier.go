/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

type MethodIdentifier struct {

	/**
	 * Returns the interface identifier of the interface this method belongs
	 */
	interfaceIdentifier InterfaceIdentifier
	/**
	 * method name
	 */
	methodName string
}

/**
 * Creates a method identifier given an interface identifier
 * and a method name.
 *
 * @param interfaceIdentifier   interface identifier
 * @param method  method name
 * @throws GenericError if <code>interfaceIdentifier</code> or
 *         <code>method</code> is <code>null</code>
 */
func NewMethodIdentifier(iface InterfaceIdentifier, methodName string) MethodIdentifier {

	var methodIdentifier = MethodIdentifier{interfaceIdentifier: iface, methodName: methodName}
	return methodIdentifier
}

func (m MethodIdentifier) MethodNameDelimiter() string {
	return "."
}

/**
* Returns the name of this method.
*
* @return  method name
 */
func (m MethodIdentifier) Name() string {
	return (m.methodName)
}

/**
* @return the fully-qualified name ([interface].[method]).
 */
func (m MethodIdentifier) FullyQualifiedName() string {
	return (m.interfaceIdentifier.interfaceName) + m.MethodNameDelimiter() + m.Name()
}

/**
* Returns the identifier of the interface containing this method.
*
* @return  interface identifier of the interface containing this method
 */
func (m MethodIdentifier) InterfaceIdentifier() InterfaceIdentifier {
	return m.interfaceIdentifier
}
