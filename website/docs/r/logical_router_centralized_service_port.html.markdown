---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_logical_router_centralized_service_port"
description: A resource that can be used to configure logical router centralized service port in NSX.
---

# nsxt_logical_router_centralized_service_port

This resource provides a means to define a centralized service port on a logical router to connect a logical tier0 or tier1 router to a logical switch. This allows the router to be used for E-W load balancing

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_logical_router_centralized_service_port" "cs_port" {
  description                   = "Centralized service port provisioned by Terraform"
  display_name                  = "CSP1"
  logical_router_id             = nsxt_logical_tier1_router.rtr1.id
  linked_logical_switch_port_id = nsxt_logical_port.logical_port1.id
  ip_address                    = "1.1.0.1/24"

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `logical_router_id` - (Required) Identifier for logical Tier-0 or Tier-1 router on which this port is created
* `linked_logical_switch_port_id` - (Required) Identifier for port on logical switch to connect to
* `ip_address` - (Required) Logical router port subnet (ip_address / prefix length)
* `urpf_mode` - (Optional) Unicast Reverse Path Forwarding mode. Accepted values are "NONE" and "STRICT" which is the default value.
* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this port.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical router centralized service port.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing logical router centralized service port can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_logical_router_centralized_service_port.cs_port UUID
```

The above command imports the logical router centralized service port named `cs_port` with the NSX id `UUID`.
