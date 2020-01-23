---
layout: "nsxt"
page_title: "NSXT: policy_tier1_gateway"
sidebar_current: "docs-nsxt-datasource-policy-tier1-gateway"
description: A policy Tier-1 gateway data source.
---

# nsxt_policy_tier1_gateway

This data source provides information about policy Tier-1s configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_tier1_gateway" "tier1_router" {
  display_name = "tier1_gw"
}
```

## Argument Reference

* `id` - (Optional) The ID of Tier-1 gateway to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Tier-1 gateway to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `edge_cluster_path` - The path of the Edge cluster where this Tier-1 gateway is placed.

* `path` - The NSX path of the policy resource.
