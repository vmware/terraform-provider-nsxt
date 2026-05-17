---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_transit_gateway"
description: A resource to configure a Transit Gateway.
---

# nsxt_policy_transit_gateway

This resource provides a method for the management of a Transit Gateway.

This resource is applicable to NSX Policy Manager.

In NSX 9.0, users cannot create TGWs. Terraform users can import and update the default TGW.
In NSX 9.1 onwards, users can create and delete TGWs

## Example Usage

```hcl
resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = "dev"
  }

  display_name    = "test"
  description     = "Terraform provisioned TransitGateway"
  transit_subnets = ["10.203.4.0/24"]

  redistribution_config {
    rule {
      types = ["PUBLIC", "TGW_STATIC_ROUTE"]
    }
    rule {
      types          = ["TGW_PRIVATE"]
      route_map_path = nsxt_policy_transit_gateway_route_map.example.path
    }
  }

  bgp_config {
    local_as_num = "65001"
    enabled      = true
    ecmp         = true

    graceful_restart_mode              = "HELPER_ONLY"
    graceful_restart_timer             = 180
    graceful_restart_stale_route_timer = 600
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `transit_subnets` - (Optional) Array of IPV4 CIDRs for internal VPC attachment networks.
* `centralized_config` - (Optional) Singleton block for high-availability and edge cluster for centralized connectivity (gateway connections, VPN). Sent in the same H-API transaction as the transit gateway (as a child object, like security policy rules). Available since NSX 9.1.0.
    * `failover_mode` - (Optional) Failover mode for active/standby centralized transit gateway (`PREEMPTIVE` or `NON_PREEMPTIVE`). Available since NSX 9.2.0.
    * `ha_mode` - (Optional) High-availability mode. Values: `ACTIVE_ACTIVE`, `ACTIVE_STANDBY`. Default is `ACTIVE_ACTIVE`.
    * `edge_cluster_paths` - (Optional) Policy paths of edge clusters. At most one item. Must be authorized for the project.
* `advanced_config` - (Optional) Advanced settings for the transit gateway. Applied via hierarchical API as the `TransitGatewayRoutingConfig` singleton child (same transaction as the transit gateway when possible). Read uses the routing-config GET API.
    * `forwarding_up_timer` - (Optional/Computed) Time in seconds the transit gateway waits before marking forwarding as active after a failover event. Valid range is `0` to `100`. Default is `5`.
* `redistribution_config` - (Optional) Route redistribution configuration. Applied as part of the `TransitGatewayRoutingConfig` singleton (same API object as `advanced_config`).
    * `rule` - (Required) List of redistribution rules. Minimum 1, maximum 5.
        * `types` - (Required) List of route types to redistribute. Supported values: `PUBLIC` (public networks from connected VPC subnets, NAT rules, and VPC static routes), `TGW_PRIVATE` (transit gateway private networks; only redistributed to BGP neighbors whose centralized network attachment has `allow_private` set to `true`), `TGW_STATIC_ROUTE` (user-configured static routes on the transit gateway), `CONNECTED_SUBNET` (centralized network attachment interface subnets).
        * `route_map_path` - (Optional) Policy path of a route map to apply for filtering routes in this rule. The route map must exist under the same transit gateway.
* `bgp_config` - (Optional/Computed) BGP routing configuration. Applied in the same H-API transaction as the transit gateway (as a `ChildBgpRoutingConfig` child object). Read uses the BGP GET API (`/orgs/{org}/projects/{proj}/transit-gateways/{tgw}/bgp`).
    * `revision` - (Computed) Revision of the BGP config object.
    * `path` - (Computed) NSX policy path of the BGP config.
    * `tag` - (Optional) A list of scope + tag pairs.
    * `enabled` - (Optional) Flag to enable BGP. Default is `true`.
    * `ecmp` - (Optional) Flag to enable ECMP. Default is `true`.
    * `inter_sr_ibgp` - (Optional/Computed) Enable inter SR IBGP configuration.
    * `local_as_num` - (Optional/Computed) BGP AS number in ASPLAIN or ASDOT format (e.g. `"65001"` or `"1.10"`).
    * `multipath_relax` - (Optional/Computed) Flag to enable BGP multipath relax option.
    * `graceful_restart_mode` - (Optional) BGP graceful restart mode. Values: `DISABLE`, `GR_AND_HELPER`, `HELPER_ONLY`. Default is `HELPER_ONLY`.
    * `graceful_restart_timer` - (Optional) BGP graceful restart timer in seconds. Valid range: 1–3600. Default is `180`.
    * `graceful_restart_stale_route_timer` - (Optional) BGP stale route timer in seconds. Valid range: 1–3600. Default is `600`.
    * `route_aggregation` - (Optional) List of routes to be aggregated (up to 1000 entries).
        * `prefix` - (Optional) CIDR of the aggregate address.
        * `summary_only` - (Optional) Send only the summarized route. Default is `true`.
* `span` - (Optional) Span configuration. Note that one of `cluster_based_span` and `zone_based_span` is required. Available since NSX 9.1.0.
    * `cluster_based_span` - (Optional) Span based on vSphere Clusters.
        * `span_path` - (Required) Policy path of the network span object.
    * `zone_based_span` - (Optional) Span based on zones.
        * `zone_external_ids` - (Optional) An array of Zone object's external IDs.
        * `use_all_zones` - (Optional) Flag to indicate that TransitGateway is associated with all project zones.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_transit_gateway.test PATH
```

The above command imports Transit Gateway named `test` with the NSX path `PATH`.
