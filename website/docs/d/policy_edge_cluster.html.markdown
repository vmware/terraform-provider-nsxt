---
layout: "nsxt"
page_title: "NSXT: policy_edge cluster"
sidebar_current: "docs-nsxt-datasource-policy-edge-cluster"
description: A policy Edge Cluster data source.
---

# nsxt_policy_edge_cluster

This data source provides information about policy edge clusters configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_edge_cluster" "ec" {
  display_name = "ec"
}
```

## Argument Reference

* `id` - (Optional) The ID of the edge cluster to retrieve.

* `display_name` - (Optional) The Display Name prefix of the edge cluster to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
