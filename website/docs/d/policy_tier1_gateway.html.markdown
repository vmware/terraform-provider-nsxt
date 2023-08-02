---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: policy_tier1_gateway"
description: A policy Tier-1 gateway data source.
---

# nsxt_policy_tier1_gateway

This data source provides information about policy Tier-1s configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_tier1_gateway" "tier1_router" {
  display_name = "tier1_gw"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_tier1_gateway" "demotier1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demotier1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Tier-1 gateway to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Tier-1 gateway to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `edge_cluster_path` - The path of the Edge cluster where this Tier-1 gateway is placed.
* `path` - The NSX path of the policy resource.
