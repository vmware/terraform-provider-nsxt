---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: policy_distributed_flood_protection_profile_binding"
description: A resource to configure Policy Distributed Flood Protection Profile BindingMap on NSX Policy manager.
---

# nsxt_policy_distributed_flood_protection_profile_binding

This resource provides a method for the management of a Distributed Flood Protection Profile BindingMap.

This resource is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_distributed_flood_protection_profile_binding" "test" {
  display_name    = "test"
  description     = "test"
  profile_path    = nsxt_policy_distributed_flood_protection_profile.test.path
  group_path      = nsxt_policy_group.test.path
  sequence_number = 3

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_distributed_flood_protection_profile_binding" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name    = "test"
  description     = "test"
  profile_path    = nsxt_policy_distributed_flood_protection_profile.test.path
  group_path      = nsxt_policy_group.test.path
  sequence_number = 3

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `profile_path` - (Required) The path of the flood protection profile to be binded.
* `group_path` - (Required) The path of the group to bind with the profile.
* `sequence_number` - (Optional) Sequence number of this profile binding map.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing Distributed Flood Protection Profile BindingMap can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_distributed_flood_protection_profile_binding.dfppb POLICY_PATH
```
The above command imports the Distributed Flood Protection Profile BindingMap named `dfppb` with the policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
