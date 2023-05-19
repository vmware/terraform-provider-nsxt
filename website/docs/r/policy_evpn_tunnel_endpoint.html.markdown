---
subcategory: "EVPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_evpn_tunnel_endpoint"
description: A resource to configure EVPN Tunnel Endpoint in NSX Policy manager.
---

# nsxt_policy_evpn_tunnel_endpoint

This resource provides a method for the management of EVPN Tunnel Endpoint.

This resource is applicable to NSX Policy Manager.
This resource is supported with NSX 3.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_policy_evpn_tunnel_endpoint" "endpoint1" {
  display_name = "endpoint1"
  description  = "terraform provisioned endpoint"

  external_interface_path = nsxt_policy_tier0_gateway_interface.if1.path
  edge_node_path          = data.nsxt_policy_edge_node.node1.path

  local_address = "10.2.0.12"
  mtu           = 1500

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `external_interface_path` - (Required) Policy path for External Interface on Tier0 Gateway.
* `edge_node_path` - (Required) Edge node path.
* `local_address` - (Required) Local IPv4 IP address.
* `mtu` - (Optional) Maximal Transmission Unit.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `gateway_id` - Tier0 Gateway ID on which EVPN Tunnel is configured.
* `locale_service_id` - Tier0 Gateway Locale Service ID on which EVPN Tunnel is configured.

## Importing

An existing EVPN Tunnel Endpoint can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_evpn_tunnel_endpoint.endpoint1 GW-ID/LOCALE-SERVICE-ID/INTERFACE-ID/ID
```

The above command imports EVPN Tunnel Endpoint named `endpoint1` with the NSX Policy ID `ID`, on Tier0 Gateway GW-ID and Locale Service LOCALE-SERVICE-ID with external interface INTERFACE-ID.
