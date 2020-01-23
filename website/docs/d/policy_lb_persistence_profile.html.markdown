---
layout: "nsxt"
page_title: "NSXT: policy_lb_persistence_profile"
sidebar_current: "docs-nsxt-datasource-policy-lb-persistence-profile"
description: Policy Load Balancer Persistence Profile data source.
---

# nsxt_policy_lb_persistence_profile

This data source provides information about policy Load Balancer Persistence Profiles configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_lb_persistence_profile" "test" {
  display_name = "policy-lb-persistence-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Load Balanacer Persistence Profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Load Balancer Persistence Profile to retrieve.
* `type` - (Optional) The Load Balancer Persistence Profile type. One of `ANY`, `SOURCE_IP`, `COOKIE` or `GENERIC`. Defaults to `ANY`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
