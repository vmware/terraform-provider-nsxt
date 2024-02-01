/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var segmentSecurityProfileSchema = map[string]*extendedSchema{
	"nsx_id":       getExtendedSchema(getNsxIDSchema()),
	"path":         getExtendedSchema(getPathSchema()),
	"display_name": getExtendedSchema(getDisplayNameSchema()),
	"description":  getExtendedSchema(getDescriptionSchema()),
	"revision":     getExtendedSchema(getRevisionSchema()),
	"tag":          getExtendedSchema(getTagsSchema()),
	"context":      getExtendedSchema(getContextSchema(false, false)),
	"bpdu_filter_allow": {
		s: schema.Schema{
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.IsMACAddress,
			},
			Optional: true,
		},
		m: metadata{
			skip: true,
		},
	},
	"bpdu_filter_enable": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		m: metadata{
			schemaType:   "bool",
			sdkFieldName: "BpduFilterEnable",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
	"dhcp_client_block_enabled": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		m: metadata{
			schemaType:   "bool",
			sdkFieldName: "DhcpClientBlockEnabled",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
	"dhcp_server_block_enabled": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		m: metadata{
			schemaType:   "bool",
			sdkFieldName: "DhcpServerBlockEnabled",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
	"dhcp_client_block_v6_enabled": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		m: metadata{
			schemaType:   "bool",
			sdkFieldName: "DhcpClientBlockV6Enabled",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
	"dhcp_server_block_v6_enabled": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		m: metadata{
			schemaType:   "bool",
			sdkFieldName: "DhcpServerBlockV6Enabled",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
	"non_ip_traffic_block_enabled": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		m: metadata{
			schemaType:   "bool",
			sdkFieldName: "NonIpTrafficBlockEnabled",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
	"ra_guard_enabled": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		m: metadata{
			schemaType:   "bool",
			sdkFieldName: "RaGuardEnabled",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
	"rate_limits_enabled": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		m: metadata{
			schemaType:   "bool",
			sdkFieldName: "RateLimitsEnabled",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
	"rate_limit": {
		s: schema.Schema{
			Type: schema.TypeList,
			Elem: &extendedResource{
				Schema: map[string]*extendedSchema{
					"rx_broadcast": {
						s: schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						m: metadata{
							schemaType:   "int",
							sdkFieldName: "RxBroadcast",
							testData: testdata{
								createValue: "100",
								updateValue: "1000",
							},
						},
					},
					"rx_multicast": {
						s: schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						m: metadata{
							schemaType:   "int",
							sdkFieldName: "RxMulticast",
							testData: testdata{
								createValue: "100",
								updateValue: "1000",
							},
						},
					},
					"tx_broadcast": {
						s: schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						m: metadata{
							schemaType:   "int",
							sdkFieldName: "TxBroadcast",
							testData: testdata{
								createValue: "100",
								updateValue: "1000",
							},
						},
					},
					"tx_multicast": {
						s: schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						m: metadata{
							schemaType:   "int",
							sdkFieldName: "TxMulticast",
							testData: testdata{
								createValue: "100",
								updateValue: "1000",
							},
						},
					},
				},
			},
			Optional: true,
			Computed: true,
		},
		m: metadata{
			schemaType:   "struct",
			sdkFieldName: "RateLimits",
			reflectType:  reflect.TypeOf(model.TrafficRateLimits{}),
		},
	},
}

func resourceNsxtPolicySegmentSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicySegmentSecurityProfileCreate,
		Read:   resourceNsxtPolicySegmentSecurityProfileRead,
		Update: resourceNsxtPolicySegmentSecurityProfileUpdate,
		Delete: resourceNsxtPolicySegmentSecurityProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtPolicyPathResourceImporter,
		},
		Schema: getSchemaFromExtendedSchema(segmentSecurityProfileSchema),
	}
}

func resourceNsxtPolicySegmentSecurityProfileExists(context utl.SessionContext, id string, connector client.Connector) (bool, error) {
	client := infra.NewSegmentSecurityProfilesClient(context, connector)
	_, err := client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicySegmentSecurityProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	// example of attribute that we prefer to handle manually in the code
	// it is marked with `skip` in metadata
	bpduFilterAllow := getStringListFromSchemaSet(d, "bpdu_filter_allow")

	obj := model.SegmentSecurityProfile{
		DisplayName:     &displayName,
		Description:     &description,
		Tags:            tags,
		BpduFilterAllow: bpduFilterAllow,
	}
	elem := reflect.ValueOf(&obj).Elem()
	schemaToStruct(elem, d, segmentSecurityProfileSchema, "", nil)

	log.Printf("[INFO] Sending SegmentSecurityProfile with ID %s", id)
	client := infra.NewSegmentSecurityProfilesClient(getSessionContext(d, m), connector)
	return client.Patch(id, obj, nil)
}

func resourceNsxtPolicySegmentSecurityProfileCreate(d *schema.ResourceData, m interface{}) error {

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtPolicySegmentSecurityProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtPolicySegmentSecurityProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("SegmentSecurityProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicySegmentSecurityProfileRead(d, m)
}

func resourceNsxtPolicySegmentSecurityProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SegmentSecurityProfile ID")
	}

	client := infra.NewSegmentSecurityProfilesClient(getSessionContext(d, m), connector)
	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "SegmentSecurityProfile", id, err)
	}

	elem := reflect.ValueOf(&obj).Elem()
	structToSchema(elem, d, segmentSecurityProfileSchema, "", nil)

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	return nil
}

func resourceNsxtPolicySegmentSecurityProfileUpdate(d *schema.ResourceData, m interface{}) error {

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SegmentSecurityProfile ID")
	}

	err := resourceNsxtPolicySegmentSecurityProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("SegmentSecurityProfile", id, err)
	}

	return resourceNsxtPolicySegmentSecurityProfileRead(d, m)
}

func resourceNsxtPolicySegmentSecurityProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SegmentSecurityProfile ID")
	}

	connector := getPolicyConnector(m)
	client := infra.NewSegmentSecurityProfilesClient(getSessionContext(d, m), connector)
	err := client.Delete(id, nil)

	if err != nil {
		return handleDeleteError("SegmentSecurityProfile", id, err)
	}

	return nil
}
