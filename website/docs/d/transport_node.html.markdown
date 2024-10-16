---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: transport_node"
description: An Transport Node data source.
---

# nsxt_transport_node

This data source provides information about Transport Node configured on NSX.

## Example Usage

```hcl
data "nsxt_transport_node" "test_node" {
  display_name = "edgenode1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Transport Node to retrieve
* `display_name` - (Optional) The Display Name of the Transport Node to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
