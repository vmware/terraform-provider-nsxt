---
subcategory: "Segments"
layout: "nsxt"
page_title: "NSXT: policy_segment"
description: Policy Segment data source.
---

# nsxt_policy_segment

This data source provides information about policy Segment configured on NSX.
This data source is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_segment" "test" {
  display_name = "segment1"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_segment" "demoseg" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demoseg"
}
```

## Argument Reference

* `id` - (Optional) The ID of Segment to retrieve. If ID is specified, no additional argument should be configured.
* `display_name` - (Optional) The Display Name prefix of the Segment to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
