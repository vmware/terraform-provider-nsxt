---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_vm_tags"
sidebar_current: "docs-nsxt-resource-policy-vm-tags"
description: A resource to configure tags for a Virtual Machine in NSX Policy.
---

# nsxt_policy_vm_tags

  This resource provides a means to configure tags that are applied to objects such as Virtual Machines. A Virtual Machine is not directly managed by NSX however, NSX allows attachment of tags to a virtual machine. This tagging enables tag based grouping of objects. Deletion of `nsxt_policy_vm_tags` resource will remove all tags from the Virtual Machine and is equivalent to update operation with empty tag set.

## Example Usage

```hcl
resource "nsxt_policy_vm_tags" "vm1_tags" {
  instance_id = "${vsphere_virtual_machine.vm1.id}"

  tag {
    scope = "color"
    tag   = "blue"
  }

  tag {
    scope = "env"
    tag   = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the Virtual Machine. Can be the instance UUID or BIOS UUID.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Virtual Machine.

## Importing

An existing Tags collection can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_vm_tags.vm1_tags ID
```

The above would import NSX Virtual Machine tags as a resource named `vm1_tags` with the NSX ID `ID`, where ID is external ID of the Virtual Machine.
