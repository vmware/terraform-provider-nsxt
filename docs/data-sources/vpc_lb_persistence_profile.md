---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_lb_persistence_profile"
description: VPC Load Balancer Persistence Profile data source.
---

# nsxt_vpc_lb_persistence_profile

This data source provides information about VPC Load Balancer Persistence Profiles configured on NSX.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

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

data "nsxt_vpc_lb_persistence_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  display_name = "vpc-lb-persistence-profile1"
}
```

## Argument Reference

* `context` - (Required) The context which the object belongs to. Requires a VPC context.
    * `project_id` - (Required) The ID of the project which the object belongs to.
    * `vpc_id` - (Required) The ID of the VPC which the object belongs to.
* `id` - (Optional) The ID of Profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of Profile to retrieve.
* `type` - (Optional) The Load Balancer Persistence Profile type. One of `ANY`, `SOURCE_IP`, `COOKIE` or `GENERIC`. Defaults to `ANY`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
