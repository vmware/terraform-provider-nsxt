// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliBfdProfilesClient = func(sessionContext utl.SessionContext, connector vapiProtocolClient.Connector) *infra.BfdProfileClientContext {
	return infra.NewBfdProfilesClient(sessionContext, connector)
}

func dataSourceNsxtPolicyBfdProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyBfdProfileRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
		},
	}
}

func dataSourceNsxtPolicyBfdProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := cliBfdProfilesClient(getSessionContext(d, m), connector)

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	if objID != "" {
		obj, err := client.Get(objID)
		if isNotFoundError(err) {
			return fmt.Errorf("BfdProfile with ID %s was not found", objID)
		}
		if err != nil {
			return fmt.Errorf("error reading BfdProfile %s: %v", objID, err)
		}
		d.SetId(*obj.Id)
		d.Set("display_name", obj.DisplayName)
		d.Set("description", obj.Description)
		d.Set("path", obj.Path)
		return nil
	}

	if objName == "" {
		return fmt.Errorf("error obtaining BfdProfile: id or display_name must be specified")
	}

	inc := false
	objList, err := client.List(nil, &inc, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("error listing BfdProfiles: %v", err)
	}

	var perfectMatch, prefixMatch []string
	type bfdEntry struct{ id, name, desc, path string }
	entries := map[string]bfdEntry{}
	for _, obj := range objList.Results {
		if obj.DisplayName == nil || obj.Id == nil {
			continue
		}
		name := *obj.DisplayName
		id := *obj.Id
		desc := ""
		path := ""
		if obj.Description != nil {
			desc = *obj.Description
		}
		if obj.Path != nil {
			path = *obj.Path
		}
		entries[id] = bfdEntry{id: id, name: name, desc: desc, path: path}
		if name == objName {
			perfectMatch = append(perfectMatch, id)
		} else if strings.HasPrefix(name, objName) {
			prefixMatch = append(prefixMatch, id)
		}
	}

	var matchID string
	if len(perfectMatch) == 1 {
		matchID = perfectMatch[0]
	} else if len(perfectMatch) > 1 {
		return fmt.Errorf("found multiple BfdProfiles with display_name '%s'", objName)
	} else if len(prefixMatch) == 1 {
		matchID = prefixMatch[0]
	} else if len(prefixMatch) > 1 {
		return fmt.Errorf("found multiple BfdProfiles with display_name starting with '%s'", objName)
	} else {
		return fmt.Errorf("BfdProfile with display_name '%s' was not found", objName)
	}

	e := entries[matchID]
	d.SetId(e.id)
	d.Set("display_name", e.name)
	d.Set("description", e.desc)
	d.Set("path", e.path)
	return nil
}
