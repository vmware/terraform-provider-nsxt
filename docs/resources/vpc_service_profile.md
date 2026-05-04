---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_service_profile"
description: A resource to configure a VPC Service Profile.
---

# nsxt_vpc_service_profile

This resource provides a method for the management of a VPC Service Profile.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards. IPv6-related arguments (`dhcpv6_config`, `dns_forwarder_config`, `ipv6_profile_paths`, `service_subnet_cidrs`) require NSX 9.2.0 or later.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_vpc_service_profile" "vpc1_service_profile" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }

  display_name = "vpc1"
  description  = "Terraform provisioned Vpc Service Profile"

  mac_discovery_profile = nsxt_policy_mac_discovery_profile.for_vpc1.path
  spoof_guard_profile   = nsxt_policy_spoof_guard_profile.for_vpc1.path
  ip_discovery_profile  = nsxt_policy_ip_discovery_profile.for_vpc1.path
  qos_profile           = nsxt_policy_qos_profile.for_vpc1.path

  dhcp_config {
    dhcp_server_config {
      ntp_servers = ["20.2.60.5"]
      lease_time  = 50840

      dns_client_config {
        dns_server_ips = ["10.204.2.20"]
      }

      advanced_config {
        is_distributed_dhcp = false
      }
    }
  }

  dhcpv6_config {
    dhcpv6_server_config {
      lease_time     = 86400
      preferred_time = 69120
      ntp_servers    = ["2001:db8::1"]
      sntp_servers   = ["2001:db8::2"]
      dns_client_config {
        dns_server_ips = ["2001:db8::53"]
      }
      advanced_config {
        is_distributed_dhcp = false
      }
    }
  }

  # Up to two paths: one IPv6 DAD profile and one NDRA profile (see data sources nsxt_policy_ipv6_dad_profile / nsxt_policy_ipv6_ndra_profile).
  ipv6_profile_paths = ["/orgs/default/projects/myproj/infra/ipv6-dad-profiles/dad1", "/orgs/default/projects/myproj/infra/ipv6-ndra-profiles/ndra1"]

  service_subnet_cidrs = ["fd12:3456::/64"]

  dns_forwarder_config {
    log_level = "INFO"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `mac_discovery_profile` - (Optional) Policy path for Mac Discovery Profile
* `spoof_guard_profile` - (Optional) Policy path for Spoof Guard Profile
* `ip_discovery_profile` - (Optional) Policy path for IP Discovery Profile
* `security_profile` - (Optional) Policy path for Security Profile
* `qos_profile` - (Optional) Policy path for QoS profile
* `dhcp_config` - (Required) DHCP configuration for this profile
    * `dhcp_server_config` - (Optional / Computed) DHCP server configuration for this profile. NSX may supply defaults when this block is omitted or only partially specified.
        * `ntp_servers` - (Optional) List of NTP servers
        * `dns_client_config` - (Optional) DNS Client configuration
            * `dns_server_ips` - (Optional) List of IP addresses of the DNS servers which need to be configured on the workload VMs
        * `lease_time` - (Optional) DHCP lease time in seconds.
        * `advanced_config` - (Optional) VPC DHCP advanced configuration
            * `is_distributed_dhcp` - DHCP server's IP allocation model based on workloads subnet port id. Can be `false` only when Edge cluster is available, in
        which case edge cluster in VPC connectivity profile must be configured. This is the traditional DHCP server that dynamically allocates IP per VM's MAC.
        If value is `true`, edge cluster will not be required. This is a DHCP server that dynamically assigns IP per VM port.
    * `dhcp_relay_config` - (Optional) DHCP Relay configuration
        * `server_addresses` - (Optional) List of DHCP server IP addresses for DHCP relay configuration. Both IPv4 and IPv6 addresses are supported.
* `dhcpv6_config` - (Optional) DHCPv6 profile template for VPC subnets (relay and/or server settings). Supported with NSX v9.2.0 and above.
    * `dhcpv6_relay_config` - (Optional) DHCPv6 relay configuration
        * `server_addresses` - (Optional) List of DHCPv6 server addresses for relay
    * `dhcpv6_server_config` - (Optional) DHCPv6 server defaults for subnets using this profile
        * `dns_client_config` - (Optional) Same structure as `dhcp_config.dhcp_server_config.dns_client_config` (DNS server IPs for clients)
        * `lease_time` - (Optional) IPv6 lease time in seconds. Default is 86400.
        * `preferred_time` - (Optional) Preferred lifetime in seconds for delegated prefixes. Minimum 48. If omitted, the server may derive a value from `lease_time`.
        * `ntp_servers` - (Optional) NTP servers as FQDN or IPv6 address
        * `sntp_servers` - (Optional) SNTP server IPv6 addresses
        * `advanced_config` - (Optional) Same semantics as IPv4 `dhcp_config.dhcp_server_config.advanced_config`
            * `is_distributed_dhcp` - (Optional / Computed) Distributed DHCP mode for DHCPv6
* `dns_forwarder_config` - (Optional) VPC DNS forwarder settings. Supported with NSX v9.2.0 and above.
    * `cache_size` - (Optional) DNS answer cache size
    * `conditional_forwarder_zone_paths` - (Optional) Policy paths of conditional DNS forwarder zones
    * `default_forwarder_zone_path` - (Optional) Policy path of the default forwarder zone
    * `log_level` - (Optional) One of `DEBUG`, `INFO`, `WARNING`, `ERROR`, `FATAL`
* `ipv6_profile_paths` - (Optional / Computed) List of policy paths for IPv6 DAD and/or NDRA profiles used on VPC subnets (up to two paths: one DAD profile and one NDRA profile). Supported with NSX v9.2.0 and above. NSX may populate this when omitted from the configuration.
* `service_subnet_cidrs` - (Optional / Computed) Service subnet CIDRs in `ip/prefix` form (IPv4 and/or IPv6). Minimum IPv4 size /28; must not overlap private, private TGW, or external IP blocks. Supported with NSX v9.2.0 and above. The manager may return CIDRs on read when this argument is omitted; Terraform keeps them in state without forcing a change.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_service_profile.test PATH
```

The above command imports VPC Service Profile named `test` with the NSX policy path `PATH`.
