---
layout: "nsxt"
page_title: "NSXT: nsxt_logical_router_link_port_on_tier1"
sidebar_current: "docs-nsxt-resource-logical-router-link-port-on-tier1"
description: |-
  Provides a resource to configure Logical Router link Port on Tier-1 Router on NSX-T Manager.
---

# nsxt_logical_router_link_port_on_tier1

Provides a resource to configure Logical Router link Port on Tier-1 Router on NSX-T Manager.

## Example Usage

```hcl
resource "nsxt_logical_router_link_port_on_tier1" "TIER1_PORT1" {
    description = "TIER1_PORT1 provisioned by Terraform"
    display_name = "TIER1_PORT1"
    logical_router_id =  "${nsxt_logical_tier1_router.RTR1.id}"
    linked_logical_router_port_id = "${nsxt_logical_router_link_port_on_tier0.TIER0_PORT1.id}"
    service_binding {
        target_id = "${nsxt_dhcp_relay_service.DRS1.id}"
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

* `logical_router_id` - (Required) Identifier for logical tier-1 router on which this port is created.
* `linked_logical_switch_port_id` - (Required) Identifier for port on logical Tier-0 router to connect to.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical switch.
* `service_binding` - (Optional) A list of services for this port

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical switch.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
