---
subcategory: "Gateways and Routing"
page_title: "NSXT: nsxt_policy_route_controller_bgp_neighbor"
description: A resource to configure a Route Controller BGP Neighbor.
---

# nsxt_policy_route_controller_bgp_neighbor

This resource provides a method for the management of a Route Controller BGP Neighbor.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_route_controller_bgp_neighbor" "test" {
  display_name     = "test"
  description      = "Terraform provisioned Route Controller BGP Neighbor"
  parent_path      = "${nsxt_policy_route_controller.rc.path}/bgp"
  neighbor_address = "192.168.100.1"
  remote_as_num    = "65001"
  enabled          = true
  hold_down_time   = 180
  keep_alive_time  = 60

  route_filtering {
    address_family = "IPV4"
    enabled        = true
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `parent_path` - (Required) Policy path of the parent Route Controller BGP config (e.g. `<route_controller_path>/bgp`).
* `neighbor_address` - (Required) Neighbor IP Address.
* `remote_as_num` - (Required) ASN of the neighbor in ASPLAIN or ASDOT format.
* `enabled` - (Optional) Flag to enable or disable BGP peering with this neighbor. Defaults to `true`.
* `allow_as_in` - (Optional) Flag to enable allow_as_in option for BGP neighbor. Defaults to `false`.
* `gateway_ips` - (Optional) List of next-hop gateway IP addresses used to reach non-directly connected BGP peers.
* `graceful_restart_mode` - (Optional) BGP Graceful Restart Configuration Mode. One of `DISABLE`, `GR_AND_HELPER`, `HELPER_ONLY`. Defaults to `HELPER_ONLY`.
* `hold_down_time` - (Optional) Wait time in seconds before declaring peer dead. Must be between 1 and 65535. Defaults to `180`.
* `keep_alive_time` - (Optional) Interval between keep alive messages sent to peer. Must be between 1 and 65535. Defaults to `60`.
* `maximum_hop_limit` - (Optional) Maximum number of hops allowed to reach BGP neighbor. Must be between 1 and 255. Defaults to `1`.
* `password` - (Optional) Password for BGP neighbor authentication. Sensitive value. Maximum 32 characters.
* `source_addresses` - (Optional) List of source IP addresses for BGP peering. Maximum 8 entries.
* `bfd_config` - (Optional) BFD configuration for failure detection. The following arguments are supported:
    * `enabled` - (Optional) Flag to enable/disable BFD configuration. Defaults to `false`.
    * `interval` - (Optional) Time interval between heartbeat packets in milliseconds. Must be between 50 and 60000. Defaults to `500`.
    * `multiple` - (Optional) Number of times heartbeat packet is missed before BFD declares the neighbor is down. Must be between 2 and 16. Defaults to `3`.
* `route_filtering` - (Optional) List of address families and route filtering configuration. Maximum 2 entries. The following arguments are supported:
    * `address_family` - (Required) Address family type. One of `IPV4`, `IPV6`, `L2VPN_EVPN`.
    * `enabled` - (Optional) Flag to enable/disable address family. Defaults to `true`.
    * `in_route_filter` - (Optional) Policy path of prefix-list or route map for IN direction.
    * `out_route_filter` - (Optional) Policy path of prefix-list or route map for OUT direction.
    * `maximum_routes` - (Optional) Maximum number of routes for the address family. Must be between 1 and 1000000.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing Route Controller BGP Neighbor can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```shell
terraform import nsxt_policy_route_controller_bgp_neighbor.test <path>
```

The above command imports the Route Controller BGP Neighbor named `test` using the policy path as import ID.
The import ID should be the full policy path of the BGP neighbor, for example:
`/infra/route-controllers/<controller-id>/bgp/neighbors/<neighbor-id>`
