---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_bgp_neighbor"
description: A resource to configure a BGP Neighbor.
---

# nsxt_policy_bgp_neighbor

This resource provides a method for the management of a BGP Neighbor.

This resource is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl

resource "nsxt_policy_gateway_prefix_list" "test" {
  display_name = "prefix_list"
  gateway_path = "nsxt_policy_tier0_gateway.testresource.path"

  prefix {
    action  = "DENY"
    ge      = "20"
    le      = "23"
    network = "4.4.0.0/20"
  }
}

resource "nsxt_policy_bgp_neighbor" "test" {
  display_name          = "tfbpg"
  description           = "Terraform provisioned BgpNeighborConfig"
  bgp_path              = nsxt_policy_tier0_gateway.testresource.bgp_config.0.path
  allow_as_in           = true
  graceful_restart_mode = "HELPER_ONLY"
  hold_down_time        = 300
  keep_alive_time       = 100
  neighbor_address      = "12.12.11.23"
  password              = "passw0rd"
  remote_as_num         = "60000"
  source_addresses      = nsxt_policy_tier0_gateway_interface.testresource.ip_addresses

  bfd_config {
    enabled  = true
    interval = 1000
    multiple = 4
  }

  route_filtering {
    address_family   = "IPV4"
    maximum_routes   = 20
    in_route_filter  = nsxt_policy_gateway_prefix_list.test.path
    out_route_filter = nsxt_policy_gateway_prefix_list.test.path
  }
}
```

~> **NOTE:** If bgp neighbor configuration depends on gateway interface, please add `depends_on` clause in `nsxt_policy_bgp_neighbor` resource in order to ensure correct order of creation/deletion.


## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `bgp_path` - (Required) The policy path to the BGP configuration for this neighbor.
* `allow_as_in` - (Optional) Flag to enable allowas_in option for BGP neighbor. Defaults to `false`.
* `graceful_restart_mode` - (Optional) BGP Graceful Restart Configuration Mode. One of `DISABLE`, `GR_AND_HELPER` or `HELPER_ONLY`.
* `hold_down_time` - (Optional) Wait time in seconds before declaring peer dead. Defaults to `180`.
* `keep_alive_time` - (Optional) Interval between keep alive messages sent to peer. Defaults to `60`.
* `maximum_hop_limit` - (Optional) Maximum number of hops allowed to reach BGP neighbor. Defaults to `1`.
* `neighbor_address` - (Required) Neighbor IP Address.
* `password` - (Optional) Password for BGP neighbor authentication. Set to the empty string to clear out the password.
* `remote_as_num` - (Required) ASN of the neighbor in ASPLAIN/ASDOT Format.
* `source_addresses` - (Optional) A list of up to 8 source IP Addresses for BGP peering. `ip_addresses` field of an existing `nsxt_policy_tier0_gateway_interface` can be used here.
* `bfd_config` - (Optional) The BFD configuration.
  * `enabled` - (Optional) A boolean flag to enable/disable BFD. Defaults to `false`.
  * `interval` - (Optional) Time interval between heartbeat packets in milliseconds. Defaults to `500`.
  * `multiple` - (Optional) Number of times heartbeat packet is missed before BFD declares the neighbor is down. Defaults to `3`.
* `route_filtering` - (Optional) Up to 2 route filters for the neighbor. Note that prior to NSX version 3.0.0, only 1 element is supported.
  * `address_family` - (Required) Address family type. Must be one of `L2VPN_EVPN`, `IPV4` or `IPV6`. Note the `L2VPN_EVPN` property is only available starting with NSX version 3.0.0.
  * `enabled`- (Optional) A boolean flag to enable/disable address family. Defaults to `false`.
  * `in_route_filter`- (Optional) Path of prefix-list or route map to filter routes for IN direction.
  * `out_route_filter`- (Optional) Path of prefix-list or route map to filter routes for OUT direction.
  * `maximum_routes` - (Optional) Maximum number of routes for the address family. Note this property is only available starting with NSX version 3.0.0.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing BGP Neighbor can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_bgp_neighbor.test T0_ID/LOCALE_SERVICE_ID/NEIGHBOR_ID
```

The above command imports BGP Neighbor named `test` with the NSX BGP Neighbor ID `NEIGHBOR_ID` from the Tier-0 `T0_ID` and Locale Service `LOCALE_SERVICE_ID`.

~> **NOTE:** BGP neigbor password configured on NSX will not be returned by NSX for security reasons, and hence not imported to terraform state.
