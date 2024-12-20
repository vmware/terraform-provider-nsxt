---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_transit_gateway_attachment"
description: A resource to configure a Transit Gateway Attachment.
---

# nsxt_policy_transit_gateway_attachment

This resource provides a method for the management of a Transit Gateway Attachment.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "test_proj" {
  display_name = "test_project"
}

data "nsxt_policy_transit_gateway" "test_tgw" {
  context {
    project_id = nsxt_policy_project.test_proj.id
  }
  id = "default"
}

resource "nsxt_policy_gateway_connection" "test_gw_conn" {
  display_name = "test_gw_conn"
  tier0_path   = "/infra/tier-0s/test-t0"
}

resource "nsxt_policy_transit_gateway_attachment" "test_tgw_att" {
  display_name    = "test"
  parent_path     = nsxt_policy_transit_gateway.test_tgw.path
  connection_path = nsxt_policy_gateway_connection.test_gw_conn.path
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required) The path of the parent to bind with the profile. This is a policy path of a transit gateway.
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `connection_path` - (Required) Policy path of the desired transit gateway external connection.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_transit_gateway_attachment.test PATH
```

The above command imports Transit Gateway Attachment named `test` with the policy path `PATH`.
