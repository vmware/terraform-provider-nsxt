---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_edge_transport_node_rtep"
description: A resource to configure an Policy Edge Transport Node RTEP (remote tunnel endpoint).
---

# nsxt_policy_edge_transport_node_rtep

This resource provides a method for the management of a Policy Edge Transport Node RTEP (remote tunnel endpoint).
This resource is supported with NSX 9.0.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_edge_transport_node" "test_node" {
  display_name = "edgenode1"
}

data "nsxt_policy_ip_pool" "ip_pool1" {
  display_name = "ip-pool1"
}

resource "nsxt_policy_edge_transport_node_rtep" "test_rtep" {
  edge_transport_node_path = data.nsxt_edge_node.test_node.path
  host_switch_name         = "someSwitch"
  ip_assignment {
    static_ipv4_pool = data.nsxt_policy_ip_pool.ip_pool1.path
  }
  named_teaming_policy = "tp123"
  vlan                 = 500
}
```

## Argument Reference

The following arguments are supported:

* `edge_transport_node_path` - (Required) Policy path for Policy Edge Transport Node to associate with remote tunnel endpoint.
* `host_switch_name` - (Required) The host switch name to be used for the remote tunnel endpoint.
* `ip_assignment` - (Required) - Specification for IPs to be used with host switch virtual tunnel endpoints. Should contain exatly one of the below:
  * `static_ipv4_list` - (Optional) IP assignment specification value for Static IPv4 List.
    * `default_gateway` - (Required) Gateway IP.
    * `ip_addresses` - (Required) List of IPV4 addresses for edge transport node host switch virtual tunnel endpoints.
    * `subnet_mask` - (Required) Subnet mask.
  * `static_ipv4_pool` - (Optional) IP assignment specification for Static IPv4 Pool. Input can be MP ip pool UUID or policy path of IP pool.
* `vlan` - (Required) VLAN id for remote tunnel endpoint.
* `named_teaming_policy` - (Optional) The named teaming policy to be used by the remote tunnel endpoint.

## Importing

An existing Edge Transport Node RTEP can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_edge_transport_node_rtep.test POLICY_PATH:SWITCH_ID
```
The above command imports Policy Edge Transport Node RTEP named `test` with the NSX Policy Transport Node path `POLICY_PATH` and host_switch_name `SWITCH_ID`.

**NOTE:** The Policy Edge Transport Node path and host_switch_name are separated by a colon. 
