---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: policy_distributed_vlan_connection"
description: Policy distributed VLAN connection data source.
---

# nsxt_policy_distributed_vlan_connection

This data source provides information about Distributed VLAN Connection on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_distributed_vlan_connection" "test" {
  display_name = "test"
}
```

## Argument Reference

* `id` - (Optional) The ID of transit gateway to retrieve.
* `display_name` - (Optional) The Display Name prefix of the transit gateway to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
