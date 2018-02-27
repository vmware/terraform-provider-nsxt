package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/common"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
	"strings"
)

// Note - this resource is time consuming to configure since
// there is no API to search VMs by host ID

func resourceNsxtVMTags() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVMTagsCreate,
		Read:   resourceNsxtVMTagsRead,
		Update: resourceNsxtVMTagsUpdate,
		Delete: resourceNsxtVMTagsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Instance id",
				Required:    true,
			},
			"tag": getTagsSchema(),
		},
	}
}

func getVMList(nsxClient *api.APIClient) (*manager.VirtualMachineListResult, error) {

	localVarOptionals := make(map[string]interface{})
	vmList, resp, err := nsxClient.FabricApi.ListVirtualMachines(nsxClient.Context, localVarOptionals)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status returned during operation on VM Tags: %v", resp.StatusCode)
	}

	return &vmList, nil

}

func findVMIDByLocalID(nsxClient *api.APIClient, localID string) (*manager.VirtualMachine, error) {
	// VM provider (for example, vsphere provider) can supply host id(BIOS id) for the VM
	// NSX API needs NSX id for the update tags API call
	vmList, err := getVMList(nsxClient)
	if err != nil {
		return nil, err
	}

	for _, elem := range vmList.Results {
		for _, computeID := range elem.ComputeIds {
			log.Printf("[DEBUG] Inspecting compute id %s for VM %s", computeID, elem.LocalIdOnHost)
			slices := strings.Split(computeID, ":")
			if len(slices) == 2 && (slices[0] == "biosUuid") && (slices[1] == localID) {
				return &elem, nil
			}
		}
	}

	return nil, fmt.Errorf("Failed to find Virtual Machine with local id %s in inventory", localID)
}

func findVMByExternalID(nsxClient *api.APIClient, instanceID string) (*manager.VirtualMachine, error) {
	vmList, err := getVMList(nsxClient)
	if err != nil {
		return nil, err
	}

	for _, elem := range vmList.Results {
		if elem.ExternalId == instanceID {
			return &elem, nil
		}
	}

	return nil, fmt.Errorf("Failed to find Virtual Machine with id %s in inventory", instanceID)
}

func updateTags(nsxClient *api.APIClient, id string, tags []common.Tag) error {
	log.Printf("[DEBUG] Updating tags for %s", id)

	tagsUpdate := manager.VirtualMachineTagUpdate{
		ExternalId: id,
		Tags:       tags,
	}

	resp, err := nsxClient.FabricApi.UpdateVirtualMachineTagsUpdateTags(nsxClient.Context, tagsUpdate)

	if err != nil {
		return fmt.Errorf("Error during operation on Tags for VM %s: %v", id, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected status returned during operation on Tags for VM %d: %v", resp.StatusCode, id)
	}

	return nil
}

func resourceNsxtVMTagsCreate(d *schema.ResourceData, m interface{}) error {
	instanceID := d.Get("instance_id").(string)

	nsxClient := m.(*api.APIClient)
	vm, err := findVMIDByLocalID(nsxClient, instanceID)
	if err != nil {
		return fmt.Errorf("Error during VM retrieval: %v", err)
	}

	tags := getTagsFromSchema(d)

	err = updateTags(nsxClient, vm.ExternalId, tags)

	if err != nil {
		return err
	}

	d.SetId(vm.ExternalId)

	return resourceNsxtVMTagsRead(d, m)
}

func resourceNsxtVMTagsRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	vm, err := findVMByExternalID(nsxClient, id)
	if err != nil {
		return fmt.Errorf("Error during VM retrieval: %v", err)
	}

	setTagsInSchema(d, vm.Tags)

	return nil
}

func resourceNsxtVMTagsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtVMTagsCreate(d, m)
}

func resourceNsxtVMTagsDelete(d *schema.ResourceData, m interface{}) error {
	instanceID := d.Get("instance_id").(string)

	nsxClient := m.(*api.APIClient)
	vm, err := findVMIDByLocalID(nsxClient, instanceID)
	if err != nil {
		return fmt.Errorf("Error during VM retrieval: %v", err)
	}

	tags := make([]common.Tag, 0)
	err = updateTags(nsxClient, vm.ExternalId, tags)

	if err != nil {
		return err
	}

	return nil
}
