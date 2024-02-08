---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_host_transport_node_collection_state"
description: A host transport node collection state data source.
---

# nsxt_policy_host_transport_node_collection_state

This data source provides information about state of host transport node collection configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_host_transport_node_collection_state" "host_transport_node_collection_state" {
  id        = data.nsxt_policy_host_transport_node_collection.id
  site_path = data.nsxt_policy_site.paris.path
}
```

## Argument Reference

* `id` - (Optional) The ID of host transport node collection to retrieve state.
* `display_name` - (Optional) The Display Name prefix of the host transport node collection to retrieve state.
* `site_path` - (Optional) The path of the site which the Transport Node Collection belongs to. `path` field of the existing `nsxt_policy_site` can be used here.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `aggregate_progress_percentage` - Aggregate percentage of compute collection deployment.
* `cluster_level_error` - Errors which needs cluster level to resolution.
* `state` - Application state of transport node profile on compute collection.
* `validation_errors` - Errors while applying transport node profile on discovered node.
* `vlcm_transition_error` - Errors while enabling vLCM on the compute collection.

