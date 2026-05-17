---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_dhcp_v6_static_binding"
description: A resource to configure a DHCP IPv6 Static Binding on a VPC subnet.
---

# nsxt_vpc_dhcp_v6_static_binding

This resource provides a method for the management of a DhcpV6StaticBindingConfig under a VPC subnet (`.../vpcs/<vpc-id>/subnets/<subnet-id>/dhcp-static-binding-configs/<binding-id>`).

This resource is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards. The parent subnet must be configured for IPv6 DHCP as required by NSX (for example dual-stack addressing and DHCPv6 server mode on the subnet).

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

resource "nsxt_vpc_dhcp_v6_static_binding" "test" {
  parent_path     = data.nsxt_vpc_subnet.test.path
  display_name    = "test"
  description     = "Terraform provisioned DhcpV6StaticBindingConfig"
  mac_address     = "10:0e:00:11:22:02"
  lease_time      = 86400
  ip_addresses    = ["2001:db8:1::101"]
  domain_names    = ["client.example.org"]
  dns_nameservers = ["2001:db8:1::53"]
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required) Policy path of parent VpcSubnet object, typically reference to `path` in `nsxt_vpc_subnet` data source or resource.
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `dns_nameservers` - (Optional) List of DNS nameserver addresses or hostnames for the client.
* `domain_names` - (Optional) List of domain names assigned to the client host.
* `ip_addresses` - (Optional) List of IPv6 addresses assigned to the client.
* `lease_time` - (Optional) DHCP lease time in seconds. Minimum 60. Default is 86400.
* `preferred_time` - (Optional) Preferred lifetime in seconds for the prefix. Minimum 48. If omitted, the server may derive a value from `lease_time`.
* `mac_address` - (Optional) MAC address of the client host.
* `ntp_servers` - (Optional) List of NTP servers as FQDN or IPv6 address.
* `sntp_servers` - (Optional) List of SNTP server IPv6 addresses.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_dhcp_v6_static_binding.test PATH
```

The above command imports DhcpV6StaticBindingConfig named `test` with the policy path `PATH`.
