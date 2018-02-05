---
layout: "nsxt"
page_title: "NSXT: nsxt_ip_set"
sidebar_current: "docs-nsxt-resource-ip-set"
description: |-
  Provides a resource to configure IP set on NSX-T manager
---

# nsxt_ip_set

Provides a resource to configure IP set on NSX-T manager

## Example Usage

```hcl
resource "nsxt_ip_set" "IS" {
  description = "IS provisioned by Terraform"
  display_name = "IS"
  tags = [{ scope = "color"
            tag = "red" }
  ]
  ip_addresses = ["1.1.1.1", "2.2.2.2"]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) Defaults to ID if not set.
* `tags` - (Optional) A list of scope + tag pairs to associate with this ip_set.
* `ip_addresses` - (Optional) IP addresses.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the ip_set.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `system_owned` - A boolean that indicates whether this resource is system-owned and thus read-only.
