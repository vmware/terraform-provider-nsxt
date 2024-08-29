package nsxt

import (
	"log"

	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/vpcs/subnets"
)

func resourceNsxtVpcExternalAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcExternalAddressCreate,
		Read:   resourceNsxtVpcExternalAddressRead,
		Update: resourceNsxtVpcExternalAddressUpdate,
		Delete: resourceNsxtVpcExternalAddressDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: map[string]*schema.Schema{
			"parent_path":                getPolicyPathSchema(true, true, "port path for address binding"),
			"allocated_external_ip_path": getPolicyPathSchema(true, false, "allocated ip path"),
			"external_ip_address": {
				Type:        schema.TypeString,
				Description: "Computed external IP address",
				Computed:    true,
			},
		},
	}
}

func updatePort(d *schema.ResourceData, m interface{}, deleteFlow bool) error {
	portPath := d.Get("parent_path").(string)
	addressPath := d.Get("allocated_external_ip_path").(string)
	log.Printf("[DEBUG] Updating external address binding for port %s", portPath)

	parents, pathErr := parseStandardPolicyPathVerifySize(portPath, 5)
	if pathErr != nil {
		return pathErr
	}

	// Get port in order to Update
	connector := getPolicyConnector(m)
	portClient := subnets.NewPortsClient(connector)
	port, err := portClient.Get(parents[0], parents[1], parents[2], parents[3], parents[4])
	if err != nil {
		return err
	}

	if deleteFlow {
		port.ExternalAddressBinding = nil
	} else {
		port.ExternalAddressBinding = &model.ExternalAddressBinding{
			AllocatedExternalIpPath: &addressPath,
		}
	}

	_, err = portClient.Update(parents[0], parents[1], parents[2], parents[3], parents[4], port)

	return err
}

func resourceNsxtVpcExternalAddressCreate(d *schema.ResourceData, m interface{}) error {
	err := updatePort(d, m, false)
	if err != nil {
		return handleCreateError("External Address", "", err)
	}
	d.SetId(newUUID())
	return resourceNsxtVpcExternalAddressRead(d, m)
}

func resourceNsxtVpcExternalAddressRead(d *schema.ResourceData, m interface{}) error {
	portPath := d.Get("parent_path").(string)

	parents, pathErr := parseStandardPolicyPathVerifySize(portPath, 5)
	if pathErr != nil {
		return pathErr
	}

	connector := getPolicyConnector(m)
	portClient := subnets.NewPortsClient(connector)
	port, err := portClient.Get(parents[0], parents[1], parents[2], parents[3], parents[4])
	if err != nil {
		return handleReadError(d, "External Address", "", err)
	}

	if port.ExternalAddressBinding == nil {
		d.Set("allocated_external_ip_path", "")
		d.Set("external_ip_address", "")
		return nil
	}

	d.Set("allocated_external_ip_path", port.ExternalAddressBinding.AllocatedExternalIpPath)
	d.Set("external_ip_address", port.ExternalAddressBinding.ExternalIpAddress)
	return nil
}

func resourceNsxtVpcExternalAddressUpdate(d *schema.ResourceData, m interface{}) error {
	err := updatePort(d, m, false)
	if err != nil {
		return handleUpdateError("External Address", "", err)
	}
	return resourceNsxtVpcExternalAddressRead(d, m)
}

func resourceNsxtVpcExternalAddressDelete(d *schema.ResourceData, m interface{}) error {
	return updatePort(d, m, true)
}
