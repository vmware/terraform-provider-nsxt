---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_vtep_ha_host_switch_profile"
description: A VTEP HA host switch profile data source.
---

# nsxt_policy_vtep_ha_host_switch_profile

This data source provides information about VTEP HA host switch profile configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_vtep_ha_host_switch_profile" "vtep_ha_host_switch_profile" {
  display_name = "vtep_ha_host_switch_profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of VTEP HA host switch profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the VTEP HA host switch profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
