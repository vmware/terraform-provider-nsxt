---
subcategory: "Gateways and Routing"
page_title: "NSXT: policy_gateway_prefix_list"
description: Policy Gateway Prefix List data source.
---

# nsxt_policy_gateway_prefix_list

This data source provides information about a policy Gateway Prefix List configured on NSX.

This data source is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "gw1" {
  display_name = "gw1"
}

data "nsxt_policy_gateway_prefix_list" "prefix_list1" {
  gateway_path = data.nsxt_policy_tier0_gateway.gw1.path
  display_name = "t0_prefix_list"
}
```

## Argument Reference

* `id` - (Optional) The ID of the Gateway Prefix List to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Gateway Prefix List to retrieve.
* `gateway_path` - (Optional) Policy path of the Gateway where the Prefix List is configured.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
