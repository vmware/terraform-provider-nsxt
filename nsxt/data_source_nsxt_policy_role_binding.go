// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-nsxt/api/aaa"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyRoleBinding() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRoleBindingRead,

		Schema: map[string]*schema.Schema{
			"id":                   getDataSourceIDSchema(),
			"display_name":         getDataSourceDisplayNameSchema(),
			"description":          getDataSourceDescriptionSchema(),
			"revision":             getRevisionSchema(),
			"tag":                  getTagsSchema(),
			"name":                 {Type: schema.TypeString, Description: "User/Group's name", Optional: true, Computed: true},
			"type":                 {Type: schema.TypeString, Description: "Indicates the type of the user", Optional: true, Computed: true, ValidateFunc: validation.StringInSlice(roleBindingUserTypes, false)},
			"user_id":              {Type: schema.TypeString, Description: "Local user's numeric id", Computed: true},
			"identity_source_id":   {Type: schema.TypeString, Description: "ID of the external identity source", Optional: true, Computed: true},
			"identity_source_type": {Type: schema.TypeString, Description: "Type of the external identity source", Optional: true, Computed: true, ValidateFunc: validation.StringInSlice(roleBindingIdentitySourceTypes, false)},
			"roles_for_path": {
				Type:        schema.TypeSet,
				Description: "List of roles that are associated with the user, limiting them to a path",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path":  {Type: schema.TypeString, Description: "Path of the entity in parent hierarchy.", Computed: true},
						"roles": {Type: schema.TypeSet, Description: "Applicable roles", Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyRoleBindingRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	rbClient := aaa.NewRoleBindingsClient(sessionContext, connector)

	id := d.Get("id").(string)
	var obj nsxModel.RoleBinding
	var err error

	if id != "" {
		obj, err = rbClient.Get(id, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleDataSourceReadError(d, "RoleBinding", id, err)
		}
	} else {
		username := d.Get("name").(string)
		userType := d.Get("type").(string)
		if username == "" || userType == "" {
			return fmt.Errorf("either id or both name and type must be set to look up a role binding")
		}
		identitySourceID := d.Get("identity_source_id").(string)
		var identitySourceIdParam *string
		if identitySourceID != "" {
			identitySourceIdParam = &identitySourceID
		}
		userTypeParam := &userType
		rbList, listErr := rbClient.List(nil, identitySourceIdParam, nil, nil, &username, nil, nil, nil, nil, nil, nil, userTypeParam)
		if listErr != nil {
			return fmt.Errorf("failed to list RoleBindings for %s: %w", username, listErr)
		}
		var found *nsxModel.RoleBinding
		for i := range rbList.Results {
			rb := &rbList.Results[i]
			if *rb.Name != username || *rb.Type_ != userType {
				continue
			}
			if identitySourceID != "" && (rb.IdentitySourceId == nil || *rb.IdentitySourceId != identitySourceID) {
				continue
			}
			found = rb
			break
		}
		if found == nil {
			return fmt.Errorf("role binding not found for name=%q type=%q", username, userType)
		}
		obj = *found
	}

	bindingID := *obj.Id
	d.SetId(bindingID)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("revision", obj.Revision)
	d.Set("name", obj.Name)
	d.Set("type", obj.Type_)
	if obj.IdentitySourceId != nil {
		d.Set("identity_source_id", obj.IdentitySourceId)
	}
	if obj.IdentitySourceType != nil {
		d.Set("identity_source_type", obj.IdentitySourceType)
	}
	if obj.UserId != nil {
		d.Set("user_id", obj.UserId)
	}
	setRolesForPathInSchema(d, obj.RolesForPaths)

	return nil
}
