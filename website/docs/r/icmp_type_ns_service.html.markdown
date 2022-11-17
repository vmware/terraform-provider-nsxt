---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_icmp_type_ns_service"
description: A resource that can be used to configure an ICMP based networking and security service in NSX.
---

# nsxt_icmp_type_ns_service

This resource provides a way to configure a networking and security service which can be used within NSX. This specific service is for the ICMP protocol.

## Example Usage

```hcl
resource "nsxt_icmp_type_ns_service" "ns_service_icmp" {
  description  = "S1 provisioned by Terraform"
  display_name = "S1"
  protocol     = "ICMPv4"
  icmp_type    = "5"
  icmp_code    = "1"

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
* `protocol` - (Required) Version of ICMP protocol ICMPv4 or ICMPv6.
* `icmp_type` - (Optional) ICMP message type.
* `icmp_code` - (Optional) ICMP message code
* `tag` - (Optional) A list of scope + tag pairs to associate with this service.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the NS service.
* `default_service` - The default NSServices are created in the system by default. These NSServices can't be modified/deleted.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing ICMP type NS Service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_icmp_type_ns_service.x id
```

The above service imports the ICMP type network and security service named `x` with the NSX id `id`.
