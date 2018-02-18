---
layout: "nsxt"
page_title: "NSXT: nsxt_dhcp_relay_profile"
sidebar_current: "docs-nsxt-resource-dhcp-relay-profile"
description: |-
  Provides a resource to configure dhcp relay profile on NSX-T manager
---

# nsxt_dhcp_relay_profile

Provides a resource to configure dhcp relay profile on NSX-T manager

## Example Usage

```hcl
resource "nsxt_dhcp_relay_profile" "DRP" {
  description = "DRP provisioned by Terraform"
  display_name = "DRP"
    tag {
        scope = "color"
        tag = "red"
    }
  server_addresses = ["1.1.1.1"]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this dhcp_relay_profile.
* `server_addresses` - (Required) IP addresses of the DHCP relay servers.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the dhcp_relay_profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `system_owned` - A boolean that indicates whether this resource is system-owned and thus read-only.
