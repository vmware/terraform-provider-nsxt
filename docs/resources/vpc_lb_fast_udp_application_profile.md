---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_lb_fast_udp_application_profile"
description: A resource to configure a VPC Load Balancer Fast UDP Application Profile.
---

# nsxt_vpc_lb_fast_udp_application_profile

This resource provides a method for the management of a VPC Load Balancer Fast UDP Application Profile.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

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

resource "nsxt_vpc_lb_fast_udp_application_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  display_name           = "test"
  description            = "Terraform provisioned profile"
  idle_timeout           = 120
  flow_mirroring_enabled = true
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
* `idle_timeout` - (Optional) Timeout in seconds to specify how long UDP connection can remain open. Default is 300.
* `flow_mirroring_enabled` - (Optional) If enabled, all the flows to the bounded virtual server are mirrored to the standby node. Default is `false`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_lb_fast_udp_application_profile.test PATH
```

The above command imports the VPC LB Fast UDP application profile named `test` with the policy path `PATH`, which must be a VPC path of the format `/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-app-profiles/[id]`. The project and VPC context is derived from the path.
