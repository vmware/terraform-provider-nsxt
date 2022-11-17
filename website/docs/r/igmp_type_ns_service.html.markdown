---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_igmp_type_ns_service"
description: A resource that can be used to configure an IGMP based networking and security service in NSX.
---

# nsxt_igmp_type_ns_service

This resource provides a way to configure a networking and security service which can be used within NSX. This specific service is for the IGMP protocol.

## Example Usage

```hcl
resource "nsxt_igmp_type_ns_service" "ns_service_igmp" {
  description  = "S1 provisioned by Terraform"
  display_name = "S1"

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
* `tag` - (Optional) A list of scope + tag pairs to associate with this service.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the NS service.
* `default_service` - The default NSServices are created in the system by default. These NSServices can't be modified/deleted.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing IGMP type NS Service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_igmp_type_ns_service.ns_service_igmp UUID
```

The above command imports the IGMP based networking and security service named `ns_service_igmp` with the NSX id `UUID`.
