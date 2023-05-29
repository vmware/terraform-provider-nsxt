---
subcategory: "Segments"
layout: "nsxt"
page_title: "NSXT: policy_bridge_profile"
description: Policy Bridge Profile data source.
---

# nsxt_policy_bridge_profile

This data source provides information about Edge Bridge Profile configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_bridge_profile" "test" {
  display_name = "profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Profile to retrieve. If ID is specified, no additional argument should be configured.
* `display_name` - (Optional) The Display Name prefix of the Profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
