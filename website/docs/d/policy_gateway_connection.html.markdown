---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: policy_gateway_connection"
description: Policy gateway connection data source.
---

# nsxt_policy_gateway_connection

This data source provides information about Gateway Connection on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_gateway_connection" "test" {
  display_name = "test"
  tier0_path   = data.nsxt_policy_tier0_gateway.path
}
```

## Argument Reference

* `id` - (Optional) The ID of transit gateway to retrieve.
* `display_name` - (Optional) The Display Name prefix of the transit gateway to retrieve.
* `tier0_path` - (Optional) Path of Tier0 for this connection

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
