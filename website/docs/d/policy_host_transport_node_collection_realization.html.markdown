---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_host_transport_node_collection_realization"
description: Host transport node collection realization information.
---

# nsxt_policy_host_transport_node_collection_realization

This data source provides information about the realization of host transport
node collection configured on NSX. This data source will fail if transport node collection
fails to realize. It is recommended to use this data source if further configuration
depends on host transport node collection realization.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_host_transport_node_collection_realization" "test" {
  path = data.nsxt_policy_host_transport_node_collection.path
}
```

## Argument Reference

* `path` - (Optional) Policy path of Transport Node Collection.
* `delay` - (Optional) Delay (in seconds) before realization polling is started. Default is set to 1.
* `timeout` - (Optional) Timeout (in seconds) for realization polling. Default is set to 1200.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - Application state of transport node profile on compute collection. Transitional state is "in_progress".

