---
subcategory: "Policy - Load Balancer"
layout: "nsxt"
page_title: "NSXT: policy_lb_service"
description: Policy Load Balancer Service data source.
---

# nsxt_policy_lb_service

This data source provides information about Policy Service for Load Balancer configured on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_lb_service" "test" {
  display_name = "myservice"
}
```

## Argument Reference

* `id` - (Optional) The ID of Service to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Service to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
