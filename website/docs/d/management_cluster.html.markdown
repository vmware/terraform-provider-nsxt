---
subcategory: "Manager"
layout: "nsxt"
page_title: "NSXT: management_cluster"
description: A NSX-T management cluster data source.
---

# nsxt_management_cluster

This data source provides information about the NSX-T management cluster.

## Example Usage

```hcl
data "nsxt_management_cluster" "cluster" {}
```

## Attributes Reference

* `id` - Unique identifier of this cluster.
* `node_sha256_thumbprint` - SHA256 of certificate thumbprint of this manager node.
