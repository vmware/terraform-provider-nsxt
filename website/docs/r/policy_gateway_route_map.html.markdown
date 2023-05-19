---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_gateway_route_map"
description: A resource to configure Route Map on Tier0 Gateway.
---

# nsxt_policy_gateway_route_map

This resource provides a method for the management of Route Map on Tier0 Gateway.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_gateway_route_map" "test" {
  display_name = "test"
  description  = "Terraform provisioned list"
  gateway_path = nsxt_policy_tier0_gateway.test.path

  entry {
    action = "PERMIT"
    community_list_match {
      criteria       = "11:*"
      match_operator = "MATCH_COMMUNITY_REGEX"
    }

    community_list_match {
      criteria       = "11:*"
      match_operator = "MATCH_LARGE_COMMUNITY_REGEX"
    }
  }

  entry {
    action              = "PERMIT"
    prefix_list_matches = [nsxt_policy_gateway_prefix_list.test.path]

    set {
      local_preference = 1122
      med              = 120
      weight           = 12
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `gateway_path` - (Required) Policy path of relevant Tier0 Gateway.
* `entry` - (Required) List of entries for the Route Map.
  * `action` - (Optional) Action for the route map entry, either `PERMIT` or `DENY`, with default being `PERMIT`.
  * `community_list_match` - (Optional) List of Prefix List match criteria for route map. Cannot be configured together with `prefix_list_matches`. If configured together, `prefix_list_matches` will be ignored.
    * `criteria` - (Required) Community list path or a regular expression.
    * `match_operator` - (Required) Match operator for the criteria, one of `MATCH_ANY`, `MATCH_ALL`, `MATCH_EXACT`, `MATCH_COMMUNITY_REGEX`, `MATCH_LARGE_COMMUNITY_REGEX`. Only last two operators can be used together with regular expression criteria.
  * `prefix_list_matches` - (Optional) List of policy paths for Prefix Lists configured on this Gateway. Cannot be configured together with `community_list_match`. If configured together, `prefix_list_matches` will be ignored.
  * `set` - (Optional) Set criteria for route map entry.
    * `as_path_prepend` - (Optional) Autonomous System (AS) path prepend to influence route selection.
    * `community` - (Optional) BGP regular or large community for matching routes.
    * `local_preference` - (Optional) Local preference indicates the degree of preference for one BGP route over other BGP routes.
    * `med` - (Optional) Multi Exit Descriminator (lower value is preferred over higher value).
    * `prefer_global_v6_next_hop` - (Optional)  Indicator whether to prefer IPv6 global address over link-local as the next hop.
    * `weight` - (Optional) Weight is used to select a route when multiple routes are available to the same network.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_gateway_route_map.test GW-ID/ID
```

The above command imports Tier0 Gateway Route Map named `test` with the NSX Route Map ID `ID` on Tier0 Gateway `GW-ID`.
