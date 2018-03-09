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
resource "nsxt_ns_group" "group2" {
  description  = "NG provisioned by Terraform"
  display_name = "NG"

  member {
    target_type = "NSGroup"
    value       = "${nsxt_ns_group.group1.id}"
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
* `member` - (Optional) Reference to the direct/static members of the NSGroup. Can be ID based expressions only. VirtualMachine cannot be added as a static member. target_type can be: NSGroup, IPSet, LogicalPort, LogicalSwitch, MACSet
* `membership_criteria` - (Optional) List of tag or ID expressions which define the membership criteria for this NSGroup. An object must satisfy at least one of these expressions to qualify as a member of this group.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the ns_group.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing NS group can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_ns_group.group2 UUID
```

The above would import the NS group named `group2` with the nsx id `UUID`
