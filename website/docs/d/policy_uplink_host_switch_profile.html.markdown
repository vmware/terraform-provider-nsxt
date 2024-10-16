---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_uplink_host_switch_profile"
description: An uplink host switch profile data source.
---

# nsxt_policy_uplink_host_switch_profile

This data source provides information about uplink host switch profile configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_uplink_host_switch_profile" "uplink_host_switch_profile" {
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
* `realized_id` - The id of realized pool object. This id should be used in `nsxt_edge_transport_node` resource.
