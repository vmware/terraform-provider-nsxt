---
subcategory: "VPN"
layout: "nsxt"
page_title: "NSXT: policy_l2_vpn_service"
description: Policy L2 VPN Service.
---

# nsxt_policy_vpn_service

This data source provides information about policy l2 vpn service configured on NSX.

This data source is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_l2_vpn_service" "test" {
  display_name = "l2_vpn_service1"
}
```

## Argument Reference

* `id` - (Optional) The ID of L2 VPN Service to retrieve.
* `display_name` - (Optional) The Display Name of the L2 VPN Service.
* `gateway_path` - (Optional) Gateway Path for this Service.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
