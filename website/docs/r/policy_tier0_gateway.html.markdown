---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_tier0_gateway"
sidebar_current: "docs-nsxt-resource-policy-tier0-gateway"
description: A resource to configure a Tier-0 gateway in NSX Policy manager.
---

# nsxt_policy_tier0_gateway

This resource provides a method for the management of a Tier-0 gateway.

## Example Usage

```hcl
resource "nsxt_policy_tier0_gateway" "tier0_gw" {
  description               = "Tier-0 provisioned by Terraform"
  display_name              = "Tier0-gw1"
  nsx_id                    = "predefined_id"
  failover_mode             = "PREEMPTIVE"
  default_rule_logging      = false
  enable_firewall           = true
  force_whitelisting        = false
  ha_mode                   = "ACTIVE_STANDBY"
  internal_transit_subnets  = ["102.64.0.0/16"]
  transit_subnets           = ["101.64.0.0/16"]
  edge_cluster_path         = data.nsxt_policy_edge_cluster.EC.path

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

  vrf_config {
    gateway_path       = data.nsxt_policy_tier0_gateway.vrf.path
    route_distinuisher = "62000:10"
    evpn_transit_vni   = 76001
    route_target {
      auto_mode      = false
      import_targets = ["62000:2"]
      export_targets = ["62000:3", "10.2.2.0:3"]
    }
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
* `edge_cluster_path` - (Optional) The path of the edge cluster where the Tier-0 is placed. Must be specified when `bgp_config` is enabled.
* `failover_mode` - (Optional) This failover mode determines, whether the preferred service router instance for given logical router will preempt the peer. Accepted values are PREEMPTIVE/NON_PREEMPTIVE.
* `default_rule_logging` - (Optional) Boolean flag indicating if the default rule logging will be enabled or not. The default value is false.
* `enable_firewall` - (Optional) Boolean flag indicating if the edge firewall will be enabled or not. The default value is true.
* `force_whitelisting` - (Optional) Boolean flag indicating if white-listing will be forced or not. The default value is false.
* `ipv6_ndra_profile_path` - (Optional) Policy path to IPv6 NDRA profile.
* `ipv6_dad_profile_path` - (Optional) Policy path to IPv6 DAD profile.
* `ha_mode` - (Optional) High-availability Mode for Tier-0. Valid values are `ACTIVE_ACTIVE` and `ACTIVE_STANDBY`.
* `internal_transit_subnets` - (Optional) Internal transit subnets in CIDR format. At most 1 CIDR.
* `transit_subnets` - (Optional) Transit subnets in CIDR format.
* `dhcp_config_path` - (Optional) Policy path to DHCP server or relay configuration to use for this gateway.
* `bgp_config` - (Optional) The BGP configuration for the Tier-0 gateway. When enabled a valid `edge_cluster_path` must be set on the Tier-0 gateway.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this Tier-0 gateway's BGP configuration.
  * `ecmp` - (Optional) A boolean flag to enable/disable ECMP. Default is `true`.
  * `enabled` - (Optional) A boolean flag to enable/disable BGP. Default is `true`.
  * `inter_sr_ibgp` - (Optional) A boolean flag to enable/disable inter SR IBGP configuration. Default is `true`.
  * `local_as_num` - (Optional) BGP AS number in ASPLAIN/ASDOT Format. Default is `65000`.
  * `multipath_relax` - (Optional) A boolean flag to enable/disable multipath relax for BGP. Default is `true`.
  * `graceful_restart_mode` - (Optional) Setting to control BGP graceful restart mode, one of `DISABLE`, `GR_AND_HELPER`, `HELPER_ONLY`.
  * `graceful_restart_timer` - (Optional) BGP graceful restart timer. Default is `180`.
  * `graceful_restart_stale_route_timer` - (Optional) BGP stale route timer. Default is `600`.
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

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_tier0_gateway.tier0_gw ID
```

The above command imports the policy Tier-0 gateway named `tier0_gw` with the NSX Policy ID `ID`.
