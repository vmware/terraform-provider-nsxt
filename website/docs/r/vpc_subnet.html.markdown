---
subcategory: "FIXME"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_subnet"
description: A resource to configure a VpcSubnet.
---

# nsxt_policy_subnet

This resource provides a method for the management of a VpcSubnet.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_subnet" "test" {
    display_name      = "test"
    description       = "Terraform provisioned VpcSubnet"
    ipv4_subnet_size = FILL VALUE FOR schema.TypeInt
ip_addresses = FILL VALUE FOR schema.TypeString
access_mode = "Private"

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `advanced_config` - (Optional) None
  * `static_ip_allocation` - (Optional) None
    * `enabled` - (Optional) Enable ip and mac addresse allocation for VPC Subnet ports from static ip pool. To enable this,
dhcp pool shall be empty and static ip pool shall own all available ip addresses.

* `ipv4_subnet_size` - (Optional) If IP Addresses are not provided, this field will be used to carve out the ips
from respective ip block defined in the parent VPC. The default is 64.
If ip_addresses field is provided then ipv4_subnet_size field is ignored.
This field cannot be modified after creating a VPC Subnet.

* `ip_addresses` - (Optional) If not provided, Ip assignment will be done based on VPC CIDRs
This represents the VPC Subnet that is associated with tier.
If IPv4 CIDR is given, ipv4_subnet_size property is ignored.
For IPv6 CIDR, supported prefix length is /64.

* `access_mode` - (Optional) There are three kinds of Access Types supported for an Application.
Private  - VPC Subnet is accessible only within the application and its IPs are allocated from
           private IP address pool from VPC configuration unless specified explicitly by user.
Public   - VPC Subnet is accessible from external networks and its IPs are allocated from public IP
           address pool from VPC configuration unless specified explicitly by user.
Isolated - VPC Subnet is not accessible from other VPC Subnets within the same VPC.

* `dhcp_config` - (Optional) None
  * `dhcp_relay_config_path` - (Optional) Policy path of DHCP-relay-config. If configured then all the subnets will be configured with the DHCP relay server.
If not specified, then the local DHCP server will be configured for all connected subnets.

  * `dns_client_config` - (Optional) None
    * `dns_server_ips` - (Optional) IPs of the DNS servers which need to be configured on the workload VMs

  * `enable_dhcp` - (Optional) If activated, the DHCP server will be configured based on IP address type.
If deactivated then neither DHCP server nor relay shall be configured.

  * `static_pool_config` - (Optional) None
    * `ipv4_pool_size` - (Optional) Number of IPs to be reserved in static ip pool. Maximum allowed value is 'subnet size - 4'.
If dhcp is enabled then by default static ipv4 pool size will be zero and all available IPs will be reserved in
local dhcp pool.
If dhcp is deactivated then by default all IPs will be reserved in static ip pool.



## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_subnet.test UUID
```

The above command imports VpcSubnet named `test` with the NSX ID `UUID`.
