---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_passive_monitor"
description: |-
  Provides a resource to configure load balancer passive monitor on NSX-T manager
---

# nsxt_lb_passive_monitor

Provides a resource to configure lb passive monitor on NSX-T manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_lb_passive_monitor" "lb_passive_monitor" {
  description  = "lb_passive_monitor provisioned by Terraform"
  display_name = "lb_passive_monitor"
  max_fails    = 3
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
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb passive monitor.
* `max_fails` - (Optional) When consecutive failures reach this value, the member is considered temporarily unavailable for a configurable period.
* `timeout` - (Optional) After this timeout period, the member is probed again.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb_passive_monitor.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb passive monitor can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_passive_monitor.lb_passive_monitor UUID
```

The above would import the lb passive monitor named `lb_passive_monitor` with the nsx id `UUID`
