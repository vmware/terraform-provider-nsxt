---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_gateway_redistribution_config"
description: A resource to configure Route Redistribution on Tier-0 gateway in NSX Policy manager.
---

# nsxt_policy_gateway_redistribution_config

This resource provides a method for the management of a Tier-0 Gateway Route Redistribution config.

This resource is applicable to NSX Global Manager and NSX Policy Manager.
This resource is supported with NSX 3.0.0 onwards.

# Example Usage

```hcl
resource "nsxt_policy_gateway_redistribution_config" "test" {
  gateway_path = data.nsxt_policy_tier0_gateway.gw1.path
  bgp_enabled  = true
  ospf_enabled = true

  rule {
    name  = "test-rule-1"
    types = ["TIER1_CONNECTED", "TIER1_LB_VIP"]
  }

  rule {
    name           = "test-rule-2"
    types          = ["TIER0_EVPN_TEP_IP"]
    route_map_path = nsxt_policy_gateway_route_map.map1.path
  }
}
```

# Global Manager Example Usage

```hcl
resource "nsxt_policy_gateway_redistribution_config" "test" {
  gateway_path = data.nsxt_policy_tier0_gateway.gw1.path
  site_path    = data.nsxt_policy_site.paris.path

  bgp_enabled = true

  rule {
    name  = "test-rule-1"
    types = ["TIER1_CONNECTED", "TIER1_LB_VIP"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_path` - (Required) Policy path to Tier0 Gateway.
* `site_path` - (Optional) Policy path to Global Manager site (domain). This attribute is required for NSX Global Manager and not applicable otherwise.
* `bgp_enabled` - (Optional) Enable route redistribution for BGP. Defaults to `true`.
* `ospf_enabled` - (Optional) Enable route redistribution for OSPF. Defaults to `false`. Applicable from NSX 3.1.0 onwards.
* `rule` - (Optional) List of redistribution rules.
  * `name` - (Optional) Rule name.
  * `route_map_path` - (Optional) Route map to be associated with the redistribution rule.
  * `types` - (Optional) List of redistribution types, possible values are: `TIER0_STATIC`, `TIER0_CONNECTED`, `TIER0_EXTERNAL_INTERFACE`, `TIER0_SEGMENT`, `TIER0_ROUTER_LINK`, `TIER0_SERVICE_INTERFACE`, `TIER0_LOOPBACK_INTERFACE`, `TIER0_DNS_FORWARDER_IP`, `TIER0_IPSEC_LOCAL_IP`, `TIER0_NAT`, `TIER0_EVPN_TEP_IP`, `TIER1_NAT`, `TIER1_STATIC`, `TIER1_LB_VIP`, `TIER1_LB_SNAT`, `TIER1_DNS_FORWARDER_IP`, `TIER1_CONNECTED`, `TIER1_SERVICE_INTERFACE`, `TIER1_SEGMENT`, `TIER1_IPSEC_LOCAL_ENDPOINT`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `gateway_id` - ID of the Tier-0 Gateway
* `locale_service_id` - ID of the Tier-0 Gateway locale service.

## Importing

An existing policy Tier-0 Gateway Redistribution config can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_gateway_redistribution_config.havip GW-ID/LOCALE-SERVICE-ID
```

The above command imports the policy Tier-0 gateway Redistribution config named `havip` on Tier0 Gateway `GW-ID`, under locale service `LOCALE-SERVICE-ID`.
