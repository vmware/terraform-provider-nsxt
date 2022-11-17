---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_vlan_logical_switch"
description: A resource to configure vlan logical switch in NSX.
---

# nsxt_vlan_logical_switch

This resource provides a method to create vlan logical switch in NSX. Virtual machines can then be connected to the appropriate logical switch for the desired topology and network connectivity.

## Example Usage

```hcl
resource "nsxt_vlan_logical_switch" "switch1" {
  admin_state       = "UP"
  description       = "LS1 provisioned by Terraform"
  display_name      = "LS1"
  transport_zone_id = data.nsxt_transport_zone.vlan_transport_zone.id
  vlan              = 2

  tag {
    scope = "color"
    tag   = "blue"
  }

  switching_profile_id {
    key   = data.nsxt_switching_profile.qos_profiles.resource_type
    value = data.nsxt_switching_profile.qos_profiles.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `transport_zone_id` - (Required) Transport Zone ID for the logical switch.
* `admin_state` - (Optional) Admin state for the logical switch. Accepted values - 'UP' or 'DOWN'. The default value is 'UP'.
* `vlan` - (Required) Vlan for the logical switch.
* `switching_profile_id` - (Optional) List of IDs of switching profiles (of various types) to be associated with this switch. Default switching profiles will be used if not specified.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `ip_pool_id` - (Optional) Ip Pool ID to be associated with the logical switch.
* `mac_pool_id` - (Optional) Mac Pool ID to be associated with the logical switch.
* `address_binding` - (Optional) List of Address Bindings for the logical switch. This setting allows to provide bindings between IP address, mac Address and vlan.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical switch.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical switch.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing X can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/docs/import

```
terraform import nsxt_vlan_logical_switch.switch1 UUID
```

The above command imports the logical switch named `switch1` with the NSX id `UUID`.
