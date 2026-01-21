---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_external_address"
description: A resource to configure a VPC External Address Binding on Port.
---

# nsxt_vpc_external_address

This resource provides a method to configure External Address Binding on Subnet Port under VPC.
Only single resource should be configured for any given port.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

```hcl
data "nsxt_vpc_subnet_port" "test" {
  subnet_path = nsxt_vpc_subnet.test.path
  vm_id       = var.vm_id
}

resource "nsxt_vpc_external_address" "test" {
  parent_path                = data.nsxt_vpc_subnet_port.test.path
  allocated_external_ip_path = data.nsxt_vpc_ip_address_allocation.test.path
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required) Policy path for the Subnet Port.
* `allocated_external_ip_path` - (Required) Policy path for an external IP address allocation object.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `external_ip_address` - The actual IP address that was allocated with `allocated_external_ip_path`.

## Importing

An existing VPC External Address Binding can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_external_address.external_address1 PATH
```

The above command imports the VPC External Address Binding named `external_address` with the NSX Policy path of the parent port `PATH`.
