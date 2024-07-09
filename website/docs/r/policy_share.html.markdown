---
subcategory: "Multitenancy"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_share"
description: A resource to configure a Share.
---

# nsxt_policy_share

This resource provides a method for the management of a Share.
This resource is applicable to NSX Policy Manager.
This resource is supported with NSX 4.1.1 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "test" {
  display_name = "project1"
}

resource "nsxt_policy_share" "test" {
  display_name = "test"
  description  = "Terraform provisioned Share"
  shared_with  = [data.nsxt_policy_project.test.path]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `shared_with` - (Required) Path of the context
* `sharing_strategy` - (Optional) Sharing Strategy. Accepted values: 'NONE_DESCENDANTS', 'ALL_DESCENDANTS'. Default: 'NONE_DESCENDANTS'.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_share.test PATH
```

The above command imports Share named `test` with the policy path `PATH`.
