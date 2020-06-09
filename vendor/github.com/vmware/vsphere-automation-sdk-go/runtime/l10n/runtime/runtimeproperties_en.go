/* Copyright © 2019-2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package runtime

var RuntimeProperties_EN = []byte(
	`vapi.connection=Could not connect to '{host}'

vapi.bindings.typeconverter.unexpected.runtime.value=Expected a value of type '{expectedType}', but received a value of type '{actualType}'
vapi.bindings.typeconverter.invalid.type=Unexpected binding type '{bindingType}'
vapi.bindings.typeconverter.nil.type=Binding type cannot be nil
vapi.bindings.typeconverter.uri.invalid=Error parsing URI value '{uriValue}': '{errorMessage}'
vapi.bindings.typeconverter.struct.field.invalid=Invalid field '{fieldName}' in structure '{structName}'
vapi.bindings.typeconverter.struct.field.missing=Field '{fieldName}' missing from structure '{structName}'
vapi.bindings.typeconverter.list.entry.invalid=The value at index '{index}' of the list is invalid
vapi.bindings.typeconverter.dict.key.invalid=Invalid value for key
vapi.bindings.typeconverter.dict.value.invalid=Invalid value for key '{key}'
vapi.bindings.typeconverter.dict.missing.key=Key %s missing from dict
vapi.bindings.typeconverter.datetime.invalid=Error parsing datetime object '{dateTime}' to '{vapiFormat}' format: '{errorMessage}'
vapi.bindings.typeconverter.value.nil=Value of element cannot be nil
vapi.bindings.typeconverter.set.invalid=Invalid element in set

vapi.bindings.input.param.invalid=Invalid value for parameter '{paramName}'

vapi.bindings.stub.rest_metadata.unavailable=REST metadata not available for invocation
vapi.bindings.stub.rest_metadata.type.mismatch=Connection metadata type invalid

vapi.data.opaque.definition.null.value=Expected non-nil value got nil value
vapi.data.validate.mismatch=Type mismatch - expected an object of type '{expectedType}', but got '{actualType}'
vapi.data.list.invalid.entry=The value '{value}' at index '{index}' in the list is invalid
vapi.data.structref.not.resolved=Structure reference of type '{referenceType}' is not resolved
vapi.data.structref.structure.not.defined=Cannot resolve structure reference, because structure '{referenceType}' is not defined

vapi.data.structure.field.missing=Field '{fieldName}' missing from structure '{structName}'
vapi.data.structure.field.invalid=Invalid field '{fieldName}' in structure '{structName}'
vapi.data.structure.name.mismatch=Name mismatch for structure. Expected '{expectedName}', but got '{actualName}'
vapi.data.structure.union.extra=Structure '{structName}' has a union with a field '{fieldName}' not required for this case
vapi.data.structure.union.missing=Structure '{structName}' has a union that is missing a required field '{fieldName}' for this case
vapi.data.structure.dynamic.invalid=Expected valid fields of structure '{structName}'
vapi.data.structure.isoneof.value.invalid=Invalid value '{value}' for field '{fieldName}' marked with IsOneOf

vapi.data.serializers.rest.datavalue.error=Error serializing DataValue of type '{type}'
vapi.data.serializers.rest.unhandled.datavalue.formurl=Error in encoding Datavalue to FormUrl type, Unhandled DataValue type detected '{datavalue}'
vapi.data.serializers.rest.formurl.field.error=FormUrl encoding failed for Field '{field}', Error: '{msg}'
vapi.data.serializers.rest.metadata.value.nil=Rest metadata cannot be nil
vapi.data.serializers.rest.nested.invalid.args=Parameter {param} doesnot have valid binding to be returned.
vapi.data.serializers.rest.unsupported_data_value=Data value of type '{type}' is not supported to be serialized as part of REST request URL
vapi.data.serializers.unsupported.json.type=Could not determine appropriate DataValue for json type '{jsonType}'
vapi.data.invalid.json.number=json.Number neither Int64 nor Float64
vapi.data.serializers.json.marshall.error=Error converting the object to byte stream: '{errorMessage}'

vapi.introspection.operation.service.not_found=Service '{serviceName}' not found
vapi.introspection.service.not_found=Service '{serviceName}' not found
vapi.introspection.operation.not_found=Operation '{operationName}' is not found in service '{serviceName}'

vapi.metadata.parser.failure = metamodel or privilege Parser failed to ingest the given json files : {msg}

vapi.method.input.invalid.interface=Invalid interface '{serviceId}'
vapi.method.input.invalid.method=Invalid method '{methodId}'
vapi.method.input.invalid=Invalid input to method '{methodId}'
vapi.method.input.invalid.definition=Invalid input definition to method '{methodId}'
vapi.method.status.errors.invalid=Invalid error '{errorName}' reported from method '{methodName}'

vapi.security.rest.context.error=Security context creation failed, Failure: {err}
vapi.security.authentication.scheme.invalid=Authentication scheme '{scheme}' is invalid
vapi.security.authentication.invalid=Unable to authenticate user
vapi.security.authentication.failed=Invalid authentication result
vapi.security.authentication.exception=Exception in invoking authentication handler %s
vapi.security.authentication.metadata.invalid=Cannot parse authentication metadata for operation %s of service %s
vapi.security.authentication.scheme=Expected one of the following schemes {allowedSchemes}, but got {providedScheme}
vapi.security.authentication.certificate.invalid=Unable to verify server certificate

vapi.security.authorization.exception=Exception in invoking authorization handler {msg}
vapi.security.authorization.invalid=Unable to authorize user
vapi.security.authorization.invalid_with_error=Unable to authorize user. Error occured: {err}
vapi.security.authorization.internal_server_error=Internal server error occured on authorization: {err}

vapi.security.sso.digest.invalid=Invalid digest.
vapi.security.sso.hash.invalid=Invalid hash algorithm.
vapi.security.sso.pubkey.invalid=Invalid public key.
vapi.security.sso.pvtkey.invalid=Invalid private key.
vapi.security.sso.samltoken.invalid=Invalid saml token.
vapi.security.sso.signature.invalid=Invalid signature.
vapi.security.sso.signature.algorithm.invalid=Invalid signature algorithm.

vapi.protocol.server.error=Server error occured: {err}

vapi.protocol.server.rest.param.invalid_value=Parsing failed because of '{errMsg}'
vapi.protocol.server.rest.param.invalid_type=Invalid value for request parameter '{paramName}'. Expected a value of type '{expectedType}', but parsing failed because of '{errMsg}'
com.vmware.vapi.rest.unsupported_property=Unsupported property with name: {arg}.
com.vmware.vapi.rest.unsupported_media_type=Unsupported media type.
vapi.protocol.server.rest.param.name_type_map_not_found='{paramName}' is not found in parameter name to type map
vapi.protocol.server.rest.param.unsupport_type=Request parameter data type '{dataType}' not supported
vapi.protocol.server.rest.param.body_parse_error=Error when parsing request body, '{errMsg}'
vapi.protocol.server.rest.param.bodyfield_parse_invalid=Body Field data is not of type Map
vapi.protocol.server.rest.param.bodyfield.unexpected_field=Body Field contains unexpected field : '{msg}'
vapi.protocol.server.rest.param.bodyfield.missing_field=Body Field does not contain required field : '{msg}'
vapi.protocol.server.rest.param.body.unsupport_type=Request body data type '{dataType}' not supported
vapi.protocol.server.rest.param.body.unexpected=Request body data type expected '{required}' but got '{type}'
vapi.protocol.server.rest.param.internal_server_error= Request process params failure: {msg}
vapi.protocol.server.rest.response.unsupport_type=Response data type '{dataType}' not supported for field '{fieldName}'
vapi.protocol.server.rest.response.not_structure=Response result is not a structure type
vapi.protocol.server.rest.response.not_error=Response result is not an error type
vapi.protocol.server.rest.response.invalid=Response result is nil
vapi.protocol.server.rest.response.error_not_structure=Response error is not a structure type
vapi.protocol.server.rest.response.result_failed=Method execution failed, do not set response header
vapi.protocol.server.rest.response.unsupport_http_status=Http status '{httpStatus}' not supported
vapi.protocol.server.rest.response.body_parse_error=Error when parsing response body
vapi.protocol.server.rest.error.not_supported=Http status '{errorName}' not supported
vapi.protocol.client.request.not_structure=Input request is not a structure type
vapi.protocol.client.response.error=Error reading server response
vapi.protocol.client.request.error=Error completing client request '{errMsg}'
vapi.protocol.client.response.unmarshall.error=Error unmarshalling server response

vapi.server.timedout = Request Timed out
vapi.server.unavailable = Service not available
vapi.server.response.error= Error reading server response '{errMsg}'


#unused by golang
vapi.bindings.skeleton.task.invalidstate=Service did not set the task state
vapi.bindings.stub.jsonrpc.unsupported=JSON-RPC connector not supported for this invocation
vapi.bindings.typeconverter.blob.base64.decode.error=Error in base64 decoding '{errMsg}'
vapi.bindings.typeconverter.datetime.deserialize.day.invalid=Datetime string has an invalid day field= '%s'
vapi.bindings.typeconverter.datetime.deserialize.hour.invalid=Datetime string has an invalid hour field= '%s'
vapi.bindings.typeconverter.datetime.deserialize.invalid.format=Datetime string '%s' does not match expected pattern '%s'
vapi.bindings.typeconverter.datetime.deserialize.invalid.time=Datetime string '%s' is invalid= '%s'
vapi.bindings.typeconverter.datetime.deserialize.minute.invalid=Datetime string has an invalid minute field= '%s'
vapi.bindings.typeconverter.datetime.deserialize.month.invalid=Datetime string has an invalid month field= '%s'
vapi.bindings.typeconverter.datetime.deserialize.second.invalid=Datetime string has an invalid second field= '%s'
vapi.bindings.typeconverter.datetime.serialize.invalid.tz=Datetime object '%s' should be in UTC timezone
vapi.bindings.typeconverter.dict.missing.key=Key %s missing from dict
vapi.bindings.typeconverter.enum.invalid.enum.type=Invalid enum type= %s
vapi.bindings.typeconverter.invalid=Invalid Type to convert from, expected %s, instead got %s
vapi.bindings.typeconverter.map.duplicate.key=List contains structure with duplicate key '%s'. Cannot convert to a Map
vapi.bindings.typeconverter.set.duplicate.element=List contains duplicate element '%s'. Cannot convert to a set
vapi.bindings.typeconverter.unexpected.enum.type=Expected enumeration of type %s, but received %s
vapi.bindings.typeconverter.unexpected.error.type=Expected VapiError instance, but received %s
vapi.bindings.typeconverter.unexpected.struct.class=Expected class of type %s, but received %s
vapi.bindings.typeconverter.unexpected.struct.type=Expected VapiStruct instance or python dictionary, but received %s
vapi.bindings.typeconverter.voiddef.expect.null=Expected null for void definition, but found object of type %s
vapi.data.definition.list.mismatch=The List content type does not match the content type of this list
vapi.data.definition.mismatch=The DataDefinition %s does not match this DataDefinition %s
vapi.data.dynamicstruct.validate.mismatch=Type mismatch, expected an object of %s, but got %s
vapi.data.invalid.double.inf=Double value %s is outside the limits of 64 bit binary floating point number
vapi.data.invalid=Invalid Type, expected %s, instead got %s
vapi.data.list.add=Invalid entry for List= %s
vapi.data.optional.construct=Cannot construct an OptionalValue with the given DataValue
vapi.data.optional.getvalue.mismatch=Request for value of optional as a %s, but it is of type %s
vapi.data.optional.getvalue.unset=Request for value of optional when no value has been set
vapi.data.optional.mismatch=The other Optional data definition does not match data definition of this Optional
vapi.data.optional.validate=The given OptionalValue does not match the OptionalDefinition
vapi.data.serializers.python.unsupported.python.type=Unsupported python type '%s' provided for field '%s' in a dynamic structure
vapi.data.serializers.invalid.type=Unsupported python type '%s' provided
vapi.data.serializers.security_context.unsupported=Security Context with scheme %s is not supported for this operation
vapi.data.serializers.rest.marshall.error=Error serializing REST requests
vapi.data.structref.already.resolved=Structure reference of type %s is already resolved
vapi.data.structref.resolve.type.mismatch=Structure reference of type %s cannot point to structure definition %s
vapi.data.structure.field.extra=Structure value has a field not in the definition = %s
vapi.data.structure.field.unexpected=Found unexpected field '%s' in structure '%s'
vapi.data.structure.getfield.mismatch=Unable to get field %s, request for field as a %s instead of %s
vapi.data.structure.getfield.unknown=Unable to get field %s, no field of that name found
vapi.data.structure.setfield.invalid=Attempt to set an invalid field %s on the structure.
vapi.data.structure.setfield.null=Cannot set a field (%s) of a struct to be null
vapi.decimal.canonicalization=Invalid decimal data provided
vapi.introspection.invalid.type=Type %s is invalid
vapi.json.read.field.extra=JSON object has extra field(s)= '%s'
vapi.json.read.field.missing=JSON object has missing field(s)= '%s'
vapi.json.read.notimplemented=JSON parsing of type '%s' not implemented
vapi.json.write.notimplemented=JSON generation of type '%s' not implemented
vapi.method.invoke.exception=Error in method invocation %s
vapi.method.notimplemented=Vapi method %s in %s is not implemented
vapi.method.output.invalid.definition=Invalid output definition
vapi.method.output.invalid.null=There was no output provided from the command.
vapi.method.output.invalid=Invalid output provided from method %s.
vapi.provider.interface.duplicate=Cannot register interface %s, an interface with that name already exists.
vapi.provider.interposer.already.added=Interposer %s is already registered with the aggregator
vapi.provider.interposer.already.removed=Interposer %s is already unregistered from the aggregator
vapi.signature.canonicalization=Cannot canonicalize the provided data
vapi.task.invalid.error=Invalid error %s reported from method %s.
vapi.task.invalid.result=Invalid output %s provided from method %s.`)
