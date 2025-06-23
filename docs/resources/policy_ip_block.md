---
subcategory: "IPAM"
page_title: "NSXT: nsxt_policy_ip_block"
description: A resource to configure IP Block on NSX Policy.
---

# nsxt_policy_ip_block

This resource provides a means to configure IP Block on NSX Policy.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_ip_block" "block1" {
  display_name = "ip-block1"
  cidr         = "192.168.1.0/24"
  visibility   = "PRIVATE"

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

## Example Usage - CIDR List

```hcl
resource "nsxt_policy_ip_block" "block1" {
  display_name = "ip-block1"
  cidr_list    = ["192.168.1.0/24"]
  visibility   = "PRIVATE"

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

resource "nsxt_policy_ip_block" "block1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "ip-block1"
  cidr         = "192.168.1.0/24"
  visibility   = "PRIVATE"

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

* `display_name` - (Required) The display name for the IP Block.
* `description` - (Optional) Description of the resource.
* `cidr` - (Optional) Network address and the prefix length which will be associated with a layer-2 broadcast domain. Either cidr or cidr_list should be set.
* `cidr_list` - (Optional) Array of contiguous IP address spaces represented by network address and prefix length. This attribute is supported starting from NSX version 9.1.0. Either cidr or cidr_list should be set.
* `is_subnet_exclusive` - (Optional) If this property is set to true, then this block is reserved for direct vlan extension use case. This attribute is supported starting from NSX version 9.1.0.
* `range_list` - (Optional) Represents list of IP address ranges in the form of start and end IPs. This attribute is supported with NSX 9.1.0 onwards.
  * `start` - (Required) The start IP address for the allocation range.
  * `end` - (Required) The end IP address for the allocation range.
* `reserved_ips` - (Optional) Represents list of reserved IP address in the form of start and end IPs. This attribute is supported with NSX 9.1.0 onwards.
  * `start` - (Required) The start IP address for the allocation range.
  * `end` - (Required) The end IP address for the allocation range.
* `visibility` - (Optional) Visibility of the IP Block. Valid options are `PRIVATE`, `EXTERNAL` or unset. Visibility cannot be changed once the block is associated with other resources.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this IP Block.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP Block.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the resource.

## Importing

An existing IP Block can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_ip_block.block1 ID
```

The above would import NSX IP Block as a resource named `block1` with the NSX id `ID`, where `ID` is NSX ID of the IP Block.

```shell
terraform import nsxt_policy_ip_block.block1 POLICY_PATH
```

The above would import NSX IP Block as a resource named `block1` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
