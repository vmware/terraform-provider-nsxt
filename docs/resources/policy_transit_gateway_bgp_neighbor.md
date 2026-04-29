---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_transit_gateway_bgp_neighbor"
description: A resource to configure a BGP Neighbor on a Transit Gateway.
---

# nsxt_policy_transit_gateway_bgp_neighbor

This resource provides a method for the management of a BGP neighbor configuration on a Transit Gateway.

Transit gateway BGP neighbors use `source_attachment` (instead of `source_addresses`) to specify the transit gateway attachment or IPSec route-based VPN session used as the BGP peering source.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "test_proj" {
  display_name = "test_project"
}

data "nsxt_policy_transit_gateway" "test_tgw" {
  context {
    project_id = data.nsxt_policy_project.test_proj.id
  }
  id = "default"
}

resource "nsxt_policy_transit_gateway_bgp_neighbor" "example" {
  display_name     = "my-bgp-neighbor"
  description      = "BGP neighbor for transit gateway"
  parent_path      = data.nsxt_policy_transit_gateway.test_tgw.path
  neighbor_address = "192.168.20.1"
  remote_as_num    = "65001"

  source_attachment = [nsxt_policy_transit_gateway_attachment.cna_att.path]
}
```

### With route filtering and local-AS configuration

```hcl
resource "nsxt_policy_transit_gateway_bgp_neighbor" "advanced" {
  display_name     = "my-bgp-neighbor-advanced"
  parent_path      = data.nsxt_policy_transit_gateway.test_tgw.path
  neighbor_address = "192.168.20.2"
  remote_as_num    = "65002"
  hold_down_time   = 90
  keep_alive_time  = 30

  source_attachment = [nsxt_policy_transit_gateway_attachment.cna_att.path]

  route_filtering {
    address_family  = "IPV4"
    enabled         = true
    in_route_filter = "/infra/tier-0s/tier0-1/locale-services/default/prefix-lists/inbound"
  }

  neighbor_local_as_config {
    local_as_num          = "65100"
    as_path_modifier_type = "NO_PREPEND"
  }

  bfd_config {
    enabled  = true
    interval = 500
    multiple = 3
  }
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required, ForceNew) Policy path of the parent transit gateway (e.g. `/orgs/default/projects/proj1/transit-gateways/tgw1`).
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `neighbor_address` - (Required) IP address (IPv4 or IPv6) of the BGP neighbor.
* `remote_as_num` - (Required) ASN of the neighbor in ASPLAIN or ASDOT format.
* `source_attachment` - (Optional) A list containing exactly one policy path of the transit gateway attachment or IPSec route-based VPN session used as the BGP peering source.
* `allow_as_in` - (Optional) Flag to enable the allow-as-in option for this BGP neighbor. Defaults to `false`.
* `graceful_restart_mode` - (Optional) BGP graceful restart configuration mode. Accepted values: `HELPER_ONLY`, `GR_AND_HELPER`, `DISABLE`. Defaults to `HELPER_ONLY`.
* `hold_down_time` - (Optional) Wait time in seconds before declaring the peer dead. Valid range 1–65535. Defaults to `180`.
* `keep_alive_time` - (Optional) Interval in seconds between keep-alive messages sent to the peer. Valid range 1–65535. Defaults to `60`.
* `maximum_hop_limit` - (Optional) Maximum number of hops allowed to reach the BGP neighbor. Valid range 1–255. Defaults to `1`.
* `password` - (Optional, Sensitive) Password for BGP neighbor MD5 authentication. Maximum 32 characters.
* `bfd_config` - (Optional) BFD failure detection configuration block:
    * `enabled` - (Optional) Flag to enable/disable BFD. Defaults to `false`.
    * `interval` - (Optional) Heartbeat interval in milliseconds. Valid range 50–60000. Defaults to `500`.
    * `multiple` - (Optional) Number of missed heartbeats before the neighbor is declared down. Valid range 2–16. Defaults to `3`.
* `route_filtering` - (Optional) List of address family route filtering entries (up to 2):
    * `address_family` - (Required) Address family type. Accepted values: `IPV4`, `IPV6`, `L2VPN_EVPN`.
    * `enabled` - (Optional) Flag to enable/disable this address family. Defaults to `true`.
    * `in_route_filter` - (Optional) Policy path of a prefix-list or route map applied in the IN direction.
    * `out_route_filter` - (Optional) Policy path of a prefix-list or route map applied in the OUT direction.
    * `maximum_routes` - (Optional) Maximum number of routes for this address family. Valid range 1–1000000.
* `neighbor_local_as_config` - (Optional) BGP neighbor local-AS configuration block:
    * `local_as_num` - (Required) Local AS number in ASPLAIN or ASDOT format.
    * `as_path_modifier_type` - (Optional) AS_PATH modifier type. Accepted values: `NO_PREPEND`, `NO_PREPEND_REPLACE_AS`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

```shell
terraform import nsxt_policy_transit_gateway_bgp_neighbor.example PARENT_PATH/ID
```

The above command imports a Transit Gateway BGP Neighbor named `example` using the transit gateway NSX path as `PARENT_PATH` and the neighbor ID as `ID`. For example:

```shell
terraform import nsxt_policy_transit_gateway_bgp_neighbor.example /orgs/default/projects/project1/transit-gateways/tgw1/bgp-neighbor-1
```

[docs-import]: https://developer.hashicorp.com/terraform/cli/import
