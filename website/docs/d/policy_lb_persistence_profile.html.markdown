---
subcategory: "Load Balancer"
layout: "nsxt"
page_title: "NSXT: policy_lb_persistence_profile"
description: Policy Load Balancer Persistence Profile data source.
---

# nsxt_policy_lb_persistence_profile

This data source provides information about Policy Load Balancer Persistence Profiles configured on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_lb_persistence_profile" "test" {
  display_name = "policy-lb-persistence-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of Profile to retrieve.
* `type` - (Optional) The Load Balancer Persistence Profile type. One of `ANY`, `SOURCE_IP`, `COOKIE` or `GENERIC`. Defaults to `ANY`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
