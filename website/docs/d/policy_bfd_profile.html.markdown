---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: policy_bfd_profile"
description: Policy BFD Profile data source.
---

# nsxt_policy_bfd_profile

This data source provides information about policy BFD Profile configured on NSX.
This data source is applicable to NSX Global Manager, NSX Policy Manager and VMC (NSX version 3.0.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_bfd_profile" "test" {
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
