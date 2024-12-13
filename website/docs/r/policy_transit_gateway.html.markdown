---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_transit_gateway"
description: A resource to configure a Transit Gateway.
---

# nsxt_policy_transit_gateway

This resource provides a method for the management of a Transit Gateway.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = "dev"
  }

  display_name    = "test"
  description     = "Terraform provisioned TransitGateway"
  transit_subnets = ["10.203.4.0/24"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `transit_subnets` - (Optional) Array of IPV4 CIDRs for internal VPC attachment networks.



## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_transit_gateway.test PATH
```

The above command imports Transit Gateway named `test` with the NSX path `PATH`.
