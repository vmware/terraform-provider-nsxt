---
subcategory: "VPN"
page_title: "NSXT: policy_ipsec_vpn_dpd_profile"
description: Policy IPSec VPN DPD profile.
---

# nsxt_policy_ipsec_vpn_dpd_profile

This data source provides information about policy IPSec VPN Dead Peer Detection Profile configured on NSX.

This data source is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_ipsec_vpn_dpd_profile" "test" {
  display_name = "ipsec_vpn_dpd_profile"
}
```

## Argument Reference

* `context` - (Optional) Policy context for this data source.
  * `project_id` - (Optional) Project ID for multitenancy.
* `id` - (Optional) The ID of IPSec VPN DPD profile to retrieve.
* `display_name` - (Optional) The Display Name of the IPSec VPN DPD profile.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
