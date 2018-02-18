---
layout: "nsxt"
page_title: "NSXT: nsxt_l4_port_set_ns_service"
sidebar_current: "docs-nsxt-resource-l4-port-set-ns-service"
description: |-
  Provides a resource to configure NS service for L4 protocol with ports on NSX-T Manager.
---

# nsxt_l4_port_set_ns_service

Provides a resource to configure NS service for L4 protocol with ports on NSX-T Manager

## Example Usage

```hcl
resource "nsxt_l4_port_set_ns_service" "S1" {
    description = "S1 provisioned by Terraform"
    display_name = "S1"
    protocol = "TCP"
    destination_ports = [ "73", "8080", "81"]
    tag {
        scope = "color"
        tag = "blue"
    }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description.
* `destination_ports` - (Optional) Set of destination ports.
* `source_ports` - (Optional) Set of source ports.
* `protocol` - (Optional) VL4 protocol
* `tag` - (Optional) A list of scope + tag pairs to associate with this ip_set.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical switch.
* `default_service` - The default NSServices are created in the system by default. These NSServices can't be modified/deleted.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
