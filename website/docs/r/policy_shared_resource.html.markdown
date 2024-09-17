---
subcategory: "Multitenancy"
layout: "nsxt"
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

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `share_path` - (Required) Share policy path to associate the resource to.
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

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_shared_resource.test PATH
```

The above command imports Shared Resource named `test` with the policy path `PATH`.
