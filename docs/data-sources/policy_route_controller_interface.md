---
subcategory: "Gateways and Routing"
page_title: "NSXT: nsxt_policy_route_controller_interface"
description: A data source to retrieve a Route Controller Interface.
---

# nsxt_policy_route_controller_interface

This data source provides information about a Route Controller Interface configured in NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_route_controller_interface" "test" {
  display_name = "my-interface"
  parent_path  = "/infra/route-controllers/my-controller"
}
```

## Argument Reference

* `id` - (Optional) The ID of the Interface to retrieve.
* `display_name` - (Optional) The Display Name of the Interface to retrieve.
* `parent_path` - (Optional) Policy path of the parent Route Controller. Narrows the search to interfaces belonging to this Route Controller.
* `description` - (Optional) The description of the Interface to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `path` - The NSX path of the policy resource.
