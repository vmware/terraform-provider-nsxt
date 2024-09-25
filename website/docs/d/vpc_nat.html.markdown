---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: vpc_nat"
description: VPC NAT data source.
---

# nsxt_vpc_nat

This data source provides information about an NAT section configured under VPC on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "proj" {
  display_name = "demoproj"
}

data "nsxt_vpc" "vpc1" {
  context {
    project_id = data.nsxt_policy_project.proj.id
  }
  display_name = "vpc1"
}

data "nsxt_vpc_nat" "test" {
  context {
    project_id = data.nsxt_policy_project.proj.id
    vpc_id     = data.nsxt_vpc.vpc1.id
  }
  nat_type = "USER"
}
```

## Argument Reference

* `nat_type` - (Required) Type of NAT, one of `USER`, `INTERNAL`, `DEFAULT` or `NAT64`.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
    * `vpc_id` - (Required) The ID of the VPC which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - (Optional) The ID of the resource.
* `display_name` - (Optional) Display Name of the resource.
* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
