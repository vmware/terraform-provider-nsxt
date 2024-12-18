---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: vpc_subnet"
description: VPC Subnet data source.
---

# nsxt_vpc_subnet

This data source provides information about a Subnet configured under VPC on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc" "demovpc" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "vpc1"
}

data "nsxt_vpc_subnet" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  display_name = "subnet1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Subnet to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Subnet to retrieve.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
    * `vpc_id` - (Required) The ID of the VPC which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
