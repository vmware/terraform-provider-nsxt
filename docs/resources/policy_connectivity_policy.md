---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_connectivity_policy"
description: A resource to configure Transit Gateway Connectivity Policy.
---

# nsxt_policy_connectivity_policy

This resource provides a method for the management of TGW Connectivity Policy.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_connectivity_policy" "test" {
  display_name = "test"
  description  = "Terraform provisioned Connectivity Policy"

  parent_path        = data.nsxt_policy_transit_gateway.default.path
  group_path         = data.nsxt_policy_group.test.path
  connectivity_scope = "COMMUNITY"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `parent_path` - (Required) Path of parent transit gateway.
* `group_path` - (Required) Path of group that the policy is applied to.
* `connectivity_scope` - (Optional) Either `COMMUNITY` or `ISOLATED`. Default is `ISOLATED`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `internal_id` - Internal ID of the resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_connectivity_policy.test PATH
```

The above command imports Connectivity Policy named `test` with the NSX path `PATH`.
