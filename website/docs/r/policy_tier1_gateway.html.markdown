---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_tier1_gateway"
description: A resource to configure a Tier-1 gateway on NSX Policy manager.
---

# nsxt_policy_tier1_gateway

This resource provides a method for the management of a Tier-1 Gateway. A Tier-1 Gateway is often used for tenants, users and applications. There can be many Tier-1 gateways connected to a common Tier-0 provider gateway.

This resource is applicable to NSX Global Manager and NSX Policy Manager.

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

## Global manager example usage

```hcl
resource "nsxt_policy_tier1_gateway" "tier1_gw" {
  description  = "Tier-1 provisioned by Terraform"
  display_name = "Tier1-gw1"

  locale_service {
    nsx_id            = "paris"
    edge_cluster_path = data.nsxt_policy_edge_cluster.paris.path
  }

  locale_service {
    nsx_id               = "london"
    edge_cluster_path    = data.nsxt_policy_edge_cluster.london.path
    preferred_edge_paths = [data.nsxt_policy_egde_node.edge1.path]
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

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_tier1_gateway" "tier1_gw" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  description               = "Tier-1 provisioned by Terraform"
  display_name              = "Tier1-gw1"
  nsx_id                    = "predefined_id"
  edge_cluster_path         = data.nsxt_policy_project.demoproj.site_info.0.edge_cluster_paths.0
  failover_mode             = "PREEMPTIVE"
  default_rule_logging      = "false"
  enable_firewall           = "true"
  enable_standby_relocation = "false"
  tier0_path                = data.nsxt_policy_project.demoproj.tier0_gateway_paths.0
  route_advertisement_types = ["TIER1_STATIC_ROUTES", "TIER1_CONNECTED"]
  pool_allocation           = "ROUTING"

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
* `tag` - (Optional) A list of scope + tag pairs to associate with this Tier-1 gateway.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `edge_cluster_path` - (Optional) The path of the edge cluster where the Tier-1 is placed.For advanced configuration, use `locale_service` clause instead. 
* `locale_service` - (Optional) This argument is required on NSX Global Manager. Multiple locale services can be specified for multiple locations.
  * `nsx_id` - (Optional) NSX id for the locale service. It is recommended to specify this attribute in order to avoid unnecessary recreation of this object. Should be unique within the gateway.
  * `edge_cluster_path` - (Required) The path of the edge cluster where the Tier-0 is placed.
  * `preferred_edge_paths` - (Optional) Policy paths to edge nodes. Specified edge is used as preferred edge cluster member when failover mode is set to `PREEMPTIVE`.
  * `display_name` - (Optional) Display name for the locale service.
* `failover_mode` - (Optional) This failover mode determines, whether the preferred service router instance for given logical router will preempt the peer. Accepted values are PREEMPTIVE/NON_PREEMPTIVE.
* `default_rule_logging` - (Optional) Boolean flag indicating if the default rule logging will be enabled or not. The default value is false.
* `enable_firewall` - (Optional) Boolean flag indicating if the edge firewall will be enabled or not. The default value is true.
* `enable_standby_relocation` - (Optional) Boolean flag indicating if the standby relocation will be enabled or not. The default value is false.
* `force_whitelisting` - (Deprecated) Boolean flag indicating if white-listing will be forced or not. The default value is false. This argument is deprecated and will be removed. Please use `nsxt_policy_predefined_gateway_policy` resource to control default action.
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
* `route_advertisement_types` - (Optional) List of desired types of route advertisements, supported values: `TIER1_STATIC_ROUTES`, `TIER1_CONNECTED`, `TIER1_NAT`, `TIER1_LB_VIP`, `TIER1_LB_SNAT`, `TIER1_DNS_FORWARDER_IP`, `TIER1_IPSEC_LOCAL_ENDPOINT`. This field is Computed, meaning that NSX can auto-assign types. Hence, in order to revert to default behavior, set route advertisement values explicitly rather than removing this clause from configuration.
* `ingress_qos_profile_path` - (Optional) QoS Profile path for ingress traffic on link connected to Tier0 gateway.
* `egress_qos_profile_path` - (Optional) QoS Profile path for egress traffic on link connected to Tier0 gateway.
* `intersite_config` - (Optional) This clause is relevant for Global Manager only.
  * `transit_subnet` - (Optional) IPv4 subnet for inter-site transit segment connecting service routers across sites for stretched gateway. For IPv6 link local subnet is auto configured.
  * `primary_site_path` - (Optional) Primary egress site for gateway.
* `ha_mode` - (Optional) High-availability Mode for Tier-1. Valid values are `ACTIVE_ACTIVE`, `ACTIVE_STANDBY` and `NONE`.  `ACTIVE_ACTIVE` is supported with NSX version 4.0.0 and above. `NONE` mode should be used for Distributed Only, e.g when a gateway is created and has no services.
* `type` - (Optional) This setting is only applicable to VMC and it helps auto-configure router advertisements for the gateway. Valid values are `ROUTED`, `NATTED` and `ISOLATED`. For `ROUTED` and `NATTED`, `tier0_path` should be specified in configuration.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Tier-1 gateway.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing policy Tier-1 gateway can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_tier1_gateway.tier1_gw ID
```
The above command imports the policy Tier-1 gateway named `tier1_gw` with the NSX Policy ID `ID`.

```
terraform import nsxt_policy_tier1_gateway.tier1_gw POLICY_PATH
```
The above command imports the policy Tier-1 gateway named `tier1_gw` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.

~> **NOTE:** When importing Gateway, `edge_cluster_path` will be assigned rather than `locale_service`. In order to switch to `locale_service` configuration, additional apply will be required.
