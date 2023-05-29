---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: policy_edge node"
description: A policy Edge Node data source.
---

# nsxt_policy_edge_node

This data source provides information about policy edge nodes configured on NSX.

This data source is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_edge_cluster" "ec" {
  display_name = "ec"
}

data "nsxt_policy_edge_node" "node1" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.ec.path
  member_index      = 0
}
```

## Argument Reference

* `edge_cluster_path` - (Required) The path of edge cluster where to which this node belongs.
* `id` - (Optional) The ID of the edge node to retrieve.
* `display_name` - (Optional) The Display Name prefix of the edge node to retrieve.
* `member_index` - (Optional) Member index of the node in edge cluster.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
