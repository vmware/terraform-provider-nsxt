---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: policy_gateway_flood_protection_profile_binding"
description: A resource to configure Policy Gateway Flood Protection Profile BindingMap on NSX Policy manager.
---

# nsxt_policy_gateway_flood_protection_profile_binding

This resource provides a method for the management of a Gateway Flood Protection Profile BindingMap.

This resource is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "test" {
  display_name = "tier0_gw"
}

resource "nsxt_policy_gateway_flood_protection_profile_binding" "test" {
  display_name = "test"
  description  = "test"
  profile_path = nsxt_policy_gateway_flood_protection_profile.test.path
  parent_path  = data.nsxt_policy_tier0_gateway.test.path

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

data "nsxt_policy_tier1_gateway" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "tier1_gw"
}

resource "nsxt_policy_gateway_flood_protection_profile_binding" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "test"
  description  = "test"
  profile_path = nsxt_policy_gateway_flood_protection_profile.test.path
  parent_path  = data.nsxt_policy_tier1_gateway.test.path

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
* `parent_path` - (Required) The path of the parent to bind with the profile. This could be either T0 path, T1 path or locale service path.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing Gateway Flood Protection Profile BindingMap can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_gateway_flood_protection_profile_binding.gfppb POLICY_PATH
```
The above command imports the Gateway Flood Protection Profile BindingMap named `gfppb` with the policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
