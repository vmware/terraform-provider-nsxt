---
layout: "nsxt"
page_title: "NSXT: logical_tier1_router"
sidebar_current: "docs-nsxt-datasource-logical-tier1-router"
description: A logical Tier 1 router data source.
---

# nsxt_logical_tier1_router

This data source provides information about logical Tier 1 routers configured in NSX.

## Example Usage

```hcl
data "nsxt_logical_tier1_router" "tier1_router" {
  display_name = "router1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Logical Tier 1 Router to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Logical Tier 1 Router to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the logical Tier 0 router.

* `edge_cluster_id` - The id of the Edge cluster where this logical router is placed.
