---
subcategory: "Beta"
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

## Argument Reference

* `value_tupe` - (Optional) Type of VM ID the user is interested in. Possible values are `bios_id`, `external_id`, `instance_id`. Default is `bios_id`.

## Attributes Reference

* `items` - Map of IDs by Display Name.
