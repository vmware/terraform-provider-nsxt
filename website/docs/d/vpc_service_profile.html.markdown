---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: vpc_service_profile"
description: VPC Service Profile data source.
---

# nsxt_vpc_service_profile

This data source provides information about an inventory Service Profile configured under VPC on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc_service_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Service Profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Service Profile to retrieve.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
