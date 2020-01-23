---
layout: "nsxt"
page_title: "NSXT: policy_lb_server_ssl_profile"
sidebar_current: "docs-nsxt-datasource-policy-lb-server-ssl-profile"
description: Policy Load Balancer Server SSL Profile data source.
---

# nsxt_policy_lb_server_ssl_profile

This data source provides information about policy Server SSL Profile for Load Balancer configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_lb_server_ssl_profile" "test" {
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
