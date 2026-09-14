---
subcategory: "VPC"
page_title: "NSXT: policy_transit_gateway_community_list"
description: A policy BGP Community List on a Transit Gateway data source.
---

# nsxt_policy_transit_gateway_community_list

This data source provides information about a BGP community list configured on a Transit Gateway.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_transit_gateway_community_list" "test" {
  display_name = "my-community-list"
  parent_path  = nsxt_policy_transit_gateway.demo.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the Transit Gateway Community List to retrieve.
* `display_name` - (Optional) The Display Name of the Transit Gateway Community List to retrieve.
* `parent_path` - (Required) Policy path of the parent Transit Gateway. This is used to scope the search, since the underlying `CommunityList` object type is shared with Tier-0 gateway community lists.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
