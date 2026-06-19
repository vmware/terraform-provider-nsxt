---
subcategory: "Beta"
page_title: "NSXT: policy_baremetal_server"
description: A bare metal server inventory data source.
---

# nsxt_policy_baremetal_server

This data source provides information about a single bare metal server discovered by NSX. Use either `external_id` (exact match) or `display_name` (exact match returning exactly one result).

This data source is applicable to NSX Policy Manager and requires NSX-T version 9.0.0 or higher (Bare Metal Server support)

## Example Usage

```hcl
# Find server by external ID
data "nsxt_policy_baremetal_server" "baremetalserver" {
  external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
}

# Find server by display name (must return exactly one result)
data "nsxt_policy_baremetal_server" "web_server" {
  display_name = "web-server-01"
}

# Use server information in other resources
resource "nsxt_policy_baremetal_server_tags" "server_tags" {
  external_id = data.nsxt_policy_baremetal_server.baremetalserver.external_id

  tag {
    scope = "environment"
    tag   = "production"
  }
}

output "server_os" {
  value = "${data.nsxt_policy_baremetal_server.baremetalserver.os_name} ${data.nsxt_policy_baremetal_server.baremetalserver.os_version}"
}
```

## Argument Reference

* `external_id` - (Optional) External ID of the bare metal server for exact match lookup. Default is empty (must specify either external_id or display_name).
* `display_name` - (Optional) Display name of the bare metal server for exact match lookup. Default is empty (must specify either external_id or display_name).

**Note**: Either `external_id` or `display_name` must be specified. If using `display_name`, exactly one server must match.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the data source.
* `source_id` - Source ID of the bare metal server.
* `cpu_cores` - Number of CPU cores.
* `os_name` - Operating system name.
* `os_version` - Operating system version.
* `resource_type` - Resource type.
* `last_sync_time` - Last synchronization time.
* `tags` - List of tags applied to the server. Each tag contains:
    * `scope` - Tag scope.
    * `tag` - Tag value.
