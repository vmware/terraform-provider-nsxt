/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package msg

type JsonRpc20Response struct {
	//A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0"
	version string

	/**
	 *  This member is REQUIRED.
	 *	It MUST be the same as the value of the id member in the Request Object.
	 *	If there was an error in detecting the id in the Request object (e.g. Parse error/Invalid Request), it MUST be Null.
	 */
	id interface{}

	/**
	 *	This member is REQUIRED on success.
	 *	This member MUST NOT exist if there was an error invoking the method.
	 *	The value of this member is determined by the method invoked on the Server.
	 */
	result interface{}

	/**
	 *	This member is REQUIRED on error.
	 *	This member MUST NOT exist if there was no error triggered during invocation.
	 *	The value for this member MUST be an Object as defined in section 5.1.
	 */
	error map[string]interface{}
}

func NewJsonRpc20Response(version string, id interface{}, result interface{}, error map[string]interface{}) JsonRpc20Response {
	return JsonRpc20Response{version: version, id: id, result: result, error: error}
}

func (j JsonRpc20Response) Id() interface{} {
	return j.id
}

func (j JsonRpc20Response) Version() string {
	return j.version
}

func (j JsonRpc20Response) Result() interface{} {
	return j.result
}

func (j JsonRpc20Response) Error() map[string]interface{} {
	return j.error
}
