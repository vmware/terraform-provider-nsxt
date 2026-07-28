---
subcategory: "Fabric"
page_title: "NSXT: policy_edge node"
description: A policy Edge Node data source.
---

# nsxt_policy_edge_node

This data source provides information about policy edge nodes configured within a cluster, on NSX.

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

~> **NOTE:** `id` behaves inconsistently on this data source. As an input, it is used to locate the edge
node. However, once the edge node is found, NSX reports back an `id` that is actually the node's
`member_index` within the edge cluster rather than a stable, globally unique identifier. If you set `id`
to look up a specific edge node, do not rely on the resulting `id` attribute value for anything other
than the lookup - use `path`, `unique_id` or `realization_id` instead to reference this resource
unambiguously elsewhere in your configuration.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `unique_id` - A unique identifier assigned by the system for this edge node.
* `realization_id` - The ID used to realize this edge node.
