---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_static_route"
description: A resource to configure a static route in NSX.
---

# nsxt_static_route

This resource provides a means to configure static routes in NSX to determine where IP traffic is routed.

## Example Usage

```hcl
resource "nsxt_static_route" "static_route" {
  description       = "SR provisioned by Terraform"
  display_name      = "SR"
  logical_router_id = nsxt_logical_tier1_router.router1.id
  network           = "4.4.4.0/24"

  next_hop {
    ip_address              = "8.0.0.10"
    administrative_distance = "1"
    logical_router_port_id  = nsxt_logical_router_downlink_port.downlink_port.id
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
* `logical_router_id` - (Required) Logical router id.
* `network` - (Required) CIDR.
* `next_hop` - (Required) List of Next Hops, each with those arguments:
    * `administrative_distance` - (Optional) Administrative Distance for the next hop IP.
    * `ip_address` - (Optional) Next Hop IP.
    * `logical_router_port_id` - (Optional) Reference of logical router port to be used for next hop.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the static route.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `next_hop` additional arguments:
    * `bfd_enabled` - Status of bfd for this next hop where bfd_enabled = true indicate bfd is enabled for this next hop and bfd_enabled = false indicate bfd peer is disabled or not configured for this next hop.
    * `blackhole_action` - Action to be taken on matching packets for NULL routes.

## Importing

An existing static route can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_static_route.static_route logical-router-uuid/static-route-num
```

The above command imports the static route named `static_route` with the number `static-route-num` that belongs to the tier 1 logical router with the NSX id `logical-router-uuid`.
