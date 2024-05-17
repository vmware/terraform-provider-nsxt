---
subcategory: "Segments"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_segment"
description: A resource to configure a network Segment.
---

# nsxt_policy_segment

This resource provides a method for the management of Segments.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

~> **NOTE:** For best user experience, it is recommended to use `nsxt_policy_fixed_segment` resource on VMC. Flexible segments and not supported in early versions of VMC UI.

## Example Usage

```hcl
resource "nsxt_policy_segment" "segment1" {
  display_name        = "segment1"
  description         = "Terraform provisioned Segment"
  transport_zone_path = data.nsxt_policy_transport_zone.tz1.path

}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_segment" "segment1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name      = "segment1"
  description       = "Terraform provisioned Segment"
  connectivity_path = data.nsxt_policy_tier1_gateway.tier1_gw.path

  subnet {
    cidr        = "12.12.2.1/24"
    dhcp_ranges = ["12.12.2.100-12.12.2.160"]

    dhcp_v4_config {
      server_address = "12.12.2.2/24"
      lease_time     = 36000

      dhcp_option_121 {
        network  = "6.6.6.0/24"
        next_hop = "1.1.1.21"
      }

      dhcp_generic_option {
        code   = "119"
        values = ["abc"]
      }
    }
  }
}
```

## Example Usage with VCD network

