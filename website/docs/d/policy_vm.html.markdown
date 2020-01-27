---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_vm"
sidebar_current: "docs-nsxt-datasource-policy-vm"
description: A Policy Virtual Machine ID data source.
---

# nsxt_policy_vm

This data source provides information about Policy based Virtual Machine (VM) configured in NSX and allows look-up of the VM by `display_name` or the BIOS, external or instance ID exposed on the VM resource.

## Example Usage

```hcl
data "nsxt_policy_vm" "nsxt_vm1" {
  display_name = "nsxt-virtualmachine1"
}
```

## Argument Reference

* `display_name` - (Optional) The Display Name prefix of the Virtual Machine to retrieve.
* `external_id` - (Optional) The external ID of the Virtual Machine.
* `bios_id` - (Optional) The BIOS UUID of the Virtual Machine.
* `instance_id` - (Optional) The instance UUID of the Virtual Machine.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the Virtual Machine.
