---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_udp_monitor"
description: |-
  Provides a resource to configure load balancer udp monitor on NSX-T manager
---

# nsxt_lb_udp_monitor

Provides a resource to configure lb udp monitor on NSX-T manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_lb_udp_monitor" "lb_udp_monitor" {
  description  = "lb_udp_monitor provisioned by Terraform"
  display_name = "lb_udp_monitor"
  fall_count   = 3
  interval     = 5
  monitor_port = 7887
  rise_count   = 3
  timeout      = 10
  send         = "hi"
  receive      = "hello"

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
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb udp monitor.
* `fall_count` - (Optional) Number of consecutive checks must fail before marking it down.
* `interval` - (Optional) The frequency at which the system issues the monitor check (in seconds).
* `monitor_port` - (Optional) If the monitor port is specified, it would override pool member port setting for healthcheck. Port range is not supported.
* `rise_count` - (Optional) Number of consecutive checks must pass before marking it up.
* `timeout` - (Optional) Number of seconds the target has in which to respond to the monitor request.
* `receive` - (Required) Expected data, if specified, can be anywhere in the response and it has to be a string, regular expressions are not supported.
* `send` - (Required) Payload to send out to the monitored server.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb_udp_monitor.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb udp monitor can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_udp_monitor.lb_udp_monitor UUID
```

The above would import the lb udp monitor named `lb_udp_monitor` with the nsx id `UUID`
