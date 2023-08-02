---
subcategory: "Segments"
layout: "nsxt"
page_title: "NSXT: policy_qos_profile"
description: Policy QoS Profile data source.
---

# nsxt_policy_qos_profile

This data source provides information about policy Quality of Service Profile configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_qos_profile" "test" {
  display_name = "qos-profile1"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_qos_profile" "demoqosprof" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demoqosprof"
}
```

## Argument Reference

* `id` - (Optional) The ID of Profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Profile to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
