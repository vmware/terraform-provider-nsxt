---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_bgp_config"
sidebar_current: "docs-nsxt-resource-policy-bgp-config"
description: A resource to configure BGP Settings of Tier0 Gateway on NSX Global Manager.
---

# nsxt_policy_bgp_config

This resource provides a method for the management of BGP for T0 Gateway on specific site. This resource is only supported on NSX Global Manager. A single resource should be specified per T0 Gateway + Site.
## Example Usage

```hcl
resource "nsxt_policy_bgp_config" "gw1-paris" {
  site_path    = data.nsxt_policy_site.paris.path
  gateway_path = nsxt_policy_tier0_gateway.gw1.path

  enabled                = true
  inter_sr_ibgp          = true
  local_as_num           = 60001
  graceful_restart_mode  = HELPER_ONLY
  graceful_restart_timer = 2400

  route_aggregation {
    prefix       = "20.1.0.0/24"
    summary_only = false
  }
}
```

~> **NOTE:** Since BGP config is auto-created on NSX, this resource will update NSX object, but never create or delete it.

## Argument Reference

The following arguments are supported:

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
  * `tag` - (Optional) A list of scope + tag pairs to associate with this Tier-0 gateway's BGP configuration.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource. This path should be used as `bgp_path` in `nsxt_policy_bgp_neighbor` resource configuration.

## Importing

Since BGP config is autocreated by the backend, and terraform create is de-facto an update, importing the resource is not useful and thus not supported.
