---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_vpc_static_route"
description: A resource to configure a Static Routes under VPC.
---

# nsxt_vpc_static_route

This resource provides a method for the management of VPC Static Routes.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_vpc_static_route" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }

  display_name = "test"
  description  = "Terraform provisioned StaticRoutes"

  network = "3.3.3.0/24"

  next_hop {
    ip_address     = "10.230.3.1"
    admin_distance = 4
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `network` - (Required) Specify network address in CIDR format. Optionally this can be allocated IP from one of the external blocks associated with VPC. Only /32 CIDR is allowed in case IP overlaps with external blocks.
* `next_hop` - (Required) Specify next hop routes for network.
  * `ip_address` - (Optional) Next hop gateway IP address
  * `admin_distance` - (Optional) Cost associated with next hop route. Default is 1.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_vpc_static_route.test <PATH>
```

The above command imports Static Route named `test` with the NSX policy VPC path `PATH`.
