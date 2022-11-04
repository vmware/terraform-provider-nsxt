---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_service"
description: |-
  Provides a resource to configure lb service on NSX-T manager
---

# nsxt_lb_service

Provides a resource to configure lb service on NSX-T manager.
Note that lb service needs to be attached to Tier-1 router that satisfies
following preconditions:
* It needs to reside on edge cluster
* It needs to be condigured with either uplink port or centralized service port

In order to enforce correct order of create/delete, it is recommended to add
depends_on clause to lb service.

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl

data "nsxt_edge_cluster" "EC" {
  display_name = "%s"
}

data "nsxt_logical_tier0_router" "test" {
  display_name = "%s"
}

resource "nsxt_logical_router_link_port_on_tier0" "test" {
  display_name      = "port_on_tier0"
  logical_router_id = data.nsxt_logical_tier0_router.test.id
}

resource "nsxt_logical_tier1_router" "test" {
  display_name    = "test"
  edge_cluster_id = data.nsxt_edge_cluster.EC.id
}

resource "nsxt_logical_router_link_port_on_tier1" "test" {
  logical_router_id             = nsxt_logical_tier1_router.test.id
  linked_logical_router_port_id = nsxt_logical_router_link_port_on_tier0.test.id
}

resource "nsxt_lb_service" "lb_service" {
  description  = "lb_service provisioned by Terraform"
  display_name = "lb_service"

  tag {
    scope = "color"
    tag   = "red"
  }

  enabled           = true
  logical_router_id = nsxt_logical_tier1_router.test.id
  error_log_level   = "INFO"
  size              = "MEDIUM"

  depends_on = [nsxt_logical_router_link_port_on_tier1.test]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb service.
* `logical_router_id` - (Required) Tier1 logical router this service is attached to. Note that this router needs to have edge cluster configured, and have an uplink port or CSP (centralized service port).
* `enabled` - (Optional) whether the load balancer service is enabled.
* `error_log_level` - (Optional) Load balancer engine writes information about encountered issues of different severity levels to the error log. This setting is used to define the severity level of the error log.
* `size` - (Required) Size of load balancer service. Accepted values are SMALL/MEDIUM/LARGE.
* `virtual_server_ids` - (Optional) Virtual servers associated with this Load Balancer.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb_service.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_service.lb_service UUID
```

The above would import the lb service named `lb_service` with the nsx id `UUID`
