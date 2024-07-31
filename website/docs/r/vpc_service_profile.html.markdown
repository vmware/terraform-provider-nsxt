---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_service_profile"
description: A resource to configure a VPC Service Profile.
---

# nsxt_policy_service_profile

This resource provides a method for the management of a VPC Service Profile.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_service_profile" "vpc1_service_profile" {
  display_name = "vpc1"
  description  = "Terraform provisioned Vpc Service Profile"

  mac_discovery_profile = nsxt_policy_mac_discovery_profile.for_vpc1.path
  spoof_guard_profile   = nsxt_policy_spoof_guard_profile.for_vpc1.path
  ip_discovery_profile  = nsxt_policy_ip_discovery_profile.for_vpc1.path
  qos_profile           = nsxt_policy_qos_profile.for_vpc1.path

  dhcp_config {
    ntp_servers = ["20.2.60.5"]

    lease_time = 50840
    mode       = "DHCP_IP_ALLOCATION_BY_PORT"

    dns_client_config {
      dns_server_ips = ["10.204.2.20"]
    }
  }

  dns_forwarder_config {
    cache_size = 1024
    log_level  = "WARNING"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `mac_discovery_profile` - (Optional) Policy path for Mac Discovery Profile
* `spoof_guard_profile` - (Optional) Policy path for Spoof Guard Profile
* `ip_discovery_profile` - (Optional) Policy path for IP Discovery Profile
* `security_profile` - (Optional) Policy path for Security Profile
* `qos_profile` - (Optional) Policy path for QoS profile
* `dhcp_config` - (Required) DHCP configuration for this profile
  * `ntp_servers` - (Optional) List of NTP servers
  * `dns_client_config` - (Optional) DNS Client configuration
    * `dns_server_ips` - (Optional) List of IP addresses of the DNS servers which need to be configured on the workload VMs
  * `lease_time` - (Optional) DHCP lease time in seconds.

  * `mode` - (Optional) DHCP mode of the VPC Profile DHCP Config. Possible values are `DHCP_IP_ALLOCATION_BY_PORT`, `DHCP_IP_ALLOCATION_BY_MAC`, `DHCP_RELAY`, `DHCP_DEACTIVATED`. Default is `DHCP_IP_ALLOCATION_BY_PORT`.
  * `dhcp_relay_config` - (Optional) DHCP Relay configuration
    * `server_addresses` - (Optional) List of DHCP server IP addresses for DHCP relay configuration. Both IPv4 and IPv6 addresses are supported.
* `dns_forwarder_config` - (Optional) DNS Forwarder configuration
  * `cache_size` - (Optional) Cache size in KB
  * `log_level` - (Optional) Log level of the DNS forwarder. Possible values: `DEBUG`, `INFO`, `WARNING`, `ERROR`, `FATAL`
  * `conditional_forwarder_zone_paths` - (Optional) Path of conditional DNS zones
  * `default_forwarder_zone_path` - (Optional) Path of the default DNS zone


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_service_profile.test UUID
```

The above command imports VpcServiceProfile named `test` with the NSX ID `UUID`.
