---
subcategory: "User Management"
page_title: "NSXT: nsxt_policy_user_management_role_binding_path"
description: A resource to add a single role-for-path to an existing user management role binding.
---

# nsxt_policy_user_management_role_binding_path

This resource adds **one** `roles_for_path` entry (a path and its roles) to an **existing** role binding. The role binding is typically looked up with the `nsxt_policy_user_management_role_binding` data source (binding created in the NSX portal or elsewhere). Use this resource when you want to attach an additional path/roles to a user without managing the full binding in Terraform.

~> **NOTE:** This resource requires NSX version 4.0.0 or higher.

~> **NOTE:** If the same role binding is fully managed by `nsxt_policy_user_management_role_binding`, that resource's next apply will overwrite the binding and remove any path added only by this resource. The common use case is to use the **data source** to reference an existing binding and add paths with this resource.

## Example Usage

Add a path to an existing role binding (look up binding with the data source):

```hcl
data "nsxt_policy_user_management_role_binding" "existing" {
  name = "johndoe@example.com"
  type = "remote_user"
}

resource "nsxt_policy_user_management_role_binding_path" "project_scope" {
  role_binding_id = data.nsxt_policy_user_management_role_binding.existing.id
  path            = "/orgs/default/projects/my-project"
  roles           = ["project_admin"]
}
```

With identity source when multiple sources exist:

```hcl
data "nsxt_policy_user_management_role_binding" "existing" {
  name               = "johndoe@example.com"
  type               = "remote_user"
  identity_source_id = nsxt_policy_ldap_identity_source.ldap.id
}

resource "nsxt_policy_user_management_role_binding_path" "extra" {
  role_binding_id = data.nsxt_policy_user_management_role_binding.existing.id
  path            = "/orgs/default"
  roles           = ["network_engineer", "security_engineer"]
}
```

Add multiple paths to the same binding (multiple resources):

```hcl
data "nsxt_policy_user_management_role_binding" "user" {
  name = "johndoe@example.com"
  type = "remote_user"
}

resource "nsxt_policy_user_management_role_binding_path" "path1" {
  role_binding_id = data.nsxt_policy_user_management_role_binding.user.id
  path            = "/orgs/default"
  roles           = ["vpc_admin"]
}

resource "nsxt_policy_user_management_role_binding_path" "path2" {
  role_binding_id = data.nsxt_policy_user_management_role_binding.user.id
  path            = "/orgs/default/projects/proj1"
  roles           = ["project_admin"]
}
```

## Argument Reference

* `role_binding_id` - (Required) ID of the role binding to add this path to. Use the `id` from the `nsxt_policy_user_management_role_binding` data source (or from the resource if the binding is managed in Terraform).
* `path` - (Required) Path in the parent hierarchy for which these roles apply (e.g. `"/"` or `"/orgs/default"`).
* `roles` - (Required) Set of role identifiers to assign for this path.

## Attributes Reference

In addition to arguments above, no additional attributes are exported. The resource state reflects the configured `role_binding_id`, `path`, and `roles`.

## Importing

Import by the composite id: `ROLE_BINDING_ID::PATH` (the path as stored on the binding).

```shell
terraform import nsxt_policy_user_management_role_binding_path.project_scope ROLE_BINDING_ID::/orgs/default/projects/my-project
```

After import, run `terraform plan`; Terraform will read the current roles for that path from NSX and update state if needed.
