---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_policy"
description: A data source to retrieve an Intrusion Service (IDS) Policy.
---

# nsxt_policy_intrusion_service_policy

This data source provides information about an existing Intrusion Service (IDS) Policy configured on NSX.
This data source can be useful for fetching policy path to use in `nsxt_policy_intrusion_service_policy_rule` resource.

~> **NOTE:** This data source retrieves only the policy metadata (id, display_name, description, path). It does not retrieve the rules within the policy. To manage rules, use the `nsxt_policy_intrusion_service_policy_rule` resource with the policy path obtained from this data source.

This data source is applicable to NSX Policy Manager and VMC (NSX version 3.1.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_intrusion_service_policy" "ids_policy" {
  display_name = "intrusion-service-policy"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_intrusion_service_policy" "ids_policy" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "intrusion-service-policy"
}
```

## Argument Reference

* `id` - (Optional) The ID of the policy to retrieve.
* `display_name` - (Optional) The display name of the policy to retrieve.
* `domain` - (Optional) The domain of the policy. Defaults to `default`.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
