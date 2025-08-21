---
subcategory: "Fabric"
page_title: "NSXT: nsxt_policy_edge_transport_node"
description: An Edge transport node data source.
---

# nsxt_policy_edge_transport_node

This data source provides information about Edge transport node configured on NSX.
This data source is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_edge_transport_node" "edge_transport_node" {
  display_name = "edge_transport_node1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Edge transport node to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Edge transport node to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `unique_id` - A unique identifier assigned by the system.
