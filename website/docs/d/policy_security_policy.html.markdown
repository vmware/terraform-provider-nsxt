---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: policy_security_policy"
description: A policy Security Policy data source.
---

# nsxt_policy_security_policy

This data source provides information about policy Security Policues configured on NSX.
This data source can be useful for fetching policy path to use in `nsxt_policy_predefined_security_policy` resource.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_security_policy" "predefined" {
  is_default = true
  category   = "Application"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_security_policy" "predefined" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }

  is_default = true
  category   = "Application"
}
```

## Argument Reference

* `id` - (Optional) The ID of Security Policy to retrieve.
* `is_default` - (Optional) Whether this is a default policy. Default is `false`.
* `domain` - (Optional) The domain of the policy, defaults to `default`. Needs to be specified in VMC environment.
* `category` - (Optional) Category of the policy to retrieve.
* `display_name` - (Optional) The Display Name of the policy to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
