---
subcategory: "Beta"
layout: "nsxt"
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



## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_gateway_connection.test PATH
```

The above command imports GatewayConnection named `test` with the policy path `PATH`.
