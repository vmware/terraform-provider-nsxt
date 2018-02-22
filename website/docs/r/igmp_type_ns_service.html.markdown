---
layout: "nsxt"
page_title: "NSXT: nsxt_igmp_type_ns_service"
sidebar_current: "docs-nsxt-resource-igmp-type-ns-service"
description: |-
  Provides a resource to configure NS service for IGMP on NSX-T Manager.
---

# nsxt_igmp_type_ns_service

Provides a resource to configure NS service for IGMP on NSX-T Manager

## Example Usage

```hcl
resource "nsxt_igmp_type_ns_service" "ns_service_igmp" {
    description = "S1 provisioned by Terraform"
    display_name = "S1"
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
* `tag` - (Optional) A list of scope + tag pairs to associate with this service.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical switch.
* `default_service` - The default NSServices are created in the system by default. These NSServices can't be modified/deleted.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing IGMP type NS Service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import nsxt_igmp_type_ns_service.x id
```

The above would import the IGMP type NS Service named `x` with the nsx id `id`
