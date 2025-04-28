---
subcategory: "Gateways and Routing"
page_title: "NSXT: policy_tier1_gateway_interface"
description: A policy Tier-1 gateway interface data source.
---

# nsxt_policy_tier1_gateway_interface

This data source provides information about policy Tier-1 gateway interface configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_tier1_gateway_interface" "tier1_gw_interface" {
  display_name = "tier1-gw-interface"
  t1_gateway_name = "tier1-gw"
}
```

## Argument Reference

* `id` - (Optional) The ID of Tier-1 gateway to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Tier-1 gateway interface to retrieve.
* `t1_gateway_name` - (Optional) The Tier-1 gateway to retrieve to which the interface should be linked to.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `segment_path` - Policy path for segment which is connected to this Tier0 Gateway
