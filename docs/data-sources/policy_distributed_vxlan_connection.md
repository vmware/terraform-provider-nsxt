---
subcategory: "EVPN"
page_title: "NSXT: policy_distributed_vxlan_connection"
description: Policy distributed VXLAN connection data source.
---

# nsxt_policy_distributed_vxlan_connection

This data source provides information about Distributed VXLAN Connection on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_distributed_vxlan_connection" "test" {
  display_name = "test"
}
```

## Argument Reference

* `id` - (Optional) The ID of distributed VXLAN connection to retrieve.
* `display_name` - (Optional) The Display Name prefix of the distributed VXLAN connection to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
