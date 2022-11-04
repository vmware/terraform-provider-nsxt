---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_logical_router_downlink_port"
description: A resource that can be used to configure logical router downlink port in NSX.
---

# nsxt_logical_router_downlink_port

This resource provides a means to define a downlink port on a logical router to connect a logical tier1 router to a logical switch. The result of this is to provide a default gateway to virtual machines running on the logical switch.

## Example Usage

```hcl
resource "nsxt_logical_router_downlink_port" "downlink_port" {
  description                   = "DP1 provisioned by Terraform"
  display_name                  = "DP1"
  logical_router_id             = nsxt_logical_tier1_router.rtr1.id
  linked_logical_switch_port_id = nsxt_logical_port.logical_port1.id
  ip_address                    = "1.1.0.1/24"

  service_binding {
    target_id   = nsxt_dhcp_relay_service.dr_service.id
    target_type = "LogicalService"
  }

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `logical_router_id` - (Required) Identifier for logical Tier-1 router on which this port is created
* `linked_logical_switch_port_id` - (Required) Identifier for port on logical switch to connect to
* `ip_address` - (Required) Logical router port subnet (ip_address / prefix length)
* `urpf_mode` - (Optional) Unicast Reverse Path Forwarding mode. Accepted values are "NONE" and "STRICT" which is the default value.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this port.
* `service_binding` - (Optional) A list of services for this port. Currently only "LogicalService" is supported as a target_type, and a DHCP relay service ID as target_id

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical router downlink port.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `mac_address` - The MAC address assigned to this port

## Importing

An existing logical router downlink port can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_logical_router_downlink_port.downlink_port UUID
```

The above command imports the logical router downlink port named `downlink_port` with the NSX id `UUID`.
