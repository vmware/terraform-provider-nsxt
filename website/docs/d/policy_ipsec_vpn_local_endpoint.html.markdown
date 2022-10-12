---
subcategory: "Policy - Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: policy_local_endpoint"
description: Policy Local Endpoint data source.
---

# nsxt_policy_local_endpoint

This data source provides information about policy local endpoint configured on NSX.

This data source is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_local_endpoint" "test" {
  display_name = "local_endpoint1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Local Endpoint to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Local Endpoint to retrieve. With VMC on AWS, the user can either choose the "Private IP1" or "Public IP1" values.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
