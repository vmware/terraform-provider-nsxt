---
subcategory: "Policy - Firewall"
layout: "nsxt"
page_title: "NSXT: policy_intrusion_service_profile"
description: Policy Intrusion Service Profile data source.
---

# nsxt_policy_intrusion_service_profile

This data source provides information about policy IDS Profile configured on NSX.
This data source is applicable to NSX Policy Manager and VMC (NSX version 3.0.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}
```

## Argument Reference

* `id` - (Optional) The ID of Profile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
