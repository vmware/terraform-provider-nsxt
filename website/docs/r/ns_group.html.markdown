---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_ns_group"
description: A resource to configure a networking and security group in NSX.
---

# nsxt_ns_group

This resource provides a method to create and manage a network and security (NS) group in NSX. A NS group is used to group other objects into collections for application of other settings.

## Example Usage

```hcl
resource "nsxt_ns_group" "group2" {
  description  = "NG provisioned by Terraform"
  display_name = "NG"

  member {
    target_type = "NSGroup"
    value       = nsxt_ns_group.group1.id
  }

  membership_criteria {
    target_type = "LogicalPort"
    scope       = "XXX"
    tag         = "YYY"
  }

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this NS group.
* `member` - (Optional) Reference to the direct/static members of the NSGroup. Can be ID based expressions only. VirtualMachine cannot be added as a static member.
  * `target_type` - (Required) Static member type, one of: NSGroup, IPSet, LogicalPort, LogicalSwitch, MACSet
  * `value` - (Required) Member ID
* `membership_criteria` - (Optional) List of tag or ID expressions which define the membership criteria for this NSGroup. An object must satisfy at least one of these expressions to qualify as a member of this group.
  * `target_type` - (Required) Dynamic member type, one of: LogicalPort, LogicalSwitch, VirtualMachine.
  * `scope` - (Optional) Tag scope for matching dynamic members.
  * `tag` - (Optional) Tag value for matching dynamic members.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the NS group.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing networking and security group can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_ns_group.group2 UUID
```

The above command imports the networking and security group named `group2` with the NSX id `UUID`.
