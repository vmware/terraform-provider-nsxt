---
subcategory: "Gateways and Routing"
page_title: "NSXT: policy_tier0_gateway_interface"
description: A policy Tier-0 gateway interface data source.
---

# nsxt_policy_tier0_gateway_interface

This data source provides information about policy Tier-0 gateway interface configured on NSX.

This data source is applicable to NSX Policy Manager and NSX Global Manager.

## Example Usage

```hcl
data "nsxt_policy_tier0_gateway_interface" "tier0_gw_interface" {
  display_name    = "tier0-gw-interface"
  t0_gateway_path = "tier0-gw"
}
```

## Argument Reference

* `id` - (Optional) The ID of Tier-0 gateway to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Tier-0 gateway interface to retrieve.
* `t0_gateway_path` - (Optional) The Tier-0 gateway to retrieve to which the interface should be linked to.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `edge_cluster_path` - The path of the Edge cluster where this Tier-0 gateway is placed. This attribute is not set for NSX Global Manager, where gateway can spawn across multiple sites.
* `path` - The NSX path of the policy resource.
* `segment_path` - Policy path for segment which is connected to this Tier0 Gateway
