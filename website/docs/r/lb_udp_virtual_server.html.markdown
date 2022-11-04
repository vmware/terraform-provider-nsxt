---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_udp_virtual_server"
description: |-
  Provides a resource to configure Load Balancer UDP Virtual Server on NSX-T manager
---

# nsxt_lb_udp_virtual_server

Provides a resource to configure Load Balancer UDP Virtual Server on NSX-T manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_lb_fast_udp_application_profile" "timeout_60" {
  idle_timeout = 60
}

resource "nsxt_lb_source_ip_persistence_profile" "ip_profile" {
  display_name = "source1"
}

resource "nsxt_lb_pool" "pool1" {
  algorithm = "LEAST_CONNECTION"
  member {
    ip_address = "3.0.0.1"
    port       = "443"
  }
  member {
    ip_address = "3.0.0.2"
    port       = "443"
  }
}

resource "nsxt_lb_pool" "sorry_pool" {
  member {
    ip_address = "3.0.0.15"
    port       = "443"
  }
}

resource "nsxt_lb_udp_virtual_server" "lb_virtual_server" {
  description                = "lb_virtual_server provisioned by terraform"
  display_name               = "virtual server 1"
  access_log_enabled         = true
  application_profile_id     = nsxt_lb_fast_udp_application_profile.timeout_60.id
  enabled                    = true
  ip_address                 = "10.0.0.2"
  ports                      = ["443"]
  default_pool_member_ports  = ["8888"]
  max_concurrent_connections = 50
  max_new_connection_rate    = 20
  persistence_profile_id     = nsxt_lb_source_ip_persistence_profile.ip_profile.id
  pool_id                    = nsxt_lb_pool.pool1.id
  sorry_pool_id              = nsxt_lb_pool.sorry_pool.id

  tag {
    scope = "color"
    tag   = "green"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `enabled` - (Optional) Whether the virtual server is enabled. Default is true.
* `ip_address` - (Required) Virtual server IP address.
* `ports` - (Required) List of virtual server port.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb udp virtual server.
* `access_log_enabled` - (Optional) Whether access log is enabled. Default is false.
* `application_profile_id` - (Required) The application profile defines the application protocol characteristics.
* `default_pool_member_ports` - (Optional) List of default pool member ports.
* `max_concurrent_connections` - (Optional) To ensure one virtual server does not over consume resources, affecting other applications hosted on the same LBS, connections to a virtual server can be capped. If it is not specified, it means that connections are unlimited.
* `max_new_connection_rate` - (Optional) To ensure one virtual server does not over consume resources, connections to a member can be rate limited. If it is not specified, it means that connection rate is unlimited.
* `persistence_profile_id` - (Optional) Persistence profile is used to allow related client connections to be sent to the same backend server. Only source ip persistence profile is accepted.
* `pool_id` - (Optional) Pool of backend servers. Server pool consists of one or more servers, also referred to as pool members, that are similarly configured and are running the same application.
* `sorry_pool_id` - (Optional) When load balancer can not select a backend server to serve the request in default pool or pool in rules, the request would be served by sorry server pool.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb udp virtual server.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb udp virtual server can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_udp_virtual_server.lb_udp_virtual_server UUID
```

The above would import the lb udp virtual server named `lb_udp_virtual_server` with the nsx id `UUID`
