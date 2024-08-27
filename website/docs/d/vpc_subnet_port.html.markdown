---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: vpc_subnet_port"
description: VPC Subnet Port data source.
---

# nsxt_vpc_subnet_port

This data source provides information about Subnet Port configured under VPC on NSX.

This data source is applicable to NSX Policy Manager.

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

data "nsxt_vpc_subnet_port" "test" {
  subnet_path = data.nsxt_vpc_subnet.test.path
  vm_id       = data.vpshere_virtual_machine.vm1.id
}
```

## Argument Reference

* `subnet_path` - (Required) Policy path of Subnet for the port.
* `vm_id` - (Required) Policy path of VM connected to the port.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `display_name` - The Display Name of the resource.
