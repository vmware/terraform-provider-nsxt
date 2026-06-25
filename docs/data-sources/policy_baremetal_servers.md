---
subcategory: "Beta"
page_title: "NSXT: policy_baremetal_servers"
description: List of bare metal servers discovered by NSX with advanced filtering.
---

# nsxt_policy_baremetal_servers

This data source provides information about bare metal servers discovered by NSX with advanced server-side filtering capabilities including tag-based search, OS filtering, and regex support.

This data source is applicable to NSX Policy Manager and requires NSX-T version 9.0.0 or higher (Bare Metal Server support)

## Example Usage

```hcl
# Get all bare metal servers
data "nsxt_policy_baremetal_servers" "all_servers" {}

# Filter by production environment using server-side tag filtering
data "nsxt_policy_baremetal_servers" "prod" {
  tag_scope = "environment"
  tag       = "production"
}

# Filter by OS name
data "nsxt_policy_baremetal_servers" "redhat" {
  os_name = "redhat"
}

# Filter by display name using regex
data "nsxt_policy_baremetal_servers" "web_servers" {
  display_name = "web-.*"
}

# Multiple filters
data "nsxt_policy_baremetal_servers" "prod_web" {
  display_name = "web-.*"
  tag_scope    = "environment"
  tag          = "production"
  os_name      = "ubuntu"
}

output "prod_external_ids" {
  value = [for s in data.nsxt_policy_baremetal_servers.prod.results : s.external_id]
}
```

## Argument Reference

* `display_name` - (Optional) Regex pattern to filter bare metal servers by display name. Default is empty (no filtering).
* `source_id` - (Optional) Filter by source (bare metal controller) ID. Default is empty (no filtering).
* `os_name` - (Optional) Filter by OS name (case-insensitive substring match). Default is empty (no filtering).
* `os_version` - (Optional) Filter by OS version (case-insensitive substring match). Default is empty (no filtering).
* `tag_scope` - (Optional) Filter by tag scope. This filter is applied server-side for better performance. Default is empty (no filtering).
* `tag` - (Optional) Filter by tag value. This filter is applied server-side for better performance. Default is empty (no filtering).

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the data source.
* `results` - List of bare metal servers matching the criteria. Each server contains:
    * `external_id` - External ID of the bare metal server.
    * `display_name` - Display name of the bare metal server.
    * `source_id` - Source ID of the bare metal server.
    * `cpu_cores` - Number of CPU cores.
    * `os_name` - Operating system name.
    * `os_version` - Operating system version.
    * `resource_type` - Resource type.
    * `last_sync_time` - Last synchronization time.
    * `tags` - List of tags applied to the server. Each tag contains:
        * `scope` - Tag scope.
        * `tag` - Tag value.
