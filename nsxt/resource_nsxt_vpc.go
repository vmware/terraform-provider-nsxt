/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	clientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/vpcs"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var vpcIPAddressTypeValues = []string{
	model.Vpc_IP_ADDRESS_TYPE_IPV4,
}

var vpcSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchema(true, false, false)),
	"private_ips": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: &metadata.ExtendedSchema{
				Schema: schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateIPCidr(),
				},
				Metadata: metadata.Metadata{
					SchemaType: "string",
				},
			},
			Optional: true,
			ForceNew: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "PrivateIps",
		},
	},
	"vpc_service_profile": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validatePolicyPath(),
			Optional:     true,
			Computed:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "VpcServiceProfile",
			OmitIfEmpty:  true,
		},
	},
	"load_balancer_vpc_endpoint": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"enabled": {
						Schema: schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "bool",
							SdkFieldName: "Enabled",
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "struct",
			SdkFieldName: "LoadBalancerVpcEndpoint",
			ReflectType:  reflect.TypeOf(model.LoadBalancerVPCEndpoint{}),
		},
	},
	"ip_address_type": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(vpcIPAddressTypeValues, false),
			Optional:     true,
			Default:      model.Vpc_IP_ADDRESS_TYPE_IPV4,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "IpAddressType",
		},
	},
	"short_id": {
		Schema: schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "ShortId",
			OmitIfEmpty:  true,
		},
	},
	"default_attachment": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"vpc_connectivity_profile": {
						Schema: schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validatePolicyPath(),
							Required:     true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "VpcConnectivityProfile",
							Skip:         true,
						},
					},
				},
			},
			Optional: true,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType: "struct",
			Skip:       true,
		},
	},
}

// VPC resource needs dedicated importer since its path is VPC path,
// but VPC does not need to be set in context
func nsxtVpcImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	if isPolicyPath(importID) {
		pathSegs := strings.Split(importID, "/")
		if len(pathSegs) != 7 || pathSegs[1] != "orgs" || pathSegs[3] != "projects" || pathSegs[5] != "vpcs" {
			return []*schema.ResourceData{d}, fmt.Errorf("invalid VPC path: %s", importID)
		}
		ctxMap := make(map[string]interface{})
		ctxMap["project_id"] = pathSegs[4]
		d.Set("context", []interface{}{ctxMap})
		d.SetId(pathSegs[6])
		return []*schema.ResourceData{d}, nil
	}
	return []*schema.ResourceData{d}, ErrNotAPolicyPath
}

func resourceNsxtVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcCreate,
		Read:   resourceNsxtVpcRead,
		Update: resourceNsxtVpcUpdate,
		Delete: resourceNsxtVpcDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVpcImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(vpcSchema),
	}
}

func resourceNsxtVpcExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error
	parents := getVpcParentsFromContext(sessionContext)
	client := clientLayer.NewVpcsClient(connector)
	_, err = client.Get(parents[0], parents[1], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID2(d, m, resourceNsxtVpcExists)
	if err != nil {
		return err
	}

	parents := getVpcParentsFromContext(getSessionContext(d, m))
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.Vpc{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating Vpc with ID %s", id)

	client := clientLayer.NewVpcsClient(connector)
	err = client.Patch(parents[0], parents[1], id, obj)
	if err != nil {
		return handleCreateError("Vpc", id, err)
	}

	if d.HasChanges("default_attachment") {
		err = patchVpcAttachment(d, connector, id, parents)
		if err != nil {
			return handleCreateError("Vpc Attachment", id, err)
		}
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcRead(d, m)
}

func resourceNsxtVpcRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Vpc ID")
	}

	client := clientLayer.NewVpcsClient(connector)
	parents := getVpcParentsFromContext(getSessionContext(d, m))
	obj, err := client.Get(parents[0], parents[1], id)
	if err != nil {
		return handleReadError(d, "Vpc", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	err = metadata.StructToSchema(elem, d, vpcSchema, "", nil)
	if err != nil {
		return handleReadError(d, "Vpc", id, err)
	}

	attachmentClient := vpcs.NewAttachmentsClient(connector)
	attachment, attErr := attachmentClient.Get(parents[0], parents[1], id, "default")
	if attErr != nil {
		return handleReadError(d, "Vpc Attachment", id, err)
	}

	var attachmentList []map[string]interface{}
	attachmentMap := make(map[string]interface{})
	attachmentMap["vpc_connectivity_profile"] = attachment.VpcConnectivityProfile
	attachmentList = append(attachmentList, attachmentMap)
	d.Set("default_attachment", attachmentList)
	return nil

}

func resourceNsxtVpcUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Vpc ID")
	}

	parents := getVpcParentsFromContext(getSessionContext(d, m))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.Vpc{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcSchema, "", nil); err != nil {
		return err
	}
	client := clientLayer.NewVpcsClient(connector)
	_, err := client.Update(parents[0], parents[1], id, obj)
	if err != nil {
		// Trigger partial update to avoid terraform updating state based on failed intent
		// TODO - move this into handleUpdateError
		d.Partial(true)
		return handleUpdateError("Vpc", id, err)
	}
	if d.HasChanges("default_attachment") {
		err = patchVpcAttachment(d, connector, id, parents)
		if err != nil {
			d.Partial(true)
			return handleUpdateError("Vpc Attachment", id, err)
		}
	}

	return resourceNsxtVpcRead(d, m)
}

func resourceNsxtVpcDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Vpc ID")
	}

	connector := getPolicyConnector(m)
	parents := getVpcParentsFromContext(getSessionContext(d, m))

	client := clientLayer.NewVpcsClient(connector)
	err := client.Delete(parents[0], parents[1], id)

	if err != nil {
		return handleDeleteError("Vpc", id, err)
	}

	return nil
}

func patchVpcAttachment(d *schema.ResourceData, connector client.Connector, vpcID string, parents []string) error {
	// Initialize to default profile in case default_attachment was removed
	profilePath := fmt.Sprintf("/orgs/%s/projects/%s/connectivity-profiles/default", parents[0], parents[1])
	defaultAttachment := d.Get("default_attachment")

	if defaultAttachment != nil {
		attachmentMap := defaultAttachment.([]interface{})
		if len(attachmentMap) > 0 {
			profilePath = attachmentMap[0].(map[string]interface{})["vpc_connectivity_profile"].(string)
		}
	}
	client := vpcs.NewAttachmentsClient(connector)
	obj := model.VpcAttachment{
		VpcConnectivityProfile: &profilePath,
	}
	return client.Patch(parents[0], parents[1], vpcID, "default", obj)
}
