---
layout: "nsxt"
page_title: "NSXT: nsxt_logical_switch"
sidebar_current: "docs-nsxt-resource-logical-switch"
description: |-
  Provides a resource to configure Logical Switch (LS) on NSX-T Manager.
---

# nsxt_logical_switch

Provides a resource to configure Logical Switch (LS) on NSX-T Manager.

## Example Usage

```hcl
resource "nsxt_logical_switch" "switch1" {
  admin_state       = "UP"
  description       = "LS1 provisioned by Terraform"
  display_name      = "LS1"
  transport_zone_id = "${data.nsxt_transport_zone.transport_zone.id}"
  replication_mode  = "MTEP"

  tag {
    scope = "color"
    tag   = "blue"
  }

  switching_profile_id {
    key   = "${data.nsxt_switching_profile.qos_profiles.resource_type}"
    value = "${data.nsxt_switching_profile.qos_profiles.id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `transport_zone_id` - (Required) Transport Zone ID for the logical switch.
* `admin_state` - (Optional) Admin state for the logical switch. Accepted values - 'UP' or 'DOWN'. The default value is 'UP'.
* `replication_mode` - (Optional) Replication mode of the Logical Switch. Accepted values - 'MTEP' (Hierarchical Two-Tier replication) and 'SOURCE' (Head Replication), with 'MTEP' being the default value. Applies to overlay logical switches.
* `switching_profile_id` - (Optional) List of IDs of switching profiles (of various types) to be associated with this switch. Default switching profiles will be used if not specified.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `ip_pool_id` - (Optional) Ip Pool ID to be associated with the logical switch.
* `mac_pool_id` - (Optional) Mac Pool ID to be associated with the logical switch.
* `vlan` - (Optional) Vlan for vlan logical switch. If not specified, this switch is overlay logical switch.
* `vni` - (Optional) Vni for the logical switch.
* `address_binding` - (Optional) List of Address Bindings for the logical switch. This setting allows to provide bindings between IP address, mac Address and vlan.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical switch.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical switch.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing X can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import nsxt_X.x id
```

The above would import the X named `x` with the nsx id `id`
