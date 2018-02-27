---
layout: "nsxt"
page_title: "NSXT: nsxt_vm_tags"
sidebar_current: "docs-nsxt-resource-vm-tags"
description: |-
  Provides a resource to configure Tags for VM for on NSX-T Manager.
---

# nsxt_vm_tags
  Provides a resource to configure Tags for Virtual Machine for on NSX-T Manager. The Virtual Machine is not managed by NSX-T, but NSX-T allows attachment of tags to the VM. This enables tag-based grouping. Deletion of nsxt_vm_tags resource will remove all tags from the VM, and is equivalent to update operation with empty tag set.


## Example Usage

```hcl
resource "nsxt_vm_tags" "vm1_tags" {
    instance_id = "${vsphere_virtual_machine.vm1.id}"
    tag {
        scope = "color"
        tag = "blue"
    }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) BIOS Id of the Virtual Machine.
* `tag` - (Required) A list of scope + tag pairs to associate with this VM.

## Importing

An existing Tags collection can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import nsxt_vm_tags.x id
```

The above would import nsxt vm tags resource named `x` with the nsx id `id`, where id is external ID (not the BIOS id) of the VM.
