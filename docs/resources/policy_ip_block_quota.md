---
subcategory: "Multitenancy"
page_title: "NSXT: nsxt_policy_ip_block_quota"
description: A resource to configure IP Block on NSX Policy.
---

# nsxt_policy_ip_block_quota

This resource provides a means to configure IP Block Quota on NSX Policy.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_ip_block_quota" "small" {
  display_name = "small-quota"
  description  = "terraform provisioned quota"

  quota {
    ip_block_paths        = [nsxt_policy_ip_block.all.path]
    ip_block_visibility   = "EXTERNAL"
    ip_block_address_type = "IPV4"

    single_ip_cidrs = 4
    other_cidrs {
      total_count = 4
    }
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
resource "nsxt_policy_ip_block_quota" "small" {
  context {
    project_id = "dev"
  }

  display_name = "small-quota"
  description  = "terraform provisioned quota"

  quota {
    ip_block_paths        = [nsxt_policy_ip_block.dev.path]
    ip_block_visibility   = "EXTERNAL"
    ip_block_address_type = "IPV6"

    single_ip_cidrs = -1
    other_cidrs {
      mask        = "/64"
      total_count = 4
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `context` - (Optional) The context which the object belongs to
          * `project_id` - (Required) The ID of the project which the object belongs to
* `display_name` - (Required) The display name for the IP Block.
* `description` - (Optional) Description of the resource.
* `quota` - (Required) Quota specification
    * `ip_block_paths` - (Optional) List of IP blocks that this quota applies to.
    * `ip_block_address_type` - (Required) One of `IPV4`, `IPV6`. A quota will be applied on block of same address type. One v4 block and another v6 block cannot be specified within the same quota.
    * `ip_block_visibility` - (Required) One of `EXTERNAL`, `PRIVATE`. A quota will be applied on blocks with same visibility. Private and External blocks cannot be specified within the same block.
    * `single_ip_cidrs` - (Optional) Quota for single IP CIDRs allowed. Default is -1 (unlimited).
    * `other_cidrs` - (Required) Quota for other cidrs
        * `mask` - (Optional) Largest mask size that is allowed, format: `/[size]`
        * `total_count` - (Optional) Number of CIDRs that can be allocated from the block. Default is -1 (unlimited).
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this IP Block.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the resource.

## Importing

An existing IP Block can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_ip_block_quota.quota1 POLICY_PATH
```

The above would import IP Block Quota as a resource named `quota1` with policy path `POLICY_PATH`.
