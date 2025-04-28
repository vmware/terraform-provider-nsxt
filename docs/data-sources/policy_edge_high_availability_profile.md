---
subcategory: "Fabric"
page_title: "NSXT: nsxt_policy_edge_high_availability_profile"
description: A policy Edge high availability profile data source.
---

# nsxt_policy_edge_high_availability_profile

This data source provides information about policy Edge high availability profile configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_edge_high_availability_profile" "edge_ha_profile" {
  display_name = "edge_ha_profile1"
  site_path    = data.nsxt_policy_site.paris.path
}
```

## Argument Reference

* `id` - (Optional) The ID of policy Edge high availability profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Edge high availability profile to retrieve.
* `site_path` - (Optional) The path of the site which the Edge high availability profile belongs to.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `unique_id` - A unique identifier assigned by the system.
