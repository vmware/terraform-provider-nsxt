---
layout: "nsxt"
page_title: "NSXT: nsxt_dhcp_relay_service"
sidebar_current: "docs-nsxt-resource-dhcp-relay-service"
description: |-
  Provides a resource to configure DHCP relay service on NSX-T manager
---

# nsxt_dhcp_relay_service

Provides a resource to configure DHCP relay service on NSX-T manager
The DHCP relay service uses a DHCP relay profile, and later consumed by a router
link port.

## Example Usage

```hcl
resource "nsxt_dhcp_relay_profile" "DRP" {
    display_name = "DRP"
    server_addresses = ["1.1.1.1"]
}

resource "nsxt_dhcp_relay_service" "DRS" {
    description = "DRS provisioned by Terraform"
    display_name = "DRS"
    tag {
        scope = "color"
        tag = "blue"
    }
    dhcp_relay_profile_id = "${nsxt_dhcp_relay_profile.DRP.id}"
}

resource "nsxt_logical_router_downlink_port" "LRDP" {
    display_name = "logical_router_downlink_port"
    linked_logical_switch_port_id = "${nsxt_logical_port.PORT1.id}"
    logical_router_id = "${nsxt_logical_tier1_router.RTR1.id}"
    subnet {
        ip_addresses = ["8.0.0.1"],
        prefix_length = 24
    }
    service_binding {
        target_id = "${nsxt_dhcp_relay_service.DRS.id}"
        target_type = "LogicalService"
    }
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this dhcp_relay_service.
* `dhcp_relay_profile_id` - (Required) dhcp relay profile referenced by the dhcp relay service.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the dhcp_relay_service.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
