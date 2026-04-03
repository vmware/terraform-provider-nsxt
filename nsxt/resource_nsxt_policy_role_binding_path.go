// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

// roleBindingPathIDSeparator is used in the resource id: binding_id + separator + path
const roleBindingPathIDSeparator = "::"

func resourceNsxtPolicyUserManagementRoleBindingPath() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyUserManagementRoleBindingPathCreate,
		Read:   resourceNsxtPolicyUserManagementRoleBindingPathRead,
		Update: resourceNsxtPolicyUserManagementRoleBindingPathUpdate,
		Delete: resourceNsxtPolicyUserManagementRoleBindingPathDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role_binding_id": {
				Type:        schema.TypeString,
				Description: "ID of the role binding to add this path to. Use the id from nsxt_policy_user_management_role_binding or the data source.",
				Required:    true,
				ForceNew:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "Path in the parent hierarchy for which these roles apply (e.g. \"/\" or \"/orgs/default\").",
				Required:    true,
				ForceNew:    true,
			},
			"roles": {
				Type:        schema.TypeSet,
				Description: "List of role identifiers to assign for this path.",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func parseRoleBindingPathID(id string) (bindingID, path string, err error) {
	idx := strings.Index(id, roleBindingPathIDSeparator)
	if idx < 0 {
		return "", "", fmt.Errorf("invalid role binding path id: %q", id)
	}
	return id[:idx], id[idx+len(roleBindingPathIDSeparator):], nil
}

func roleBindingPathResourceID(bindingID, path string) string {
	return bindingID + roleBindingPathIDSeparator + path
}

// mergeRolesForPathWithPath merges the existing binding's RolesForPaths with the given path and roles.
// If the path already exists, its roles are replaced with newRoles. Otherwise the path is appended.
func mergeRolesForPathWithPath(existing []nsxModel.RolesForPath, path string, newRoles []string) []nsxModel.RolesForPath {
	nsxRoles := make([]nsxModel.Role, 0, len(newRoles))
	for _, r := range newRoles {
		role := r
		nsxRoles = append(nsxRoles, nsxModel.Role{Role: &role})
	}
	pathPtr := strings.Clone(path)
	newEntry := nsxModel.RolesForPath{
		Path:  &pathPtr,
		Roles: nsxRoles,
	}

	out := make([]nsxModel.RolesForPath, 0, len(existing)+1)
	found := false
	for _, p := range existing {
		if p.Path != nil && *p.Path == path {
			out = append(out, newEntry)
			found = true
			continue
		}
		out = append(out, p)
	}
	if !found {
		out = append(out, newEntry)
	}
	return out
}

// rolesForPathWithoutPath returns RolesForPath list with the given path removed.
// The path to remove is sent with DeletePath: true and one role (required by API).
func rolesForPathWithoutPath(existing []nsxModel.RolesForPath, pathToRemove string, oneRoleForDelete string) []nsxModel.RolesForPath {
	boolTrue := true
	out := make([]nsxModel.RolesForPath, 0, len(existing))
	for _, p := range existing {
		if p.Path != nil && *p.Path == pathToRemove {
			path := strings.Clone(pathToRemove)
			role := oneRoleForDelete
			out = append(out, nsxModel.RolesForPath{
				Path:       &path,
				DeletePath: &boolTrue,
				Roles:      []nsxModel.Role{{Role: &role}},
			})
			continue
		}
		out = append(out, p)
	}
	return out
}

func resourceNsxtPolicyUserManagementRoleBindingPathCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	rbClient := cliRoleBindingsClient(sessionContext, connector)

	bindingID := d.Get("role_binding_id").(string)
	path := d.Get("path").(string)
	roles := interface2StringList(d.Get("roles").(*schema.Set).List())
	if len(roles) == 0 {
		return fmt.Errorf("at least one role is required")
	}

	obj, err := rbClient.Get(bindingID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleReadError(d, "RoleBinding", bindingID, err)
	}

	merged := mergeRolesForPathWithPath(obj.RolesForPaths, path, roles)
	obj.RolesForPaths = merged
	boolTrue := true
	obj.ReadRolesForPaths = &boolTrue

	log.Printf("[INFO] Adding role_for_path to RoleBinding %s: path=%q roles=%v", bindingID, path, roles)
	_, err = rbClient.Update(bindingID, obj)
	if err != nil {
		return handleUpdateError("RoleBinding", bindingID, err)
	}

	d.SetId(roleBindingPathResourceID(bindingID, path))
	return resourceNsxtPolicyUserManagementRoleBindingPathRead(d, m)
}

