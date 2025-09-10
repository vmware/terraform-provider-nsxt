---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_subnet"
description: A resource to configure a VpcSubnet.
---

# nsxt_vpc_subnet

This resource provides a method for the management of a Vpc Subnet.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc" "demovpc" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "vpc1"
}

// This will create an isolated subnet, which requires
// a CIDR to be specified in ip_addresses
resource "nsxt_vpc_subnet" "test_isolated" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }

  display_name = "test-subnet-isolated"
  description  = "Test Isolated VPC subnet"
  ip_addresses = ["192.168.240.0/24"]
  access_mode  = "Isolated"
}

// This will create a private subnet. A CIDR will be
// automatically cut from the VPC's private IP ranges.
resource "nsxt_vpc_subnet" "test_private" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }

  display_name     = "test-subnet-private"
  description      = "Test Private VPC subnet"
  ipv4_subnet_size = 32
  access_mode      = "Private"
}
```

~> **NOTE:** In some cases, subnet creation will depend on VPC attachment. If both resources are being created within same apply,
  explicit `depends_on` meta argument needs to be added to enforce this dependency.

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `ipv4_subnet_size` - (Optional) If IP Addresses are not provided, this field will be used to carve out the ips
  from respective ip block defined in the parent VPC. The default is 64. Conflicts with `ip_addresses`.
* `ip_addresses` - (Optional) If not provided, Ip assignment will be done based on VPC CIDRs. Conflicts with `ipv4_subnet_size`. This argument is required when access_mode is set to `Isolated`
* `access_mode` - (Optional) Subnet access mode, one of `Private`, `Public`, `Isolated` or `Private_TGW`. Default is `Private`
* `advanced_config` - (Optional) Advanced Configuration for the Subnet
    * `gateway_addresses` - (Optional) List of Gateway IP Addresses per address family, in CIDR format
    * `connectivity_state` - (Optional) Connectivity state for the subnet, one of `CONNECTED`, `DISCONNECTED`
    * `dhcp_server_addresses` - (Optional) List of DHCP server addresses per address family, in CIDR format
    * `static_ip_allocation` - (Optional) Static IP allocation configuration
        * `enabled` - (Optional) Enable ip and mac address allocation for VPC Subnet ports from static ip pool. To enable this,
          dhcp pool shall be empty and static ip pool shall own all available ip addresses.
    * `extra_configs` - (Optional) List of vendor specific configuration key/value pairs
        * `config_pair` - (Required)
            * `key` - (Required) key for vendor-specific configuration
            * `value` - (Required) value for vendor-specific configuration
    * `enable_vlan_extension` - (Optional) Enabling VLAN connection for the subnet. The user must configure the exclusive IP block for this subnet. This IP block must be leveraged by a distributed VLAN connection. The default value for this will be false. This attribute is supported with NSX 9.1.0 onwards.
* `dhcp_config` - (Optional) DHCP configuration block
    * `dns_server_preference` - (Optional) DNS server IP preference. Select the preference between the DNS server IPs (from the DHCP config in VPC service profile), and VPC DNS forwarder IP. The preferred DNS server IP config will be attempted first when the system selects the DNS server to forward DNS requests. This can be one of `PROFILE_DNS_SERVERS_PREFERRED_OVER_DNS_FORWARDER`, `DNS_FORWARDER_PREFERRED_OVER_PROFILE_DNS_SERVERS`. The default value is `PROFILE_DNS_SERVERS_PREFERRED_OVER_DNS_FORWARDER`.
    * `mode` - (Optional) The operational mode of DHCP within the subnet, can be one of `DHCP_SERVER`, `DHCP_RELAY`, `DHCP_DEACTIVATED`.
       Default is `DHCP_DEACTIVATED`
    * `dhcp_server_additional_config` - (Optional) Additional DHCP server config
        * `options` - (Optional) DHCPv4 options block
            * `option121` - (Optional) Specification for DHCP option 121
                * `static_route` - (Optional) Static route
                    * `network` - (Optional) Destination network in CIDR format
                        * `next_hop` - (Optional) IP Address for next hop of the route
                    * `other` - (Optional) DHCP option in generic format
                        * `code` - (Optional) Code of DHCP option
                        * `values` - (Optional) List of values in string format
        * `reserved_ip_ranges` - (Optional) Specifies IP ranges that are reserved and excluded from being assigned by the DHCP server to clients.
         This is a list of IP ranges or IP addresses.
* `ip_blocks` - (Optional) List of IP block path for subnet IP allocation

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful
  for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_subnet.test PATH
```

The above command imports VpcSubnet named `test` with the policy path `PATH`.
