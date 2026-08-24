---
subcategory: "Fabric"
page_title: "NSXT: manager_cluster_node"
description: A NSX manager cluster node data source.
---

# nsxt_manager_cluster_node

This data source provides information about a specific NSX manager cluster node.

## Example Usage

```hcl
data "nsxt_manager_cluster_node" "node1" {
  display_name = "nsx-manager-node1"
}
```

## Argument Reference

* `id` - (Optional) Unique ID of the manager cluster node to retrieve.
* `display_name` - (Optional) Display name prefix of the manager cluster node to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - Description of this resource.
* `appliance_mgmt_listen_address` - The IP and port for the appliance management API service on this node.
