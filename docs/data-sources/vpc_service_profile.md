---
subcategory: "VPC"
page_title: "NSXT: vpc_service_profile"
description: VPC Service Profile data source.
---

# nsxt_vpc_service_profile

This data source provides information about an inventory Service Profile configured under VPC on NSX.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc_service_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "profile1"
}
```

## Example Usage: Fetching Default VPC Service Profile Attribute

```hcl
data "nsxt_vpc_service_profile" "test" {
  context {
    project_id = "default"
  }
  is_default = true
}

output "nsxt_vpc_service_profile_test" {
  value = data.nsxt_vpc_service_profile.test.display_name
}
```

## Argument Reference

* `id` - (Optional) The ID of Service Profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Service Profile to retrieve.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `is_default` - Specifies if the VPC Service Profile is set as the default.
