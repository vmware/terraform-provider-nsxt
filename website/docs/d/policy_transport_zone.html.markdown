---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_transport_zone"
description: A Policy Transport Zone data source.
---

# nsxt_policy_transport_zone

This data source provides information about Policy based Transport Zones (TZ) configured in NSX. A Transport Zone defines the scope to which a network can extend in NSX. For example an overlay based Transport Zone is associated with both hypervisors and logical switches and defines which hypervisors will be able to serve the defined logical switch. Virtual machines on the hypervisor associated with a Transport Zone can be attached to logical switches in that same Transport Zone.

This data source is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_transport_zone" "overlay_transport_zone" {
  display_name = "1-transportzone-87"
}
```

```hcl
data "nsxt_policy_transport_zone" "vlan_transport_zone" {
  transport_type = "VLAN_BACKED"
  is_default     = true
}
```

Note: This usage is for Global Manager only.
```hcl
data "nsxt_policy_site" "paris" {
  display_name = "Paris"
}

data "nsxt_policy_transport_zone" "overlay_transport_zone" {
  display_name = "1-transportzone-87"
  site_path    = data.nsxt_policy_site.paris.path
}
```

## Argument Reference

* `id` - (Optional) The ID of Transport Zone to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Transport Zone to retrieve.
* `transport_type` - (Optional) Transport type of requested Transport Zone, one of `OVERLAY_STANDARD`, `OVERLAY_ENS`, `OVERLAY_BACKED`, `VLAN_BACKED` and `UNKNOWN`.
* `is_default` - (Optional) May be set together with `transport_type` in order to retrieve default Transport Zone for this transport type.
* `site_path` - (Optional) The path of the site which the Transport Zone belongs to, this configuration is required for global manager only. `path` field of the existing `nsxt_policy_site` can be used here.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the Transport Zone.
* `is_default` - A boolean flag indicating if this Transport Zone is the default.
* `transport_type` - The transport type of this transport zone.
* `path` - The NSX path of the policy resource.
* `realized_id` - The id of realized transport zone object. This id should be used in `nsxt_edge_transport_node` resource.
