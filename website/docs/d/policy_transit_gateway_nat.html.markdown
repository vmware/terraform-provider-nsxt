---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: policy_transit_gateway_nat"
description: Transit Gateway NAT data source.
---

# nsxt_policy_transit_gateway_nat

This data source provides information about an NAT section configured under Transit Gateway on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "proj" {
  display_name = "demoproj"
}

data "nsxt_policy_transit_gateway" "tgw1" {
  context {
    project_id = data.nsxt_policy_project.proj.id
  }
  display_name = "TGW1"
}

data "nsxt_policy_transit_gateway_nat" "test" {
  transit_gateway_path = data.nsxt_policy_transit_gateway.tgw1.path
}
```

## Argument Reference

* `transit_gateway_path` - (Required) Policy path of parent Transit Gateway
* `type` - (Optional) Type of NAT, one of `USER`, `DEFAULT`. Default is `USER`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - (Optional) The ID of the resource.
* `display_name` - (Optional) Display Name of the resource.
* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
