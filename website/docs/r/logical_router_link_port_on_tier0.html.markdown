---
layout: "nsxt"
page_title: "NSXT: nsxt_logical_router_link_port_on_tier0"
sidebar_current: "docs-nsxt-resource-logical-router-link-port-on-tier0"
description: |-
  Provides a resource to configure Logical Router link Port on Tier-0 Router on NSX-T Manager.
---

# nsxt_logical_router_link_port_on_tier0

Provides a resource to configure Logical Router link Port on Tier-0 Router on NSX-T Manager.

## Example Usage

```hcl
resource "nsxt_logical_router_link_port_on_tier0" "link_port_tier0" {
    description = "TIER0_PORT1 provisioned by Terraform"
    display_name = "TIER0_PORT1"
    logical_router_id =  "${data.nsxt_logical_tier0_router.rtr1.id}"
    service_binding {
        target_id = "${nsxt_dhcp_relay_service.dr_service.id}"
        target_type = "LogicalService"
    }
    tag {
        scope = "color"
        tag = "blue"
    }
}
```

## Argument Reference

The following arguments are supported:

* `logical_router_id` - (Required) Identifier for logical Tier0 router on which this port is created.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical switch.
* `service_binding` - (Optional) A list of services for this port

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical switch.
* `linked_logical_switch_port_id` - Identifier for port on logical router to connect to.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing logical router link port on Tier-0 can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import nsxt_logical_router_link_port_on_tier0.x id
```

The above would import the logical router link port on Tier-0 named `x` with the nsx id `id`
