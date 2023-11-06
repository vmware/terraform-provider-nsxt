/* Copyright © 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/aaa"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// Only support local user at the moment
var roleBindingUserTypes = [](string){
	nsxModel.RoleBinding_TYPE_LOCAL_USER,
	nsxModel.RoleBinding_TYPE_REMOTE_USER,
	nsxModel.RoleBinding_TYPE_REMOTE_GROUP,
	nsxModel.RoleBinding_TYPE_PRINCIPAL_IDENTITY,
}

var roleBindingIdentitySourceTypes = [](string){
	nsxModel.RoleBinding_IDENTITY_SOURCE_TYPE_VIDM,
	nsxModel.RoleBinding_IDENTITY_SOURCE_TYPE_LDAP,
	nsxModel.RoleBinding_IDENTITY_SOURCE_TYPE_OIDC,
	nsxModel.RoleBinding_IDENTITY_SOURCE_TYPE_CSP,
}

func resourceNsxtPolicyUserManagementRoleBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyUserManagementRoleBindingCreate,
		Read:   resourceNsxtPolicyUserManagementRoleBindingRead,
		Update: resourceNsxtPolicyUserManagementRoleBindingUpdate,
		Delete: resourceNsxtPolicyUserManagementRoleBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"name": {
				Type:        schema.TypeString,
				Description: "User/Group's name",
				Required:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "Indicates the type of the user",
				Required:     true,
				ValidateFunc: validation.StringInSlice(roleBindingUserTypes, false),
			},
			"user_id": {
				Type:        schema.TypeString,
				Description: "Local user's numeric id",
				Computed:    true,
			},
			"identity_source_id": {
				Type:        schema.TypeString,
				Description: "ID of the external identity source",
				Optional:    true,
			},
			"identity_source_type": {
				Type:         schema.TypeString,
				Description:  "ID of the external identity source",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(roleBindingIdentitySourceTypes, false),
			},
			"roles_for_path": getRolesForPathSchema(false),
		},
	}
}

// getRolesForPathSchema return schema for RolesForPath, which is shared between role bindings and PI
func getRolesForPathSchema(forceNew bool) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of roles that are associated with the user, limiting them to a path",
		Required:    true,
		ForceNew:    forceNew,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"path": {
					Type:        schema.TypeString,
					Description: "Path of the entity in parent hierarchy.",
					Required:    true,
				},
				"role": {
					Type:        schema.TypeList,
					Description: "Applicable roles",
					Required:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"role": {
								Type:        schema.TypeString,
								Description: "Short identifier for the role",
								Required:    true,
								ValidateFunc: validation.StringMatch(
									regexp.MustCompile(
										`^[_a-z0-9-]+$`),
									"Must be a valid role identifier matching: ^[_a-z0-9-]+$"),
							},
							"role_display_name": {
								Type:        schema.TypeString,
								Description: "Display name for role",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

type rolesForPath map[string]rolesPerPath
type rolesPerPath map[string]bool

func getRolesForPathFromSchema(d *schema.ResourceData) rolesForPath {
	rolesForPathMap := make(rolesForPath)
	rolesForPathInput := d.Get("roles_for_path").([]interface{})
	for _, rolesPerPathInput := range rolesForPathInput {
		data := rolesPerPathInput.(map[string]interface{})
		path := data["path"].(string)
		roles := data["role"].([]interface{})
		rolesPerPathMap := make(rolesPerPath)
		for _, role := range roles {
			roleData := role.(map[string]interface{})
			roleInput := roleData["role"].(string)
			rolesPerPathMap[roleInput] = true
		}
		rolesForPathMap[path] = rolesPerPathMap
	}

	return rolesForPathMap
}

func setRolesForPathInSchema(d *schema.ResourceData, nsxRolesForPathList []nsxModel.RolesForPath) {
	var rolesForPathList []map[string]interface{}
	for _, nsxRolesForPath := range nsxRolesForPathList {
		elem := make(map[string]interface{})
		elem["path"] = nsxRolesForPath.Path
		var roles []map[string]interface{}
		for _, nsxRole := range nsxRolesForPath.Roles {
			rElem := make(map[string]interface{})
			rElem["role"] = nsxRole.Role
			rElem["role_display_name"] = nsxRole.RoleDisplayName
			roles = append(roles, rElem)
		}
		elem["role"] = roles
		rolesForPathList = append(rolesForPathList, elem)
	}
	err := d.Set("roles_for_path", rolesForPathList)
	if err != nil {
		log.Printf("[WARNING] Failed to set roles_for_path in schema: %v", err)
	}
}

func getRolesForPathList(d *schema.ResourceData) []nsxModel.RolesForPath {
	boolTrue := true
	rolesPerPathMap := getRolesForPathFromSchema(d)
	nsxRolesForPaths := make([]nsxModel.RolesForPath, 0)

	for k, v := range rolesPerPathMap {
		path := k
		nsxRoles := make([]nsxModel.Role, 0, len(v))
		for k := range v {
			roleID := k
			nsxRoles = append(nsxRoles, nsxModel.Role{
				Role: &roleID,
			})
		}
		nsxRolesForPaths = append(nsxRolesForPaths, nsxModel.RolesForPath{
			Path:  &path,
			Roles: nsxRoles,
		})
	}

	// Handle deletion of entire paths
	if d.HasChange("roles_for_path") {
		o, _ := d.GetChange("roles_for_path")
		oldRoles := o.([]interface{})
		for _, oldRole := range oldRoles {
			data := oldRole.(map[string]interface{})
			path := data["path"].(string)
			if _, ok := rolesPerPathMap[path]; ok {
				continue
			}
			// Add one role in the list to make NSX happy
			roles := data["role"].([]interface{})
			if len(roles) == 0 {
				continue
			}
			roleData := roles[0].(map[string]interface{})
			roleID := roleData["role"].(string)
			nsxRolesForPaths = append(nsxRolesForPaths, nsxModel.RolesForPath{
				Path:       &path,
				DeletePath: &boolTrue,
				Roles:      []nsxModel.Role{{Role: &roleID}},
			})
		}
	}
	return nsxRolesForPaths
}

func getRoleBindingObject(d *schema.ResourceData) *nsxModel.RoleBinding {
	boolTrue := true
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	name := d.Get("name").(string)
	identitySrcID := d.Get("identity_source_id").(string)
	identitySrcType := d.Get("identity_source_type").(string)
	roleBindingType := d.Get("type").(string)
	nsxRolesForPaths := getRolesForPathList(d)

	obj := nsxModel.RoleBinding{
		DisplayName:        &displayName,
		Description:        &description,
		Tags:               tags,
		Name:               &name,
		IdentitySourceId:   &identitySrcID,
		IdentitySourceType: &identitySrcType,
		Type_:              &roleBindingType,
		ReadRolesForPaths:  &boolTrue,
		RolesForPaths:      nsxRolesForPaths,
	}
	return &obj
}

func resourceNsxtPolicyUserManagementRoleBindingCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	roleBindingType := d.Get("type").(string)
	username := d.Get("name").(string)
	if roleBindingType == nsxModel.RoleBinding_TYPE_LOCAL_USER {
		return fmt.Errorf("creation of RoleBinding for %s is not allowed as it's created by NSX. "+
			"import the binding first", roleBindingType)
	}

	// Create the resource using POST
	log.Printf("[INFO] Creating RoleBinding for %s %s", roleBindingType, username)
	obj := getRoleBindingObject(d)
	rbClient := aaa.NewRoleBindingsClient(connector)
	rbObj, err := rbClient.Create(*obj)
	if err != nil {
		return handleCreateError("RoleBinding", username, err)
	}

	id := *rbObj.Id
	d.SetId(id)

	return resourceNsxtPolicyUserManagementRoleBindingRead(d, m)
}

func resourceNsxtPolicyUserManagementRoleBindingRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	rbClient := aaa.NewRoleBindingsClient(connector)

	id := d.Id()
	if id == "" {
		err := fmt.Errorf("error obtaining RoleBinding ID")
		return err
	}

	obj, err := rbClient.Get(id, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleReadError(d, "RoleBinding", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("revision", obj.Revision)
	d.Set("name", obj.Name)
	d.Set("identity_source_id", obj.IdentitySourceId)
	d.Set("identity_source_type", obj.IdentitySourceType)
	d.Set("type", obj.Type_)
	if obj.UserId != nil {
		d.Set("user_id", obj.UserId)
	}
	setRolesForPathInSchema(d, obj.RolesForPaths)

	return nil
}

func resourceNsxtPolicyUserManagementRoleBindingUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		err := fmt.Errorf("error obtaining RoleBinding ID")
		return err
	}

	log.Printf("[INFO] Updateing RoleBinding with ID %s", id)
	connector := getPolicyConnector(m)
	rbClient := aaa.NewRoleBindingsClient(connector)
	obj := getRoleBindingObject(d)
	_, err := rbClient.Update(id, *obj)
	if err != nil {
		return handleCreateError("RoleBinding", id, err)
	}

	return resourceNsxtPolicyUserManagementRoleBindingRead(d, m)
}

func resourceNsxtPolicyUserManagementRoleBindingDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	rbClient := aaa.NewRoleBindingsClient(connector)

	id := d.Id()
	if id == "" {
		err := fmt.Errorf("error obtaining RoleBinding ID")
		return err
	}
	roleBindingType := d.Get("type").(string)
	if roleBindingType == nsxModel.RoleBinding_TYPE_LOCAL_USER {
		return fmt.Errorf("role binding for %s %s can not be deleted", roleBindingType, id)
	}

	log.Printf("[INFO] Deleting RoleBinding with ID %s", id)
	err := rbClient.Delete(id, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleDeleteError("RoleBinding", id, err)
	}

	return nil
}
