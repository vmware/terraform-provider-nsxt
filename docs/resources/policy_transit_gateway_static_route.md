---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_transit_gateway_static_route"
description: A resource to configure a static route
---

# nsxt_policy_transit_gateway_static_route

This resource provides a method for the management of a transit gateway StaticRoute.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "proj" {
  display_name = "demoproj"
}

data "nsxt_policy_transit_gateway" "tgw1" {
  context {
    project_id = data.nsxt_policy_project.proj.id
  }
  display_name = "TGW1"
}

resource "nsxt_policy_transit_gateway_static_route" "test" {
  display_name         = "test"
  description          = "Terraform provisioned StaticRoute for transit gateway"
  parent_path          = data.nsxt_policy_transit_gateway.tgw1.path
  enabled_on_secondary = false
  network              = "3.3.3.0/24"
  next_hop {
    ip_address     = "192.168.1.1"
    admin_distance = 10
    scope          = [nsxt_policy_transit_gateway_attachment.test.path]
  }
}

```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `parent_path` - (Required) Path of parent object
* `network` - (Required) Specify network address in CIDR format.
* `next_hop` - (Required) List of next hops (only one supported),each with those arguments:
  * `administrative_distance` - (Optional) Administrative Distance for the next hop IP.
  * `ip_address` - (Optional) Next Hop IP.
  * `scope` - (Required) Set of policy object paths where the rule is applied.
* `enabled_on_secondary` - (Optional) Flag to plumb route on secondary site.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_transit_gateway_static_route.test PATH
```

The above command imports StaticRoute named `test` with the NSX path `PATH`.
