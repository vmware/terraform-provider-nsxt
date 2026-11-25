---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_lb_passive_monitor_profile"
description: A resource to configure a VPC LBPassiveMonitorProfile.
---

# nsxt_vpc_lb_passive_monitor_profile

This resource provides a method for the management of a VPC LBPassiveMonitorProfile.

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

resource "nsxt_vpc_lb_passive_monitor_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  display_name = "test"
  description  = "Terraform provisioned LBPassiveMonitorProfile"
  max_fails    = 2
  timeout      = 2

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
* `max_fails` - (Optional) Number of consecutive failures before a member is considered temporarily unavailable.
* `timeout` - (Optional) After this timeout period, the member is tried again for a new connection to see if it is available.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_lb_passive_monitor_profile.test PATH
```

The above command imports the VPC LB passive monitor profile named `test` with the policy path `PATH`, which must be a VPC path of the format `/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-monitor-profiles/[id]`. The project and VPC context is derived from the path.
