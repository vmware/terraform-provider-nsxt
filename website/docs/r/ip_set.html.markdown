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
resource "nsxt_ip_set" "ip_set" {
    description = "IS provisioned by Terraform"
    display_name = "IS"
    tag {
        scope = "color"
        tag = "blue"
    }
    ip_addresses = ["1.1.1.1", "2.2.2.2"]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this ip_set.
* `ip_addresses` - (Optional) IP addresses.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the ip_set.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing IP set can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import nsxt_ip_set.x id
```

The above would import the IP set named `x` with the nsx id `id`
