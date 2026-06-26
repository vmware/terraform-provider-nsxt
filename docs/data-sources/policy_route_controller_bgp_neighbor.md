---
subcategory: "Gateways and Routing"
page_title: "NSXT: nsxt_policy_route_controller_bgp_neighbor"
description: A data source to retrieve a Route Controller BGP Neighbor.
---

# nsxt_policy_route_controller_bgp_neighbor

This data source provides information about a Route Controller BGP Neighbor configured in NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_bgp_neighbor" "test" {
  display_name = "my-bgp-neighbor"
  parent_path  = "/infra/route-controllers/my-controller/bgp"
}
```

## Argument Reference

* `id` - (Optional) The ID of the BGP Neighbor to retrieve.
* `display_name` - (Optional) The Display Name of the BGP Neighbor to retrieve.
* `parent_path` - (Optional) Policy path of the parent Route Controller BGP config. Narrows the search to neighbors belonging to this BGP config.
* `description` - (Optional) The description of the BGP Neighbor to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `path` - The NSX path of the policy resource.
