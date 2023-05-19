---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_tier0_gateway_ha_vip_config"
description: A resource to configure HA Vip config on Tier-0 gateway in NSX Policy manager.
---

# nsxt_policy_tier0_gateway_ha_vip_config

This resource provides a method for the management of a Tier-0 gateway HA VIP config. Note that this configuration can be defined only for Active-Standby Tier0 gateway.

~> **NOTE:** All HA VIP configuration for particular Gateway should be defined within single resource clause - use multiple `config` sections to define HA on multiple interfaces. Defining multiple resources for same Gateway may result in apply conflicts.

# Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "gw1" {
  display_name = "gw1"
}

resource "nsxt_policy_tier0_gateway_interface" "if1" {
  display_name = "segment0_interface"
  description  = "connection to segment1"
  type         = "EXTERNAL"
  gateway_path = data.nsxt_policy_tier0_gateway.gw1.path
  segment_path = nsxt_policy_vlan_segment.segment1.path
  subnets      = ["12.12.2.1/24"]
}

resource "nsxt_policy_tier0_gateway_interface" "if2" {
  display_name = "segment0_interface"
  description  = "connection to segment2"
  type         = "EXTERNAL"
  gateway_path = data.nsxt_policy_tier0_gateway.gw1.path
  segment_path = nsxt_policy_vlan_segment.segment2.path
  subnets      = ["12.12.2.2/24"]
}

resource "nsxt_policy_tier0_gateway_ha_vip_config" "ha-vip" {
  config {
    enabled                  = true
    external_interface_paths = [nsxt_policy_tier0_gateway_interface.if1.path, nsxt_policy_tier0_gateway_interface.if2.path]
    vip_subnets              = ["12.12.2.3/24"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `config` - (Required) List of HA vip configurations (all belonging to the same Tier0 locale-service) containing:
  * `enabled` - (Optional) Flag indicating if this HA VIP config is enabled. True by default.
  * `vip_subnets` - (Required) 1 or 2 Ip Addresses/Prefixes in CIDR format, which will be used as floating IP addresses.
  * `external_interface_paths` - (Required) Paths of 2 external interfaces belonging to the same Tier0 gateway locale-service, which are to be paired to provide redundancy. Floating IP will be owned by one of these interfaces depending upon which edge node is active.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `tier0_id` - ID of the Tier-0 Gateway
* `locale_service_id` - ID of the Tier-0 Gateway locale service.

## Importing

An existing policy Tier-0 Gateway HA Vip config can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_tier0_gateway_ha_vip_config.havip GW-ID/LOCALE-SERVICE-ID
```

The above command imports the policy Tier-0 gateway HA Vip config named `havip` on Tier0 Gateway `GW-ID`, under locale service `LOCALE-SERVICE-ID`.
