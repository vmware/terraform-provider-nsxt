---
subcategory: "VPC"
page_title: "NSXT: policy_transit_gateway_route_map"
description: A policy Route Map on a Transit Gateway data source.
---

# nsxt_policy_transit_gateway_route_map

This data source provides information about a route map configured on a Transit Gateway.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_transit_gateway_route_map" "test" {
  display_name = "my-route-map"
  parent_path  = nsxt_policy_transit_gateway.demo.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the Transit Gateway Route Map to retrieve.
* `display_name` - (Optional) The Display Name of the Transit Gateway Route Map to retrieve.
* `parent_path` - (Required) Policy path of the parent Transit Gateway.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
