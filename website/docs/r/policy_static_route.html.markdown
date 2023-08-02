---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_static_route"
description: A resource to configure Static Routes in NSX Policy manager.
---

# nsxt_policy_static_route

This resource provides a method for the management of a Static Route.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_static_route" "route1" {
  display_name = "sroute"
  gateway_path = nsxt_policy_tier0_gateway.tier0_gw.path
  network      = "13.1.1.0/24"

  next_hop {
    admin_distance = "2"
    ip_address     = "11.10.10.1"
  }

  next_hop {
    admin_distance = "4"
    ip_address     = "12.10.10.1"
  }

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_static_route" "route1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "sroute"
  gateway_path = nsxt_policy_tier1_gateway.tier1_gw.path
  network      = "13.1.1.0/24"

  next_hop {
    admin_distance = "2"
    ip_address     = "11.10.10.1"
  }

  next_hop {
    admin_distance = "4"
    ip_address     = "12.10.10.1"
  }

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
* `tag` - (Optional) A list of scope + tag pairs to associate with this Tier-0 gateway.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `network` - (Required) The network address in CIDR format for the route.
* `gateway_path` (Required) The NSX Policy path to the Tier0 or Tier1 Gateway for this Static Route.
* `next_hop` - (Required) One or more next hops for the static route.
  * `admin_distance` - (Optional) The cost associated with the next hop. Valid values are 1 - 255 and the default is 1.
  * `ip_address` - (Optional) The gateway address of the next hop.
  * `interface` - (Optional) The policy path to the interface associated with the static route.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing policy Static Route can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_static_route.route1 GWID/ID
```
The above command imports the policy Static Route named `route1` for the NSX Tier0 or Tier1 Gateway `GWID` with the NSX Policy ID `ID`.

```
terraform import nsxt_policy_static_route.route1 POLICY_PATH
```
The above command imports the policy Static Route named `route1` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
