---
subcategory: "Policy - Segments"
layout: "nsxt"
page_title: "NSXT: policy_qos_profile"
description: Policy QosProfile data source.
---

# nsxt_policy_qos_profile

This data source provides information about policy Quality of Service Profile configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_qos_profile" "test" {
  display_name = "qos-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Profile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
