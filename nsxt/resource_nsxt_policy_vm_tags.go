/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/realized_state"
	updateClient "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
	"strings"
)

var (
	nsxtPolicyBiosUUIDKey     = "biosUuid"
	nsxtPolicyInstanceUUIDKey = "instanceUuid"
)

func resourceNsxtPolicyVMTags() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyVMTagsCreate,
		Read:   resourceNsxtPolicyVMTagsRead,
		Update: resourceNsxtPolicyVMTagsUpdate,
		Delete: resourceNsxtPolicyVMTagsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Instance id",
				Required:    true,
			},
			"tag": getRequiredTagsSchema(),
		},
	}
}

func listAllPolicyVirtualMachines(connector *client.RestConnector, m interface{}) ([]model.VirtualMachine, error) {
	client := realized_state.NewDefaultVirtualMachinesClient(connector)
	var results []model.VirtualMachine
	boolFalse := false
	var cursor *string
	total := 0

	for {
		// NOTE: the search API does not return resource_type of VirtualMachine
		enforcementPointPath := getPolicyEnforcementPointPath(m)
		vms, err := client.List(cursor, &enforcementPointPath, &boolFalse, nil, nil, &boolFalse, nil)
		if err != nil {
			return results, err
		}
		results = append(results, vms.Results...)
		if total == 0 && vms.ResultCount != nil {
			// first response
			total = int(*vms.ResultCount)
		}
		cursor = vms.Cursor
		if len(results) >= total {
			return results, nil
		}
	}
}

func findNsxtPolicyVMByNamePrefix(connector *client.RestConnector, namePrefix string, m interface{}) ([]model.VirtualMachine, []model.VirtualMachine, error) {
	var perfectMatch, prefixMatch []model.VirtualMachine

	allVMs, err := listAllPolicyVirtualMachines(connector, m)
	if err != nil {
		log.Printf("[INFO] Error reading Virtual Machines when looking for name: %s", namePrefix)
		return perfectMatch, prefixMatch, err
	}

	for _, vmResult := range allVMs {
		if *vmResult.DisplayName == namePrefix {
			perfectMatch = append(perfectMatch, vmResult)
		} else if strings.HasPrefix(*vmResult.DisplayName, namePrefix) {
			prefixMatch = append(prefixMatch, vmResult)
		}
	}
	return perfectMatch, prefixMatch, nil
}

func findNsxtPolicyVMByID(connector *client.RestConnector, vmID string, m interface{}) (model.VirtualMachine, error) {
	var virtualMachineStruct model.VirtualMachine

	allVMs, err := listAllPolicyVirtualMachines(connector, m)
	if err != nil {
		log.Printf("[INFO] Error reading Virtual Machines when looking for ID: %s", vmID)
		return virtualMachineStruct, err
	}

	for _, vmResult := range allVMs {
		if (vmResult.ExternalId != nil) && *vmResult.ExternalId == vmID {
			return vmResult, nil
		}
		computeIDMap := collectSeparatedStringListToMap(vmResult.ComputeIds, ":")
		if computeIDMap[nsxtPolicyBiosUUIDKey] == vmID {
			return vmResult, nil
		} else if computeIDMap[nsxtPolicyInstanceUUIDKey] == vmID {
			return vmResult, nil
		}
	}
	return virtualMachineStruct, fmt.Errorf("Could not find Virtual Machine with ID: %s", vmID)
}

func updateNsxtPolicyVMTags(connector *client.RestConnector, externalID string, tags []model.Tag, m interface{}) error {
	client := updateClient.NewDefaultVirtualMachinesClient(connector)

	tagUpdate := model.VirtualMachineTagsUpdate{
		Tags:             tags,
		VirtualMachineId: &externalID,
	}
	return client.Updatetags(getPolicyEnforcementPoint(m), tagUpdate)
}

func resourceNsxtPolicyVMTagsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	vmID := d.Id()
	if vmID == "" {
		return fmt.Errorf("Error obtaining Virtual Machine ID")
	}

	vm, err := findNsxtPolicyVMByID(connector, vmID, m)
	if err != nil {
		return fmt.Errorf("Error during Virtual Machine retrieval: %v", err)
	}

	setPolicyTagsInSchema(d, vm.Tags)

	if d.Get("instance_id") == "" {
		// for import
		d.Set("instance_id", vm.ExternalId)
	}

	return nil
}

func resourceNsxtPolicyVMTagsCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	instanceID := d.Get("instance_id").(string)

	vm, err := findNsxtPolicyVMByID(connector, instanceID, m)
	if err != nil {
		return fmt.Errorf("Error finding Virtual Machine: %v", err)
	}

	tags := getPolicyTagsFromSchema(d)
	if tags == nil {
		tags = make([]model.Tag, 0)
	}
	err = updateNsxtPolicyVMTags(connector, *vm.ExternalId, tags, m)
	if err != nil {
		return handleCreateError("Virtual Machine Tag", *vm.ExternalId, err)
	}

	d.SetId(*vm.ExternalId)

	return resourceNsxtPolicyVMTagsRead(d, m)
}

func resourceNsxtPolicyVMTagsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyVMTagsCreate(d, m)
}

func resourceNsxtPolicyVMTagsDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	instanceID := d.Get("instance_id").(string)

	vm, err := findNsxtPolicyVMByID(connector, instanceID, m)
	if err != nil {
		return fmt.Errorf("Error finding Virtual Machine: %v", err)
	}

	tags := make([]model.Tag, 0)
	err = updateNsxtPolicyVMTags(connector, *vm.ExternalId, tags, m)

	if err != nil {
		return handleDeleteError("Virtual Machine Tag", *vm.ExternalId, err)
	}

	return err
}
