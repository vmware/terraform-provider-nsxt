---
subcategory: "Load Balancer"
layout: "nsxt"
page_title: "NSXT: policy_lb_client_ssl_profile"
description: Policy Load Balancer Client SSL Profile data source.
---

# nsxt_policy_lb_client_ssl_profile

This data source provides information about Policy Client SSL Profile for Load Balancer configured on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_lb_client_ssl_profile" "test" {
  display_name = "myprofile"
}
```

## Argument Reference

* `id` - (Optional) The ID of Profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
