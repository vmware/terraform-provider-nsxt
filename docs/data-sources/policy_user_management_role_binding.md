---
subcategory: "User Management"
page_title: "NSXT: nsxt_policy_user_management_role_binding"
description: A data source to read user management role bindings. Use this data source to look up role bindings managed by the nsxt_policy_user_management_role_binding resource or created outside Terraform.
---

# nsxt_policy_user_management_role_binding

This data source provides information about a user management role binding configured on NSX. It corresponds to the `nsxt_policy_user_management_role_binding` resource.

~> **NOTE:** This data source requires NSX version 4.0.0 or higher.

## Example Usage - Look up by ID

```hcl
data "nsxt_policy_user_management_role_binding" "rb" {
  id = "ROLE_BINDING_ID"
}
```

## Example Usage - Look up by name and type

```hcl
data "nsxt_policy_user_management_role_binding" "rb" {
  name = "johndoe@example.com"
  type = "remote_user"
}
```

## Example Usage - With identity source (remote user/group)

When multiple identity sources exist, specify `identity_source_id` to disambiguate:

```hcl
data "nsxt_policy_user_management_role_binding" "rb" {
  name               = "johndoe@example.com"
  type               = "remote_user"
  identity_source_id = nsxt_policy_ldap_identity_source.ldap.id
}
```

## Example Usage - Look up a role binding created by the resource

```hcl
resource "nsxt_policy_user_management_role_binding" "test" {
  display_name         = "johndoe@example.com"
  name                 = "johndoe@example.com"
  type                 = "remote_user"
  identity_source_type = "LDAP"

  roles_for_path {
    path  = "/"
    roles = ["auditor"]
  }
  roles_for_path {
    path  = "/orgs/default"
    roles = ["vpc_admin"]
  }
}

# Look up by ID
data "nsxt_policy_user_management_role_binding" "by_id" {
  id = nsxt_policy_user_management_role_binding.test.id
}

# Look up by name and type
data "nsxt_policy_user_management_role_binding" "by_name_type" {
  name               = nsxt_policy_user_management_role_binding.test.name
  type               = nsxt_policy_user_management_role_binding.test.type
  identity_source_id = nsxt_policy_user_management_role_binding.test.identity_source_id
}
```

## Argument Reference

* `id` - (Optional) The ID of the role binding to retrieve.
* `name` - (Optional) User/Group's name. Required when `id` is not set.
* `type` - (Optional) Type of the user. Required when `id` is not set. Valid options: `local_user`, `remote_user`, `remote_group`.
* `identity_source_id` - (Optional) ID of the external identity source. Used to disambiguate when looking up by name and type.

Either `id` or both `name` and `type` must be set.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `display_name` - Display name of the resource.
* `description` - Description of the resource.
* `revision` - Current revision number of the object as seen by the NSX API.
* `tag` - Tags associated with the resource.
* `user_id` - Local user's numeric id (only for `local_user`).
* `identity_source_type` - Type of the external identity source (for remote user/group).
* `roles_for_path` - List of roles associated with the user per path.
    * `path` - Path in the parent hierarchy.
    * `roles` - List of role identifiers for that path.
