---
subcategory: "VPC"
page_title: "NSXT: policy_transit_gateway_attachment"
description: A policy Transit Gateway Attachment data source.
---

# nsxt_policy_transit_gateway_attachment

This data source provides information about a policy Transit Gateway Attachment attached to a Policy Transit Gateway.

This data source is applicable to NSX Policy Manager and requires NSX 9.0.0 or higher.

## Example Usage

```hcl
data "nsxt_policy_transit_gateway_attachment" "test" {
  display_name = "my-tgw-attachment"
  parent_path  = data.nsxt_policy_transit_gateway.test.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the Transit Gateway Attachment to retrieve.
* `display_name` - (Optional) The Display Name of the Transit Gateway Attachment to retrieve.
* `parent_path` - (Required) Policy path of the parent Transit Gateway.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
