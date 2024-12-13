---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_vpc_dhcp_v4_static_binding"
description: A resource to configure a DHCP IPv4 Static Binding.
---

# nsxt_vpc_dhcp_v4_static_binding

This resource provides a method for the management of a DhcpV4StaticBindingConfig.
This resource is applicable to NSX Policy Manager.

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

data "nsxt_vpc_subnet" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  display_name = "subnet1"
}

resource "nsxt_vpc_dhcp_v4_static_binding" "test" {
  parent_path     = data.nsxt_vpc_subnet.test.path
  display_name    = "test"
  description     = "Terraform provisioned DhcpV4StaticBindingConfig"
  gateway_address = "192.168.240.1"
  host_name       = "host.example.org"
  mac_address     = "10:0e:00:11:22:02"
  lease_time      = 162
  ip_address      = "192.168.240.41"
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required) Policy path of parent VpcSubnet object, typically reference to `path` in `nsxt_vpc_subnet` data source or resource.
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `gateway_address` - (Optional) When not specified, gateway address is auto-assigned from segment configuration.
* `host_name` - (Optional) Hostname to assign to the host.
* `mac_address` - (Optional) MAC address of the host.
* `lease_time` - (Optional) DHCP lease time in seconds.
* `ip_address` - (Optional) IP assigned to host. The IP address must belong to the subnet, if any, configured on Segment.
* `options` - (Optional) DHCPv4 options block
  * `option121` - (Optional) Specification for DHCP option 121
    * `static_route` - (Required) Classless static route of DHCP option 121.
      * `next_hop` - (Required) IP address of next hop of the route.
      * `network` - (Required) Destination network in CIDR format.
  * `other` - (Optional) To define DHCP options other than option 121 in generic format.
    * `code` - (Required) Code of the dhcp option.
    * `values` - (Optional) Value of the option.
  
Please note, only the following options can be defined in generic
format. Those other options will be accepted without validation
but will not take effect.
--------------------------
Code    Name
--------------------------
    2   Time Offset
    6   Domain Name Server
    13  Boot File Size
    19  Forward On/Off
    26  MTU Interface
    28  Broadcast Address
    35  ARP Timeout
    40  NIS Domain
    41  NIS Servers
    42  NTP Servers
    44  NETBIOS Name Srv
    45  NETBIOS Dist Srv
    46  NETBIOS Node Type
    47  NETBIOS Scope
    58  Renewal Time
    59  Rebinding Time
    64  NIS+-Domain-Name
    65  NIS+-Server-Addr
    66  TFTP Server-Name (used by PXE)
    67  Bootfile-Name (used by PXE)
    117 Name Service Search
    119 Domain Search
    150 TFTP server address (used by PXE)
    209 PXE Configuration File
    210 PXE Path Prefix
    211 PXE Reboot Time

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_vpc_dhcp_v4_static_binding.test PATH
```

The above command imports DhcpV4StaticBindingConfig named `test` with the policy path `PATH`.
