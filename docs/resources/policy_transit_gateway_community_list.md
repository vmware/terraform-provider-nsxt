---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_transit_gateway_community_list"
description: A resource to configure a BGP Community List on a Transit Gateway.
---

# nsxt_policy_transit_gateway_community_list

This resource provides a method for the management of a BGP community list on a Transit Gateway.

Transit gateway community lists use the same `CommunityList` schema as Tier-0 gateway community lists and are referenced by route map entries for BGP route policy.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "test_proj" {
  display_name = "test_project"
}

data "nsxt_policy_transit_gateway" "test_tgw" {
  context {
    project_id = data.nsxt_policy_project.test_proj.id
  }
  id = "default"
}

resource "nsxt_policy_transit_gateway_community_list" "example" {
  display_name = "my-community-list"
  description  = "BGP community list for route policy"
  parent_path  = data.nsxt_policy_transit_gateway.test_tgw.path

  communities = [
    "65000:100",
    "65000:200",
    "no-export",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required, ForceNew) Policy path of the parent Transit Gateway (e.g. `/orgs/default/projects/proj1/transit-gateways/tgw1`).
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `communities` - (Required) Set of BGP community string entries. Each entry must be a standard BGP community (`<ASN>:<value>`), a well-known community (`internet`, `no-export`, `no-advertise`, `local-AS`), or a large community (`<ASN>:<local-1>:<local-2>`).

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

```shell
terraform import nsxt_policy_transit_gateway_community_list.example PARENT_PATH/ID
```

The above command imports a Transit Gateway Community List named `example` with the transit gateway NSX path as `PARENT_PATH` and the community list ID as `ID`. For example:

```shell
terraform import nsxt_policy_transit_gateway_community_list.example /orgs/default/projects/project1/transit-gateways/tgw1/community-list-1
```

[docs-import]: https://developer.hashicorp.com/terraform/cli/import
