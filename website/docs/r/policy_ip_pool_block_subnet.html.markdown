---
subcategory: "IPAM"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ip_pool_block_subnet"
description: A resource to configure IP Pool Block Subnets in NSX Policy.
---

# nsxt_policy_ip_pool_block_subnet

This resource provides a means to configure IP Pool Block Subnets on NSX Policy.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_ip_pool_block_subnet" "block_subnet1" {
  display_name        = "block-subnet1"
  pool_path           = nsxt_policy_ip_pool.pool1.path
  block_path          = nsxt_policy_ip_block.block1.path
  size                = 8
  auto_assign_gateway = false

  tag {
    scope = "color"
    tag   = "blue"
  }

  tag {
    scope = "env"
    tag   = "test"
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_ip_pool_block_subnet" "block_subnet1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name        = "block-subnet1"
  pool_path           = nsxt_policy_ip_pool.pool1.path
  block_path          = nsxt_policy_ip_block.block1.path
  size                = 8
  auto_assign_gateway = false

  tag {
    scope = "color"
    tag   = "blue"
  }

  tag {
    scope = "env"
    tag   = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The display name for the Block Subnet.
* `pool_path` - (Required) The Policy path to the IP Pool for this Block Subnet.
* `block_path` - (Required) The Policy path to the IP Block for this Block Subnet.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `size` - (Required) The size of this Block Subnet. Must be a power of 2
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Block Subnet.
* `auto_assign_gateway` - (Optional) A boolean flag to toggle auto-assignment of the Gateway IP for this Subnet

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of this Block Subnet.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the resource.

## Importing

An existing Block can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ip_pool_block_subnet.block_subnet1 pool-id/subnet-id
```
The above would import NSX Block Subnet as a resource named `block_subnet1` with the NSX ID `subnet-id` in the IP Pool `pool-id`, where `subnet-id` is NSX ID of Block Subnet and `pool-id` is the IP Pool ID the Subnet is in.

```
terraform import nsxt_policy_ip_pool_block_subnet.block_subnet1 POLICY_PATH
```
The above would import NSX Block Subnet as a resource named `block_subnet1` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
