/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkerrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/realized_state"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

// ErrNotAPolicyPath - Define an ignorable error for  policy path importer - to indicate that the given path is not a
// policy path and may be processed as an id - which is handy for legacy import method
var ErrNotAPolicyPath = errors.New("specified import identifier is not a policy path")

func getOrGenerateID2(d *schema.ResourceData, m interface{}, presenceChecker func(utl.SessionContext, string, client.Connector) (bool, error)) (string, error) {
	connector := getPolicyConnector(m)

	id := d.Get("nsx_id").(string)
	if id == "" {
		return newUUID(), nil
	}

	exists, err := presenceChecker(getSessionContext(d, m), id, connector)
	if err != nil {
		return "", err
	}

	if exists {
		return "", fmt.Errorf("Resource with id %s already exists", id)
	}

	return id, nil
}

func getOrGenerateID(d *schema.ResourceData, m interface{}, presenceChecker func(string, client.Connector, bool) (bool, error)) (string, error) {
	connector := getPolicyConnector(m)
	isGlobalManager := isPolicyGlobalManager(m)

	id := d.Get("nsx_id").(string)
	if id == "" {
		return newUUID(), nil
	}

	exists, err := presenceChecker(id, connector, isGlobalManager)
	if err != nil {
		return "", err
	}

	if exists {
		return "", fmt.Errorf("Resource with id %s already exists", id)
	}

	return id, nil
}

func newUUID() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func getPolicyTagsFromSet(tagSet *schema.Set) []model.Tag {
	tags := tagSet.List()
	var tagList []model.Tag
	for _, tag := range tags {
		data := tag.(map[string]interface{})
		tagScope := data["scope"].(string)
		tagTag := data["tag"].(string)
		elem := model.Tag{
			Scope: &tagScope,
			Tag:   &tagTag}

		tagList = append(tagList, elem)
	}
	return tagList
}

func initPolicyTagsSet(tags []model.Tag) []map[string]interface{} {
	var tagList []map[string]interface{}
	for _, tag := range tags {
		elem := make(map[string]interface{})
		elem["scope"] = tag.Scope
		elem["tag"] = tag.Tag
		tagList = append(tagList, elem)
	}
	return tagList
}

func getCustomizedPolicyTagsFromSchema(d *schema.ResourceData, schemaName string) []model.Tag {
	tags := d.Get(schemaName).(*schema.Set).List()
	tagList := make([]model.Tag, 0)
	for _, tag := range tags {
		data := tag.(map[string]interface{})
		tagScope := data["scope"].(string)
		tagTag := data["tag"].(string)
		elem := model.Tag{
			Scope: &tagScope,
			Tag:   &tagTag}

		tagList = append(tagList, elem)
	}
	return tagList
}

func setCustomizedPolicyTagsInSchema(d *schema.ResourceData, tags []model.Tag, schemaName string) {
	var tagList []map[string]interface{}
	for _, tag := range tags {
		elem := make(map[string]interface{})
		elem["scope"] = tag.Scope
		elem["tag"] = tag.Tag
		tagList = append(tagList, elem)
	}
	err := d.Set(schemaName, tagList)
	if err != nil {
		log.Printf("[WARNING] Failed to set tag in schema: %v", err)
	}
}

func getPolicyTagsFromSchema(d *schema.ResourceData) []model.Tag {
	return getCustomizedPolicyTagsFromSchema(d, "tag")
}

func setPolicyTagsInSchema(d *schema.ResourceData, tags []model.Tag) {
	setCustomizedPolicyTagsInSchema(d, tags, "tag")
}

func getPathListFromMap(data map[string]interface{}, attrName string) []string {
	pathList := interface2StringList(data[attrName].(*schema.Set).List())
	if len(pathList) == 0 {
		// Convert empty value to "ANY"
		pathList = append(pathList, "ANY")
	}

	return pathList
}

func setPathListInMap(data map[string]interface{}, attrName string, pathList []string) {
	if len(pathList) == 1 && pathList[0] == "ANY" {
		data[attrName] = nil
	} else {
		data[attrName] = pathList
	}
}

func getDomainFromResourcePath(rPath string) string {
	return getResourceIDFromResourcePath(rPath, "domains")
}

func getResourceIDFromResourcePath(rPath string, rType string) string {
	segments := strings.Split(rPath, "/")
	for i, seg := range segments {
		if seg == rType && i+1 < len(segments) {
			return segments[i+1]
		}
	}
	return ""
}

func nsxtDomainResourceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importDomain := defaultDomain
	importID := d.Id()
	s := strings.Split(importID, "/")
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err == nil {
		for i, seg := range s {
			if seg == "domains" {
				d.Set("domain", s[i+1])
			}
		}
		return rd, nil
	} else if !errors.Is(err, ErrNotAPolicyPath) {
		return rd, err
	}
	if len(s) == 2 {
		importDomain = s[0]
		d.SetId(s[1])
	} else {
		d.SetId(s[0])
	}

	d.Set("domain", importDomain)

	return []*schema.ResourceData{d}, nil
}

func getParameterFromPolicyPath(startDelimiter, endDelimiter, policyPath string) (string, error) {
	startIndex := strings.Index(policyPath, startDelimiter)
	endIndex := strings.Index(policyPath, endDelimiter)
	if startIndex < 0 || endIndex < 0 || (startIndex+len(startDelimiter)) > endIndex {
		return "", fmt.Errorf("failed to parse policy path %s, delimited by '%s' and '%s'", policyPath, startDelimiter, endDelimiter)
	}
	return policyPath[startIndex+len(startDelimiter) : endIndex], nil
}

func nsxtPolicyPathResourceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if errors.Is(err, ErrNotAPolicyPath) {
		return rd, nil
	} else if err != nil {
		return rd, err
	}
	return rd, nil
}

func nsxtPolicyPathResourceImporterHelper(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	if isPolicyPath(importID) {
		pathSegs := strings.Split(importID, "/")
		if strings.Contains(pathSegs[1], "infra") {
			d.SetId(pathSegs[len(pathSegs)-1])
		} else if pathSegs[1] == "orgs" && pathSegs[3] == "projects" {
			if len(pathSegs) < 5 {
				return nil, fmt.Errorf("invalid policy multitenancy path %s", importID)
			}
			// pathSegs[2] should contain the organization. Once we support multiple organization, it should be
			// assigned into the context as well
			contexts := make([]interface{}, 1)
			ctxMap := make(map[string]interface{})
			ctxMap["project_id"] = pathSegs[4]
			contexts[0] = ctxMap
			d.Set("context", contexts)
			d.SetId(pathSegs[len(pathSegs)-1])
		}
		return []*schema.ResourceData{d}, nil
	}
	return []*schema.ResourceData{d}, ErrNotAPolicyPath
}

func isPolicyPath(policyPath string) bool {
	pathSegs := strings.Split(policyPath, "/")
	if len(pathSegs) < 4 {
		return false
	} else if pathSegs[0] != "" || pathSegs[len(pathSegs)-1] == "" {
		return false
	} else if !strings.Contains(pathSegs[1], "infra") && pathSegs[1] != "orgs" {
		// must be infra, global-infra or orgs as of now
		return false
	}
	return true
}

func getPolicyIDFromPath(path string) string {
	tokens := strings.Split(path, "/")
	return tokens[len(tokens)-1]
}

func interfaceListToStringList(interfaces []interface{}) []string {
	var strList []string
	for _, elem := range interfaces {
		strList = append(strList, elem.(string))
	}
	return strList
}

func policyResourceNotSupportedError() error {
	return fmt.Errorf("This NSX policy resource is not supported with given provider settings")
}

func collectSeparatedStringListToMap(stringList []string, separator string) map[string]string {
	strMap := make(map[string]string)
	for _, elem := range stringList {
		segs := strings.Split(elem, separator)
		if len(segs) > 1 {
			strMap[segs[0]] = segs[1]
		}

	}
	return strMap
}

func stringListToCommaSeparatedString(stringList []string) *string {
	if len(stringList) > 0 {
		var str string
		for i, seg := range stringList {
			str += seg
			if i < len(stringList)-1 {
				str += ","
			}
		}
		return &str
	}
	return nil
}

func commaSeparatedStringToStringList(commaString string) []string {
	var strList []string
	for _, seg := range strings.Split(commaString, ",") {
		if seg != "" {
			strList = append(strList, seg)
		}
	}
	return strList
}

func nsxtPolicyWaitForRealizationStateConf(connector client.Connector, d *schema.ResourceData, realizedEntityPath string) *resource.StateChangeConf {
	client := realized_state.NewRealizedEntitiesClient(connector)
	pendingStates := []string{"UNKNOWN", "UNREALIZED"}
	targetStates := []string{"REALIZED", "ERROR"}
	stateConf := &resource.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {

			realizationResult, realizationError := client.List(realizedEntityPath, nil)
			if realizationError == nil {
				// Find the right entry
				for _, objInList := range realizationResult.Results {
					if objInList.State != nil {
						return objInList, *objInList.State, nil
					}
				}
				// Realization info not found yet
				return nil, "UNKNOWN", nil
			}
			return nil, "", realizationError
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 1 * time.Second,
		Delay:      1 * time.Second,
	}

	return stateConf
}

func getPolicyEnforcementPointPath(m interface{}) string {
	return "/infra/sites/default/enforcement-points/" + getPolicyEnforcementPoint(m)
}

func getGlobalPolicyEnforcementPointPathWithLocation(m interface{}, location string) string {
	return "/global-infra/sites/" + location + "/enforcement-points/" + getPolicyEnforcementPoint(m)
}

func convertModelBindingType(obj interface{}, sourceType bindings.BindingType, destType bindings.BindingType) (interface{}, error) {
	converter := bindings.NewTypeConverter()
	dataValue, err := converter.ConvertToVapi(obj, sourceType)
	if err != nil {
		return nil, err[0]
	}

	gmObj, err := converter.ConvertToGolang(dataValue, destType)
	if err != nil {
		return nil, err[0]
	}

	return gmObj, nil
}

func retryUponPreconditionFailed(readAndUpdate func() error, maxRetryAttempts int) error {
	// This retry specific to Precondition Error, and solution
	// here required refreshing the object, and updating revision
	// in request body. This can not be solved with SDK-based retry
	// functionality since it always retries with same request.
	var err error
	for i := 0; i <= maxRetryAttempts; i++ {
		err = readAndUpdate()
		if err == nil {
			return nil
		}

		if _, ok := err.(sdkerrors.InvalidRequest); !ok {
			// other type of error
			return err
		}

		log.Printf("[INFO] Refreshing object and repeating operation, attempt %d", i+1)
	}

	return err
}
