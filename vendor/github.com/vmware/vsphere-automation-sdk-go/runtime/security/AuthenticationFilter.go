/* Copyright © 2019, 2021-2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import (
	"encoding/json"
	"fmt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"io/ioutil"
	"reflect"
	"strings"
)

// Enforces authentication schemes specified in authentication metadata file.
type AuthenticationFilter struct {
	authHandlers []AuthenticationHandler
	provider     core.APIProvider
	schemes      map[string]Scheme   // maps schemeId to scheme mapping
	packages     map[string][]string // maps packageId to scheme Id mapping
	services     map[string][]string // maps serviceId to scheme Id mapping
	operations   map[string][]string // maps operationId to scheme Id mapping
}

var anonOps = []string{
	"com.vmware.vapi.std.introspection.operation.list",
	"com.vmware.vapi.std.introspection.operation.get",
	"com.vmware.vapi.std.introspection.provider.get",
	"com.vmware.vapi.std.introspection.service.list",
	"com.vmware.vapi.std.introspection.service.get"}

func NewAuthenticationFilter(authHandlers []AuthenticationHandler, provider core.APIProvider, authnMetadataFilePath []string) (*AuthenticationFilter, error) {
	aFilter := &AuthenticationFilter{authHandlers: authHandlers, provider: provider, schemes: map[string]Scheme{},
		packages: map[string][]string{}, services: map[string][]string{}, operations: map[string][]string{}}
	err := aFilter.loadAuthnMetadataFromFile(authnMetadataFilePath...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return aFilter, nil
}

// read authentication metadata json files
func (a *AuthenticationFilter) loadAuthnMetadataFromFile(authnMetadataFilepaths ...string) error {
	for _, authnMetadataFilepath := range authnMetadataFilepaths {
		authnData, err := ioutil.ReadFile(authnMetadataFilepath)
		if err != nil {
			log.Error("Error reading authentication metadata file %s", authnMetadataFilepath)
			return err
		}
		err = a.loadAuthnMetadata(authnData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AuthenticationFilter) loadAuthnMetadata(authnData []byte) error {
	var authnMetadata AuthenticationMetadata
	err := json.Unmarshal(authnData, &authnMetadata)
	if err != nil {
		return err
	}
	componentName := authnMetadata.Authentication.Component.Name
	for key, val := range authnMetadata.Authentication.Component.Schemes {
		a.schemes[componentName+":"+key] = val
	}
	for key, val := range authnMetadata.Authentication.Component.Operations {
		loadAuthnKeys(key, val, componentName, a.operations)
	}
	for key, val := range authnMetadata.Authentication.Component.Packages {
		err = checkAnonymousScheme(key, val, componentName, a)
		if err != nil {
			return err
		}
		loadAuthnKeys(key, val, componentName, a.packages)
	}
	for key, val := range authnMetadata.Authentication.Component.Services {
		err = checkAnonymousScheme(key, val, componentName, a)
		if err != nil {
			return err
		}
		loadAuthnKeys(key, val, componentName, a.services)
	}

	return nil
}

func checkAnonymousScheme(key string, val interface{}, componentName string, filter *AuthenticationFilter) error {
	if valStr, ok := val.(string); ok {
		if isNoAuthScheme(componentName+":"+valStr, filter) {
			return fmt.Errorf("Invalid authentication metadata for %s. "+
				"Anonymous scheme should be assigned to individual operations.", key)
		}
	} else if valSlice, ok := val.([]interface{}); ok {
		for _, scheme := range valSlice {
			if isNoAuthScheme(componentName+":"+scheme.(string), filter) {
				return fmt.Errorf("Invalid authentication metadata for %s. "+
					"Anonymous scheme should be assigned to individual operations.", key)
			}
		}
	} else {
		return fmt.Errorf("Expected string or json array but found %s", reflect.TypeOf(val))
	}
	return nil
}

func isNoAuthScheme(schemeName string, filter *AuthenticationFilter) bool {
	if schemeInfo, ok := filter.schemes[schemeName]; ok {
		return schemeInfo.AuthenticationScheme == NO_AUTH
	}
	return false
}

func loadAuthnKeys(key string, val interface{}, componentName string, maps map[string][]string) {
	if valStr, ok := val.(string); ok {
		maps[key] = []string{componentName + ":" + valStr}
	} else if valSlice, ok := val.([]interface{}); ok {
		schemeResult := []string{}
		for _, scheme := range valSlice {
			schemeResult = append(schemeResult, componentName+":"+scheme.(string))
		}
		maps[key] = schemeResult
	} else {
		log.Errorf("Expected string or json array but found %s", reflect.TypeOf(val))
	}
}

// finds closestpackage for a given serviceId
// ex: if serviceId is com.vmware.vmc.svc, com.vmware.vmc is closest package compared to com.vmware
func (a *AuthenticationFilter) findClosestPackage(serviceId string) string {

	runes := []rune(serviceId)
	lastIndex := strings.LastIndex(serviceId, ".")
	if lastIndex < 0 {
		lastIndex = len(serviceId)
	}
	servicePackageName := string(runes[0:lastIndex])
	var matchlength = 0
	var closestPackage = ""
	for packageName, _ := range a.packages {
		if strings.HasPrefix(servicePackageName, packageName) {
			if len(packageName) > matchlength {
				matchlength = len(packageName)
				closestPackage = packageName
			}
		}
	}
	return closestPackage
}

// returns defined auth scheme for closest package of serviceId
func (a *AuthenticationFilter) packageSpecificScheme(serviceId string) []string {
	closestPackage := a.findClosestPackage(serviceId)
	if schemeNames, ok := a.packages[closestPackage]; ok {
		result := []string{}
		for _, schemeName := range schemeNames {
			if schemeInfo, ok := a.schemes[schemeName]; ok {
				result = append(result, schemeInfo.AuthenticationScheme)
			}
		}
		return result
	}
	log.Debugf("Could not find package specific auth scheme for %s", serviceId)
	return nil
}

// returns auth scheme for serviceID if found otherwise returns default value.
func (a *AuthenticationFilter) serviceSpecificScheme(serviceID string) []string {

	if schemeNames, ok := a.services[serviceID]; ok {
		result := []string{}
		for _, schemeName := range schemeNames {
			if schemeInfo, ok := a.schemes[schemeName]; ok {
				result = append(result, schemeInfo.AuthenticationScheme)
			}
		}
		return result
	}
	log.Debugf("Service specific authentication scheme for %s not found.", serviceID)
	return nil
}

// returns defined auth scheme for given operationID
func (a *AuthenticationFilter) operationSpecificScheme(operationFqn string) []string {
	if schemeNames, ok := a.operations[operationFqn]; ok {
		result := []string{}
		for _, schemeName := range schemeNames {
			if schemeInfo, ok := a.schemes[schemeName]; ok {
				result = append(result, schemeInfo.AuthenticationScheme)
			}
		}
		return result
	}
	log.Debugf("Operation specific authentication scheme for %s not found.", operationFqn)
	return nil
}

// compares two arrays x & y and check whether they have same content
func equals(x []string, y []string) bool {
	if len(x) != len(y) {
		return false
	}

	elemMap := make(map[string]int)
	for _, val := range x {
		elemMap[val]++
	}
	for _, val := range y {
		elemMap[val]--
		if elemMap[val] == 0 {
			delete(elemMap, val)
		}
	}
	if len(elemMap) > 0 {
		return false
	}
	return true
}

// Removes the duplicate elements in the slice
func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// returns list of allowed schemes for combination of serviceID and operationID
func (a *AuthenticationFilter) allowedSchemes(serviceID string, operationID string) []string {
	allowedSchemes := []string{}
	emptySlice := []string{""}
	methodIdentifier := core.NewMethodIdentifier(core.NewInterfaceIdentifier(serviceID), operationID)
	operationFqn := methodIdentifier.FullyQualifiedName()
	opScheme := a.operationSpecificScheme(operationFqn)
	if !equals(opScheme, emptySlice) {
		allowedSchemes = append(allowedSchemes, opScheme...)
	}
	serviceScheme := a.serviceSpecificScheme(serviceID)
	if !equals(serviceScheme, emptySlice) {
		allowedSchemes = append(allowedSchemes, serviceScheme...)
	}
	packageScheme := a.packageSpecificScheme(serviceID)
	if !equals(packageScheme, emptySlice) {
		allowedSchemes = append(allowedSchemes, packageScheme...)
	}
	if len(allowedSchemes) == 0 {
		for _, anonOp := range anonOps {
			if anonOp == operationFqn {
				allowedSchemes = append(allowedSchemes, NO_AUTH)
				break
			}
		}
	}
	return unique(allowedSchemes)
}

func (a *AuthenticationFilter) Invoke(serviceID string, operationId string,
	input data.DataValue, ctx *core.ExecutionContext) core.MethodResult {

	var isNoAuthAllowed bool = false
	var authnSchemeFound bool = false

	// Get required Authn Schemes for ServiceId and OperationID
	allowedAuthnSchemes := a.allowedSchemes(serviceID, operationId)

	// Checks for NO_AUTH
	if len(allowedAuthnSchemes) == 0 {
		errorVal := bindings.CreateErrorValueFromMessageId(bindings.OP_NOT_FOUND_ERROR_DEF,
			"vapi.authentication.metadata.required", nil)
		return core.NewMethodResult(nil, errorVal)
	} else {
		for _, allowedAuthnScheme := range allowedAuthnSchemes {
			if allowedAuthnScheme == NO_AUTH {
				isNoAuthAllowed = true
			}
		}
	}

	securityCtx := ctx.SecurityContext()
	if securityCtx == nil && !isNoAuthAllowed {
		errorVal := bindings.CreateErrorValueFromMessageId(bindings.UNAUTHENTICATED_ERROR_DEF,
			"vapi.security.authentication.invalid", nil)
		return core.NewMethodResult(nil, errorVal)
	}

	authScheme := securityCtx.Property(AUTHENTICATION_SCHEME_ID)
	if authScheme == nil {
		if isNoAuthAllowed {
			log.Debugf("Authn scheme Id is not provided but NO AUTH is allowed hence invoking the operation")
			return a.invokeProvider(serviceID, operationId, input, ctx)
		}
		errorVal := bindings.CreateErrorValueFromMessageId(bindings.UNAUTHENTICATED_ERROR_DEF,
			"vapi.security.authentication.invalid", nil)
		return core.NewMethodResult(nil, errorVal)
	}

	if requestScheme, ok := authScheme.(string); ok {

		if requestScheme == NO_AUTH && isNoAuthAllowed {
			log.Debugf("Provided requestScheme is NO_AUTH hence invoking api provider")
			return a.invokeProvider(serviceID, operationId, input, ctx)
		}

		for _, allowedAuthnScheme := range allowedAuthnSchemes {
			if allowedAuthnScheme == requestScheme {
				log.Debugf("Required Authentication scheme is provided")
				authnSchemeFound = true
				break
			}
		}

		if !authnSchemeFound && !isNoAuthAllowed {
			log.Debugf("Provided Authentication Scheme is not valid to invoke this operation ", operationId)
			args := map[string]string{
				"allowedSchemes": strings.Join(allowedAuthnSchemes, ","),
				"providedScheme": requestScheme,
			}
			errorVal := bindings.CreateErrorValueFromMessageId(bindings.UNAUTHENTICATED_ERROR_DEF,
				"vapi.security.authentication.scheme", args)
			return core.NewMethodResult(nil, errorVal)
		}

	} else {
		log.Debugf("Invalid Authentication Scheme present in request, invalid type : %s", reflect.TypeOf(authScheme).String())
		args := map[string]string{"type": reflect.TypeOf(authScheme).String()}
		errorVal := bindings.CreateErrorValueFromMessageId(bindings.UNAUTHENTICATED_ERROR_DEF,
			"vapi.security.authentication.scheme.invalid", args)
		return core.NewMethodResult(nil, errorVal)
	}

	var authnResult *UserIdentity
	var authError error
	for _, authHandlers := range a.authHandlers {
		authnResult, authError = authHandlers.Authenticate(securityCtx)
		if authError != nil {
			log.Errorf("Authentication failed: %v", authError)
			errorVal := bindings.CreateErrorValueFromMessageId(bindings.UNAUTHENTICATED_ERROR_DEF,
				"vapi.security.authentication.invalid", nil)
			return core.NewMethodResult(nil, errorVal)

		}
		if authnResult != nil {
			securityCtx.SetProperty(AUTHN_IDENTITY, authnResult)
			break
		}
	}

	if authnResult == nil && !isNoAuthAllowed {
		log.Info("Could not find supporting authentication handler")
		errorVal := bindings.CreateErrorValueFromMessageId(bindings.UNAUTHENTICATED_ERROR_DEF,
			"vapi.security.authentication.invalid", nil)
		return core.NewMethodResult(nil, errorVal)
	}

	return a.invokeProvider(serviceID, operationId, input, ctx)
}

func (a *AuthenticationFilter) invokeProvider(serviceID string, operationId string,
	input data.DataValue, ctx *core.ExecutionContext) core.MethodResult {
	return a.provider.Invoke(serviceID, operationId, input, ctx)
}
