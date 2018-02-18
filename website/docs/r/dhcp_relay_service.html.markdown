---
layout: "nsxt"
page_title: "NSXT: nsxt_dhcp_relay_service"
sidebar_current: "docs-nsxt-resource-dhcp-relay-service"
description: |-
  Provides a resource to configure dhcp relay service on NSX-T manager
---

# nsxt_dhcp_relay_service

Provides a resource to configure dhcp relay service on NSX-T manager

## Example Usage

```hcl
resource "nsxt_dhcp_relay_service" "DRS" {
    description = "DRS provisioned by Terraform"
    display_name = "DRS"
    tag {
        scope = "color"
        tag = "blue"
    }
    dhcp_relay_profile_id = ...
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this dhcp_relay_service.
* `dhcp_relay_profile_id` - (Required) dhcp relay profile referenced by the dhcp relay service.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the dhcp_relay_service.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `system_owned` - A boolean that indicates whether this resource is system-owned and thus read-only.
