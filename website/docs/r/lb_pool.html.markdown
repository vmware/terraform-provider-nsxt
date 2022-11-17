---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_pool"
description: |-
  Provides a resource to configure Load Balancer Pool on NSX-T manager
---

# nsxt_lb_pool

Provides a resource to configure Load Balancer Pool on NSX-T manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_lb_icmp_monitor" "lb_icmp_monitor" {
  display_name = "lb_icmp_monitor"
  fall_count   = 3
  interval     = 5
}

resource "nsxt_lb_passive_monitor" "lb_passive_monitor" {
  display_name = "lb_passive_monitor"
  max_fails    = 3
  timeout      = 10
}

resource "nsxt_lb_pool" "lb_pool" {
  description              = "lb_pool provisioned by Terraform"
  display_name             = "lb_pool"
  algorithm                = "WEIGHTED_ROUND_ROBIN"
  min_active_members       = 1
  tcp_multiplexing_enabled = false
  tcp_multiplexing_number  = 3
  active_monitor_id        = nsxt_lb_icmp_monitor.lb_icmp_monitor.id
  passive_monitor_id       = nsxt_lb_passive_monitor.lb_passive_monitor.id

  member {
    admin_state                = "ENABLED"
    backup_member              = "false"
    display_name               = "1st-member"
    ip_address                 = "1.1.1.1"
    max_concurrent_connections = "1"
    port                       = "87"
    weight                     = "1"
  }

  tag {
    scope = "color"
    tag   = "red"
  }
}

resource "nsxt_lb_pool" "lb_pool_with_dynamic_membership" {
  description              = "lb_pool provisioned by Terraform"
  display_name             = "dynamic_lb_pool"
  algorithm                = "LEAST_CONNECTION"
  min_active_members       = 1
  tcp_multiplexing_enabled = false
  tcp_multiplexing_number  = 3
  active_monitor_id        = nsxt_lb_icmp_monitor.lb_icmp_monitor.id
  passive_monitor_id       = nsxt_lb_passive_monitor.lb_passive_monitor.id

  snat_translation {
    type = "SNAT_IP_POOL"
    ip   = "1.1.1.1"
  }

  member_group {
    ip_version_filter  = "IPV4"
    limit_ip_list_size = true
    max_ip_list_size   = "4"
    port               = "80"

    grouping_object {
      target_type = "NSGroup"
      target_id   = nsxt_ns_group.group1.id
    }
  }

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `description` - (Optional) Description of this resource.
* `active_monitor_id` - (Optional) Active health monitor Id. If one is not set, the active healthchecks will be disabled.
* `algorithm` - (Optional) Load balancing algorithm controls how the incoming connections are distributed among the members. Supported algorithms are: ROUND_ROBIN, WEIGHTED_ROUND_ROBIN, LEAST_CONNECTION, WEIGHTED_LEAST_CONNECTION, IP_HASH.
* `member` - (Optional) Server pool consists of one or more pool members. Each pool member is identified, typically, by an IP address and a port. Each member has the following arguments:
  * `admin_state` - (Optional) Pool member admin state. Possible values: ENABLED, DISABLED and GRACEFUL_DISABLED
  * `backup_member` - (Optional) A boolean flag which reflects whether this is a backup pool member. Backup servers are typically configured with a sorry page indicating to the user that the application is currently unavailable. While the pool is active (a specified minimum number of pool members are active) BACKUP members are skipped during server selection. When the pool is inactive, incoming connections are sent to only the BACKUP member(s).
  * `display_name` - (Optional) The display name of this resource. pool member name.
  * `ip_address` - (Required) Pool member IP address.
  * `max_concurrent_connections` - (Optional) To ensure members are not overloaded, connections to a member can be capped by the load balancer. When a member reaches this limit, it is skipped during server selection. If it is not specified, it means that connections are unlimited.
  * `port` - (Optional) If port is specified, all connections will be sent to this port. Only single port is supported. If unset, the same port the client connected to will be used, it could be overrode by default_pool_member_port setting in virtual server. The port should not specified for port range case.
  * `weight` - (Optional) Pool member weight is used for WEIGHTED_ROUND_ROBIN balancing algorithm. The weight value would be ignored in other algorithms.
* `member_group` - (Optional) Dynamic pool members for the loadbalancing pool. When member group is defined, members setting should not be specified. The member_group has the following arguments:
  * `grouping_object` - (Required) Grouping object of type NSGroup which will be used as dynamic pool members. The IP list of the grouping object would be used as pool member IP setting.
  * `ip_version_filter` - (Optional) Ip version filter is used to filter IPv4 or IPv6 addresses from the grouping object. If the filter is not specified, both IPv4 and IPv6 addresses would be used as server IPs. Supported filtering is "IPV4" and "IPV6" ("IPV4" is the default one)
  * `limit_ip_list_size` - (Optional) Limits the max number of pool members. If false, allows the dynamic pool to grow up to the load balancer max pool member capacity.
  * `max_ip_list_size` - (Optional) Should only be specified if limit_ip_list_size is set to true. Limits the max number of pool members to the specified value.
  * `port` - (Optional) If port is specified, all connections will be sent to this port. If unset, the same port the client connected to will be used, it could be overridden by default_pool_member_ports setting in virtual server. The port should not specified for multiple ports case.
* `min_active_members` - (Optional) The minimum number of members for the pool to be considered active. This value is 1 by default.
* `passive_monitor_id` - (Optional) Passive health monitor Id. If one is not set, the passive healthchecks will be disabled.
* `snat_translation - (Optional) SNAT translation configuration for the pool.
  * `type` - (Optional) Type of SNAT performed to ensure reverse traffic from the server can be received and processed by the loadbalancer. Supported types are: SNAT_AUTO_MAP, SNAT_IP_POOL and TRANSPARENT
  * `ip` - (Required for snat_translation of type SNAT_IP_POOL) Ip address or Ip range for SNAT of type SNAT_IP_POOL.
* `tcp_multiplexing_enabled` - (Optional) TCP multiplexing allows the same TCP connection between load balancer and the backend server to be used for sending multiple client requests from different client TCP connections. Disabled by default.
* `tcp_multiplexing_number` - (Optional) The maximum number of TCP connections per pool that are idly kept alive for sending future client requests. The default value for this is 6.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb pool.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb pool.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb pool can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_pool.lb_pool UUID
```

The above would import the lb pool named `lb_pool` with the nsx id `UUID`
