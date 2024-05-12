---
subcategory: "Multitenancy"
layout: "nsxt"
page_title: "NSXT: policy_vpc"
description: Policy VPC data source.
---

# nsxt_policy_vpc

This data source provides information about policy VPC configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_vpc" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "vpc1"
}
```

## Argument Reference

* `id` - (Optional) The ID of VPC to retrieve. If ID is specified, no additional argument should be configured.
* `display_name` - (Optional) The Display Name prefix of the VPC to retrieve.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `short_id` - Defaults to id if id is less than equal to 8 characters or defaults to random generated id if not set.
* `site_info` - Information related to sites applicable for given Project.
    * `edge_cluster_paths` - PolicyPath of the edge cluster or label.
    * `site_path` - This represents the path of the site which is managed by Global Manager. For the local manager, if set, this needs to point to 'default'.
* `default_gateway_path` - PolicyPath of Tier0 or Tier0 VRF gateway or label path referencing to Tier0 or Tier0 VRF.
