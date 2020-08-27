---
subcategory: "Policy - Segments"
layout: "nsxt"
page_title: "NSXT: policy_ip_discovery_profile"
description: Policy IP Discovery Profile data source.
---

# nsxt_policy_ip_discovery_profile

This data source provides information about Policy IP Discovery Profile configured on NSX.

This data source is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_ip_discovery_profile" "test" {
  display_name = "ip-discovery-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Profile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
