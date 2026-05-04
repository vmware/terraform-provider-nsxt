---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_transit_gateway_prefix_list"
description: A resource to configure a Prefix List on a Transit Gateway.
---

# nsxt_policy_transit_gateway_prefix_list

This resource provides a method for the management of a prefix list on a Transit Gateway.

Transit gateway prefix lists use the same `PrefixList` schema as Tier-0 gateways and are used for BGP route filtering.

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

resource "nsxt_policy_transit_gateway_prefix_list" "example" {
  display_name = "my-prefix-list"
  description  = "Prefix list for BGP route filtering"
  parent_path  = data.nsxt_policy_transit_gateway.test_tgw.path

  prefix {
    action  = "PERMIT"
    network = "10.0.0.0/8"
  }

  prefix {
    action  = "DENY"
    network = "192.168.0.0/16"
    ge      = 24
    le      = 32
  }

  prefix {
    action = "DENY"
  }
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required, ForceNew) Policy path of the parent Transit Gateway (e.g. `/orgs/default/projects/proj1/transit-gateways/tgw1`).
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `prefix` - (Required) Ordered list of network prefix entries. At least one entry is required. Each entry supports:
    * `action` - (Optional) Action for the prefix list entry. Accepted values: `PERMIT`, `DENY`. Defaults to `PERMIT`.
    * `network` - (Optional) Network prefix in CIDR format. If not set, matches ANY network.
    * `ge` - (Optional) Minimum prefix length (greater than or equal to). Valid range 0–128. Defaults to `0` (not applied).
    * `le` - (Optional) Maximum prefix length (less than or equal to). Valid range 0–128. Defaults to `0` (not applied).

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

```shell
terraform import nsxt_policy_transit_gateway_prefix_list.example PARENT_PATH/ID
```

The above command imports a Transit Gateway Prefix List named `example` with the transit gateway NSX path as `PARENT_PATH` and the prefix list ID as `ID`. For example:

```shell
terraform import nsxt_policy_transit_gateway_prefix_list.example /orgs/default/projects/project1/transit-gateways/tgw1/prefix-list-1
```

[docs-import]: https://developer.hashicorp.com/terraform/cli/import
