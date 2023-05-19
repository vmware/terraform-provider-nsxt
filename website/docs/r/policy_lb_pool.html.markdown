---
subcategory: "Load Balancer"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_lb_pool"
description: A resource to configure Load Balancer Pool.
---

# nsxt_policy_lb_pool

This resource provides a method for the management of Load Balancer Pool.

This resource is applicable to NSX Policy Manager.
 
## Example Usage

```hcl
resource "nsxt_policy_lb_pool" "test" {
  display_name         = "test"
  description          = "Terraform provisioned LB Pool"
  algorithm            = "IP_HASH"
  min_active_members   = 2
  active_monitor_path  = "/infra/lb-monitor-profiles/default-icmp-lb-monitor"
  passive_monitor_path = "/infra/lb-monitor-profiles/default-passive-lb-monitor"
  member {
    admin_state                = "ENABLED"
    backup_member              = false
    display_name               = "member1"
    ip_address                 = "5.5.5.5"
    max_concurrent_connections = 12
    port                       = "77"
    weight                     = 1
  }
  snat {
    type = "AUTOMAP"
  }
  tcp_multiplexing_enabled = true
  tcp_multiplexing_number  = 8
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `algorithm` - (Optional) Load balancing algorithm, one of `ROUND_ROBIN`, `WEIGHTED_ROUND_ROBIN`, `LEAST_CONNECTION`, `WEIGHTED_LEAST_CONNECTION`, `IP_HASH`. Default is `ROUND_ROBIN`.
* `member_group` - (Optional) Grouping specification for pool members. When `member_group` is set, `member` should not be specified.
  * `group_path` - (Required) Path for policy group.
  * `allow_ipv4` - (Optional) Use IPv4 addresses from the grouping object, default is `true`.
  * `allow_ipv6` - (Optional) Use IPv6 addresses from the grouping object, default is `false`. Note: this setting is only supported for pools that contain IPv6 addresses.
  * `max_ip_list_size` - (Optional) Maximum number of IPs to use from the grouping object.
  * `port` - (Optional) If port is specified, all connections will be redirected to this port.
* `member`- (Optional) Members of the pool. When `member' is set, `member_group` should not be specified.
  * `ip_address` - (Required) Member IP address.
  * `admin_state` - (Optional) One of `ENABLED`, `DISABLED`, `GRACEFUL_DISABLED`. Default is `ENABLED`.
  * `backup_member` - (Optional) Whether this member is a backup member.
  * `display_name` - (Optional) Display name of the member.
  * `max_concurrent_connections` - (Optional) To ensure members are not overloaded, connections to a member can be capped by this setting.
  * `port` - (Optional) If port is specified, all connections will be redirected to this port.
  * `weight` - (Optional) Pool member weight is used for WEIGHTED algorithms.
* `min_active_members` - (Optional) A pool is considered active if there are at least certain minimum number of members.
* `active_monitor_path` - (Optional) Active monitor to be associated with this pool.
* `passive_monitor_path` - (Optional) Passive monitor to be associated with this pool.
* `snat` - (Optional) Source NAT may be required to ensure traffic from the server destined to the client is received by the load balancer.
  * `type` - (Optional) SNAT type, one of `AUTOMAP`, `DISABLED`, `IPPOOL`. Default is `AUTOMAP`.
  * `ip_pool_addresses` - (Optional) List of IP ranges or IP CIDRs to use for IPPOOL SNAT type.
* `tcp_multiplexing_enabled` - (Optional) Enable TCP multiplexing within the pool.
* `tcp_multiplexing_number` - (Optional) The maximum number of TCP connections per pool that are idly kept alive for sending future client requests.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Security Policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing pool can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_lb_pool.test ID
```

The above command imports LBPool named `test` with the NSX LBPool ID `ID`.
