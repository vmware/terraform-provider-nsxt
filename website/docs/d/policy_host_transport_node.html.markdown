---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_host_transport_node"
description: A host transport node data source.
---

# nsxt_policy_host_transport_node

This data source provides information about host transport node configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_host_transport_node" "host_transport_node" {
  display_name = "host_transport_node1"
}
```

## Argument Reference

* `id` - (Optional) The ID of host transport node to retrieve.
* `display_name` - (Optional) The Display Name prefix of the host transport node to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `unique_id` - A unique identifier assigned by the system.
