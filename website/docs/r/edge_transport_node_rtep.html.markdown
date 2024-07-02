---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_edge_transport_node_rtep"
description: A resource to configure an Edge Transport Node RTEP (remote tunnel endpoint).
---

# nsxt_edge_transport_node_rtep

This resource provides a method for the management of an Edge Transport Node RTEP (remote tunnel endpoint).
This resource is supported with NSX 4.1.0 onwards.

## Example Usage

```hcl
data "nsxt_edge_transport_node" "test_node" {
  display_name = "edgenode1"
}

resource "nsxt_edge_transport_node_rtep" "test_rtep" {
  edge_id = data.nsxt_edge_transport_node.test_node.id

  host_switch_name = "someSwitch"
  ip_assignment {
    assigned_by_dhcp = true
  }
  named_teaming_policy = "tp123"
  rtep_vlan            = 500
}
```

## Argument Reference

The following arguments are supported:

* `edge_id` - (Required) Edge ID to associate with remote tunnel endpoint.
* `host_switch_name` - (Required) The host switch name to be used for the remote tunnel endpoint.
* `ip_assignment` - (Required) - Specification for IPs to be used with host switch virtual tunnel endpoints. Should contain exatly one of the below:
    * `assigned_by_dhcp` - (Optional) Enables DHCP assignment.
    * `static_ip` - (Optional) IP assignment specification for Static IP List.
        * `ip_addresses` - (Required) List of IPs for transport node host switch virtual tunnel endpoints.
        * `subnet_mask` - (Required) Subnet mask.
        * `default_gateway` - (Required) Gateway IP.
    * `static_ip_pool` - (Optional) IP assignment specification for Static IP Pool.
* `rtep_vlan` - (Required) VLAN id for remote tunnel endpoint.
* `named_teaming_policy` - (Optional) The named teaming policy to be used by the remote tunnel endpoint.


## Importing

An existing Edge Transport Node RTEP can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_edge_transport_node_rtep.test UUID
```
The above command imports Edge Transport Node RTEP named `test` with the NSX Transport Node ID `UUID`.
