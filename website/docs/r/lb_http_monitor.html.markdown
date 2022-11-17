---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_http_monitor"
description: |-
  Provides a resource to configure lb http monitor on NSX-T manager
---

# nsxt_lb_http_monitor

Provides a resource to configure Load Balancer HTTP Monitor on NSX-T manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl

resource "nsxt_lb_http_monitor" "lb_http_monitor" {
  description           = "lb_http_monitor provisioned by Terraform"
  display_name          = "lb_http_monitor"
  fall_count            = 2
  interval              = 5
  monitor_port          = 8080
  rise_count            = 5
  timeout               = 10
  request_body          = "ping"
  request_method        = "HEAD"
  request_url           = "/index.html"
  request_version       = "HTTP_VERSION_1_1"
  response_body         = "pong"
  response_status_codes = [200, 304]

  tag {
    scope = "color"
    tag   = "red"
  }

  request_header {
    name  = "X-healthcheck"
    value = "NSX"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb http monitor.
* `fall_count` - (Optional) Number of consecutive checks that must fail before marking it down.
* `interval` - (Optional) The frequency at which the system issues the monitor check (in seconds).
* `monitor_port` - (Optional) If the monitor port is specified, it would override pool member port setting for healthcheck. A port range is not supported.
* `rise_count` - (Optional) Number of consecutive checks that must pass before marking it up.
* `timeout` - (Optional) Number of seconds the target has to respond to the monitor request.
* `request_body` - (Optional) String to send as HTTP health check request body. Valid only for certain HTTP methods like POST.
* `request_header` - (Optional) HTTP request headers.
* `request_method` - (Optional) Health check method for HTTP monitor type. Valid values are GET, HEAD, PUT, POST and OPTIONS.
* `request_url` - (Optional) URL used for HTTP monitor.
* `request_version` - (Optional) HTTP request version. Valid values are HTTP_VERSION_1_0 and HTTP_VERSION_1_1.
* `response_body` - (Optional) If response body is specified, healthcheck HTTP response body is matched against the specified string and server is considered healthy only if there is a match (regular expressions not supported). If response body string is not specified, HTTP healthcheck is considered successful if the HTTP response status code is among configured values.
* `response_status_codes` - (Optional) HTTP response status code should be a valid HTTP status code.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb_http_monitor.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb http monitor can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_http_monitor.lb_http_monitor UUID
```

The above would import the lb http monitor named `lb_http_monitor` with the nsx id `UUID`
