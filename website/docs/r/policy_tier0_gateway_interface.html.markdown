---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_tier0_gateway_interface"
description: A resource to configure an Interface on Tier-0 gateway on NSX Policy manager.
---

# nsxt_policy_tier0_gateway_interface

This resource provides a method for the management of a Tier-0 gateway Interface. Note that edge cluster must be configured on Tier-0 Gateway in order to configure interfaces on it.

~> **NOTE:** When configuring VRF-Lite interfaces, please specify explicit dependency on parent Tier-0 Gateway interface(see VRF interface example below).  This will ensure correct order of object deletion.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

# Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "gw1" {
  display_name = "gw1"
}

data "nsxt_policy_ipv6_ndra_profile" "slaac" {
  display_name = "slaac"
}

resource "nsxt_policy_vlan_segment" "segment0" {
  display_name = "segment0"
  vlan_ids     = [12]
}

resource "nsxt_policy_tier0_gateway_interface" "if1" {
  display_name           = "segment0_interface"
  description            = "connection to segment0"
  type                   = "SERVICE"
  gateway_path           = data.nsxt_policy_tier0_gateway.gw1.path
  segment_path           = nsxt_policy_vlan_segment.segment0.path
  subnets                = ["12.12.2.13/24"]
  mtu                    = 1500
  ipv6_ndra_profile_path = data.nsxt_policy_ipv6_ndra_profile.slaac.path
}
```

# VRF Interface Example Usage

```hcl
resource "nsxt_policy_tier0_gateway_interface" "red_vrf_uplink1" {
  display_name   = "Uplink-01"
  type           = "EXTERNAL"
  edge_node_path = data.nsxt_policy_edge_node.edge_node_1.path
  gateway_path   = nsxt_policy_tier0_gateway.red_vrf.path
  segment_path   = nsxt_policy_vlan_segment.vrf_trunk_1.path
  access_vlan_id = 112
  subnets        = ["192.168.112.254/24"]
  mtu            = 1500

  depends_on = [nsxt_policy_tier0_gateway_interface.parent_uplink1]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `type` - (Optional) Type of this interface, one of `SERVICE`, `EXTERNAL`, `LOOPBACK`. Default is `EXTERNAL`
* `gateway_path` - (Required) Policy path for the Tier-0 Gateway.
* `segment_path` - (Optional) Policy path for segment to be connected with this Tier1 Gateway. This argemnt is required for interfaces of type `SERVICE` and `EXTERNAL`.
* `subnets` - (Required) list of Ip Addresses/Prefixes in CIDR format, to be associated with this interface.
* `edge_node_path` - (Optional) Path of edge node for this interface, relevant for interfaces of type `EXTERNAL`.
* `mtu` - (Optional) Maximum Transmission Unit for this interface.
* `ipv6_ndra_profile_path` - (Optional) IPv6 NDRA profile to be associated with this interface.
* `dhcp_relay_path` - (Optional) DHCP relay path to be associated with this interface.
* `enable_pim` - (Optional) Flag to enable Protocol Independent Multicast, relevant only for interfaces of type `EXTERNAL`. This attribute will always be `false` for other interface types. This attribute is supported with NSX 3.0.0 onwards, and only for local managers.
* `access_vlan_id`- (Optional) Access VLAN ID, relevant only for VRF interfaces. This attribute is supported with NSX 3.0.0 onwards.
* `urpf_mode` - (Optional) Unicast Reverse Path Forwarding mode, one of `NONE`, `STRICT`. Default is `STRICT`. This attribute is supported with NSX 3.0.0 onwards.
* `site_path` - (Required for global manager only) Path of the site the Tier0 edge cluster belongs to. This configuration is required for global manager only. `path` field of the existing `nsxt_policy_site` can be used here.
* `ospf` - (Optional) OSPF configuration block - supported for `EXTERNAL` interface only. Not supported on Global Manager.
  * `enabled` - (Optional) Flag to enable/disable OSPF for this interface. Default is `true`.
  * `enable_bfd` - (Optional) Flag that controls whether OSPF will register for BFD event. Default is `false`.
  * `bfd_profile_path` - (Optional) Policy path to BFD profile. Relevant only if BFD is enabled.
  * `area_path` - (Required) Policy path to OSPF area defined on this Tier0 Gateway.
  * `network_type` - (Optional) OSPF network type, one of `BROADCAST` and `P2P`. Default is `BROADCAST`.
  * `hello_interval` - (Optional) Interval between OSPF Hello Packets, in seconds. Defaults to 10.
  * `dead_interval` - (Optional) Interval to wait before declaring OSPF peer as down, in seconds. Defaults to 40. Must be at least 3 times greater than `hello_interval`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `ip_addresses` - list of Ip Addresses picked from each subnet in `subnets` field. This attribute can serve as `source_addresses` field of `nsxt_policy_bgp_neighbor` resource.

## Importing

An existing policy Tier-0 Gateway Interface can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_tier0_gateway_interface.interface1 GW-ID/LOCALE-SERVICE-ID/ID
```

The above command imports the policy Tier-0 gateway interface named `interface1` with the NSX Policy ID `ID` on Tier0 Gateway `GW-ID`, under locale service `LOCALE-SERVICE-ID`.
