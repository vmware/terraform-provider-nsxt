// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	clientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var transitGatewayHaModeValues = []string{
	model.TransitGatewayHighAvailabilityConfig_HA_MODE_ACTIVE,
	model.TransitGatewayHighAvailabilityConfig_HA_MODE_STANDBY,
}

var transitGatewaySchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchemaExtended(true, false, false, true)),
	"transit_subnets": {
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
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "TransitSubnets",
		},
	},
	"is_default": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "IsDefault",
		},
	},
	"high_availability_config": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"ha_mode": {
						Schema: schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice(transitGatewayHaModeValues, false),
							Optional:     true,
							Default:      model.TransitGatewayHighAvailabilityConfig_HA_MODE_ACTIVE,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "HaMode",
						},
					},
					"edge_cluster_paths": {
						Schema: schema.Schema{
							Type: schema.TypeList,
							Elem: &metadata.ExtendedSchema{
								Schema: schema.Schema{
									Type:         schema.TypeString,
									ValidateFunc: validatePolicyPath(),
								},
								Metadata: metadata.Metadata{
									SchemaType: "string",
								},
							},
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "list",
							SdkFieldName: "EdgeClusterPaths",
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "struct",
			SdkFieldName: "HighAvailabilityConfig",
			ReflectType:  reflect.TypeOf(model.TransitGatewayHighAvailabilityConfig{}),
		},
	},
	"span": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"cluster_based_span": {
						Schema: schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							ExactlyOneOf: []string{
								"span.0.cluster_based_span",
								"span.0.zone_based_span",
							},
							Optional: true,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"span_path": {
										Schema: schema.Schema{
											Type:         schema.TypeString,
											Optional:     true,
											Computed:     true,
											ValidateFunc: validatePolicyPath(),
										},
										Metadata: metadata.Metadata{
											Skip: true,
										},
									},
								},
							},
						},
						Metadata: metadata.Metadata{
							Skip: true,
						},
					},
					"zone_based_span": {
						Schema: schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"zone_external_ids": {
										Schema: schema.Schema{
											Type:     schema.TypeList,
											Required: true,
											MinItems: 0,
											MaxItems: 10,
											Elem: &metadata.ExtendedSchema{
												Schema: schema.Schema{
													Type: schema.TypeString,
												},
												Metadata: metadata.Metadata{
													Skip: true,
												},
											},
										},
										Metadata: metadata.Metadata{
											Skip: true,
										},
									},
								},
							},
						},
						Metadata: metadata.Metadata{
							Skip: true,
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			Skip: true,
		},
	},
}

func resourceNsxtPolicyTransitGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayCreate,
		Read:   resourceNsxtPolicyTransitGatewayRead,
		Update: resourceNsxtPolicyTransitGatewayUpdate,
		Delete: resourceNsxtPolicyTransitGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathResourceImporter(transitGatewayPathExample),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(transitGatewaySchema),
	}
}

func resourceNsxtPolicyTransitGatewayExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error
	parents := getVpcParentsFromContext(sessionContext)
	client := clientLayer.NewTransitGatewaysClient(connector)
	_, err = client.Get(parents[0], parents[1], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func getSpanFromSchema(iSpan interface{}) (*data.StructValue, error) {
	if len(iSpan.([]interface{})) == 0 {
		return nil, nil
	}
	converter := bindings.NewTypeConverter()

	// We're limiting to one span of any kind in the schema
	span := iSpan.([]interface{})[0].(map[string]interface{})
	if len(span["cluster_based_span"].([]interface{})) > 0 {
		cbs := span["cluster_based_span"].([]interface{})[0].(map[string]interface{})
		spanPath := cbs["span_path"].(string)
		clusterBasedSpan := model.ClusterBasedSpan{
			SpanPath: &spanPath,
			Type_:    model.BaseSpan_TYPE_CLUSTERBASEDSPAN,
		}
		dataValue, errs := converter.ConvertToVapi(clusterBasedSpan, model.ClusterBasedSpanBindingType())
		if errs != nil {
			return nil, errs[0]
		}
		return dataValue.(*data.StructValue), nil
	}
	if len(span["zone_based_span"].([]interface{})) > 0 {
		zbs := span["zone_based_span"].([]interface{})[0].(map[string]interface{})
		zoneExternalIds := interfaceListToStringList(zbs["zone_external_ids"].([]interface{}))
		zoneBasedSpan := model.ZoneBasedSpan{
			ZoneExternalIds: zoneExternalIds,
			Type_:           model.BaseSpan_TYPE_ZONEBASEDSPAN,
		}
		dataValue, errs := converter.ConvertToVapi(zoneBasedSpan, model.ZoneBasedSpanBindingType())
		if errs != nil {
			return nil, errs[0]
		}
		return dataValue.(*data.StructValue), nil
	}

	return nil, nil
}

func resourceNsxtPolicyTransitGatewayCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyTransitGatewayExists)
	if err != nil {
		return err
	}

	parents := getVpcParentsFromContext(getSessionContext(d, m))
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags, tagErr := getValidatedTagsFromSchema(d)
	if tagErr != nil {
		return tagErr
	}

	obj := model.TransitGateway{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, transitGatewaySchema, "", nil); err != nil {
		return err
	}

	if util.NsxVersionHigherOrEqual("9.1.0") {
		obj.Span, err = getSpanFromSchema(d.Get("span"))
		if err != nil {
			return handleCreateError("TransitGateway", id, err)
		}

	}

	log.Printf("[INFO] Creating TransitGateway with ID %s", id)

	client := clientLayer.NewTransitGatewaysClient(connector)
	err = client.Patch(parents[0], parents[1], id, obj)
	if err != nil {
		return handleCreateError("TransitGateway", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTransitGatewayRead(d, m)
}

func resourceNsxtPolicyTransitGatewayRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining TransitGateway ID")
	}

	client := clientLayer.NewTransitGatewaysClient(connector)
	parents := getVpcParentsFromContext(getSessionContext(d, m))
	obj, err := client.Get(parents[0], parents[1], id)
	if err != nil {
		return handleReadError(d, "TransitGateway", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	if util.NsxVersionHigherOrEqual("9.1.0") && obj.Span != nil {
		span, err := setSpanFromSchema(obj.Span)

		if err != nil {
			return handleReadError(d, "TransitGateway", id, err)
		}
		d.Set("span", span)
	}

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, transitGatewaySchema, "", nil)
}

func setSpanFromSchema(span *data.StructValue) (interface{}, error) {
	converter := bindings.NewTypeConverter()
	base, errs := converter.ConvertToGolang(span, model.BaseSpanBindingType())

	schemaSpan := make(map[string]interface{})
	if errs != nil {
		return nil, errs[0]
	}
	switch base.(model.BaseSpan).Type_ {
	case model.BaseSpan_TYPE_CLUSTERBASEDSPAN:
		cbs, errs := converter.ConvertToGolang(span, model.ClusterBasedSpanBindingType())
		if errs != nil {
			return nil, errs[0]
		}
		clusterBasedSpan := cbs.(model.ClusterBasedSpan)

		elem := make(map[string]interface{})
		elem["span_path"] = clusterBasedSpan.SpanPath
		schemaSpan["cluster_based_span"] = []interface{}{elem}
		return []interface{}{schemaSpan}, nil

	case model.BaseSpan_TYPE_ZONEBASEDSPAN:
		zbs, errs := converter.ConvertToGolang(span, model.ZoneBasedSpanBindingType())
		if errs != nil {
			return nil, errs[0]
		}
		zoneBasedSpan := zbs.(model.ZoneBasedSpan)

		elem := make(map[string]interface{})
		elem["zone_external_ids"] = zoneBasedSpan.ZoneExternalIds
		schemaSpan["zone_based_span"] = []interface{}{elem}
		return []interface{}{schemaSpan}, nil
	}
	return nil, nil
}

func resourceNsxtPolicyTransitGatewayUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining TransitGateway ID")
	}

	parents := getVpcParentsFromContext(getSessionContext(d, m))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags, tagErr := getValidatedTagsFromSchema(d)
	if tagErr != nil {
		return tagErr
	}

	revision := int64(d.Get("revision").(int))

	obj := model.TransitGateway{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, transitGatewaySchema, "", nil); err != nil {
		return err
	}

	var err error
	if util.NsxVersionHigherOrEqual("9.1.0") {
		obj.Span, err = getSpanFromSchema(d.Get("span"))
		if err != nil {
			return handleCreateError("TransitGateway", id, err)
		}
	}

	client := clientLayer.NewTransitGatewaysClient(connector)
	_, err = client.Update(parents[0], parents[1], id, obj)
	if err != nil {
		return handleUpdateError("TransitGateway", id, err)
	}

	return resourceNsxtPolicyTransitGatewayRead(d, m)
}

func resourceNsxtPolicyTransitGatewayDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining TransitGateway ID")
	}

	connector := getPolicyConnector(m)
	parents := getVpcParentsFromContext(getSessionContext(d, m))

	client := clientLayer.NewTransitGatewaysClient(connector)
	err := client.Delete(parents[0], parents[1], id)

	if err != nil {
		return handleDeleteError("TransitGateway", id, err)
	}

	return nil
}
