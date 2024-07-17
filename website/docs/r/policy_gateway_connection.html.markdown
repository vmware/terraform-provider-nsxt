---
subcategory: "FIXME"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_gateway_connection"
description: A resource to configure a GatewayConnection.
---

# nsxt_policy_gateway_connection

This resource provides a method for the management of a GatewayConnection.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_gateway_connection" "test" {
    display_name      = "test"
    description       = "Terraform provisioned GatewayConnection"
    advertise_outbound_route_filter = FILL VALUE FOR schema.TypeString
tier0_path = FILL VALUE FOR schema.TypeString
aggregate_routes = FILL VALUE FOR schema.TypeString

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `advertise_outbound_route_filter` - (Optional) Path of a prefixlist object that will have Transit gateway to tier-0 gateway advertise route filter.

* `tier0_path` - (Optional) Tier-0 gateway object path

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
terraform import nsxt_policy_gateway_connection.test UUID
```

The above command imports GatewayConnection named `test` with the NSX ID `UUID`.
