---
subcategory: "Gateways and Routing"
page_title: "NSXT: policy_tier0_gateway_interface"
description: A policy Tier-0 gateway interface data source.
---

# nsxt_policy_gateway_interface

This data source provides information about policy Tier-0 & Tier-1 gateway interface configured on NSX.

This data source is applicable to NSX Policy Manager and NSX Global Manager.

## Tier0 gateway interface example

```hcl
data "nsxt_policy_tier0_gateway" "t0_gw" {
  display_name      = "t0Gateway"
}

data "nsxt_policy_gateway_interface" "tier0_gw_interface" {
  display_name    = "gw-interface1"
  gateway_path = data.nsxt_policy_tier0_gateway.t0_gw.path
}
```

## Tier1 gateway interface example

```hcl
data "nsxt_policy_tier1_gateway" "t1_gw" {
  display_name      = "t1Gateway"
}

data "nsxt_policy_gateway_interface" "tier1_gw_interface" {
  display_name    = "gw-interface2"
  gateway_path = data.nsxt_policy_tier1_gateway.t1_gw.path
}
```

## Argument Reference

* `display_name` - (Required) The Display Name prefix of the gateway interface to retrieve.
* `gateway_path` - (Required) The path of the gateway to retrieve which the interface should be linked to.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the interface.
* `description` - The description of the resource.
* `edge_cluster_path` - The path of the Edge cluster where this gateway is placed. This attribute is not set for NSX Global Manager, where gateway can span across multiple sites. This attribute is set only for Tier0 gateways.
* `path` - The NSX path of the policy resource.
* `segment_path` - Policy path for segment which is connected to this Tier0 Gateway
