---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_host_transport_node_collection_realization"
description: Host transport node collection realization information.
---

# nsxt_policy_host_transport_node_collection_realization

This data source provides information about the realization of host transport
node collection configured on NSX. This data source will wait until realization is
determined as either success or error. It is recommended to use this data source
if further configuration depends on host transport node collection realization.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_host_transport_node_collection_realization" "test" {
  id = data.nsxt_policy_host_transport_node_collection.id
}
```

## Argument Reference

* `id` - (Optional) The ID of host transport node collection to retrieve information.
* `delay` - (Optional) Delay (in seconds) before realization polling is started. Default is set to 1.
* `timeout` - (Optional) Timeout (in seconds) for realization polling. Default is set to 1200.
* `site_path` - (Optional) The path of the site which the Transport Node Collection belongs to. `path` field of the existing `nsxt_policy_site` can be used here.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - Application state of transport node profile on compute collection. Transitional state is "in_progress".
* `aggregate_progress_percentage` - Aggregate percentage of compute collection deployment.
* `cluster_level_error` - Errors which needs cluster level to resolution.
* `validation_errors` - Errors while applying transport node profile on discovered node.
* `vlcm_transition_error` - Errors while enabling vLCM on the compute collection.

