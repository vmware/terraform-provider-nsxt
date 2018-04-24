---
layout: "nsxt"
page_title: "NSXT: nsxt_ip_block"
sidebar_current: "docs-nsxt-resource-ip-block"
description: |-
  Provides a resource to configure ip block on NSX-T manager
---

# nsxt_ip_block

Provides a resource to configure ip block on NSX-T manager

## Example Usage

```hcl
resource "nsxt_ip_block" "ip_block" {
  description = "ip_block provisioned by Terraform"
  display_name = "ip_block"
  cidr         = "2.1.1.0/24"

  subnet {
    display_name = "sub1"
    description  = "subnet"
    size         = "16"
  }

  tag = {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `description` - (Optional) Description of this resource.
* `cidr` - (Required) Represents network address and the prefix length which will be associated with a layer-2 broadcast domain.
* `tag` - (Optional) A list of scope + tag pairs to associate with this ip block.
* `subnet` - (Optional) List of subnets. Please note that a subnet can be added or removed (unless in use), but not updated. Each subnet has the following arguments:
  * `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
  * `description` - (Optional) Description of this resource.
  * `size` - (Required) Represents the size or number of ip addresses in the subnet. The size of all subnets in a block should be the same, and a power of 2

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the ip_block.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing ip block can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_ip_block.ip_block UUID
```

The above would import the ip block named `ip_block` with the nsx id `UUID`
