---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_lb_service"
description: A resource to configure a VPC Load Balancer Service.
---

# nsxt_vpc_lb_service

This resource provides a method for the management of a Load Balancer Service within a VPC.

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

resource "nsxt_vpc_lb_service" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }

  display_name    = "test"
  description     = "Terraform provisioned LB Service"
  size            = "SMALL"
  enabled         = true
  error_log_level = "ERROR"
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
* `size` - (Optional) Load Balancer Service size, one of `SMALL`, `MEDIUM`, `LARGE`, `XLARGE`, `DLB`. Default is `SMALL`.
* `enabled` - (Optional) Flag to enable the service.
* `relax_scale_validation` - (Optional) Whether the scale validations of virtual servers, pools, pool members and rules are relaxed for the service. If not set, NSX applies its default of `true` on creation, and the value configured on NSX is preserved on later updates.
* `error_log_level` - (Optional) Log level for the service, one of `DEBUG`, `INFO`, `WARNING`, `ERROR`, `CRITICAL`, `ALERT`, `EMERGENCY`. Default is `INFO`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `connectivity_path` - Policy path of the VPC which the service is connected to. It is assigned by NSX.

## Importing

An existing service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_lb_service.test PATH
```

The above command imports the VPC LBService named `test` with the policy path `PATH`, which must be a VPC path of the format `/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lbs/[id]`. The project and VPC context is derived from the path.
