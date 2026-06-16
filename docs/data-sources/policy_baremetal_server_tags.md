---
subcategory: "Beta"
page_title: "NSXT: policy_baremetal_server_tags"
description: A Bare Metal Server Tags data source.
---

# nsxt_policy_baremetal_server_tags

This data source provides information about tags applied to a Bare Metal Server in NSX-T Policy.

This data source is applicable to NSX Policy Manager and requires NSX-T version 9.0.0 or higher (Bare Metal Server support)

## Example Usage

```hcl
data "nsxt_policy_baremetal_server_tags" "test" {
  external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
}
```

## Example Usage - Using with BMS Server Discovery

```hcl
# First discover BMS servers
data "nsxt_policy_baremetal_servers" "production_servers" {
  display_name = "prod-server"
}

# Then get tags for a specific server
data "nsxt_policy_baremetal_server_tags" "server_tags" {
  external_id = data.nsxt_policy_baremetal_servers.production_servers.results[0].external_id
}

# Use the tags in a local value or output
output "server_environment" {
  value = [
    for tag in data.nsxt_policy_baremetal_server_tags.server_tags.tag :
    tag.tag if tag.scope == "environment"
  ][0]
}
```

## Argument Reference

* `external_id` - (Required) External ID of the Bare Metal Server to retrieve tags for.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the data source (same as `external_id`).
* `tag` - Set of tags applied to the Bare Metal Server. Each tag has the following attributes:
    * `scope` - Tag scope.
    * `tag` - Tag value.

## Notes

* This data source retrieves tags that have been applied to BMS servers using the `nsxt_policy_baremetal_server_tags` resource or other NSX-T mechanisms.
* Tags are used for dynamic group membership and policy application in NSX-T.
* The `external_id` must be the external ID of a Bare Metal Server that exists in NSX-T inventory.
