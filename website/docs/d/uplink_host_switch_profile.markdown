---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_uplink_host_switch_profile"
description: A host switch profile data source.
---

# nsxt_uplink_host_switch_profile

This data source provides information about uplink host switch profile configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_uplink_host_switch_profile" "uplink_host_switch_profile" {
  display_name = "uplink_host_switch_profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of uplink host switch profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the uplink host switch profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
