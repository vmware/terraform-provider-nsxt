---
subcategory: "Grouping and Tagging"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_vm"
description: A Discovered Policy Virtual Machine data source.
---

# nsxt_policy_vm

This data source provides information about Policy based Virtual Machine (VM) listed in NSX inventory, and allows look-up of the VM by `display_name` or the BIOS, external or instance ID exposed on the VM resource.

This data source is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_vm" "nsxt_vm1" {
  display_name = "nsxt-virtualmachine1"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_vm" "nsxt_vm1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "nsxt-virtualmachine1"
}
```

## Argument Reference

* `display_name` - (Optional) The Display Name prefix of the Virtual Machine to retrieve.
* `external_id` - (Optional) The external ID of the Virtual Machine.
* `bios_id` - (Optional) The BIOS UUID of the Virtual Machine.
* `instance_id` - (Optional) The instance UUID of the Virtual Machine.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the Virtual Machine.
* `tag` - A list of scope + tag pairs associated with this Virtual Machine.
