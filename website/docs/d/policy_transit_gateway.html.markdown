---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: policy_transit_gateway"
description: Policy transit gateway data source.
---

# nsxt_policy_transit_gateway

This data source provides information about an inventory transit gateway on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "tgw1"
}
```

## Argument Reference

* `id` - (Optional) The ID of transit gateway to retrieve.
* `display_name` - (Optional) The Display Name prefix of the transit gateway to retrieve.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
