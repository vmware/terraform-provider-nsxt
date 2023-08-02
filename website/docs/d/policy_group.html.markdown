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

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_group" "demogroup" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demogroup"
}
```

## Argument Reference

* `id` - (Optional) The ID of Group to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Group to retrieve.
* `domain` - (Optional) The domain this Group belongs to. For VMware Cloud on AWS use `cgw`. For Global Manager, please use site id for this field. If not specified, this field is default to `default`. 
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
