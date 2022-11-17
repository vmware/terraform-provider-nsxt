---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_ip_protocol_ns_service"
description: A resource that can be used to configure an IP protocol based networking and security service in NSX.
---

# nsxt_ip_protocol_ns_service

This resource provides a way to configure a networking and security service which can be used within NSX. This specific service is for the IP protocol.

## Example Usage

```hcl
resource "nsxt_ip_protocol_ns_service" "ns_service_ip" {
  description  = "S1 provisioned by Terraform"
  display_name = "S1"
  protocol     = "10"

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description.
* `protocol` - (Required) IP protocol number (0-255)
* `tag` - (Optional) A list of scope + tag pairs to associate with this service.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the NS service.
* `default_service` - The default NSServices are created in the system by default. These NSServices can't be modified/deleted.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing IP protocol NS service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_ip_protocol_ns_service.ns_service_ip UUID
```

The above command imports the IP protocol based networking and security service named `ns_service_ip` with the NSX id `UUID`.
