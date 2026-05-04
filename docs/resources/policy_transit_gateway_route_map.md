---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_transit_gateway_route_map"
description: A resource to configure a Route Map on a Transit Gateway.
---

# nsxt_policy_transit_gateway_route_map

This resource provides a method for the management of a route map on a Transit Gateway.

Transit gateway route maps use the same `RouteMapEntry` schema as Tier-0 gateway route maps and are used for BGP route policy.

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

resource "nsxt_policy_transit_gateway_route_map" "example" {
  display_name = "my-route-map"
  description  = "Route map for BGP policy"
  parent_path  = data.nsxt_policy_transit_gateway.test_tgw.path

  entry {
    action              = "PERMIT"
    prefix_list_matches = [nsxt_policy_transit_gateway_prefix_list.allow.path]

    set {
      local_preference = 200
      med              = 100
    }
  }

  entry {
    action = "DENY"
  }
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required, ForceNew) Policy path of the parent Transit Gateway (e.g. `/orgs/default/projects/proj1/transit-gateways/tgw1`).
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `entry` - (Required) Ordered list of route map entries. Each entry supports:
    * `action` - (Optional) Action for the entry. Accepted values: `PERMIT`, `DENY`. Defaults to `PERMIT`.
    * `prefix_list_matches` - (Optional) Set of policy paths for prefix lists to match against.
    * `community_list_match` - (Optional) List of BGP community match criteria. Each item supports:
        * `criteria` - (Required) Community list path or a regular expression.
        * `match_operator` - (Required) Match operator. Accepted values: `ANY`, `ALL`, `EXACT`, `COMMUNITY_REGEX`, `LARGE_COMMUNITY_REGEX`.
    * `set` - (Optional) Set clause applied to matching routes (max 1 block):
        * `as_path_prepend` - (Optional) AS path string to prepend (influences route selection).
        * `community` - (Optional) BGP regular or large community value to set on matching routes.
        * `local_preference` - (Optional, Computed) Local preference value.
        * `med` - (Optional, Computed) Multi-Exit Discriminator value; lower is preferred.
        * `prefer_global_v6_next_hop` - (Optional) When `true`, prefer global IPv6 address over link-local as the next hop.
        * `weight` - (Optional) Weight used to prefer one route over another to the same destination.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

```shell
terraform import nsxt_policy_transit_gateway_route_map.example PARENT_PATH/ID
```

The above command imports a Transit Gateway Route Map named `example` with the transit gateway NSX path as `PARENT_PATH` and the route map ID as `ID`. For example:

```shell
terraform import nsxt_policy_transit_gateway_route_map.example /orgs/default/projects/project1/transit-gateways/tgw1/route-map-1
```

[docs-import]: https://developer.hashicorp.com/terraform/cli/import
