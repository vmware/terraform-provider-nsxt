---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: ns_group"
description: A networking and security group data source.
---

# nsxt_ns_group

This data source provides information about a network and security (NS) group in NSX. A NS group is used to group other objects into collections for application of other settings.

## Example Usage

```hcl
data "nsxt_ns_group" "ns_group_1" {
  display_name = "test group"
}
```

## Argument Reference

* `id` - (Optional) The ID of NS group to retrieve
* `display_name` - (Optional) The Display Name of the NS group to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the NS group.
