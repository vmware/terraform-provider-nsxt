---
subcategory: "Policy - IPAM"
layout: "nsxt"
page_title: "NSXT: policy_ip_block"
description: Policy IP Block data source.
---

# nsxt_policy_ip_block

This data source provides information about Policy IP Blocks configured on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_ip_block" "test" {
  display_name = "ipblock1"
}
```

## Argument Reference

* `id` - (Optional) The ID of IP Block to retrieve.
* `display_name` - (Optional) The Display Name prefix of the IP Block to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
