---
subcategory: "Gateways and Routing"
page_title: "NSXT: policy_tier0_gateway"
description: A policy Tier-0 gateway data source.
---

# nsxt_policy_tier0_gateway

This data source provides information about policy Tier-0 gateways configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "tier0_gw_gateway" {
  display_name = "tier0-gw"
}
```

## Example Usage - Global infra

```hcl
data "nsxt_policy_tier0_gateway" "tier0_gw_gateway_global" {
  context {
    from_global = true
  }
  display_name = "tier0-gw"
}
```

## Argument Reference

* `id` - (Optional) The ID of Tier-0 gateway to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Tier-0 gateway to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `from_global` - (Optional) Set to True if the data source will need to search Tier-0 gateway created in a global manager instance (/global-infra)

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `edge_cluster_path` - The path of the Edge cluster where this Tier-0 gateway is placed. This attribute is not set for NSX Global Manager, where gateway can span across multiple sites.
* `path` - The NSX path of the policy resource.
* `is_global` - Set to true to get global objects from the local NSX manager.
* `vrf_config` - VRF config for VRF Tier0. This clause is supported with NSX 3.0.0 onwards.
    * `gateway_path` - Default Tier0 path. Cannot be modified after realization.
    * `evpn_transit_vni` - L3 VNI associated with the VRF for overlay traffic. VNI must be unique and belong to configured VNI pool.
    * `route_distinguisher` - Route distinguisher. Format: `<ASN>:<number>` or `<IPAddress>:<number>`.
    * `route_target` - Only one target is supported.
        * `auto_mode` - When true, import and export targets should not be specified.
        * `address_family` - Address family, currently only `L2VPN_EVPN` is supported, which is the default.
        * `import_targets` - List of import route targets. Format: `<ASN>:<number>`.
        * `export_targets` - List of export route targets. Format: `<ASN>:<number>`.
