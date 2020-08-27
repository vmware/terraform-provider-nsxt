---
subcategory: "Policy - IPAM"
layout: "nsxt"
page_title: "NSXT: policy_ip_pool"
description: Policy IP Pool Config data source.
---

# nsxt_policy_ip_pool

This data source provides information about policy IP Pools configured on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_ip_pool" "test" {
  display_name = "ippool1"
}
```

## Argument Reference

* `id` - (Optional) The ID of IP Pool Config to retrieve.
* `display_name` - (Optional) The Display Name prefix of the IP Pool Config to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
