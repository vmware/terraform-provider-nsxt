---
subcategory: "VPN"
page_title: "NSXT: policy_ipsec_vpn_tunnel_profile"
description: Policy IPSec VPN Tunnel Profile.
---

# nsxt_policy_ipsec_vpn_tunnel_profile

This data source provides information about policy IPSec VPN Tunnel Profile configured on NSX.

This data source is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_ipsec_vpn_tunnel_profile" "test" {
  display_name = "ipsec_vpn_tunnel_profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of IPSec VPN Tunnel Profile to retrieve.
* `display_name` - (Optional) The Display Name of the IPSec VPN Tunnel Profile.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
