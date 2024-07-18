---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_transport_zone"
description: A resource to configure Policy Transport Zone.
---

# nsxt_policy_transport_zone

This resource provides a method for the management of Policy based Transport Zones (TZ). A Transport Zone defines the scope to which a network can extend in NSX. For example an overlay based Transport Zone is associated with both hypervisors and segments and defines which hypervisors will be able to serve the defined segment. Virtual machines on the hypervisor associated with a Transport Zone can be attached to segments in that same Transport Zone.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_transport_zone" "overlay_transport_zone" {
  display_name   = "1-transportzone-87"
  transport_type = "OVERLAY_BACKED"
}
```

```hcl
resource "nsxt_policy_transport_zone" "vlan_transport_zone" {
  display_name                = "1-transportzone-87"
  description                 = "VLAN transport zone"
  transport_type              = "VLAN_BACKED"
  is_default                  = true
  uplink_teaming_policy_names = ["teaming-1"]
  site_path                   = "/infra/sites/default"
  enforcement_point           = "default"

  tag {
    scope = "app"
    tag   = "web"
  }
}
```

## Argument Reference

* `display_name` - (Optional) The Display Name of the Transport Zone.
* `description` - (Optional) Description of the Transport Zone.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `transport_type` - (Required) Transport type of requested Transport Zone, one of `OVERLAY_STANDARD`, `OVERLAY_ENS`, `OVERLAY_BACKED`, `VLAN_BACKED` and `UNKNOWN`.
* `is_default` - (Optional) Set this Transport Zone as the default zone of given `transport_type`. Default value is `false`. When setting a Transport Zone with `is_default`: `true`, no existing Transport Zone of same `transport_type` should be set as default.
* `uplink_teaming_policy_names` - (Optional) The names of switching uplink teaming policies that all transport nodes in this transport zone support. Uplinkin teaming policies are only valid for `VLAN_BACKED` transport zones. 
* `site_path` - (Optional) The path of the site which the Transport Zone belongs to. `path` field of the existing `nsxt_policy_site` can be used here.
* `enforcement_point` - (Optional) The ID of enforcement point under given `site_path` to manage the Transport Zone.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `realized_id` - Realized ID for the transport zone. For reference in fabric resources (such as `transport_node`), `realized_id` should be used rather than `id`.

## Importing

An existing Transport Zone can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import
```
terraform import nsxt_policy_transport_zone.overlay-tz POLICY_PATH
```
The above command imports the Transport Zone named `overlay-tz` with the policy path `POLICY_PATH`.
