---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_tier0_gateway"
description: A resource to configure a Tier-0 gateway on NSX Policy manager.
---

# nsxt_policy_tier0_gateway

This resource provides a method for the management of a Tier-0 gateway or VRF-Lite gateway.

This resource is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_tier0_gateway" "tier0_gw" {
  description              = "Tier-0 provisioned by Terraform"
  display_name             = "Tier0-gw1"
  failover_mode            = "PREEMPTIVE"
  default_rule_logging     = false
  enable_firewall          = true
  force_whitelisting       = false
  ha_mode                  = "ACTIVE_STANDBY"
  internal_transit_subnets = ["102.64.0.0/16"]
  transit_subnets          = ["101.64.0.0/16"]
  vrf_transit_subnets      = ["100.64.0.0/16"]
  edge_cluster_path        = data.nsxt_policy_edge_cluster.EC.path
  rd_admin_address         = "192.168.0.2"

  bgp_config {
    local_as_num    = "60000"
    multipath_relax = false

    route_aggregation {
      prefix = "12.10.10.0/24"
    }

    route_aggregation {
      prefix = "12.11.10.0/24"
    }
  }

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## VRF-Lite Example Usage

```hcl
resource "nsxt_policy_tier0_gateway" "vrf-blue" {
  description              = "Tier-0 VRF provisioned by Terraform"
  display_name             = "Tier0-vrf"
  failover_mode            = "PREEMPTIVE"
  default_rule_logging     = false
  enable_firewall          = true
  ha_mode                  = "ACTIVE_STANDBY"
  internal_transit_subnets = ["102.64.0.0/16"]
  transit_subnets          = ["101.64.0.0/16"]
  edge_cluster_path        = data.nsxt_policy_edge_cluster.EC.path

  bgp_config {
    ecmp = true

    route_aggregation {
      prefix = "12.10.10.0/24"
    }
  }

  vrf_config {
    gateway_path        = data.nsxt_policy_tier0_gateway.parent.path
    route_distinguisher = "62000:10"
    evpn_transit_vni    = 76001
    route_target {
      auto_mode      = false
      import_targets = ["62000:2"]
      export_targets = ["62000:3", "10.2.2.0:3"]
    }
  }
}
```


## Global manager example usage
```hcl
resource "nsxt_policy_tier0_gateway" "tier0_gw" {
  description   = "Tier-0 provisioned by Terraform"
  display_name  = "Tier0-gw1"
  failover_mode = "PREEMPTIVE"

  locale_service {
    nsx_id            = "paris"
    edge_cluster_path = data.nsxt_policy_edge_cluster.paris.path
  }

  locale_service {
    nsx_id               = "london"
    edge_cluster_path    = data.nsxt_policy_edge_cluster.london.path
    preferred_edge_paths = [data.nsxt_policy_edge_node.edge1.path]
  }

  intersite_config {
    primary_site_path = data.nsxt_policy_site.paris.path
  }

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```


## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Tier-0 gateway.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `edge_cluster_path` - (Optional) The path of the edge cluster where the Tier-0 is placed.For advanced configuration and on Global Manager, use `locale_service` clause instead. Note that for some configurations (such as BGP) setting edge cluster is required.
* `locale_service` - (Optional) This is required on NSX Global Manager. Multiple locale services can be specified for multiple locations.
  * `nsx_id` - (Optional) NSX id for the locale service. It is recommended to specify this attribute in order to avoid unnecessary recreation of this object. Should be unique within the gateway.
  * `edge_cluster_path` - (Required) The path of the edge cluster where the Tier-0 is placed.
  * `preferred_edge_paths` - (Optional) Policy paths to edge nodes. Specified edge is used as preferred edge cluster member when failover mode is set to `PREEMPTIVE`.
  * `display_name` - (Optional) Display name for the locale service.
