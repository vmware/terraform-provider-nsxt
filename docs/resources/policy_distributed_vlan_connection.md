---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_distributed_vlan_connection"
description: A resource to configure a Distributed Vlan Connection.
---

# nsxt_policy_distributed_vlan_connection

This resource provides a method for the management of a Distributed Vlan Connection.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_distributed_vlan_connection" "test" {
  display_name      = "test"
  description       = "Terraform provisioned Distributed Vlan Connection"
  gateway_addresses = ["192.168.2.1/24", "192.168.3.1/24"]
  vlan_id           = 12
}
```

## Example Usage - With VLAN Extension (NSX 9.1.0+)

```hcl
resource "nsxt_policy_ip_block" "vlan_block" {
  display_name     = "vlan-extension-block"
  cidrs            = ["10.66.66.0/24"]
  visibility       = "EXTERNAL"
  subnet_exclusive = true
}

resource "nsxt_policy_distributed_vlan_connection" "test" {
  display_name                = "test"
  description                 = "Terraform provisioned Distributed Vlan Connection"
  gateway_addresses           = ["192.168.2.1/24"]
  vlan_id                     = 12
  subnet_extension_connection = "ENABLED_L2"
  associated_ip_block_paths   = [nsxt_policy_ip_block.vlan_block.path]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `vlan_id` - (Required) Vlan id for external gateway traffic.
* `gateway_addresses` - (Optional) List of gateway addresses in CIDR format. Only one gateway IP CIDR is allowed for subnet exclusive configuration.
* `subnet_extension_connection` - (Optional) Controls the connectivity mode for VPC Subnets referencing this distributed VLAN connection. This can be one of `DISABLED`, `ENABLED_L2`, or `ENABLED_L2_AND_L3`. Default is `DISABLED`. This attribute is supported with NSX 9.1.0 onwards.
* `associated_ip_block_paths` - (Optional) List of IP address block(s) that are associated with the distributed vlan connection. This attribute is supported with NSX 9.1.0 onwards.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_distributed_vlan_connection.test PATH
```

The above command imports Distributed Vlan Connection named `test` with the policy path `PATH`.
