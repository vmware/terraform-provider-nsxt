---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_project_centralized_network_attachment"
description: A resource to configure a Project Centralized Network Attachment.
---

# nsxt_policy_project_centralized_network_attachment

This resource provides a method for the management of a Centralized Network Attachment (CNA) within a project. A CNA connects a centralized transit gateway to a VLAN or overlay VPC subnet, enabling external or internal network connectivity with BGP/static routing support.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

API path: `/orgs/{org}/projects/{project}/centralized-network-attachments/{id}`

## Example Usage

```hcl
resource "nsxt_policy_project_centralized_network_attachment" "example" {
  context {
    project_id = "dev"
  }

  display_name = "example-cna"
  description  = "Terraform provisioned CNA"
  subnet_path  = nsxt_vpc_subnet.external.path

  advertise_outbound_networks {
    allow_private         = true
    allow_external_blocks = ["10.0.0.0/8", "172.16.0.0/12"]
  }

  tag {
    scope = "env"
    tag   = "prod"
  }
}
```

### With manual interface IP assignment

```hcl
resource "nsxt_policy_project_centralized_network_attachment" "manual_ip" {
  context {
    project_id = "dev"
  }

  display_name = "cna-manual-ip"
  subnet_path  = nsxt_vpc_subnet.external.path

  # One entry per address family (IPv4 or IPv6), max 2.
  interface_subnet {
    interface_ip_address = ["192.168.10.5", "192.168.10.6"]
    prefix_length        = 24
    ha_vip_ip_address    = "192.168.10.100"
  }
}
```

## Argument Reference

The following arguments are supported:

* `context` - (Required) Project context for this resource.
    * `project_id` - (Required) ID of the project this CNA belongs to. Immutable after creation.
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `subnet_path` - (Required, Immutable) Policy path of the VPC subnet (overlay or VLAN-backed) to which this CNA is connected. Cannot be changed after creation.
* `advertise_outbound_networks` - (Optional) Outbound route advertisement configuration.
    * `allow_private` - (Optional) When `true`, disables VPC auto-SNAT and EIP translation on the interface, and enables `TGW_PRIVATE` route redistribution on this connection. Default: `false`.
    * `allow_external_blocks` - (Optional) List of external IP block CIDRs used as an advertisement filter for prefixes redistributed from the transit gateway.
* `interface_subnet` - (Optional) Manual IP assignment for per-node edge interfaces. Supports 1–2 entries (one IPv4 and/or one IPv6 subnet). CNAs with manual IP assignment cannot be shared to other centralized transit gateways.
    * `interface_ip_address` - (Required) List of IP addresses assigned to each edge node interface. The count must match the number of edge nodes hosting the centralized transit gateway.
    * `prefix_length` - (Required) Prefix length of the subnet. Valid range: 1–32 for IPv4, 1–128 for IPv6.
    * `ha_vip_ip_address` - (Optional) Floating VIP assigned to the active edge node. Only applicable when the centralized transit gateway uses `ACTIVE_STANDBY` HA mode.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_project_centralized_network_attachment.example PATH
```

The above command imports the Centralized Network Attachment named `example` using its full NSX policy path, for example:

```text
/orgs/default/projects/dev/centralized-network-attachments/my-cna-id
```
