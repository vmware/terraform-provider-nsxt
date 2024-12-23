---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: vpc_connectivity_profile"
description: VPC Connectivity Profile data source.
---

# nsxt_vpc_connectivity_profile

This data source provides information about an inventory Connectivity Profile configured under VPC on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc_connectivity_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Connectivity Profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Connectivity Profile to retrieve.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
