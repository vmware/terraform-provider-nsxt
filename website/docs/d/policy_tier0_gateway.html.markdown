---
layout: "nsxt"
page_title: "NSXT: policy_tier0_gateway"
sidebar_current: "docs-nsxt-datasource-policy-tier0-gateway"
description: A policy Tier-0 gateway data source.
---

# nsxt_policy_tier0_gateway

This data source provides information about policy Tier-0 gateways configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "tier0_gw_gateway" {
  display_name = "tier0-gw"
}
```

## Argument Reference

* `id` - (Optional) The ID of Tier-0 gateway to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Tier-0 gateway to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `edge_cluster_path` - The path of the Edge cluster where this Tier-0 gateway is placed.

* `path` - The NSX path of the policy resource.
