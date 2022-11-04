---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_tcp_monitor"
description: |-
  Provides a resource to configure load balancer tcp monitor on NSX-T manager
---

# nsxt_lb_tcp_monitor

Provides a resource to configure lb tcp monitor on NSX-T manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_lb_tcp_monitor" "lb_tcp_monitor" {
  description  = "lb_tcp_monitor provisioned by Terraform"
  display_name = "lb_tcp_monitor"
  fall_count   = 3
  interval     = 5
  monitor_port = 7887
  rise_count   = 3
  timeout      = 10

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
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb tcp monitor.
* `fall_count` - (Optional) Number of consecutive checks must fail before marking it down.
* `interval` - (Optional) The frequency at which the system issues the monitor check (in seconds).
* `monitor_port` - (Optional) If the monitor port is specified, it would override pool member port setting for healthcheck. Port range is not supported.
* `rise_count` - (Optional) Number of consecutive checks must pass before marking it up.
* `timeout` - (Optional) Number of seconds the target has in which to respond to the monitor request.
* `receive` - (Optional) Expected data, if specified, can be anywhere in the response and it has to be a string, regular expressions are not supported.
* `send` - (Optional) Payload to send out to the monitored server. If both send and receive are not specified, then just a TCP connection is established (3-way handshake) to validate server is healthy, no data is sent.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb_tcp_monitor.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb tcp monitor can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_tcp_monitor.lb_tcp_monitor UUID
```

The above would import the lb tcp monitor named `lb_tcp_monitor` with the nsx id `UUID`
