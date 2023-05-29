---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: edge_cluster"
description: An Edge Cluster data source.
---

# nsxt_edge_cluster

This data source provides information about Edge clusters configured in NSX. An Edge cluster is a collection of Edge nodes which can be deployed as either VM form-factor or bare-metal form-factor machines for connectivity between overlay logical switches and non-NSX underlay networking for north/south layer 2 or layer 3 connectivity. Each T0 router will be placed on one ore more Edge nodes in an Edge cluster therefore this data source is needed for the creation of T0 logical routers.

## Example Usage

```hcl
data "nsxt_edge_cluster" "edge_cluster1" {
  display_name = "edgecluster"
}
```

## Argument Reference

* `id` - (Optional) The ID of Edge Cluster to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Edge Cluster to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the edge cluster.
* `deployment_type` - This field could show deployment_type of members. It would return UNKNOWN if there is no members, and return VIRTUAL_MACHINE|PHYSICAL_MACHINE if all Edge members are VIRTUAL_MACHINE|PHYSICAL_MACHINE.
* `member_node_type` - An Edge cluster is homogeneous collection of NSX transport nodes used for north/south connectivity between NSX logical networking and physical networking. Hence all transport nodes of the cluster must be of same type. This field shows the type of transport node,
