package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func resourceLogicalPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalPortCreate,
		Read:   resourceLogicalPortRead,
		Update: resourceLogicalPortUpdate,
		Delete: resourceLogicalPortDelete,

		Schema: map[string]*schema.Schema{
			"Revision": GetRevisionSchema(),
			"SystemOwned": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Indicates system owned resource",
				Optional:    true,
				Computed:    true,
			},
			"DisplayName": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"Description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"LogicalSwitchId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // Cannot change the logical switch of a logical port
			},
			"AdminState": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"SwitchingProfileIds": GetSwitchingProfileIdsSchema(),
			"Tags":                GetTagsSchema(),
			//TODO: add attachments
		},
	}
}

func resourceLogicalPortCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*nsxt.APIClient)

	name := d.Get("DisplayName").(string)
	description := d.Get("Description").(string)
	ls_id := d.Get("LogicalSwitchId").(string)
	admin_state := d.Get("AdminState").(string)
	profilesList := GetSwitchingProfileIdsFromSchema(d)
	tagList := GetTagsFromSchema(d)

	lp := manager.LogicalPort{
		DisplayName:         name,
		Description:         description,
		LogicalSwitchId:     ls_id,
		AdminState:          admin_state,
		SwitchingProfileIds: profilesList,
		Tags:                tagList}

	lp, _, err := nsxClient.LogicalSwitchingApi.CreateLogicalPort(nsxClient.Context, lp)

	if err != nil {
		return fmt.Errorf("Error while creating logical port %s: %v\n", lp.DisplayName, err)
	}

	resource_id := lp.Id
	d.SetId(resource_id)

	return resourceLogicalPortRead(d, m)
}

func resourceLogicalPortRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*nsxt.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical port ID from state during read")
	}
	logical_port, resp, err := nsxClient.LogicalSwitchingApi.GetLogicalPort(nsxClient.Context, id)

	if err != nil {
		return fmt.Errorf("Error while reading logical port %s: %v\n", id, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Logical port %s was not found\n", id)
		d.SetId("")
		return nil
	}

	d.Set("Revision", logical_port.Revision)
	d.Set("SystemOwned", logical_port.SystemOwned)
	d.Set("DisplayName", logical_port.DisplayName)
	d.Set("Description", logical_port.Description)
	d.Set("LogicalSwitchId", logical_port.LogicalSwitchId)
	d.Set("AdminState", logical_port.AdminState)
	SetSwitchingProfileIdsInSchema(d, logical_port.SwitchingProfileIds)
	SetTagsInSchema(d, logical_port.Tags)
	d.Set("Revision", logical_port.Revision)

	return nil
}

func resourceLogicalPortUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*nsxt.APIClient)

	lp_id := d.Id()
	name := d.Get("DisplayName").(string)
	description := d.Get("Description").(string)
	ls_id := d.Get("LogicalSwitchId").(string)
	admin_state := d.Get("AdminState").(string)
	profilesList := GetSwitchingProfileIdsFromSchema(d)
	tagList := GetTagsFromSchema(d)
	revision := int64(d.Get("Revision").(int))

	lp := manager.LogicalPort{DisplayName: name,
		Description:         description,
		LogicalSwitchId:     ls_id,
		AdminState:          admin_state,
		SwitchingProfileIds: profilesList,
		Tags:                tagList,
		Revision:            revision}

	lp, resp, err := nsxClient.LogicalSwitchingApi.UpdateLogicalPort(nsxClient.Context, lp_id, lp)
	if err != nil {
		return fmt.Errorf("Error while updating logical port %s: %v\n", lp_id, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Logical port %s was not found\n", lp_id)
		d.SetId("")
		return nil
	}
	fmt.Printf("Data %s\nresp %s", lp, resp)
	return resourceLogicalPortRead(d, m)
}

func resourceLogicalPortDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*nsxt.APIClient)

	lp_id := d.Id()
	if lp_id == "" {
		return fmt.Errorf("Error obtaining logical port ID from state during delete")
	}
	//TODO: add optional detach param
	localVarOptionals := make(map[string]interface{})

	resp, err := nsxClient.LogicalSwitchingApi.DeleteLogicalPort(nsxClient.Context, lp_id, localVarOptionals)

	if err != nil {
		return fmt.Errorf("Error while deleting logical port %s: %v\n", lp_id, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Logical port %s was not found\n", lp_id)
		d.SetId("")
		return nil
	}

	return nil
}
