---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_logical_router_link_port_on_tier1"
description: A resource to configure a logical router link port on a Tier-1 router in NSX.
---

# nsxt_logical_router_link_port_on_tier1

This resource provides the ability to configure a logical router link port on a tier 1 logical router. This port can then be used to connect the tier 1 logical router to another logical router.

## Example Usage

```hcl
resource "nsxt_logical_router_link_port_on_tier1" "link_port_tier1" {
  description                   = "TIER1_PORT1 provisioned by Terraform"
  display_name                  = "TIER1_PORT1"
  logical_router_id             = nsxt_logical_tier1_router.rtr1.id
  linked_logical_router_port_id = nsxt_logical_router_link_port_on_tier0.link_port_tier0.id

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `logical_router_id` - (Required) Identifier for logical tier-1 router on which this port is created.
* `linked_logical_switch_port_id` - (Required) Identifier for port on logical Tier-0 router to connect to.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this port.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical router link port.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing logical router link port on Tier-1 can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_logical_router_link_port_on_tier1.link_port_tier1 UUID
```

The above command imports the logical router link port on the tier 1 router named `link_port_tier1` with the NSX id `UUID`.
