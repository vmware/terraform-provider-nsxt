---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: policy_intrusion_service_profile"
description: Policy Intrusion Service Profile data source.
---

# nsxt_policy_intrusion_service_profile

This data source provides information about policy IDS Profile configured on NSX.
This data source is applicable to NSX Policy Manager and VMC (NSX version 3.0.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_intrusion_service_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "DefaultIDSProfile"
}
```

## Argument Reference

* `id` - (Optional) The ID of Profile to retrieve. If ID is specified, no additional argument should be configured.
* `display_name` - (Optional) The Display Name prefix of the Profile to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
