---
subcategory: "Multitenancy"
layout: "nsxt"
page_title: "NSXT: policy_project"
description: Policy Project data source.
---

# nsxt_policy_project

This data source provides information about policy Project configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "test" {
  display_name = "project1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Project to retrieve. If ID is specified, no additional argument should be configured.
* `display_name` - (Optional) The Display Name prefix of the Project to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `short_id` - Unique ID used for logging.
* `site_info` - Information related to sites applicable for given Project.
  * `edge_cluster_paths` - The edge cluster on which the networking elements for the Org are be created.
  * `site_path` - This represents the path of the site which is managed by Global Manager. For the local manager the value would be 'default'.
* `tier0_gateway_paths` - Policy paths of Tier0 gateways associated with the project.
* `external_ipv4_blocks` - Policy paths of IPv4 blocks associated with the project.
