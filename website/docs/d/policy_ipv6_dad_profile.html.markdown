---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: policy_ipv6_dad_profile"
description: Policy IPv6 DAD Profile data source.
---

# nsxt_policy_ipv6_dad_profile

This data source provides information about policy IPv6 Duplicate Address Discovery Profile configured on NSX.

This data source is applicable to NSX Global Manager, and NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_ipv6_dad_profile" "test" {
  display_name = "ipv6-dad-profile1"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_ipv6_dad_profile" "demodad6" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demodad6"
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
