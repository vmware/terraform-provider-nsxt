---
subcategory: "VPN"
page_title: "NSXT: policy_ipsec_vpn_ike_profile"
description: Policy IPSec VPN IKE Profile.
---

# nsxt_policy_ipsec_vpn_ike_profile

This data source provides information about policy IPSec VPN IKE Profile configured on NSX.

This data source is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_ipsec_vpn_ike_profile" "test" {
  display_name = "ipsec_vpn_ike_profile1"
}
```

## Argument Reference

* `context` - (Optional) Policy context for this data source.
  * `project_id` - (Optional) Project ID for multitenancy.
* `id` - (Optional) The ID of IPSec VPN IKE Profile to retrieve.
* `display_name` - (Optional) The Display Name of the IPSec VPN IKE Profile.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
