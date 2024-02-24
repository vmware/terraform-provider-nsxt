---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_tier0_gateway_gre_tunnel"
description: A resource to configure a GRE tunnel on Tier-0 gateway on NSX Policy manager.
---

# nsxt_policy_tier0_gateway_gre_tunnel

This resource provides a method for the management of a Tier-0 gateway GRE tunnel. Note that edge cluster and an interface must be configured on Tier-0 Gateway in order to configure GRE tunnels on it.

This resource is applicable to NSX Policy Manager.

# Example Usage

```hcl
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "EDGECLUSTER1"
}

data "nsxt_policy_transport_zone" "test" {
  display_name = "tz-vlan"
}
resource "nsxt_policy_vlan_segment" "test" {

  transport_zone_path = data.nsxt_policy_transport_zone.test.path
  display_name        = "test-vlan-segment"
  vlan_ids            = [11]
  subnet {
    cidr = "10.2.2.2/24"
  }
}
data "nsxt_policy_edge_node" "EN" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  member_index      = 0
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name = "test-gateway"
  ha_mode      = "ACTIVE_STANDBY"

  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

resource "nsxt_policy_tier0_gateway_interface" "test" {
  display_name   = "test-interface"
  type           = "EXTERNAL"
  mtu            = 1500
  gateway_path   = nsxt_policy_tier0_gateway.test.path
  segment_path   = nsxt_policy_vlan_segment.test.path
  edge_node_path = data.nsxt_policy_edge_node.EN.path
  subnets        = ["192.168.240.10/24"]
  enable_pim     = "true"
}

data "nsxt_policy_gateway_locale_service" "test" {
  gateway_path = nsxt_policy_tier0_gateway.test.path
}

data "nsxt_policy_edge_node" "node1" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  member_index      = 0
}

resource "nsxt_policy_tier0_gateway_gre_tunnel" "test" {
  display_name        = "test-tunnel"
  description         = "Test GRE tunnel"
  locale_service_path = data.nsxt_policy_gateway_locale_service.test.path
  destination_address = "192.168.221.221"
  tunnel_address {
    source_address = nsxt_policy_tier0_gateway_interface.test.ip_addresses[0]
    edge_path      = data.nsxt_policy_edge_node.node1.path
    tunnel_interface_subnet {
      ip_addresses = ["192.168.243.243"]
      prefix_len   = 24
    }
  }
  tunnel_keepalive {
    enabled = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `locale_service_path` - (Required) Policy path of associated Gateway Locale Service on NSX.
* `destination_address` - (Required) Destination IPv4 address.
* `enabled` - (Optional) - Enable/Disable Tunnel. Default is true.
* `mtu` - (Optional) - Maximum transmission unit. Default is 1476.
* `tunnel_address` - (Required) Tunnel Address object parameter. At least one is required, maximum is 8.
  * `edge_path` - (Required) Policy edge node path.
  * `source_address` - (Required) IPv4 source address.
  * `tunnel_interface_subnet` - (Required) Interface Subnet object parameter. At least one is required, maximum is 2.
    * `ip_addresses` - (Required) List of IP addresses assigned to interface.
    * `prefix_len` - (Required) Subnet prefix length.
* `tunnel_keepalive` - (Optional) tunnel keep alive object. One is required.
  * `dead_time_multiplier` - (Optional) Dead time multiplier. Default is 3.
  * `enable_keepalive_ack` - (Optional) Enable tunnel keep alive acknowledge. Default is true.
  * `enabled` - (Optional) Enable/Disable tunnel keep alive. Default is false.
  * `keepalive_interval` - (Optional) Keep alive interval. Default is 10.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.


## Importing

An existing Tier-0 GRE tunnel can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import
```
terraform import nsxt_policy_tier0_gateway_gre_tunnel.test_tunnel POLICY_PATH
```
The above command imports the Tier-0 GRE tunnel named `test_tunnel` with the policy path `POLICY_PATH`.
