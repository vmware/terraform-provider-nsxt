---
subcategory: "User Management"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_user_management_role"
description: A resource to configure user management Roles.
---

# nsxt_policy_user_management_role

This resource provides a method for the management of Roles.

## Example Usage

```hcl
resource "nsxt_policy_user_management_role" "test_role" {
  display_name = "test_role"
  description  = "Terraform provisioned role"
  role         = "test_role"

  feature {
    feature    = "policy_grouping"
    permission = "crud"
  }

  feature {
    feature    = "vm_vm_info"
    permission = "crud"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `role` - (Required) Short identifier for the role. Must be all lower case with no spaces.	This will also be the NSX ID of this resource.
* `feature` - (Required) A list of permissions for features to be granted with this role.
    * `feature` - (Required) The ID of feature to grant permission.
    * `permission` - (Required) Type of permission to grant. Valid values are `crud`, `read`, `execute`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `feature` - The permissions for features in the list will also have the following attributes exported:
  * `feature_description` - Description of the feature.
  * `feature_name` - Name of the feature.


## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_user_management_role.test ROLE_ID
```
The above command imports Role named `test` with the role identifier `ROLE_ID`.
