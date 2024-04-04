---
subcategory: "Security"
layout: "nsxt"
page_title: "NSXT: policy_distributed_flood_protection_profile"
description: Policy Distributed Flood Protection Profile data source.
---

# nsxt_policy_distributed_flood_protection_profile

This data source provides information about policy Distributed Flood Protection Profile configured in NSX.
This data source is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_distributed_flood_protection_profile" "test" {
  display_name = "test"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_distributed_flood_protection_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "test"
}
```

## Argument Reference

* `id` - (Optional) The ID of Distributed Flood Protection Profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Distributed Flood Protection Profile to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
