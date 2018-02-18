---
layout: "nsxt"
page_title: "NSXT: nsxt_ip_protocol_ns_service"
sidebar_current: "docs-nsxt-resource-ip-protocol-ns-service"
description: |-
  Provides a resource to configure NS service for IP protocol on NSX-T Manager.
---

# nsxt_ip_protocol_ns_service

Provides a resource to configure NS service for IP protocol on NSX-T Manager

## Example Usage

```hcl
resource "nsxt_ip_protocol_ns_service" "S1" {
    description = "S1 provisioned by Terraform"
    display_name = "S1"
    protocol = "10"
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
* `protocol` - (Required) IP protocol number (0-255)
* `tag` - (Optional) A list of scope + tag pairs to associate with this ip_set.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical switch.
* `default_service` - The default NSServices are created in the system by default. These NSServices can't be modified/deleted.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
