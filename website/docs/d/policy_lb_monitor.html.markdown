---
subcategory: "Load Balancer"
layout: "nsxt"
page_title: "NSXT: policy_lb_monitor"
description: Policy Load Balancer Monitor data source.
---

# nsxt_policy_lb_monitor

This data source provides information about Policy Load Balancer Monitor configured on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_lb_monitor" "test" {
  type         = "TCP"
  display_name = "my-tcp-monitor"
}
```

## Argument Reference

* `id` - (Optional) The ID of Monitor to retrieve.
* `type` - (Optional) Type of Monitor to retrieve, one of `HTTP`, `HTTPS`, `TCP`, `UDP`, `ICMP`, `PASSIVE`, `ANY`.
* `display_name` - (Optional) The Display Name prefix of Monitor to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
