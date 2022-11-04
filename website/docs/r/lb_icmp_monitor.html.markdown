---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_icmp_monitor"
description: |-
  Provides a resource to configure load balancer icmp monitor on NSX-T manager
---

# nsxt_lb_icmp_monitor

Provides a resource to configure lb icmp monitor on NSX-T manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_lb_icmp_monitor" "lb_icmp_monitor" {
  description  = "lb_icmp_monitor provisioned by Terraform"
  display_name = "lb_icmp_monitor"
  fall_count   = 3
  interval     = 5
  monitor_port = 7887
  rise_count   = 3
  timeout      = 10
  data_length  = 56

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb icmp monitor.
* `fall_count` - (Optional) Number of consecutive checks must fail before marking it down.
* `interval` - (Optional) The frequency at which the system issues the monitor check (in seconds).
* `monitor_port` - (Optional) If the monitor port is specified, it would override pool member port setting for healthcheck. Port range is not supported.
* `rise_count` - (Optional) Number of consecutive checks must pass before marking it up.
* `timeout` - (Optional) Number of seconds the target has in which to respond to the monitor request.
* `data_length` - (Optional) The data size (in bytes) of the ICMP healthcheck packet.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb_icmp_monitor.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb icmp monitor can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_icmp_monitor.lb_icmp_monitor UUID
```

The above would import the lb icmp monitor named `lb_icmp_monitor` with the nsx id `UUID`
