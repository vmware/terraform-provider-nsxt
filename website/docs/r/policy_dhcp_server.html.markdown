---
subcategory: "DHCP"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_dhcp_server"
description: A resource to configure a DHCP Servers in NSX-T.
---

# nsxt_policy_dhcp_server

This resource provides a method for the management of a DHCP Server configurations.

This resource is applicable to NSX Global Manager, and NSX Policy Manager.

This resource is supported with NSX 3.0.0 onwards.

## Example Usage

```hcl
resource "nsxt_policy_dhcp_server" "test" {
  display_name      = "test"
  description       = "Terraform provisioned DhcpServerConfig"
  edge_cluster_path = data.nsxt_policy_edge_cluster.ec1.path
  lease_time        = 200
  server_addresses  = ["110.64.0.1/16", "2001::1234:abcd:ffff:c0a8:101/64"]
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_dhcp_server" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name      = "test"
  description       = "Terraform provisioned DhcpServerConfig"
  edge_cluster_path = data.nsxt_policy_project.demoproj.site_info.0.edge_cluster_paths.0
  lease_time        = 200
  server_addresses  = ["110.64.0.1/16", "2001::1234:abcd:ffff:c0a8:101/64"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `edge_cluster_path` - (Optional) The Policy path to the edge cluster for this DHCP Server.
* `lease_time` - (Optional) IP address lease time in seconds. Valid values from `60` to `4294967295`. Default is `86400`.
* `preferred_edge_paths` - (Optional) Policy paths to edge nodes. The first edge node is assigned as active edge, and second one as standby edge.
* `server_addresses` - (Optional) DHCP server address in CIDR format. At most 2 supported; one IPv4 and one IPv6 address. Server address can also be specified on segment subnet level.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing DHCP Server can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_dhcp_server.dhcp1 ID
```
The above command imports a DHCP Server named `dhcp1` with the NSX DHCP Server ID `ID`.

```
terraform import nsxt_policy_dhcp_server.dhcp1 POLICY_PATH
```
The above command imports a DHCP Server named `dhcp1` with the NSX DHCP Server policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
