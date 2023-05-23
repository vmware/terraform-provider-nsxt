---
subcategory: "Grouping and Tagging"
layout: "nsxt"
page_title: "NSXT: policy_group"
description: Policy Group data source.
---

# nsxt_policy_group

This data source provides information about an inventory Group configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_group" "test" {
  display_name = "group1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Group to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Group to retrieve.

* `domain` - (Optional) The domain this Group belongs to. For VMware Cloud on AWS use `cgw`. For Global Manager, please use site id for this field. If not specified, this field is default to `default`. 

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
