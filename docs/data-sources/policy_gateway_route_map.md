---
subcategory: "Gateways and Routing"
page_title: "NSXT: policy_gateway_route_map"
description: Policy Gateway Route Map data source.
---

# nsxt_policy_gateway_route_map

This data source provides information about a policy Gateway Route Map configured on NSX.

This data source is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "gw1" {
  display_name = "gw1"
}

data "nsxt_policy_gateway_route_map" "route_map1" {
  gateway_path = data.nsxt_policy_tier0_gateway.gw1.path
  display_name = "test"
}
```

## Argument Reference

* `id` - (Optional) The ID of the Gateway Route Map to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Gateway Route Map to retrieve.
* `gateway_path` - (Optional) Policy path of the Gateway where the Route Map is configured.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
