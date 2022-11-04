---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_ip_block_subnet"
description: |-
  Provides a resource to configure IP block subnet on NSX-T manager
---

# nsxt_ip_block_subnet

Provides a resource to configure IP block subnet on NSX-T manager

## Example Usage

```hcl
resource "nsxt_ip_block" "ip_block" {
  display_name = "block1"
  cidr         = "55.0.0.0/24"
}

resource "nsxt_ip_block_subnet" "ip_block_subnet" {
  description  = "ip_block_subnet provisioned by Terraform"
  display_name = "ip_block_subnet"
  block_id     = nsxt_ip_block.ip_block.id
  size         = 16

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `description` - (Optional) Description of this resource.
* `block_id` - (Required) Block id for which the subnet is created.
* `size` - (Required) Represents the size or number of IP addresses in the subnet.
* `tag` - (Optional) A list of scope + tag pairs to associate with this IP block subnet.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP block subnet.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `allocation_range` - A collection of IPv4 IP ranges used for IP allocation.
* `cidr` - Represents the size or number of IP addresses in the subnet. All subnets of the same block must have the same size, which must be a power of 2.


## Importing

An existing IP block subnet can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_ip_block_subnet.ip_block_subnet UUID
```

The above would import the IP block subnet named `ip_block_subnet` with the nsx id `UUID`
