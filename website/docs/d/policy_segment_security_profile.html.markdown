---
subcategory: "Policy - Segments"
layout: "nsxt"
page_title: "NSXT: policy_segment_security_profile"
description: Policy SegmentSecurityProfile data source.
---

# nsxt_policy_segment_security_profile

This data source provides information about policy Segment Security Profile configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_segment_security_profile" "test" {
  display_name = "segment-security-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of SegmentSecurityProfile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the SegmentSecurityProfile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
