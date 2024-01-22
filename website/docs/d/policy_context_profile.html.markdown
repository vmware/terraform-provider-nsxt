---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: policy_context_profile"
description: Policy Context Profile Profile data source.
---

# nsxt_policy_context_profile

This data source provides information about Policy Context Profile configured on NSX.

This data source is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_context_profile" "diameter" {
  display_name = "DIAMETER"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_context_profile" "demoprof" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demoprof"
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
