---
subcategory: "VPC"
page_title: "NSXT: vpc_lb_monitor"
description: VPC Load Balancer Monitor data source.
---

# nsxt_vpc_lb_monitor

This data source provides information about VPC Load Balancer Monitor configured on NSX.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

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

data "nsxt_vpc_lb_monitor" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  type         = "TCP"
  display_name = "my-tcp-monitor"
}
```

## Argument Reference

* `context` - (Required) The context which the object belongs to. Requires a VPC context.
    * `project_id` - (Required) The ID of the project which the object belongs to.
    * `vpc_id` - (Required) The ID of the VPC which the object belongs to.
* `id` - (Optional) The ID of Monitor to retrieve.
* `type` - (Optional) Type of Monitor to retrieve, one of `HTTP`, `HTTPS`, `TCP`, `UDP`, `ICMP`, `PASSIVE`, `ANY`.
* `display_name` - (Optional) The Display Name prefix of Monitor to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
