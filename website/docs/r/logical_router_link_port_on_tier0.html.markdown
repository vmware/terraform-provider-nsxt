---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_logical_router_link_port_on_tier0"
description: A resource that can be used to configure a logical router link port on Tier-0 router on NSX.
---

# nsxt_logical_router_link_port_on_tier0

This resource provides the ability to configure a logical router link port on a tier 0 logical router. This port can then be used to connect the tier 0 logical router to another logical router.

## Example Usage

```hcl
resource "nsxt_logical_router_link_port_on_tier0" "link_port_tier0" {
  description       = "TIER0_PORT1 provisioned by Terraform"
  display_name      = "TIER0_PORT1"
  logical_router_id = data.nsxt_logical_tier0_router.rtr1.id

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `logical_router_id` - (Required) Identifier for logical Tier0 router on which this port is created.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this port.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical router link port.
* `linked_logical_switch_port_id` - Identifier for port on logical router to connect to.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing logical router link port on Tier-0 can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_logical_router_link_port_on_tier0.link_port_tier0 UUID
```

The above command imports the logical router link port on the tier 0 logical router named `link_port_tier0` with the NSX id `UUID`.
