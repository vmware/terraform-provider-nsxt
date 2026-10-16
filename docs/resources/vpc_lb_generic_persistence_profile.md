---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_lb_generic_persistence_profile"
description: A resource to configure a VPC Load Balancer Generic Persistence Profile.
---

# nsxt_vpc_lb_generic_persistence_profile

This resource provides a method for the management of a LB Generic Persistence Profile.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.2 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc" "demovpc" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "vpc1"
}

resource "nsxt_vpc_lb_generic_persistence_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  display_name       = "test"
  description        = "Terraform provisioned profile"
  persistence_shared = false
  timeout            = 1800
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Required) The context which the object belongs to. Requires a VPC context.
    * `project_id` - (Required) The ID of the project which the object belongs to.
    * `vpc_id` - (Required) The ID of the VPC which the object belongs to.
* `ha_persistence_mirroring_enabled` - (Optional) If enabled, persistence entries are synchronized to HA peer.
* `timeout` - (Optional) Expiration time once all connections are complete. Default is 300.
* `persistence_shared` - (Optional) If enabled, all virtual servers with this profile will share the same persistence mechanism.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_lb_generic_persistence_profile.test PATH
```

The above command imports a VPC LB Generic Persistence Profile named `test` with the policy path `PATH`.
The path must be a VPC path of the form `/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-persistence-profiles/[id]`.
Importing by ID only is not supported, since the project and VPC are derived from the path.
