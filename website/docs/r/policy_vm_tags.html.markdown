---
subcategory: "Grouping and Tagging"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_vm_tags"
description: A resource to configure tags for a Virtual Machine in NSX Policy.
---

# nsxt_policy_vm_tags

  This resource provides a means to configure tags that are applied to objects such as Virtual Machines. A Virtual Machine is not directly managed by NSX however, NSX allows attachment of tags to a virtual machine. This tagging enables tag based grouping of objects. Deletion of `nsxt_policy_vm_tags` resource will remove all tags from the Virtual Machine and is equivalent to update operation with empty tag set.

  When updating resource, if you wish to delete existing port tags while leaving VM tags in place, please specify `port` clause with no tags.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_vm_tags" "vm1_tags" {
  instance_id = vsphere_virtual_machine.vm1.id

  tag {
    scope = "color"
    tag   = "blue"
  }

  tag {
    scope = "env"
    tag   = "test"
  }

  port {
    segment_path = nsxt_policy_segment.seg1.path
    tag {
      tag {
        scope = "color"
        tag   = "green"
      }
    }
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_vm_tags" "vm1_tags" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  instance_id = vsphere_virtual_machine.vm1.id

  tag {
    scope = "color"
    tag   = "blue"
  }

  port {
    segment_path = nsxt_policy_segment.seg1.path
    tag {
      tag {
        scope = "color"
        tag   = "green"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the Virtual Machine. Can be the instance UUID or BIOS UUID.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Virtual Machine.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `port` - (Optional) Option to tag segment port auto-created for the VM on specified segment.
  * `segment_path` - (Required) Segment where the port is to be tagged.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this segment port.

## Importing

An existing VM Tags collection can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_vm_tags.vm1_tags ID
```
The above would import NSX Virtual Machine tags as a resource named `vm1_tags` with the NSX ID `ID`, where ID is external ID of the Virtual Machine.

```
terraform import nsxt_policy_vm_tags.vm1_tags POLICY_PATH
```
The above would import NSX Virtual Machine tags as a resource named `vm1_tags` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.

Note that import of port tags is not supported.
