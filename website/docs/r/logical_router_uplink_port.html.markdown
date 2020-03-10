---
layout: "nsxt"
page_title: "NSXT: nsxt_logical_router_uplink_port"
sidebar_current: "docs-nsxt-resource-logical-router-uplink-port"
description: A resource that can be used to configure logical router uplink port in NSX.
---

# nsxt_logical_router_uplink_port

This resource provides a means to define a uplink port on a logical router to connect a logical tier0 router to a VLAN logical switch. The result of this is to provide a North-South connection.

## Example Usage

```hcl
resource "nsxt_logical_router_uplink_port" "uplink_port" {
  description                   = "UP1 provisioned by Terraform"
  display_name                  = "UP1"
  logical_router_id             = "${nsxt_logical_tier0_router.rtr0.id}"
  linked_logical_switch_port_id = "${nsxt_logical_port.logical_port1.id}"
  edge_cluster_member_index     = [1]

  subnets {
    ip_addresses  = ["104.103.102.101"]
    prefix_length = 24
  }

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `logical_router_id` - (Required) Identifier for logical Tier-0 router on which this port is created.
* `linked_logical_switch_port_id` - (Required) Identifier for port on VLAN logical switch to connect to.
* `urpf_mode` - (Optional) Unicast Reverse Path Forwarding mode. Accepted values are "NONE" and "STRICT" which is the default value.
* `edge_cluster_member_index` -  (Required) A list of member index of the edge node on the cluster.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `subnets` - (Required) Logical router port subnets.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this port.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical router downlink port.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `mac_address` - The MAC address assigned to this port
