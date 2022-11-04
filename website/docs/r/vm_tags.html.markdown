---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_vm_tags"
description: A resource to configure tags for a virtual machine in NSX.
---

# nsxt_vm_tags

  This resource provides a means to configure tags that are applied to objects such as virtual machines. A virtual machine is not directly managed by NSX however, NSX allows attachment of tags to a virtual machine. This tagging enables tag based grouping of objects. Deletion of `nsxt_vm_tags` resource will remove all tags from the virtual machine and is equivalent to update operation with empty tag set.

## Example Usage

```hcl
resource "nsxt_vm_tags" "vm1_tags" {
  instance_id = vsphere_virtual_machine.vm1.id

  tag {
    scope = "color"
    tag   = "blue"
  }

  logical_port_tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) BIOS Id of the Virtual Machine.
* `tag` - (Optional) A list of scope + tag pairs to associate with this VM.
* `logical_port_tag` - (Optional) A list of scope + tag pairs to associate with all logical ports that are automatically created for this VM.

## Importing

An existing Tags collection can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_vm_tags.vm1_tags id
```

The above would import NSX virtual machine tags as a resource named `vm1_tags` with the NSX id `id`, where id is external ID (not the BIOS id) of the virtual machine.
