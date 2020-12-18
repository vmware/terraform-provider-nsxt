---
subcategory: "Policy - FIXME"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_tier0_gateway_community_list"
description: A resource to configure a Tier0GatewayCommunityList.
---

# nsxt_policy_tier0_gateway_community_list

This resource provides a method for the management of a Tier0GatewayCommunityList.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_tier0_gateway_community_list" "test" {
    display_name      = "test"
    description       = "Terraform provisioned Tier0GatewayCommunityList"
    
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_tier0_gateway_community_list.test UUID
```

The above command imports Tier0GatewayCommunityList named `test` with the NSX Tier0GatewayCommunityList ID `UUID`.
