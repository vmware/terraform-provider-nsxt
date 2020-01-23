---
layout: "nsxt"
page_title: "NSXT: policy_spoofguard_profile"
sidebar_current: "docs-nsxt-datasource-policy-spoofguard-profile"
description: Policy SpoofGuardProfile data source.
---

# nsxt_policy_spoofguard_profile

This data source provides information about policy Spoofguard Profile configured on NSX.

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
