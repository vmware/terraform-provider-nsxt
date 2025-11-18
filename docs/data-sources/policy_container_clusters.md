---
subcategory: "Beta"
page_title: "NSXT: policy_container_clusters"
description: A policy container clusters data source. This data source builds "display name to policy paths" map representation of the whole table.
---

# nsxt_policy_container_clusters

This data source builds a "name to paths" map of the whole policy Container Clusters table. Such map can be referenced in configuration to obtain object identifier attributes by display name at a cost of single roundtrip to NSX, which improves apply and refresh
time at scale, compared to multiple instances of `nsxt_policy_container_cluster` data source.

## Example Usage

```hcl
data "nsxt_policy_container_clusters" "map" {
}

resource "nsxt_policy_security_policy_container_cluster" "antreacluster" {
  display_name           = "cluster1"
  description            = "Terraform provisioned SecurityPolicyContainerCluster"
  policy_path            = nsxt_policy_parent_security_policy.policy1.path
  container_cluster_path = data.nsxt_policy_container_clusters.map.items["cluster1"].path
}
```

## Attributes Reference

The following attributes are exported:

* `items` - Map of policy service policy paths keyed by display name.
