---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_static_route_bfd_peer"
description: A resource to configure Static Route BFD Peer on Tier0 Gateway.
---

# nsxt_policy_static_route_bfd_peer

This resource provides a method for the management of Static Route BFD Peer on Tier0 Gateway.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_static_route_bfd_peer" "test" {
  gateway_path     = nsxt_policy_tier0_gateway.test.path
  bfd_profile_path = data.nsxt_policy_bfd_profile.default.path

  display_name     = "test"
  description      = "Terraform provisioned peer"
  enabled          = true
  peer_address     = "10.2.12.4"
  source_addresses = ["20.9.14.5"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `gateway_path` - (Required) Policy path of relevant Tier0 Gateway.
* `bfd_profile_path` - (Required) Policy path of relevant BFD Profile.
* `enabled` - (Optional) A fkag to enable/disable this Peer, default is `true`.
* `peer_address` - (Required) IPv4 address of the Peer.
* `source_addresses` - (Optional) List of relevant IPv4 Tier0 external interface addresses.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_static_route_bfd_peer.test GW-ID/ID
```

The above command imports Tier0 Gateway Static Route BFD Peer named `test` with ID `ID` on Tier0 Gateway `GW-ID`.
