// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-nsxt/api/nsx"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliFailureDomainsClient = nsx.NewFailureDomainsClient

var edgeServices = []string{
	"active",
	"standby",
	"no_preference",
}

func resourceNsxtFailureDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtFailureDomainCreate,
		Read:   resourceNsxtFailureDomainRead,
		Update: resourceNsxtFailureDomainUpdate,
		Delete: resourceNsxtFailureDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"preferred_edge_services": {
				Type:         schema.TypeString,
				Description:  "Set preference for failure domain",
				Optional:     true,
				Default:      "no_preference",
				ValidateFunc: validation.StringInSlice(edgeServices, false),
			},
		},
	}
}

func resourceNsxtFailureDomainRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining FailureDomain ID")
	}

	sessionContext := utl.SessionContext{ClientType: utl.Local}
	client := cliFailureDomainsClient(sessionContext, connector)
	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "FailureDomain", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setMPTagsInSchema(d, obj.Tags)
	d.Set("revision", obj.Revision)

	preferPtr := obj.PreferredActiveEdgeServices
	preferStr := "no_preference"
	if preferPtr != nil {
		if *preferPtr {
			preferStr = "active"
		} else {
			preferStr = "standby"
		}
	}
	d.Set("preferred_edge_services", preferStr)
	return nil
}

func failureDomainSchemaToModel(d *schema.ResourceData) model.FailureDomain {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getMPTagsFromSchema(d)

	obj := model.FailureDomain{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	preferredEdgeServices := d.Get("preferred_edge_services").(string)
	if preferredEdgeServices != "no_preference" {
		activePrefer := true
		standbyPrefer := false
		if preferredEdgeServices == "active" {
			obj.PreferredActiveEdgeServices = &activePrefer
		} else {
			obj.PreferredActiveEdgeServices = &standbyPrefer
		}
	}
	return obj
}

func resourceNsxtFailureDomainCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	client := cliFailureDomainsClient(sessionContext, connector)

	failureDomain := failureDomainSchemaToModel(d)
	displayName := d.Get("display_name").(string)
	log.Printf("[INFO] Creating Failure Domain %s", displayName)
	obj, err := client.Create(failureDomain)
	if err != nil {
		return handleCreateError("Failure Domain", displayName, err)
	}
	d.SetId(*obj.Id)

	return resourceNsxtFailureDomainRead(d, m)
}

func resourceNsxtFailureDomainUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining FailureDomain ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	client := cliFailureDomainsClient(sessionContext, connector)

	failureDomain := failureDomainSchemaToModel(d)
	revision := int64(d.Get("revision").(int))
	failureDomain.Revision = &revision

	_, err := client.Update(id, failureDomain)
	if err != nil {
		return handleUpdateError("FailureDomain", id, err)
	}

	return resourceNsxtFailureDomainRead(d, m)
}

func resourceNsxtFailureDomainDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining FailureDomain ID")
	}
	connector := getPolicyConnector(m)
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	client := cliFailureDomainsClient(sessionContext, connector)
	err := client.Delete(id)
	if err != nil {
		return handleDeleteError("FailureDomain", id, err)
	}
	return nil
}
