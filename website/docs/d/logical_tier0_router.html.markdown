---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: logical_tier0_router"
description: A logical Tier 0 router data source.
---

# nsxt_logical_tier0_router

This data source provides information about Tier 0 Logical Routers configured in NSX. A Tier 0 router is used to connect NSX networking with traditional physical networking. Tier 0 routers are placed on an Edge cluster and will exist on one or more Edge node depending on deployment settings (i.e. active/active or active/passive). A Tier 0 router forwards layer 3 IP packets and typically peers with a traditional physical router using BGP or can use static routing.

## Example Usage

```hcl
data "nsxt_logical_tier0_router" "tier0_router" {
  display_name = "PLR1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Logical Router to retrieve.
* `display_name` - (Optional) The Display Name prefix of Logical Router to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the Logical Router.
* `edge_cluster_id` - The id of the Edge Cluster where this Logical Router is placed.
* `high_availability_mode` - The high availability mode of this Logical Router.
