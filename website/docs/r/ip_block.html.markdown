---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_ip_block"
description: |-
  Provides a resource to configure IP block on NSX-T manager
---

# nsxt_ip_block

Provides a resource to configure IP block on NSX-T manager

## Example Usage

```hcl
resource "nsxt_ip_block" "ip_block" {
  description  = "ip_block provisioned by Terraform"
  display_name = "ip_block"
  cidr         = "2.1.1.0/24"

  tag {
    scope = "color"
    tag   = "red"
  }
}

resource "nsxt_ip_block_subnet" "ip_block_subnet" {
  description = "ip_block_subnet"
  block_id    = nsxt_ip_block.ip_block.id
  size        = 16
}

```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `description` - (Optional) Description of this resource.
* `cidr` - (Required) Represents network address and the prefix length which will be associated with a layer-2 broadcast domain.
* `tag` - (Optional) A list of scope + tag pairs to associate with this IP block.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP block.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing IP block can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_ip_block.ip_block UUID
```

The above would import the IP block named `ip_block` with the nsx id `UUID`
