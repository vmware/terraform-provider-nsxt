---
subcategory: "VPC"
page_title: "NSXT: vpc_lb_app_profile"
description: VPC Load Balancer Application Profile data source.
---

# nsxt_vpc_lb_app_profile

This data source provides information about VPC Load Balancer Application Profile configured in NSX.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.0.2 onwards.

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

data "nsxt_vpc_lb_app_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  type         = "TCP"
  display_name = "my-tcp-profile"
}
```

## Argument Reference

* `context` - (Required) The context which the object belongs to. Requires a VPC context.
    * `project_id` - (Required) The ID of the project which the object belongs to.
    * `vpc_id` - (Required) The ID of the VPC which the object belongs to.
* `id` - (Optional) The ID of Profile to retrieve.
* `type` - (Optional) Type of Profile to retrieve, one of `HTTP`, `TCP`, `UDP`, `ANY`.
* `display_name` - (Optional) The Display Name prefix of the Profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
