---
subcategory: "User Management"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_user_management_role_binding"
description: A resource to configure user management Role Bindings.
---

# nsxt_policy_user_management_role_binding

This resource provides a method for the management of Role Bindings of users.

~> **NOTE:** This resource requires NSX version 4.0.0 or higher.

## Example Usage

```hcl
resource "nsxt_policy_user_management_role_binding" "test" {
  display_name         = "johndoe@abc.com"
  name                 = "johndoe@abc.com"
  type                 = "remote_user"
  identity_source_type = "LDAP"

  roles_for_path {
    path  = "/"
    roles = ["auditor"]
  }

  roles_for_path {
    path  = "/orgs/default"
    roles = ["org_admin", "vpc_admin"]
  }
}
```
As nsxt_policy_user_management_role_binding instances apply to nsxt_node_user and nsxt_policy_user_management_role resources, when they are created in the same Terraform configuration
users need to specify resource dependencies using the `depends_on` clause as in the following example:

```
resource "nsxt_node_user" "node_user" {
  active                    = true
  full_name                 = "John Doe"
  password                  = "Str0ng_Pwd!Wins$"
  username                  = "johndoe123"
  password_change_frequency = 90
  password_change_warning   = 30
}

resource "nsxt_policy_user_management_role_binding" "user_management_role_binding" {
  display_name         = "johndoe123"
  name                 = "johndoe123"
  type                 = "local_user"
  roles_for_path {
    path  = "/"
    roles = ["auditor"]
  }
  overwrite_local_user = true

  depends_on = [
    nsxt_node_user.node_user
  ]
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
    * `local_user` - This is a user local to NSX. These are linux users. Note: Role bindings for local users are owned by NSX. Creation and deletion is not allowed for local users' binding. For updates, import existing bindings first. Alternatively, set `overwrite_local_user` to overwrite current role bindings with the one defined in terraform.
* `identity_source_type` - (Optional) Identity source type. Applicable only to `remote_user` and `remote_group` user types. Valid options are: `VIDM`, `LDAP`, `OIDC`, `CSP`. Defaults to `VIDM` when applicable.
* `identity_source_id` - (Optional) The ID of the external identity source that holds the referenced external entity. Currently, only external `LDAP` and `OIDC` servers are allowed.
* `roles_for_path` - (Required) A list of The roles that are associated with the user, limiting them to a path. In case the path is '/', the roles apply everywhere.
    * `path` - (Required) Path of the entity in parent hierarchy.
    * `roles` - (Required) A list of identifiers for the roles to associate with the given user limited to a path.
* `overwrite_local_user` - (Optional) Flag to allow overwriting existing role bindings for local user with terraform definition. On deletion, the user's role will be reverted to auditor only. Any existing configuration will be lost.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `user_id` - Local user's numeric id. Only applicable to `local_user`.

# Build-in NSX roles

`roles_for_path.roles` accepts user created roles as well as roles native to NSX. For reference, the following is a list of native roles as of NSX 4.1.2
- `network_engineer`: Network Admin
- `support_bundle_collector`: Support Bundle Collector
- `security_op`: Security Operator
- `lb_auditor`: LB Operator
- `netx_partner_admin`: NETX Partner Admin
- `project_admin`: Project Admin
- `auditor`: Auditor
- `network_op`: Network Operator
- `enterprise_admin`: Enterprise Admin
- `lb_admin`: LB Admin
- `gi_partner_admin`: GI Partner Admin
- `vpn_admin`: VPN Admin
- `vpc_admin`: VPC Admin
- `security_engineer`: Security Admin

The permission matrix for above roles is available on [NSX documentation](https://docs.vmware.com/en/VMware-NSX/4.1/administration/GUID-26C44DE8-1854-4B06-B6DA-A2FD426CDF44.html)

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_user_management_role_binding.test ROLE_BINDING_ID
```
The above command imports Role named `test` with the role binding identifier `ROLE_BINDING_ID`.
