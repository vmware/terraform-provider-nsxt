---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: policy_edge cluster"
description: A policy Edge Cluster data source.
---

# nsxt_policy_edge_cluster

This data source provides information about policy edge cluster configured on NSX.

This data source is applicable to NSX Global Manager, NSX Policy Manager and VMC.

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
  site_path    = data.nsxt_policy_site.paris.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the edge cluster to retrieve.
* `display_name` - (Optional) The Display Name prefix of the edge cluster to retrieve.
* `site_path` - (Optional) The path of the site which the Edge Cluster belongs to, this configuration is required for global manager only. `path` field of the existing `nsxt_policy_site` can be used here. If a single edge cluster is configured on site, `id` and `display_name` can be omitted in configuration, otherwise either of these is required to specify the desired cluster.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
