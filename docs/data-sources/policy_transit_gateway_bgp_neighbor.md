---
subcategory: "VPC"
page_title: "NSXT: policy_transit_gateway_bgp_neighbor"
description: A policy BGP Neighbor on a Transit Gateway data source.
---

# nsxt_policy_transit_gateway_bgp_neighbor

This data source provides information about a BGP neighbor configuration on a Transit Gateway.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_transit_gateway_bgp_neighbor" "test" {
  display_name = "my-bgp-neighbor"
  parent_path  = nsxt_policy_transit_gateway.demo.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the Transit Gateway BGP Neighbor to retrieve.
* `display_name` - (Optional) The Display Name of the Transit Gateway BGP Neighbor to retrieve.
* `parent_path` - (Required) Policy path of the parent Transit Gateway. This is used to scope the search, since the underlying BGP neighbor object type is shared across Tier-0/Tier-1 gateways and Transit Gateways.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
