---
subcategory: "Policy - Multi Tenancy"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_project"
description: A resource to configure a Project.
---

# nsxt_policy_project

This resource provides a method for the management of a Project.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_project" "test" {
  display_name        = "test"
  description         = "Terraform provisioned Project"
  _default            = true
  short_id            = "test"
  tier0_gateway_paths = ["/infra/tier-0s/test"]

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `_default` - (Optional) The server will populate this field when returing the resource. Ignored on PUT and POST.true - the project is a default project. Default projects are non-editable, system create ones.
* `short_id` - (Optional) Defaults to id if id is less than equal to 8 characters or defaults to random generated id if not set.
* `site_info` - (Optional) Information related to sites applicable for given Project. For on-prem deployment, only 1 is allowed.
  * `edge_cluster_paths` - (Optional) The edge cluster on which the networking elements for the Org will be created.
  * `site_path` - (Optional) This represents the path of the site which is managed by Global Manager. For the local manager, if set, this needs to point to 'default'.
* `tier0_gateway_paths` - (Optional) The tier 0 has to be pre-created before Project is created. The tier 0 typically provides connectivity to external world. List of sites for Project has to be subset of sites where the tier 0 spans.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_project.test UUID
```

The above command imports Project named `test` with the NSX ID `UUID`.
