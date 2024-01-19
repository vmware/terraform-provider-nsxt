---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: discover_node"
description: An Discover Node data source.
---

# nsxt_discover_node

This data source provides information about Discover Node configured in NSX. A Discover Node can be used to create a Host Transport Node.

## Example Usage

```hcl
data "nsxt_discover_node" "test" {
  ip_address = "10.43.251.142"
}
```

## Argument Reference

* `id` - (Optional) External id of the discovered node, ex. a mo-ref from VC.
* `ip_address` - (Optional) IP Address of the discovered node.
