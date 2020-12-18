/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gm_tier0s "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/tier_0s"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)



func resourceNsxtPolicyTier0GatewayCommunityList() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTier0GatewayCommunityListCreate,
		Read:   resourceNsxtPolicyTier0GatewayCommunityListRead,
		Update: resourceNsxtPolicyTier0GatewayCommunityListUpdate,
		Delete: resourceNsxtPolicyTier0GatewayCommunityListDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyTier0GatewayCommunityListImport,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":         getNsxIDSchema(),
			"path":           getPathSchema(),
			"display_name":   getDisplayNameSchema(),
			"description":    getDescriptionSchema(),
			"revision":       getRevisionSchema(),
			"tag":            getTagsSchema(),
			"communities": {
				Type:        schema.TypeList,
				Description: "List of BGP community entries",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceNsxtPolicyTier0GatewayCommunityListExists(tier0Id string, id string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {
        var err error
        if isGlobalManager {
            client := gm_tier0s.NewDefaultCommunityListsClient(connector)
             _, err = client.Get(tier0Id, id)
        } else {
            client := tier_0s.NewDefaultCommunityListsClient(connector)
            _, err = client.Get(tier0Id, id)
        }
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyTier0GatewayCommunityListCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id := d.Get("nsx_id").(string)
	tier0Path := d.Get("gateway_path").(string)
	tier0ID := getPolicyIDFromPath(tier0Path)

	if id == "" {
		id = newUUID()
	} else {
		var err error
		if isPolicyGlobalManager(m) {
			client := gm_tier0s.NewDefaultCommunityListsClient(connector)
			_, err = client.Get(tier0ID, id)
		} else {
			client := tier_0s.NewDefaultCommunityListsClient(connector)
			_, err = client.Get(tier0ID, id)
		}
		if err == nil {
			return fmt.Errorf("Community List with ID '%s' already exists on Tier0 Gateway %s", id, tier0ID)
		} else if !isNotFoundError(err) {
			return err
		}
	}


	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	communities := d.Get("communities").(*schema.Set).List()
	var communitiesList []string
	for _, community := range communities {
		communitiesList = append(communitiesList, community.(string))
	}

	tags := getPolicyTagsFromSchema(d)
        

	obj := model.CommunityList{
		DisplayName:    &displayName,
		Description:    &description,
		Tags:           tags,
		Communities: communitiesList,
	}

	var err error
	// Create the resource using PATCH
	log.Printf("[INFO] Creating Tier0GatewayCommunityList with ID %s", id)
        if isPolicyGlobalManager(m) {
            gmObj, convErr := convertModelBindingType(obj, model.CommunityListBindingType(), gm_model.CommunityListBindingType())
            if convErr != nil {
                return convErr
            }
	    client := gm_tier0s.NewDefaultCommunityListsClient(connector)
            err = client.Patch(tier0ID, id, gmObj.(gm_model.CommunityList))
        } else {
	    client := tier_0s.NewDefaultCommunityListsClient(connector)
            err = client.Patch(tier0ID, id, obj)
        }
	if err != nil {
		return handleCreateError("Tier0GatewayCommunityList", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTier0GatewayCommunityListRead(d, m)
}

func resourceNsxtPolicyTier0GatewayCommunityListRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier0GatewayCommunityList ID")
	}
	tier0Path := d.Get("gateway_path").(string)
	tier0ID := getPolicyIDFromPath(tier0Path)

        var obj model.CommunityList
        if isPolicyGlobalManager(m) {
            client := gm_tier0s.NewDefaultCommunityListsClient(connector)
            gmObj, err := client.Get(tier0ID, id)
            if err != nil {
                return handleReadError(d, "Tier0GatewayCommunityList", id, err)
            }

            lmObj, err := convertModelBindingType(gmObj, gm_model.CommunityListBindingType(), model.CommunityListBindingType())
            if err != nil {
                return err
            }
            obj = lmObj.(model.CommunityList)
        } else {
	    	client := tier_0s.NewDefaultCommunityListsClient(connector)
            var err error
            obj, err = client.Get(tier0ID, id)
            if err != nil {
                return handleReadError(d, "Tier0GatewayCommunityList", id, err)
            }
        }

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("communities", obj.Communities)

	return nil
}

func resourceNsxtPolicyTier0GatewayCommunityListUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier0GatewayCommunityList ID")
	}
	tier0Path := d.Get("gateway_path").(string)
	tier0ID := getPolicyIDFromPath(tier0Path)

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	communities := d.Get("communities").(*schema.Set).List()
	var communitiesList []string
	for _, community := range communities {
		communitiesList = append(communitiesList, community.(string))
	}

	obj := model.CommunityList{
		DisplayName:    &displayName,
		Description:    &description,
		Tags:           tags,
		Communities: communitiesList,
	}

	var err error
	// Create the resource using PATCH
	log.Printf("[INFO] Creating Tier0GatewayCommunityList with ID %s", id)
	if isPolicyGlobalManager(m) {
		gmObj, convErr := convertModelBindingType(obj, model.CommunityListBindingType(), gm_model.CommunityListBindingType())
		if convErr != nil {
			return convErr
		}
		client := gm_tier0s.NewDefaultCommunityListsClient(connector)
		err = client.Patch(tier0ID, id, gmObj.(gm_model.CommunityList))
	} else {
		client := tier_0s.NewDefaultCommunityListsClient(connector)
		err = client.Patch(tier0ID, id, obj)
	}
	if err != nil {
		return handleCreateError("Tier0GatewayCommunityList", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTier0GatewayCommunityListRead(d, m)
}

func resourceNsxtPolicyTier0GatewayCommunityListDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier0GatewayCommunityList ID")
	}
	tier0Path := d.Get("gateway_path").(string)
	tier0ID := getPolicyIDFromPath(tier0Path)

	connector := getPolicyConnector(m)
	var err error
	if isPolicyGlobalManager(m) {
		client := gm_tier0s.NewDefaultCommunityListsClient(connector)
		err = client.Delete(tier0ID, id)
	} else {
		client := tier_0s.NewDefaultCommunityListsClient(connector)
		err = client.Delete(tier0ID, id)
	}

	if err != nil {
		return handleDeleteError("Tier0GatewayCommunityList", id, err)
	}

	return nil
}

func resourceNsxtPolicyTier0GatewayCommunityListImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	s := strings.Split(importID, "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("Please provide <gateway-id>/<community-list-id> as an input")
	}

	gwID := s[0]
	connector := getPolicyConnector(m)
	var tier0GW model.Tier0
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewDefaultTier0sClient(connector)
		gmObj, err1 := client.Get(gwID)
		if err1 != nil {
			return nil, err1
		}

		convertedObj, err2 := convertModelBindingType(gmObj, model.Tier0BindingType(), model.Tier0BindingType())
		if err2 != nil {
			return nil, err2
		}

		tier0GW = convertedObj.(model.Tier0)
	} else {
		client := infra.NewDefaultTier0sClient(connector)
		var err error
		tier0GW, err = client.Get(gwID)
		if err != nil {
			return nil, err
		}
	}

	d.Set("gateway_path", tier0GW.Path)

	d.SetId(s[1])
	return []*schema.ResourceData{d}, nil
}
