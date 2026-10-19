---
subcategory: "Multitenancy"
page_title: "NSXT: nsxt_policy_shared_resource"
description: A resource to configure a Shared Resource.
---

# nsxt_policy_shared_resource

This resource provides a method for the management of a Shared Resource.
This resource is applicable to NSX Policy Manager.
This resource is supported with NSX 4.1.1 onwards.

## Example Usage

```hcl
resource "nsxt_policy_shared_resource" "test" {
  display_name = "test"
  description  = "Terraform provisioned Shared Resource"

  share_path = nsxt_policy_share.share.path
  resource_object {
    resource_path = data.nsxt_policy_context_profile.cprof.path
  }
}
```

## Example Usage - Associate an IP Block with a Project

IP blocks must use `ASSOCIATION_ONLY` permission, `NONE_DESCENDANTS` sharing,
and `include_children = false`. Only nonexclusive EXTERNAL IP blocks can be
custom shared. Create PRIVATE IP blocks in the consuming project's context.

Given an existing project and EXTERNAL IP block:

```hcl
resource "nsxt_policy_share" "ip_block" {
  display_name     = "project-ip-block"
  sharing_strategy = "NONE_DESCENDANTS"
  shared_with      = [nsxt_policy_project.project1.path]
}

resource "nsxt_policy_shared_resource" "ip_block" {
  display_name    = "project-ip-block"
  share_path      = nsxt_policy_share.ip_block.path
  with_permission = "ASSOCIATION_ONLY"

  resource_object {
    resource_path    = nsxt_policy_ip_block.external.path
    include_children = false
  }
}
```

When a project-scoped quota references this block, set
`depends_on = [nsxt_policy_shared_resource.ip_block]` on the quota to ensure
the association exists before quota creation and is removed after quota deletion.
Alternatively, associate an external IPv4 block using the project's
`external_ipv4_blocks` attribute, which creates a system-managed share.

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `share_path` - (Required) Share policy path to associate the resource to.
* `with_permission` - (Optional, Computed) Access granted to the target of the share. Accepted values are `READ_ONLY`, `ASSOCIATION_ONLY`, and `READ_AND_ASSOCIATION`. When omitted, the backend chooses the default and the provider records it in state. Permission availability and semantics depend on the NSX version and enabled features; `READ_AND_ASSOCIATION` requires support for the NSX 9.2 sharing enhancements. IP-block association requires `ASSOCIATION_ONLY`.
* `resource_object` _ (Required) List of resources to be shared.
    * `include_children` - (Optional) Denotes if the children of the shared path are also shared.
    * `resource_path` - (Required) Path of the resource to be shared.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_shared_resource.test PATH
```

The above command imports Shared Resource named `test` with the policy path `PATH`.
