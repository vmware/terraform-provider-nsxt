---
layout: "nsxt"
page_title: "NSXT: nsxt_ns_group"
sidebar_current: "docs-nsxt-resource-ns-group"
description: |-
  Provides a resource to configure ns group on NSX-T manager
---

# nsxt_ns_group

Provides a resource to configure ns group on NSX-T manager

## Example Usage

```hcl
resource "nsxt_ns_group" "NG" {
  description = "NG provisioned by Terraform"
  display_name = "NG"
  tags = [{ scope = "color"
            tag = "red" }
  ]
  members = [{target_type = "NSGroup"
              value = "${nsxt_ns_group.GRP1.id}"},
  membership_criteria = ...
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) Defaults to ID if not set.
* `tags` - (Optional) A list of scope + tag pairs to associate with this ns_group.
* `members` - (Optional) Reference to the direct/static members of the NSGroup. Can be ID based expressions only. VirtualMachine cannot be added as a static member. target_type can be: NSGroup, IPSet, LogicalPort, LogicalSwitch, MACSet
* `membership_criteria` - (Optional) List of tag or ID expressions which define the membership criteria for this NSGroup. An object must satisfy at least one of these expressions to qualify as a member of this group.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the ns_group.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `system_owned` - A boolean that indicates whether this resource is system-owned and thus read-only.
* `member_count` - Count of the members added to this NSGroup.
