/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/realized_state"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/segments"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
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
			"tag": getTagsSchema(),
			"port": {
				Type:        schema.TypeList,
				Description: "Tag specificiation for corresponding segment port",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_path": getPolicyPathSchema(true, true, "Segment path where VM port should be tagged"),
						"tag":          getTagsSchema(),
					},
				},
			},
		},
	}
}

func listAllPolicyVirtualMachines(connector *client.RestConnector, m interface{}) ([]model.VirtualMachine, error) {
	client := realized_state.NewVirtualMachinesClient(connector)
	var results []model.VirtualMachine
	boolFalse := false
	var cursor *string
	total := 0

	enforcementPointPath := getPolicyEnforcementPointPath(m)
	for {
		// NOTE: Search API doesn't filter by realized state resources
		// NOTE: Contrary to the spec, this API does not populate cursor and result count
		// parameters, respects cursor input. Therefore we determine end of VM list by
		// looking for empty result.
		vms, err := client.List(cursor, &enforcementPointPath, &boolFalse, nil, nil, &boolFalse, nil)
		if err != nil {
			return results, err
		}
		results = append(results, vms.Results...)
		if len(vms.Results) == 0 {
			// no more results
			return results, nil
		}
		if total == 0 && vms.ResultCount != nil {
			// first response
			total = int(*vms.ResultCount)
		}
		cursor = vms.Cursor
		if cursor == nil {
			resultCount := strconv.Itoa(len(results))
			cursor = &resultCount
		}

		if (total > 0) && (len(results) >= total) {
			return results, nil
		}
	}
}

func listAllPolicySegmentPorts(connector *client.RestConnector, segmentPath string) ([]model.SegmentPort, error) {
	client := segments.NewPortsClient(connector)
	segmentID := getPolicyIDFromPath(segmentPath)
	var results []model.SegmentPort
	boolFalse := false
	var cursor *string
	total := 0

	for {
		vms, err := client.List(segmentID, cursor, &boolFalse, nil, nil, &boolFalse, nil)
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
			log.Printf("[DEBUG] Found %d ports for segment %s", len(results), segmentID)
			return results, nil
		}
	}
}

func listAllPolicyVifs(m interface{}) ([]model.VirtualNetworkInterface, error) {

	client := enforcement_points.NewVifsClient(getPolicyConnector(m))
	var results []model.VirtualNetworkInterface
	var cursor *string
	total := 0

	enforcementPointPath := getPolicyEnforcementPoint(m)
	for {
		// NOTE: Search API doesn't filter by realized state resources
		vifs, err := client.List(enforcementPointPath, cursor, nil, nil, nil, nil, nil)
		if err != nil {
			return results, err
		}
		results = append(results, vifs.Results...)
		if total == 0 && vifs.ResultCount != nil {
			// first response
			total = int(*vifs.ResultCount)
		}
		cursor = vifs.Cursor
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
	client := enforcement_points.NewVirtualMachinesClient(connector)

	tagUpdate := model.VirtualMachineTagsUpdate{
		Tags:             tags,
		VirtualMachineId: &externalID,
	}
	return client.Updatetags(getPolicyEnforcementPoint(m), tagUpdate)
}

func listPolicyVifAttachmentsForVM(m interface{}, externalID string) ([]string, error) {
	var vifAttachmentIds []string
	vifs, err := listAllPolicyVifs(m)
	if err != nil {
		return vifAttachmentIds, err
	}

	for _, vif := range vifs {
		if (vif.LportAttachmentId != nil) && (vif.OwnerVmId != nil) && *vif.OwnerVmId == externalID {
			vifAttachmentIds = append(vifAttachmentIds, *vif.LportAttachmentId)
		}
	}

	return vifAttachmentIds, nil
}

func updateNsxtPolicyVMPortTags(connector *client.RestConnector, externalID string, portTags []interface{}, m interface{}, isDelete bool) error {

	client := segments.NewPortsClient(connector)

	vifAttachmentIds, err := listPolicyVifAttachmentsForVM(m, externalID)
	if err != nil {
		return err
	}

	for _, portTag := range portTags {
		data := portTag.(map[string]interface{})
		segmentPath := data["segment_path"].(string)
		var tags []model.Tag
		if !isDelete {
			tags = getPolicyTagsFromSet(data["tag"].(*schema.Set))
		}

		ports, portsErr := listAllPolicySegmentPorts(connector, segmentPath)
		if portsErr != nil {
			return portsErr
		}
		for _, port := range ports {
			if port.Attachment == nil || port.Attachment.Id == nil {
				continue
			}

			for _, attachment := range vifAttachmentIds {
				if attachment == *port.Attachment.Id {
					port.Tags = tags
					log.Printf("[DEBUG] Updating port %s with %d tags", *port.Path, len(tags))
					segmentID := getPolicyIDFromPath(segmentPath)
					_, err = client.Update(segmentID, *port.Id, port)
					if err != nil {
						return err
					}
					break
				}
			}
		}
	}

	return nil
}

func setPolicyVMPortTagsInSchema(d *schema.ResourceData, m interface{}, externalID string) error {

	connector := getPolicyConnector(m)
	vifAttachmentIds, err := listPolicyVifAttachmentsForVM(m, externalID)
	if err != nil {
		return err
	}

	portTags := d.Get("port").([]interface{})
	var actualPortTags []map[string]interface{}
	for _, portTag := range portTags {
		data := portTag.(map[string]interface{})
		segmentPath := data["segment_path"].(string)

		ports, portsErr := listAllPolicySegmentPorts(connector, segmentPath)
		if portsErr != nil {
			return portsErr
		}
		for _, port := range ports {
			if port.Attachment == nil || port.Attachment.Id == nil {
				continue
			}

			for _, attachment := range vifAttachmentIds {
				if attachment == *port.Attachment.Id {
					tags := make(map[string]interface{})
					tags["segment_path"] = segmentPath
					tags["tag"] = initPolicyTagsSet(port.Tags)
					actualPortTags = append(actualPortTags, tags)
				}
			}
		}
	}

	d.Set("port", actualPortTags)

	return nil
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

	return setPolicyVMPortTagsInSchema(d, m, *vm.ExternalId)
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

	portTags := d.Get("port").([]interface{})
	err = updateNsxtPolicyVMPortTags(connector, *vm.ExternalId, portTags, m, false)
	if err != nil {
		return handleCreateError("Segment Port Tag", *vm.ExternalId, err)
	}

	d.SetId(*vm.ExternalId)
	d.Set("port", portTags)

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

	portTags := d.Get("port").([]interface{})
	err = updateNsxtPolicyVMPortTags(connector, *vm.ExternalId, portTags, m, true)
	if err != nil {
		return handleCreateError("Segment Port Tag", *vm.ExternalId, err)
	}

	return err
}
