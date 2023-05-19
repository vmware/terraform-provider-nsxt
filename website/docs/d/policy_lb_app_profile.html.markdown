---
subcategory: "Load Balancer"
layout: "nsxt"
page_title: "NSXT: policy_lb_app_profile"
description: Policy Load Balancer Application Profile data source.
---

# nsxt_policy_lb_app_profile

This data source provides information about policy Load Balancer Application Profile configured in NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_lb_app_profile" "test" {
  type         = "TCP"
  display_name = "my-tcp-profile"
}
```

## Argument Reference

* `id` - (Optional) The ID of Profile to retrieve.
* `type` - (Optional) Type of Profile to retrieve, one of `HTTP`, `TCP`, `UDP`, `ANY`.
* `display_name` - (Optional) The Display Name prefix of the Profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
