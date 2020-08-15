---
layout: "nsxt"
page_title: "NSXT: cluster"
sidebar_current: "docs-nsxt-datasource-cluster"
description: A cluster data source.
---

# nsxt_certificate

This data source provides information about the NSX-T management cluster.

## Example Usage

```hcl
data "nsxt_cluster" "cluster" {}
```

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - Unique identifier of this cluster.
* `node_sha256_thumbprint` - SHA256 of certificate thumbprint of this manager node.