```hcl
resource "nsxt_policy_segment" "segment1" {
  display_name        = "external"
  transport_zone_path = data.nsxt_policy_transport_zone.tz1.path

  subnet {
    cidr = "10.240.12.1/24"
  }

  ignore_tags {
    scopes = ["SYSTEM"]
  }

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `ignore_tags` - (Optional) A list of tag scopes that provider should ignore, more specifically, it should not detect drift when tags with such scope are present on NSX, and it should not overwrite them when applying its own tags. This feature is useful for external network with VCD scenario.
  * `scopes` - (Required) - List of scope values that should cause scope/tag pair to be ignored.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `connectivity_path` - (Optional) Policy path to the connecting Tier-0 or Tier-1.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `domain_name`- (Optional) DNS domain names.
* `overlay_id` - (Optional) Overlay connectivity ID for this Segment.
* `vlan_ids` - (Optional) List of VLAN IDs or ranges. Specifying vlan ids can be useful for overlay segments, f.e. for EVPN.
* `transport_zone_path` - (Optional) Policy path to the Overlay transport zone. This property is required for NSX Local Manager, and should not be specified for NSX Global Manager, where NSX will automatically assign default transport zone on each site.
* `dhcp_config_path` - (Optional) Policy path to DHCP server or relay configuration to use for subnets configured on this segment. This attribute is supported with NSX 3.0.0 onwards.
* `subnet` - (Optional) Subnet configuration block.
  * `cidr` - (Required) Gateway IP address CIDR. This argument can not be changed if DHCP is enabled for the subnet.
  * `dhcp_ranges` - (Optional) List of DHCP address ranges for dynamic IP allocation.
  * `dhcp_v4_config` - (Optional) DHCPv4 config for IPv4 subnet. This clause is supported with NSX 3.0.0 onwards.
     * `server_address` - (Optional) IP address of the DHCP server in CIDR format. This attribute is required if segment has provided dhcp_config_path and it represents a DHCP server config.
     * `dns_servers` - (Optional) List of IP addresses of DNS servers for the subnet.
     * `lease_time`  - (Optional) DHCP lease time in seconds.
     * `dhcp_option_121` - (Optional) DHCP classless static routes.
         * `network` - (Required) Destination in cidr format.
         * `next_hop` - (Required) IP address of next hop.
     * `dhcp_generic_option` - (Optional) Generic DHCP options.
         * `code` - (Required) DHCP option code. Valid values are from 0 to 255.
         * `values` - (Required) List of DHCP option values.
  * `dhcp_v6_config` - (Optional) DHCPv6 config for IPv6 subnet. This clause is supported with NSX 3.0.0 onwards.
     * `server_address` - (Optional) IP address of the DHCP server in CIDR format. This attribute is required if segment has provided dhcp_config_path and it represents a DHCP server config.
     * `dns_servers` - (Optional) List of IP addresses of DNS servers for the subnet.
     * `lease_time`  - (Optional) DHCP lease time in seconds.
     * `preferred_time` - (Optional) The time interval in seconds, in which the prefix is advertised as preferred.
     * `domain_names` - (Optional) List of domain names for this subnet.
     * `excluded_range` - (Optional) List of excluded address ranges to define dynamic ip allocation ranges.
         * `start` - (Required) IPv6 address that marks beginning of the range.
         * `end` - (Required) IPv6 address that marks end of the range.
     * `sntp_servers` - (Optional) IPv6 address of SNTP servers for the subnet.
* `l2_extension` - (Optional) Configuration for extending Segment through L2 VPN.
  * `l2vpn_paths` - (Optional) Policy paths of associated L2 VPN sessions.
  * `tunnel_id` - (Optional) The Tunnel ID that's a int value between 1 and 4093.
* `advanced_config` - (Optional) Advanced Segment configuration.
  * `address_pool_path` - (Optional) List of Policy path to IP address pools (for now only one pool is supported by NSX)
  * `connectivity` - (Optional) Connectivity configuration to manually connect (ON) or disconnect (OFF).
  * `hybrid` - (Optional) Boolean flag to identify a hybrid logical switch.
  * `local_egress` - (Optional) Boolean flag to enable local egress when used in conjunction with L2VPN.
  * `uplink_teaming_policy` - (Optional) The name of the switching uplink teaming policy for the bridge endpoint. This name corresponds to one of the switching uplink teaming policy names listed in the transport zone.
  * `urpf_mode` - (Optional) URPF mode to be applied to gateway downlink interface. One of `STRICT`, `NONE`.
* `discovery_profile` - (Optional) IP and MAC discovery profile specification for the segment.
  * `ip_discovery_profile_path` - (Optional) Path for IP discovery profile to be associated with the segment.
  * `mac_discovery_profile_path` - (Optional) Path for MAC discovery profile to be associated with the segment.
* `security_profile` - (Optional) Security profile specification for the segment.
  * `spoofguard_profile_path` - (Optional) Path for spoofguard profile to be associated with the segment.
  * `security_profile_path` - (Optional) Path for segment security profile to be associated with the segment.
* `qos_profile` - (Optional) QoS profile specification for the segment.
  * `qos_profile_path` - (Optional) Path for qos profile to be associated with the segment.
* `bridge_config` - (Optional) List of edge bridge configuration for the segment. This setting is not supported on Global Manager.
  * `profile_path` - (Required) Path for edge bridge profile to be associated with the segment.
  * `transport_zone_path` - (Required) Path for vlan transport zone for the bridge.
  * `vlan_ids` - (Required) List of VLAN IDs or ranges.
  * `uplink_teaming_policy` - (Optional) The name of the switching uplink teaming policy for the bridge endpoint.
* `metadata_proxy_paths` - (Optional) Metadata Proxy Configuration Paths.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Security Policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* In the `subnet`:
  * `network` The network CIDR for the subnet.

## Importing

An existing segment can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_segment.segment1 ID
```
The above command imports the segment  named `segment1` with the NSX Segment ID `ID`.

```
terraform import nsxt_policy_segment.segment1 POLICY_PATH
```
The above command imports the segment  named `segment1` with the policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.

~> **NOTE:** Only flexible (infra) segments can be imported here. To import fixed segment, please use `nsxt_policy_fixed_segment` resource.

~> **NOTE:** Please make sure `advanced_config` clause is present in configuration if you with to include it in import, otherwise it will be ignored with NSX 3.2 onwards. This is due to a platform change in handling advanced config in the API.
