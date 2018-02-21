---
layout: "nsxt"
page_title: "NSXT: nsxt_logical_tier1 router"
sidebar_current: "docs-nsxt-resource-logical-tier1 router"
description: |-
  Provides a resource to configure Logical Tier1 router on NSX-T Manager.
---

# nsxt_logical_tier1_router

Provides a resource to configure Logical Tier1 router on NSX-T Manager.

## Example Usage

```hcl
resource "nsxt_logical_tier1_router" "tier1_router" {
    description = "RTR1 provisioned by Terraform"
    display_name = "RTR1"
    failover_mode =  "PREEMPTIVE"
    high_availability_mode = "ACTIVE_STANDBY"
    edge_cluster_id = "${data.nsxt_edge_cluster.edge_cluster.id}"
    enable_router_advertisement = true
    advertise_connected_routes = false
    advertise_static_routes = true
    advertise_nat_routes = true
    tag {
        scope = "color"
        tag = "blue"
    }
}
```

## Argument Reference

The following arguments are supported:

* `edge_cluster_id` - (Required) Edge Cluster ID for the logical Tier1 router.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical Tier1 router.
* `failover_mode` - (Optional) This failover mode determines, whether the preferred service router instance for given logical router will preempt the peer. Note - It can be specified if and only if logical router is ACTIVE_STANDBY and NON_PREEMPTIVE mode is supported only for a Tier1 logical router. For ACTIVE_ACTIVE logical routers, this field must not be populated
* `high_availability_mode` - (Optional) High availability mode
* `enable_router_advertisement` - (Optional) Enable the router advertisement
* `advertise_connected_routes` - (Optional) Enable the router advertisement for all NSX connected routes
* `advertise_static_routes` - (Optional) Enable the router advertisement for static routes
* `advertise_nat_routes` - (Optional) Enable the router advertisement for NAT routes

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical Tier1 router.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `advertise_config_revision` - Indicates current revision number of the advertisement configuration object as seen by NSX-T API server. This attribute can be useful for debugging.
* `firewall_sections` - (Optional) The list of firewall sections for this router

## Importing

An existing logical tier1 router can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import nsxt_logical_tier1_router.x id
```

The above would import the logical tier1 router named `x` with the nsx id `id`
