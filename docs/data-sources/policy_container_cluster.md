---
subcategory: "Beta"
page_title: "NSXT: policy_container_cluster"
description: Policy Container Cluster data source.
---

# nsxt_policy_container_cluster

This data source provides information about Policy Container Cluster configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_container_cluster" "cluster" {
  display_name = "containercluster1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Container Cluster to retrieve. If ID is specified, no additional argument should be configured.
* `display_name` - (Optional) The Display Name prefix of the Container Cluster to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
