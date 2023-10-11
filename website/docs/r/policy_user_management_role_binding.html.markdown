---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_user_management_role_binding"
description: A resource to configure user management Role Bindings.
---

# nsxt_policy_user_management_role_binding

This resource provides a method for the management of Role Bindings of users.

## Example Usage

```hcl
resource "nsxt_policy_user_management_role_binding" "test" {
  display_name         = "johndoe@abc.com"
  name                 = "johndoe@abc.com"
  type                 = "remote_user"
  identity_source_type = "LDAP"

  roles_for_path {
    path = "/"
    role {
      role = "auditor"
    }
  }

  roles_for_path {
    path = "/orgs/default"
    role {
      role = "org_admin"
    }
    role {
      role = "vpc_admin"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `name` - (Required) User/Group's name.
* `type` - (Required) Indicates the type of the user. Valid options:
    * `remote_user` - This is a user which is external to NSX. 
    * `remote_group` - This is a group of users which is external to NSX.
    * `local_user` - This is a user local to NSX. These are linux users. Note: Role bindings for local users are owned by NSX. Creation and deletion is not allowed for local users' binding. For updates, import existing bindings first.
    * `principal_identity` - This is a principal identity user.
* `identity_source_type` - (Optional) Identity source type. Applicable only to `remote_user` and `remote_group` user types. Valid options are: `VIDM`, `LDAP`, `OIDC`, `CSP`. Defaults to `VIDM` when applicable.
* `identity_source_id` - (Optional) The ID of the external identity source that holds the referenced external entity. Currently, only external `LDAP` and `OIDC` servers are allowed.
* `roles_for_path` - (Required) A list of The roles that are associated with the user, limiting them to a path. In case the path is '/', the roles apply everywhere.
    * `path` - (Required) Path of the entity in parent hierarchy.
    * `role` - (Required) A list of identifiers for the roles to associate with the given user limited to a path.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `user_id` - Local user's numeric id. Only applicable to `local_user`.


## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_user_management_role_binding.test ROLE_BINDING_ID
```
The above command imports Role named `test` with the role binding identifier `ROLE_BINDING_ID`.
