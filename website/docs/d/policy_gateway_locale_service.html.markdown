---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: policy_gateway_locale_service"
description: A policy gateway locale service data source.
---

# nsxt_policy_gateway_locale_service

This data source provides information about certain locale service for Tier-0 or Tier-1 gateway configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_gateway_locale_service" "test" {
  gateway_path = data.nsxt_policy_tier0_gateway.path
  display_name = "london"
}
```

## Argument Reference

* `gateway_path` - (Required) Path for the gateway.
* `id` - (Optional) The ID of locale service gateway to retrieve.
* `display_name` - (Optional) The Display Name or prefix of locale service to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `edge_cluster_path` - The path of the Edge cluster configured on this service.
* `path` - The NSX path of the policy resource.
* `bgp_path` - Path for BGP configuration configured on this service.
