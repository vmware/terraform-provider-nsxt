---
subcategory: "VPC"
page_title: "NSXT: vpc_dhcp_v6_static_binding"
description: A DHCP IPv6 Static Binding on a VPC subnet data source.
---

# nsxt_vpc_dhcp_v6_static_binding

This data source provides information about a `DhcpV6StaticBindingConfig` configured under a VPC subnet.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_vpc_dhcp_v6_static_binding" "test" {
  display_name = "my-binding"
  parent_path  = nsxt_vpc_subnet.demo.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the DHCP IPv6 Static Binding to retrieve.
* `display_name` - (Optional) The Display Name of the DHCP IPv6 Static Binding to retrieve.
* `parent_path` - (Required) Policy path of the parent VPC Subnet.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
