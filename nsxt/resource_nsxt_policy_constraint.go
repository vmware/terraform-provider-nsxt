// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

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

	clientLayer "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var constraintTargetOwnerTypeValues = []string{
	model.Constraint_TARGET_OWNER_TYPE_GM,
	model.Constraint_TARGET_OWNER_TYPE_LM,
	model.Constraint_TARGET_OWNER_TYPE_ALL,
}

var constraintPathExample = getMultitenancyPathExample("/infra/constraints/[constraint]")

var constraintSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchema(false, false, false)),
	"instance_count": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"count": {
						Schema: schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "int",
							SdkFieldName: "Count",
						},
					},
					"operator": {
						Schema: schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "<=",
							ValidateFunc: validation.StringInSlice([]string{"<=", "<"}, false),
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "Operator",
						},
					},
					"target_resource_type": {
						Schema: schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "TargetResourceType",
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:      "list",
			SdkFieldName:    "ConstraintExpressions",
			PolymorphicType: metadata.PolymorphicTypeFlatten,
			ReflectType:     reflect.TypeOf(model.EntityInstanceCountConstraintExpression{}),
			BindingType:     model.EntityInstanceCountConstraintExpressionBindingType(),
			TypeIdentifier:  metadata.ResourceTypeTypeIdentifier,
			ResourceType:    model.ConstraintExpression_RESOURCE_TYPE_ENTITYINSTANCECOUNTCONSTRAINTEXPRESSION,
		},
	},
	"target_owner_type": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(constraintTargetOwnerTypeValues, false),
			Default:      model.Constraint_TARGET_OWNER_TYPE_ALL,
			Optional:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "TargetOwnerType",
		},
	},
	"target": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"path_prefix": {
						Schema: schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: pathPrefixDiffSupress,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "PathPrefix",
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "Targets",
			ReflectType:  reflect.TypeOf(model.ConstraintTarget{}),
		},
	},
	"message": {
		Schema: schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "Message",
		},
	},
}

func resourceNsxtPolicyConstraint() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyConstraintCreate,
		Read:   resourceNsxtPolicyConstraintRead,
		Update: resourceNsxtPolicyConstraintUpdate,
		Delete: resourceNsxtPolicyConstraintDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathOrIDResourceImporter(constraintPathExample),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(constraintSchema),
	}
}

func addTrailingSlash(s string) string {
	if strings.HasSuffix(s, "/") {
		return s
	}

	return s + "/"
}

func pathPrefixDiffSupress(k, oldVal, newVal string, d *schema.ResourceData) bool {
	return addTrailingSlash(oldVal) == addTrailingSlash(newVal)
}

func resourceNsxtPolicyConstraintExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error

	client := clientLayer.NewConstraintsClient(sessionContext, connector)
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func polishConstraintTargets(constraint *model.Constraint) {
	// To improve user experience, add trailing slash to PathPrefix if needed
	if len(constraint.Targets) == 0 {
		return
	}

	for i, target := range constraint.Targets {
		if target.PathPrefix != nil {
			withSlash := addTrailingSlash(*target.PathPrefix)
			constraint.Targets[i].PathPrefix = &withSlash
		}
	}

}

func resourceNsxtPolicyConstraintCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("policy constraint resource requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyConstraintExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.Constraint{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, constraintSchema, "", nil); err != nil {
		return err
	}

	polishConstraintTargets(&obj)

	log.Printf("[INFO] Creating Constraint with ID %s", id)

	client := clientLayer.NewConstraintsClient(getSessionContext(d, m), connector)
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("Constraint", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyConstraintRead(d, m)
}

func resourceNsxtPolicyConstraintRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Constraint ID")
	}

	client := clientLayer.NewConstraintsClient(getSessionContext(d, m), connector)

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "Constraint", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, constraintSchema, "", nil)
}

func resourceNsxtPolicyConstraintUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Constraint ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.Constraint{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, constraintSchema, "", nil); err != nil {
		return err
	}

	polishConstraintTargets(&obj)

	client := clientLayer.NewConstraintsClient(getSessionContext(d, m), connector)
	_, err := client.Update(id, obj)
	if err != nil {
		return handleUpdateError("Constraint", id, err)
	}

	return resourceNsxtPolicyConstraintRead(d, m)
}

func resourceNsxtPolicyConstraintDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Constraint ID")
	}

	connector := getPolicyConnector(m)

	client := clientLayer.NewConstraintsClient(getSessionContext(d, m), connector)
	err := client.Delete(id)

	if err != nil {
		return handleDeleteError("Constraint", id, err)
	}

	return nil
}
