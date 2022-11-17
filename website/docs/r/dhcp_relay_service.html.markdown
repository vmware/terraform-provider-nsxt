---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_dhcp_relay_service"
description: A resource that can be used to configure a DHCP relay service on NSX.
---

# nsxt_dhcp_relay_service

This resource provides a way to configure the DHCP relay service on the NSX manager.
The DHCP relay service uses a DHCP relay profile and later consumed by a router
downlink port to provide DHCP addresses to virtual machines connected to a logical switch.
Currently the DHCP relay is not supported for logical routers link ports on Tier0 or Tier1.

## Example Usage

```hcl
resource "nsxt_dhcp_relay_profile" "dr_profile" {
  description  = "DRP provisioned by Terraform"
  display_name = "DRP"

  tag {
    scope = "color"
    tag   = "red"
  }

  server_addresses = ["1.1.1.1"]
}

resource "nsxt_dhcp_relay_service" "dr_service" {
  display_name          = "DRS"
  dhcp_relay_profile_id = nsxt_dhcp_relay_profile.dr_profile.id
}

resource "nsxt_logical_router_downlink_port" "router_downlink" {
  display_name                  = "logical_router_downlink_port"
  linked_logical_switch_port_id = nsxt_logical_port.port1.id
  logical_router_id             = nsxt_logical_tier1_router.rtr1.id

  subnet {
    ip_addresses  = ["8.0.0.1"]
    prefix_length = 24
  }

  service_binding {
    target_id   = nsxt_dhcp_relay_service.dr_service.id
    target_type = "LogicalService"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this dhcp_relay_service.
* `dhcp_relay_profile_id` - (Required) DHCP relay profile referenced by the DHCP relay service.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the DHCP relay service.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing DHCP Relay service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_dhcp_relay_service.dr_service UUID
```

The above command imports the DHCP relay service named `dr_service` with the NSX id `UUID`.
