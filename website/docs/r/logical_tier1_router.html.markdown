---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_logical_tier1 router"
description: A resource to configure a logical Tier1 router in NSX.
---

# nsxt_logical_tier1_router

This resource provides a method for the management of a tier 1 logical router. A tier 1 logical router is often used for tenants, users and applications. There can be many tier 1 logical routers connected to a common tier 0 provider router.

## Example Usage

```hcl
resource "nsxt_logical_tier1_router" "tier1_router" {
  description                 = "RTR1 provisioned by Terraform"
  display_name                = "RTR1"
  failover_mode               = "PREEMPTIVE"
  edge_cluster_id             = data.nsxt_edge_cluster.edge_cluster.id
  enable_router_advertisement = true
  advertise_connected_routes  = false
  advertise_static_routes     = true
  advertise_nat_routes        = true
  advertise_lb_vip_routes     = true
  advertise_lb_snat_ip_routes = false

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `edge_cluster_id` - (Optional) Edge Cluster ID for the logical Tier1 router.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical Tier1 router.
* `failover_mode` - (Optional) This failover mode determines, whether the preferred service router instance for given logical router will preempt the peer. Note - It can be specified if and only if logical router is ACTIVE_STANDBY and NON_PREEMPTIVE mode is supported only for a Tier1 logical router. For ACTIVE_ACTIVE logical routers, this field must not be populated
* `enable_router_advertisement` - (Optional) Enable the router advertisement
* `advertise_connected_routes` - (Optional) Enable the router advertisement for all NSX connected routes
* `advertise_static_routes` - (Optional) Enable the router advertisement for static routes
* `advertise_nat_routes` - (Optional) Enable the router advertisement for NAT routes
* `advertise_lb_vip_routes` - (Optional) Enable the router advertisement for LB VIP routes
* `advertise_lb_snat_ip_routes` - (Optional) Enable the router advertisement for LB SNAT IP routes

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical Tier1 router.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `advertise_config_revision` - Indicates current revision number of the advertisement configuration object as seen by NSX-T API server. This attribute can be useful for debugging.
* `firewall_sections` - (Optional) The list of firewall sections for this router

## Importing

An existing logical tier1 router can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_logical_tier1_router.tier1_router UUID
```

The above command imports the logical tier 1 router named `tier1_router` with the NSX id `UUID`.
