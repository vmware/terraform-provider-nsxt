---
subcategory: "Beta"
page_title: "NSXT: policy_network_span"
description: A policy Network Span data source.
---

# nsxt_policy_edge_cluster

This data source provides information about policy network span configured on NSX.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.1.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_network_span" "netspan" {
  display_name = "ns1"
}
```

## Argument Reference

* `id` - (Optional) The ID of the network span to retrieve.
* `display_name` - (Optional) The Display Name prefix of the network span to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
