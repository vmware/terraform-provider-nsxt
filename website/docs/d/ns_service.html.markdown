---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: ns_service"
description: A networking and security service data source.
---

# nsxt_ns_service

This data source provides information about a network and security (NS) service configured in NSX. NS services are either factory defined in NSX or can be defined by the NSX administrator. They provide a convenience name for a port/protocol pair that is often used in fire walling or load balancing.

## Example Usage

```hcl
data "nsxt_ns_service" "ns_service_dns" {
  display_name = "DNS"
}
```

## Argument Reference

* `id` - (Optional) The ID of NS service to retrieve
* `display_name` - (Optional) The Display Name of the NS service to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the NS service.
