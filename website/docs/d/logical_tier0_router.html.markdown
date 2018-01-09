---
layout: "nsxt"
page_title: "NSXT: logical_tier0_router"
sidebar_current: "docs-nsxt-datasource-logical-tier0-router"
description: Logical Tier 0 Router data source.
---

# nsxt_logical_tier0_router

Provides information about logical tier 0 routers configured on NSX-T manager.

## Example Usage

```
data "nsxt_logical_tier0_router" "RT0" {
  display_name = "PLR1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Logical Tier 0 Router to retrieve

* `display_name` - (Optional) Display Name prefix of the Logical Tier 0 Router to retrieve

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - Description of the Logical Tier 0 Router.

* `router_type` - Type of Logical Router.

* `edge_cluster_id` - This id of the edge cluster connected to the router.

* `high_availability_mode` - The High availability mode of this router.