func resourceNsxtPolicyUserManagementRoleBindingPathRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	rbClient := cliRoleBindingsClient(sessionContext, connector)

	bindingID, path, err := parseRoleBindingPathID(d.Id())
	if err != nil {
		return err
	}

	obj, err := rbClient.Get(bindingID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleReadError(d, "RoleBinding", bindingID, err)
	}

	for _, p := range obj.RolesForPaths {
		if p.Path != nil && *p.Path == path {
			roles := make([]string, 0, len(p.Roles))
			for _, r := range p.Roles {
				if r.Role != nil {
					roles = append(roles, *r.Role)
				}
			}
			d.Set("role_binding_id", bindingID)
			d.Set("path", path)
			if err := d.Set("roles", roles); err != nil {
				return err
			}
			return nil
		}
	}

	// Path no longer present on the binding (e.g. removed elsewhere)
	d.SetId("")
	return fmt.Errorf("role binding path %q not found on binding %s", path, bindingID)
}

func resourceNsxtPolicyUserManagementRoleBindingPathUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	rbClient := cliRoleBindingsClient(sessionContext, connector)

	bindingID, path, err := parseRoleBindingPathID(d.Id())
	if err != nil {
		return err
	}
	roles := interface2StringList(d.Get("roles").(*schema.Set).List())
	if len(roles) == 0 {
		return fmt.Errorf("at least one role is required")
	}

	obj, err := rbClient.Get(bindingID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleReadError(d, "RoleBinding", bindingID, err)
	}

	merged := mergeRolesForPathWithPath(obj.RolesForPaths, path, roles)
	obj.RolesForPaths = merged
	boolTrue := true
	obj.ReadRolesForPaths = &boolTrue

	log.Printf("[INFO] Updating role_for_path on RoleBinding %s: path=%q roles=%v", bindingID, path, roles)
	_, err = rbClient.Update(bindingID, obj)
	if err != nil {
		return handleUpdateError("RoleBinding", bindingID, err)
	}

	return resourceNsxtPolicyUserManagementRoleBindingPathRead(d, m)
}

func resourceNsxtPolicyUserManagementRoleBindingPathDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	rbClient := cliRoleBindingsClient(sessionContext, connector)

	bindingID, path, err := parseRoleBindingPathID(d.Id())
	if err != nil {
		return err
	}

	obj, err := rbClient.Get(bindingID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleReadError(d, "RoleBinding", bindingID, err)
	}

	var oneRole string
	for _, p := range obj.RolesForPaths {
		if p.Path != nil && *p.Path == path && len(p.Roles) > 0 && p.Roles[0].Role != nil {
			oneRole = *p.Roles[0].Role
			break
		}
	}
	if oneRole == "" {
		// Path already gone
		return nil
	}

	without := rolesForPathWithoutPath(obj.RolesForPaths, path, oneRole)
	obj.RolesForPaths = without
	boolTrue := true
	obj.ReadRolesForPaths = &boolTrue

	log.Printf("[INFO] Removing role_for_path from RoleBinding %s: path=%q", bindingID, path)
	_, err = rbClient.Update(bindingID, obj)
	if err != nil {
		return handleUpdateError("RoleBinding", bindingID, err)
	}

	return nil
}
