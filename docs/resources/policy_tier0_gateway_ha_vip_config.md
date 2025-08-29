---
subcategory: "Gateways and Routing"
page_title: "NSXT: nsxt_policy_tier0_gateway_ha_vip_config"
description: A resource to configure HA Vip config on Tier-0 gateway in NSX Policy manager.
---

# nsxt_policy_tier0_gateway_ha_vip_config

This resource provides a method for the management of a Tier-0 gateway HA VIP config. Note that this configuration can be defined only for Active-Standby Tier0 gateway.

In order to correctly configure the HA VIP for a Tier-0 gateway, the selected external interfaces should be configured with the same subnet CIDR, and the vip_subnets setting
should specify an IP in this subnet CIDR which is not already in use by any external interface.

~> **NOTE:** All HA VIP configuration for particular Gateway should be defined within single resource clause - use multiple `config` sections to define HA on multiple interfaces. Defining multiple resources for same Gateway may result in apply conflicts.

## Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "tier0_gw" {
  display_name = "tier0_gw"
}

data "nsxt_policy_segment" "extseg_1" {
  display_name = "extseg1"
}

data "nsxt_policy_segment" "extseg_2" {
  display_name = "extseg2"
}

data "nsxt_policy_edge_cluster" "ec" {
  display_name = "ec"
}

data "nsxt_policy_edge_node" "node1" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.ec.path
  member_index      = 0
}

data "nsxt_policy_edge_node" "node2" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.ec.path
  member_index      = 1
}

resource "nsxt_policy_tier0_gateway_interface" "if1" {
  description    = "connection to segment1"
  display_name   = "seg1"
  gateway_path   = data.nsxt_policy_tier0_gateway.tier0_gw.path
  segment_path   = data.nsxt_policy_segment.extseg_1.path
  edge_node_path = data.nsxt_policy_edge_node.node1.path
  subnets        = ["12.12.1.1/24"]
  type           = "EXTERNAL"
  urpf_mode      = "STRICT"
}

resource "nsxt_policy_tier0_gateway_interface" "if2" {
  description    = "connection to segment1"
  display_name   = "seg2"
  gateway_path   = data.nsxt_policy_tier0_gateway.tier0_gw.path
  segment_path   = data.nsxt_policy_segment.extseg_2.path
  edge_node_path = data.nsxt_policy_edge_node.node2.path
  subnets        = ["12.12.1.2/24"]
  type           = "EXTERNAL"
  urpf_mode      = "STRICT"
}

resource "nsxt_policy_tier0_gateway_ha_vip_config" "ha-vip" {
  config {
    enabled                  = true
    external_interface_paths = [nsxt_policy_tier0_gateway_interface.if1.path, nsxt_policy_tier0_gateway_interface.if2.path]
    vip_subnets              = ["12.12.1.3/24"]
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

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_tier0_gateway_ha_vip_config.havip GW-ID/LOCALE-SERVICE-ID
```

The above command imports the policy Tier-0 gateway HA Vip config named `havip` on Tier0 Gateway `GW-ID`, under locale service `LOCALE-SERVICE-ID`.