* `failover_mode` - (Optional) This failover mode determines, whether the preferred service router instance for given logical router will preempt the peer. Accepted values are PREEMPTIVE/NON_PREEMPTIVE.
* `default_rule_logging` - (Optional) Boolean flag indicating if the default rule logging will be enabled or not. The default value is false.
* `enable_firewall` - (Optional) Boolean flag indicating if the edge firewall will be enabled or not. The default value is true.
* `force_whitelisting` - (Deprecated) Boolean flag indicating if white-listing will be forced or not. The default value is false. This argument is deprecated and will be removed. Please use `nsxt_policy_predefined_gateway_policy` resource to control default action.
* `ipv6_ndra_profile_path` - (Optional) Policy path to IPv6 NDRA profile.
* `ipv6_dad_profile_path` - (Optional) Policy path to IPv6 DAD profile.
* `ha_mode` - (Optional) High-availability Mode for Tier-0. Valid values are `ACTIVE_ACTIVE` and `ACTIVE_STANDBY`.
* `internal_transit_subnets` - (Optional) Internal transit subnets in CIDR format. At most 1 CIDR.
* `transit_subnets` - (Optional) Transit subnets in CIDR format.
* `vrf_transit_subnets` - (Optional) VRF transit subnets in CIDR format. Maximum 1 item allowed in the list.
* `dhcp_config_path` - (Optional) Policy path to DHCP server or relay configuration to use for this gateway.
* `rd_admin_address` - (Optional) Route distinguisher administrator address. If using EVPN service, then this attribute should be defined if auto generation of route distinguisher on VRF configuration is needed.
* `bgp_config` - (Optional) The BGP configuration for the Tier-0 gateway. When enabled a valid `edge_cluster_path` must be set on the Tier-0 gateway. This clause is not applicable for Global Manager - use `nsxt_policy_bgp_config` resource instead.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this Tier-0 gateway's BGP configuration.
  * `ecmp` - (Optional) A boolean flag to enable/disable ECMP. Default is `true`.
  * `enabled` - (Optional) A boolean flag to enable/disable BGP. Default is `true`.
  * `inter_sr_ibgp` - (Optional) A boolean flag to enable/disable inter SR IBGP configuration. Default is `true`. This setting is not applicable to VRF-Lite Gateway.
  * `local_as_num` - (Optional) BGP AS number in ASPLAIN/ASDOT Format. This setting is optional for VRF-Lite Gateways, and is required otherwise.
  * `multipath_relax` - (Optional) A boolean flag to enable/disable multipath relax for BGP. Default is `true`. This setting is not applicable to VRF-Lite Gateway.
  * `graceful_restart_mode` - (Optional) Setting to control BGP graceful restart mode, one of `DISABLE`, `GR_AND_HELPER`, `HELPER_ONLY`. This setting is not applicable to VRF-Lite Gateway.
  * `graceful_restart_timer` - (Optional) BGP graceful restart timer. Default is `180`. This setting is not applicable to VRF-Lite Gateway.
  * `graceful_restart_stale_route_timer` - (Optional) BGP stale route timer. Default is `600`. This setting is not applicable to VRF-Lite Gateway.
  * `route_aggregation`- (Optional) Zero or more route aggregations for BGP.
      * `prefix` - (Required) CIDR of aggregate address.
      * `summary_only` - (Optional) A boolean flag to enable/disable summarized route info. Default is `true`.
* `vrf_config` - (Optional) VRF config for VRF Tier0. This clause is supported with NSX 3.0.0 onwards.
  * `gateway_path` - (Required) Default Tier0 path. Cannot be modified after realization.
  * `evpn_transit_vni` - (Optional) L3 VNI associated with the VRF for overlay traffic. VNI must be unique and belong to configured VNI pool.
  * `route_distinguisher` - (Optional) Route distinguisher. Format: <ASN>:<number> or <IPAddress>:<number>.
  * `route_target` - (Optional) Only one target is supported.
      * `auto_mode` - (Optional) When true, import and export targets should not be specified.
      * `address_family` - (Optional) Address family, currently only `L2VPN_EVPN` is supported, which is the default.
      * `import_targets` - (Optional) List of import route targets. Format: <ASN>:<number>.
      * `export_targets` - (Optional) List of export route targets. Format: <ASN>:<number>.
* `intersite_config` - (Optional) This clause is relevant for Global Manager only.
  * `transit_subnet` - (Optional) IPv4 subnet for inter-site transit segment connecting service routers across sites for stretched gateway. For IPv6 link local subnet is auto configured.
  * `primary_site_path` - (Optional) Primary egress site for gateway.
  * `fallback_site_paths` - (Optional) Fallback sites to be used as new primary site on current primary site failure.
* `redistribution_config` - (Deprecated) Route redistribution properties. This setting is for local manager only and supported with NSXt 3.0.0 onwards. This setting is deprecated, please use `nsxt_policy_gateway_redistribution_config` resource instead.
  * `enabled` - (Optional) Enable route redistribution for BGP. Defaults to `true`.
  * `ospf_enabled` - (Optional) Enable route redistribution for OSPF. Defaults to `false`. Applicable from NSX 3.1.0 onwards.
  * `rule` - (Optional) List of redistribution rules.
      * `name` - (Optional) Rule name.
      * `route_map_path` - (Optional) Route map to be associated with the redistribution rule.
      * `types` - (Optional) List of redistribution types, possible values are: `TIER0_STATIC`, `TIER0_CONNECTED`, `TIER0_EXTERNAL_INTERFACE`, `TIER0_SEGMENT`, `TIER0_ROUTER_LINK`, `TIER0_SERVICE_INTERFACE`, `TIER0_LOOPBACK_INTERFACE`, `TIER0_DNS_FORWARDER_IP`, `TIER0_IPSEC_LOCAL_IP`, `TIER0_NAT`, `TIER0_EVPN_TEP_IP`, `TIER1_NAT`, `TIER1_STATIC`, `TIER1_LB_VIP`, `TIER1_LB_SNAT`, `TIER1_DNS_FORWARDER_IP`, `TIER1_CONNECTED`, `TIER1_SERVICE_INTERFACE`, `TIER1_SEGMENT`, `TIER1_IPSEC_LOCAL_ENDPOINT`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Tier-0 gateway.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `bgp_config` - The following attributes are exported for `bgp_config`:
  * `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
  * `path` - The NSX path of the policy resource.

## Importing

An existing policy Tier-0 gateway can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_tier0_gateway.tier0_gw ID
```

The above command imports the policy Tier-0 gateway named `tier0_gw` with the NSX Policy ID `ID`.

~> **NOTE:** When importing Gateway, `edge_cluster_path` will be assigned rather than `locale_service`. In order to switch to `locale_service` configuration, additional apply will be required.

~> **NOTE:** Redistribution config on Tier-0 resource is deprecated and thus will not be imported. Please import this configuration with `policy_gateway_redistribution_config` resource.

