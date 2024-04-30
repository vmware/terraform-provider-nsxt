---
subcategory: "Grouping and Tagging"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_vms"
description: A Discovered Policy Virtual Machines data source.
---

# nsxt_policy_vms

This data source provides map of all Policy based Virtual Machines (VMs) listed in NSX inventory, and allows look-up of the VM by `display_name` in the map. Value of the map would provide one of VM ID types, according to `value_type` argument.

This data source is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_vms" "all" {
  state      = "running"
  guest_os   = "ubuntu"
  value_type = "bios_id"
}

resource "nsxt_policy_vm_tags" "test" {
  instance_id = data.nsxt_policy_vms.all.items["vm-1"]

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_vms" "all" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  state      = "running"
  guest_os   = "ubuntu"
  value_type = "bios_id"
}

resource "nsxt_policy_vm_tags" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  instance_id = data.nsxt_policy_vms.all.items["vm-1"]

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

* `value_type` - (Optional) Type of VM ID the user is interested in. Possible values are `bios_id`, `external_id`, `instance_id`. Default is `bios_id`.
* `state` - (Optional) Filter results by power state of the machine.
* `guest_os` - (Optional) Filter results by operating system of the machine. The match is case insensitive and prefix-based.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

* `items` - Map of IDs by Display Name.
