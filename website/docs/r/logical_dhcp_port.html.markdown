---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_logical_dhcp_port"
description: A resource that can be used to configure a Logical DHCP Port in NSX.
---

# nsxt_logical_dhcp_port

This resource provides a resource to configure a logical port on a logical switch, and attach it to a DHCP server.

## Example Usage

```hcl
resource "nsxt_logical_dhcp_server" "logical_dhcp_server" {
  display_name    = "logical_dhcp_server"
  dhcp_profile_id = nsxt_dhcp_server_profile.PRF.id
  dhcp_server_ip  = "1.1.1.10/24"
  gateway_ip      = "1.1.1.20"
}

resource "nsxt_logical_switch" "switch" {
  display_name      = "LS1"
  admin_state       = "UP"
  transport_zone_id = data.nsxt_transport_zone.transport_zone.id
}

resource "nsxt_logical_dhcp_port" "dhcp_port" {
  admin_state       = "UP"
  description       = "LP1 provisioned by Terraform"
  display_name      = "LP1"
  logical_switch_id = nsxt_logical_switch.switch.id
  dhcp_server_id    = nsxt_logical_dhcp_server.logical_dhcp_server.id

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of this resource.
* `logical_switch_id` - (Required) Logical switch ID for the logical port.
* `dhcp_server_id` - (Required) Logical DHCP server ID for the logical port.
* `admin_state` - (Optional) Admin state for the logical port. Accepted values - 'UP' or 'DOWN'. The default value is 'UP'.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical port.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical DHCP port.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing DHCP Logical Port can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_logical_dhcp_port.dhcp_port UUID
```

The above command imports the logical DHCP port named `dhcp_port` with the NSX id `UUID`.
