---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_gateway_connection"
description: A resource to configure a GatewayConnection.
---

# nsxt_policy_gateway_connection

This resource provides a method for the management of a GatewayConnection.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "test" {
  display_name = "test-t0gw"
}

resource "nsxt_policy_gateway_connection" "test" {
  display_name     = "test"
  description      = "Terraform provisioned GatewayConnection"
  tier0_path       = data.nsxt_policy_tier0_gateway.test.path
  aggregate_routes = ["192.168.240.0/24"]
}
```

## Example Usage for 9.1.0 onwards

```hcl
data "nsxt_policy_tier0_gateway" "test" {
  display_name = "test-t0gw"
}

resource "nsxt_policy_ip_block" "block1" {
  display_name = "ip-block1"
  cidr         = "168.1.1.0/24"
  visibility   = "EXTERNAL"
}

resource "nsxt_policy_ip_block" "block2" {
  display_name = "ip-block1"
  cidr         = "172.1.1.0/24"
  visibility   = "EXTERNAL"
}

resource "nsxt_policy_gateway_connection" "test" {
  display_name = "test"
  description  = "Terraform provisioned GatewayConnection"
  tier0_path   = data.nsxt_policy_tier0_gateway.test.path

  aggregate_routes        = ["192.168.240.0/24"]
  inbound_remote_networks = ["196.1.1.0/24"]

  advertise_outbound_networks {
    allow_external_blocks = [nsxt_policy_ip_block.block1.path, nsxt_policy_ip_block.block2.path]
    allow_private         = false
  }

  nat_config {
    enable_snat = true
    ip_block    = nsxt_policy_ip_block.block2.path
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tier0_path` - (Required) Tier-0 gateway object path
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `advertise_outbound_route_filters` - (Optional) List of prefixlist object paths that will have Transit gateway to tier-0 gateway advertise route filter.
* `aggregate_routes` - (Optional) Configure aggregate TGW_PREFIXES routes on Tier-0 gateway for prefixes owned by TGW gateway.
If not specified then in-use prefixes are configured as TGW_PREFIXES routes on Tier-0 gateway.
* `inbound_remote_networks` - (Optional) List of inbound remote network routes for transit gateways. This attribute is supported with NSX 9.1.0 onwards.
* `advertise_outbound_networks` - (Optional) Configure advertisement of outbound networks. This attribute is supported with NSX 9.1.0 onwards.
  * `allow_external_blocks` - (Optional) IP block paths used in advertisement filter to advertise prefixes from transit gateway.
  * `allow_private` - (Optional) Setting to true allows tenant user to configure advertisement rules and nat. If set to true, snat can not be enabled for this connection.
* `nat_config` - (Optional) List of inbound remote network routes for transit gateways. This attribute is supported with NSX 9.1.0 onwards.
  * `enable_snat` - (Optional) Enable the provider managed SNAT rule with translated ip from ip blocks. Defaults to `false`. Note that if set to `true`, `allow_private` in `advertise_outbound_networks` will not be supported.
  * `ip_block` - (Optional) IP block path (must be part of `allow_external_blocks` for advertisement)
  * `logging_enabled` - (Optional) Enable NAT translation logging. Defaults to `false`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_gateway_connection.test PATH
```

The above command imports GatewayConnection named `test` with the policy path `PATH`.
