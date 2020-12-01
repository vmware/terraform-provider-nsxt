---
subcategory: "Policy - DHCP"
layout: "nsxt"
page_title: "NSXT: policy_dhcp_server"
description: A policy DHCP server data source.
---

# nsxt_policy_dhcp_server

This data source provides information about policy DHCP servers configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_dhcp_server" "test" {
  display_name = "dhcp-2"
}
```

## Argument Reference

* `id` - (Optional) The ID of DHCP server to retrieve.

* `display_name` - (Optional) The Display Name prefix of DHCP server to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
