---
subcategory: "VPC"
page_title: "NSXT: policy_transit_gateway_bfd_peer"
description: A policy Transit Gateway BFD Peer data source.
---

# nsxt_policy_transit_gateway_bfd_peer

This data source provides information about a BFD (Bidirectional Forwarding Detection) peer configured on a Transit Gateway.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_transit_gateway_bfd_peer" "test" {
  display_name = "my-bfd-peer"
  parent_path  = nsxt_policy_transit_gateway.demo.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the Transit Gateway BFD Peer to retrieve.
* `display_name` - (Optional) The Display Name of the Transit Gateway BFD Peer to retrieve.
* `parent_path` - (Required) Policy path of the parent Transit Gateway.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
