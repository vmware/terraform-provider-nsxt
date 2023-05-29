---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: logical_tier1_router"
description: A logical Tier 1 router data source.
---

# nsxt_logical_tier1_router

This data source provides information about Tier1 Logical Routers configured on NSX.

## Example Usage

```hcl
data "nsxt_logical_tier1_router" "tier1_router" {
  display_name = "router1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Logical Router to retrieve.
* `display_name` - (Optional) The Display Name prefix of Logical Router to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the Logical Router.
* `edge_cluster_id` - The id of the Edge cluster where this Logical Router is placed.
