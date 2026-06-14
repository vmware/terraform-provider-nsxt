---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_endpoint"
description: VPC Endpoint data source.
---

# nsxt_vpc_endpoint

This data source provides information about a VPC Endpoint on NSX.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc" "demovpc" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "consumer-vpc"
}

data "nsxt_vpc_endpoint" "example" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  display_name = "payment-ep"
}
```

## Argument Reference

* `id` - (Optional) The ID of the VPC Endpoint to retrieve.
* `display_name` - (Optional) The display name of the VPC Endpoint to retrieve.
* `context` - (Required) The context which the object belongs to.
    * `project_id` - (Required) The ID of the project which the object belongs to.
    * `vpc_id` - (Required) The ID of the VPC which the object belongs to.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `vpc_service_endpoint` - Policy path to the VPC service endpoint being consumed.
* `ip_allocation_path` - Policy path to the VPC IP address allocation.
