---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_tier1_gateway"
sidebar_current: "docs-nsxt-resource-policy-tier1-gateway"
description: A resource to configure a Tier-1 gateway in NSX Policy manager.
---

# nsxt_policy_tier1_gateway

This resource provides a method for the management of a Tier-1 gateway. A Tier-1 gateway is often used for tenants, users and applications. There can be many Tier-1 gateways connected to a common Tier-0 provider gateway.

## Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "T0" {
  display_name = "T0"
}

data "nsxt_policy_edge_cluster" "EC" {
  display_name = "EC"
}

resource "nsxt_policy_tier1_gateway" "tier1_gw" {
  description               = "Tier-1 provisioned by Terraform"
  display_name              = "Tier1-gw1"
  nsx_id                    = "predefined_id"
  edge_cluster_path         = data.nsxt_policy_edge_cluster.EC.path
  failover_mode             = "PREEMPTIVE"
  default_rule_logging      = "false"
  enable_firewall           = "true"
  enable_standby_relocation = "false"
  force_whitelisting        = "true"
  tier0_path                = data.nsxt_policy_tier0_gateway.T0.path
  route_advertisement_types = ["TIER1_STATIC_ROUTES", "TIER1_CONNECTED"]
  pool_allocation           = "ROUTING"

  tag {
    scope = "color"
    tag   = "blue"
  }

  route_advertisement_rule {
    name                      = "rule1"
    action                    = "DENY"
    subnets                   = ["20.0.0.0/24", "21.0.0.0/24"]
    prefix_operator           = "GE"
    route_advertisement_types = ["TIER1_CONNECTED"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Tier-1 gateway.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `edge_cluster_path` - (Optional) The path of the edge cluster where the Tier-1 is placed.
* `failover_mode` - (Optional) This failover mode determines, whether the preferred service router instance for given logical router will preempt the peer. Accepted values are PREEMPTIVE/NON_PREEMPTIVE.
* `default_rule_logging` - (Optional) Boolean flag indicating if the default rule logging will be enabled or not. The default value is false.
* `enable_firewall` - (Optional) Boolean flag indicating if the edge firewall will be enabled or not. The default value is true.
* `enable_standby_relocation` - (Optional) Boolean flag indicating if the standby relocation will be enabled or not. The default value is false.
* `force_whitelisting` - (Optional) Boolean flag indicating if white-listing will be forced or not. The default value is false. Setting it to `true` will create a base deny entry rule on Tier-1 firewall.
* `tier0_path` - (Optional) The path of the connected Tier0.
* `route_advertisement_types` - (Optional) Enable different types of route advertisements: TIER1_STATIC_ROUTES, TIER1_CONNECTED, TIER1_NAT, TIER1_LB_VIP, TIER1_LB_SNAT, TIER1_DNS_FORWARDER_IP, TIER1_IPSEC_LOCAL_ENDPOINT.
* `ipv6_ndra_profile_path` - (Optional) Policy path to IPv6 NDRA profile.
* `ipv6_dad_profile_path` - (Optional) Policy path to IPv6 DAD profile.
* `dhcp_config_path` - (Optional) Policy path to DHCP server or relay configuration to use for this gateway.
* `pool_allocation` - (Optional) Size of edge node allocation at for routing and load balancer service to meet performance and scalability requirements, one of `ROUTING`, `LB_SMALL`, `LB_MEDIUM`, `LB_LARGE`, `LB_XLARGE`. Default is `ROUTING`. Changing this attribute would force new resource.
* `route_advertisement_rule` - (Optional) List of rules for routes advertisement:
  * `name` - (Required) The name of the rule.
  * `action` - (Required) Action to advertise filtered routes to the connected Tier0 gateway. PERMIT (which is the default): Enables the advertisement, DENY: Disables the advertisement.
  * `subnets` - (Required) list of network CIDRs to be routed.
  * `prefix_operator` - (Optional) Prefix operator to apply on subnets. GE prefix operator (which is the default|) filters all the routes having network subset of any of the networks configured in Advertise rule. EQ prefix operator filter all the routes having network equal to any of the network configured in Advertise rule.The name of the rule.
* `route_advertisement_types` - (Optional) List of desired types of route advertisements, supported values: `TIER1_STATIC_ROUTES`, `TIER1_CONNECTED`, `TIER1_NAT`, `TIER1_LB_VIP`, `TIER1_LB_SNAT`, `TIER1_DNS_FORWARDER_IP`, `TIER1_IPSEC_LOCAL_ENDPOINT`.
* `ingress_qos_profile_path` - (Optional) QoS Profile path for ingress traffic on link connected to Tier0 gateway.
* `egress_qos_profile_path` - (Optional) QoS Profile path for egress traffic on link connected to Tier0 gateway.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Tier-1 gateway.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing policy Tier-1 gateway can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_tier1_gateway.tier1_gw ID
```

The above command imports the policy Tier-1 gateway named `tier1_gw` with the NSX Policy ID `ID`.
