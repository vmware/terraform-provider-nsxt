---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: vpc_ip_address_allocation"
description: VPC IP address allocation data source.
---

# nsxt_vpc_ip_address_allocation

This data source provides information about IP address allocation under VPC on NSX.

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

data "nsxt_vpc_ip_address_allocation" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  allocation_ips = "10.203.5.19"
}
```

## Argument Reference

* `id` - (Optional) The ID of Subnet to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Subnet to retrieve.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
    * `vpc_id` - (Required) The ID of the VPC which the object belongs to
* `allocation_ips` - (Optional) IP or CIDR that was allocated

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
