---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_vpc_external_address"
description: A resource to configure a VPC External Address Binding on Port.
---

# nsxt_vpc_external_address

This resource provides a method to configure External Address Binding on Subnet Port under VPC.
Only single resource should be configured for any given port.

This resource is applicable to NSX Policy Manager.

```hcl
resource "nsxt_vpc_ip_address_allocation" "test" {
  context {
    project_id = "Dev_project"
    vpc_id     = "dev_vpc"
  }
  display_name    = "external"
  allocation_size = 1
}

data "nsxt_policy_vm" "vm1" {
  context {
    project_id = "Dev_project"
    vpc_id     = "dev_vpc"
  }
  display_name = "myvm-1"
}

data "nsxt_vpc_subnet_port" "test" {
  subnet_path = nsxt_vpc_subnet.test.path
  vm_id       = data.nsxt_policy_vm.vm1.instance_id
}

resource "nsxt_vpc_external_address" "test" {
  parent_path                = data.nsxt_vpc_subnet_port.test.path
  allocated_external_ip_path = nsxt_vpc_ip_address_allocation.test.path
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required) Policy path for the Subnet Port.
* `allocated_external_ip_path` - (Required) Policy path for extrenal IP address allocation object.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `external_ip_address` - The actual IP address that was allocated with `allocated_external_ip_path`.

## Importing

An existing VPC External Address Binding can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_vpc_external_address.external_address1 PATH
```

The above command imports the VPC External Address Binding named `external_address` with the NSX Policy path of the parent port `PATH`.
