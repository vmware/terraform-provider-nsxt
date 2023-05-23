---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_gateway_community_list"
description: A resource to configure Community List on Tier0 Gateway.
---

# nsxt_policy_gateway_community_list

This resource provides a method for the management of Community List on Tier0 Gateway.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_gateway_community_list" "test" {
  display_name = "test"
  description  = "Terraform provisioned list"
  gateway_path = nsxt_policy_tier0_gateway.test.path
  communities  = ["65001:12", "65002:12"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `gateway_path` - (Required) Policy path of relevant Tier0 Gateway.
* `comminities`  - (Required) List of BGP communities. Valid values are well-known communities, such as `NO_EXPORT`, `NO_ADVERTISE`, `NO_EXPORT_SUBCONFED`; Standard community notation `aa:nn`; Large community notation `aa:bb:nn`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_gateway_community_list.test GW-ID/ID
```

The above command imports Tier0 Gateway Community List named `test` with the NSX Community List ID `ID` on Tier0 Gateway `GW-ID`.
