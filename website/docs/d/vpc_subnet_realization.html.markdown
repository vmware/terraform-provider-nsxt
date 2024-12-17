---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: vpc_subnet_realization"
description: State of VPC subnet realization on hypervisors.
---

# nsxt_vpc_subnet_realization

This data source provides information about the realization of VPC subnet or hypervisor.
This data source will wait until realization is complete with either success, partial success or error. 
It is recommended to use this data source in conjunction with vsphere provider, in order to ensure subnet is realized 
on hypervisor before VM is created on same network.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_vpc_subnet" "s1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }

  display_name     = "test-subnet"
  description      = "Test VPC subnet"
  ipv4_subnet_size = 32
  ip_addresses     = ["192.168.240.0/24"]
  access_mode      = "Isolated"
}

data "nsxt_vpc_subnet_realization" "s1" {
  path = nsxt_vpc_subnet.s1.path
}

# usage in vsphere provider
data "vsphere_network" "net" {
  name          = nsxt_vpc_subnet_realization.s1.network_name
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

## Argument Reference

* `path` - (Required) The policy path of the subnet.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - The realization state of the resource: `success`, `partial_success`, `orphaned`, `failed` or `error`.
* `network_name` - Network name on the hypervisor. This attribute can be used in vsphere provider in order to ensure implicit dependency on segment realization.
