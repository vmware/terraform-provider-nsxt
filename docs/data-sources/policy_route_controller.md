---
subcategory: "Gateways and Routing"
page_title: "NSXT: policy_route_controller"
description: Policy route controller data source.
---

# nsxt_policy_route_controller

This data source provides information about a Route Controller on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller" "test" {
  display_name = "test"
}
```

## Argument Reference

* `id` - (Optional) The ID of route controller to retrieve.
* `display_name` - (Optional) The Display Name prefix of the route controller to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
