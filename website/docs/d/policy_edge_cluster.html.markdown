---
layout: "nsxt"
page_title: "NSXT: policy_edge cluster"
sidebar_current: "docs-nsxt-datasource-policy-edge-cluster"
description: A policy Edge Cluster data source.
---

# nsxt_policy_edge_cluster

This data source provides information about policy edge clusters configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_edge_cluster" "ec" {
  display_name = "ec"
}
```

Note: This usage is for Global Manager only.
```hcl

data "nsxt_policy_site" "paris" {
  display_name = "Paris"
}

data "nsxt_policy_edge_cluster" "gm_ec" {
  display_name = "ec"
  site_path = data.nsxt_policy_site.paris.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the edge cluster to retrieve.

* `display_name` - (Optional) The Display Name prefix of the edge cluster to retrieve.

* `site_path` - (Optional) The path of the site which the Edge Cluster belongs to, this configuration is required for global manager only. `path` field of the existing `nsxt_policy_site` can be used here.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
