---
layout: "nsxt"
page_title: "NSXT: nsxt_logical_router_downlink_port"
sidebar_current: "docs-nsxt-resource-logical-router_downlink-port"
description: |-
  Provides a resource to configure Logical Router Downlink Port on NSX-T Manager.
---

# nsxt_logical_router_downlink_port

Provides a resource to configure Logical Router Downlink Port on NSX-T Manager.

## Example Usage

```hcl
resource "nsxt_logical_router_downlink_port" "DP1" {
    description = "DP1 provisioned by Terraform"
    display_name = "DP1"
    logical_router_id = "${nsxt_logical_router.RTR1.id}"
    linked_logical_switch_port_id = "${nsxt_logical_port.LP1.id}"
    subnet {
      ip_addresses = ["1.1.1.0"]
      prefix_length = "24"
    }
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

* `logical_router_id` - (Required) Identifier for logical Tier-1 router on which this port is created
* `linked_logical_switch_port_id` - (Required) Identifier for port on logical switch to connect to
* `subnets` - (Required) Logical router port subnets
* `mac_address` - (Optional) Mac Address
* `urpf_mode` - (Optional) Unicast Reverse Path Forwarding mode
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical switch.
* `service_binding` - (Optional) A list of services for this port

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical switch.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
