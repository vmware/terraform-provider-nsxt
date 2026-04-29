---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_transit_gateway_bfd_peer"
description: A resource to configure a Transit Gateway BFD Peer.
---

# nsxt_policy_transit_gateway_bfd_peer

This resource provides a method for the management of a BFD (Bidirectional Forwarding Detection) peer on a Transit Gateway.

BFD peers enable fast liveness detection of static route next-hop peers on transit gateway centralized network attachments.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "test_proj" {
  display_name = "test_project"
}

data "nsxt_policy_transit_gateway" "test_tgw" {
  context {
    project_id = data.nsxt_policy_project.test_proj.id
  }
  id = "default"
}

resource "nsxt_policy_transit_gateway_bfd_peer" "test" {
  display_name = "my-bfd-peer"
  parent_path  = data.nsxt_policy_transit_gateway.test_tgw.path
  peer_address = "192.168.10.1"

  source_attachment = [nsxt_policy_transit_gateway_attachment.cna_att.path]
}
```

### With BFD profile and explicit enable flag

```hcl
resource "nsxt_policy_transit_gateway_bfd_peer" "test_with_profile" {
  display_name     = "my-bfd-peer"
  description      = "BFD peer for static route next-hop"
  parent_path      = data.nsxt_policy_transit_gateway.test_tgw.path
  peer_address     = "192.168.10.1"
  bfd_profile_path = "/orgs/default/projects/project1/infra/bfd-profiles/fast"
  enabled          = true

  source_attachment = [nsxt_policy_transit_gateway_attachment.cna_att.path]
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required, ForceNew) Policy path of the parent transit gateway (e.g. `/orgs/default/projects/proj1/transit-gateways/tgw1`).
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `peer_address` - (Required) IP address (IPv4 or IPv6) of the static route next-hop peer for which BFD liveness detection is configured. Only one BFD peer configuration per peer address is allowed per transit gateway.
* `source_attachment` - (Required) A list containing exactly one policy path of the transit gateway attachment that serves as the BFD session source. The referenced attachment must have a `cna_path` pointing to a centralized network attachment.
* `bfd_profile_path` - (Optional) Policy path of the BFD profile to use for timing parameters of this BFD session. If not specified, default BFD timing parameters are used.
* `enabled` - (Optional) Flag to enable or disable this BFD peer configuration. Defaults to `true`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

```shell
terraform import nsxt_policy_transit_gateway_bfd_peer.example PARENT_PATH/ID
```

The above command imports a Transit Gateway BFD Peer named `example` with the NSX path `PARENT_PATH` and resource ID `ID`. For example:

```shell
terraform import nsxt_policy_transit_gateway_bfd_peer.example /orgs/default/projects/project1/transit-gateways/tgw1/bfd-peer-1
```

[docs-import]: https://developer.hashicorp.com/terraform/cli/import
