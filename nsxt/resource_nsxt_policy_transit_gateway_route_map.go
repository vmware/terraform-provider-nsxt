// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"maps"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	tgwrouting "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways/routing"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliTGWRouteMapsClient = tgwrouting.NewRouteMapsClient

func resourceNsxtPolicyTransitGatewayRouteMap() *schema.Resource {
	entrySchema := getPolicyRouteMapEntrySchema()
	// Local Manager rejects this GM-only field, even when false. Remove it
	// from the fresh TGW schema without changing the gateway route-map schema.
	if setSchema, ok := entrySchema.Schema["set"].Elem.(*schema.Resource); ok {
		delete(setSchema.Schema, "prefer_global_v6_next_hop")
	}

	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayRouteMapCreate,
		Read:   resourceNsxtPolicyTransitGatewayRouteMapRead,
		Update: resourceNsxtPolicyTransitGatewayRouteMapUpdate,
		Delete: resourceNsxtPolicyTransitGatewayRouteMapDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtTGWRoutingChildResourceImporter,
		},
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"parent_path":  getPolicyPathSchema(true, true, "Policy path of the parent Transit Gateway"),
			"entry": {
				Type:        schema.TypeList,
				Description: "List of Route Map entries",
				Required:    true,
				Elem:        entrySchema,
			},
		},
	}
}

func resourceNsxtPolicyTransitGatewayRouteMapExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := cliTGWRouteMapsClient(sessionContext, connector)
	if client == nil {
		return false, fmt.Errorf("unsupported client type")
	}
	_, err := client.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving TransitGatewayRouteMap", err)
}

func tgwRouteMapBuildObj(d *schema.ResourceData, id string) (model.TransitGatewayRouteMap, error) {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	schemaEntries := d.Get("entry").([]interface{})
	var entries []model.RouteMapEntry
	for entryNo, entry := range schemaEntries {
		entryData, ok := entry.(map[string]interface{})
		if !ok {
			return model.TransitGatewayRouteMap{}, fmt.Errorf("entry %d must be an object", entryNo)
		}
		// Reuse match validation without invoking the gateway-specific set builder.
		matchData := maps.Clone(entryData)
		delete(matchData, "set")
		entryObj, err := policyGatewayRouteMapBuildEntry(d, entryNo, matchData)
		if err != nil {
			return model.TransitGatewayRouteMap{}, err
		}
		if data := getElemOrEmptyMapFromMap(entryData, "set"); len(data) > 0 {
			asPathPrepend, _ := data["as_path_prepend"].(string)
			community, _ := data["community"].(string)
			weightValue, _ := data["weight"].(int)
			weight := int64(weightValue)
			entrySet := &model.RouteMapEntrySet{Weight: &weight}
			if asPathPrepend != "" {
				entrySet.AsPathPrepend = &asPathPrepend
			}
			if community != "" {
				entrySet.Community = &community
			}
			if value, ok := d.GetOk(fmt.Sprintf("entry.%d.set.0.local_preference", entryNo)); ok {
				localPreference, ok := value.(int)
				if !ok {
					return model.TransitGatewayRouteMap{}, fmt.Errorf("entry %d local_preference must be an integer", entryNo)
				}
				localPreferenceValue := int64(localPreference)
				entrySet.LocalPreference = &localPreferenceValue
			}
			if value, ok := d.GetOk(fmt.Sprintf("entry.%d.set.0.med", entryNo)); ok {
				med, ok := value.(int)
				if !ok {
					return model.TransitGatewayRouteMap{}, fmt.Errorf("entry %d med must be an integer", entryNo)
				}
				medValue := int64(med)
				entrySet.Med = &medValue
			}
			entryObj.Set = entrySet
		}
		entries = append(entries, entryObj)
	}

	return model.TransitGatewayRouteMap{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Entries:     entries,
	}, nil
}

func setTGWRouteMapEntriesInSchema(d *schema.ResourceData, entries []model.RouteMapEntry) error {
	var entryList []interface{}
	for _, entry := range entries {
		schemaEntry := map[string]interface{}{
			"action":              entry.Action,
			"prefix_list_matches": entry.PrefixListMatches,
		}
		var commList []map[string]interface{}
		for _, comm := range entry.CommunityListMatches {
			commList = append(commList, map[string]interface{}{
				"criteria":       comm.Criteria,
				"match_operator": comm.MatchOperator,
			})
		}
		schemaEntry["community_list_match"] = commList
		var setList []map[string]interface{}
		if entry.Set != nil {
			setList = append(setList, map[string]interface{}{
				"as_path_prepend":  entry.Set.AsPathPrepend,
				"community":        entry.Set.Community,
				"med":              entry.Set.Med,
				"local_preference": entry.Set.LocalPreference,
				"weight":           entry.Set.Weight,
			})
		}
		schemaEntry["set"] = setList
		entryList = append(entryList, schemaEntry)
	}
	return d.Set("entry", entryList)
}

func resourceNsxtPolicyTransitGatewayRouteMapCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTransitGatewayRouteMapExists)
	if err != nil {
		return err
	}
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj, err := tgwRouteMapBuildObj(d, id)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating TransitGatewayRouteMap with ID %s", id)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWRouteMapsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Patch(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleCreateError("TransitGatewayRouteMap", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyTransitGatewayRouteMapRead(d, m)
}

func resourceNsxtPolicyTransitGatewayRouteMapRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayRouteMap ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWRouteMapsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "TransitGatewayRouteMap", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	return setTGWRouteMapEntriesInSchema(d, obj.Entries)
}

func resourceNsxtPolicyTransitGatewayRouteMapUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayRouteMap ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj, err := tgwRouteMapBuildObj(d, id)
	if err != nil {
		return err
	}
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	log.Printf("[INFO] Updating TransitGatewayRouteMap with ID %s", id)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWRouteMapsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Patch(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleUpdateError("TransitGatewayRouteMap", id, err)
	}
	return resourceNsxtPolicyTransitGatewayRouteMapRead(d, m)
}

func resourceNsxtPolicyTransitGatewayRouteMapDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayRouteMap ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWRouteMapsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Delete(parents[0], parents[1], parents[2], id); err != nil {
		return handleDeleteError("TransitGatewayRouteMap", id, err)
	}
	return nil
}
