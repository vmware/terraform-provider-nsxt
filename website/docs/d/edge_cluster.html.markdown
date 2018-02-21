---
layout: "nsxt"
page_title: "NSXT: edge_cluster"
sidebar_current: "docs-nsxt-datasource-edge-cluster"
description: Edge Cluster data source.
---

# nsxt_edge_cluster

Provides information about edge clusters configured on NSX-T manager.

## Example Usage

```
data "nsxt_edge_cluster" "edge_cluster1" {
  display_name = "edgecluster"
}
```

## Argument Reference

* `id` - (Optional) The ID of Edge Cluster to retrieve

* `display_name` - (Optional) Display Name prefix of the Edge Cluster to retrieve

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - Description of the edge cluster.

* `deployment_type` - This field could show deployment_type of members (VIRTUAL_MACHINE|PHYSICAL_MACHINE).

* `member_node_type` - Edge cluster is homogenous collection of transport nodes. Hence all transport nodes of the cluster must be of same type. This field shows the type of transport nodes