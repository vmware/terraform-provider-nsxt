---
subcategory: "Policy - Segments"
layout: "nsxt"
page_title: "NSXT: policy_spoofguard_profile"
description: Policy SpoofGuardProfile data source.
---

# nsxt_policy_spoofguard_profile

This data source provides information about policy Spoofguard Profile configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_spoofguard_profile" "test" {
  display_name = "spoofguard-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of SpoofGuardProfile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the SpoofGuardProfile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
