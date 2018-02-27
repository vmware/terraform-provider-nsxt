---
layout: "nsxt"
page_title: "NSXT: nsxt_static_route"
sidebar_current: "docs-nsxt-resource-static-route"
description: |-
  Provides a resource to configure static route on NSX-T manager
---

# nsxt_static_route

Provides a resource to configure static route on NSX-T manager

## Example Usage

```hcl
resource "nsxt_static_route" "static_route" {
  description       = "SR provisioned by Terraform"
  display_name      = "SR"
  logical_router_id = "${nsxt_logical_tier1_router.router1.id}"
  network           = "4.4.4.0/24"

  next_hop {
    ip_address              = "8.0.0.10"
    administrative_distance = "1" 
    logical_router_port_id  = "${nsxt_logical_router_downlink_port.downlink_port.id}"
  }

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this static route.
* `logical_router_id` - (Optional) Logical router id.
* `network` - (Required) CIDR.
* `next_hop` - (Required) List of Next Hops, each with those arguments:
    * `administrative_distance` - (Optional) Administrative Distance for the next hop IP.
    * `ip_address` - (Optional) Next Hop IP.
    * `logical_router_port_id` - (Optional) Reference of logical router port to be used for next hop.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the static_route.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `next_hop` additional arguments:
    * `bfd_enabled` - Status of bfd for this next hop where bfd_enabled = true indicate bfd is enabled for this next hop and bfd_enabled = false indicate bfd peer is disabled or not configured for this next hop.
    * `blackhole_action` - Action to be taken on matching packets for NULL routes. 