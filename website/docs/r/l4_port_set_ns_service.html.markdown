---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_l4_port_set_ns_service"
description: A resource that can be used to configure a layer 4 networking and security service with ports in NSX.
---

# nsxt_l4_port_set_ns_service

This resource provides a way to configure a networking and security service which can be used within NSX. This specific service is for configuration of layer 4 ports.

## Example Usage

```hcl
resource "nsxt_l4_port_set_ns_service" "ns_service_l4" {
  description       = "S1 provisioned by Terraform"
  display_name      = "S1"
  protocol          = "TCP"
  destination_ports = ["73", "8080", "81"]

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
* `destination_ports` - (Optional) Set of destination ports.
* `source_ports` - (Optional) Set of source ports.
* `protocol` - (Required) L4 protocol. Accepted values - 'TCP' or 'UDP'.
* `tag` - (Optional) A list of scope + tag pairs to associate with this service.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the NS service.
* `default_service` - The default NSServices are created in the system by default. These NSServices can't be modified/deleted.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing L4 port set NS service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_l4_port_set_ns_service.ns_service_l4 UUID
```

The above command imports the layer 4 port based networking and security service named `ns_service_l4` with the NSX id `UUID`.
