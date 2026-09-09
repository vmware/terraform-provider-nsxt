---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_baremetal_server_interfaces"
description: List of bare metal server interfaces in NSX inventory with advanced filtering.
---

# nsxt_policy_baremetal_server_interfaces

This data source provides information about bare metal server interfaces discovered by NSX with advanced server-side filtering capabilities including tag-based search, parent server filtering, and regex support.

This data source is applicable to NSX Policy Manager and requires NSX-T version 9.0.0 or higher (Bare Metal Server support)

## Example Usage

```hcl
# Get all bare metal server interfaces
data "nsxt_policy_baremetal_server_interfaces" "all_interfaces" {}

# Filter by parent bare metal server
data "nsxt_policy_baremetal_server_interfaces" "bm1_interfaces" {
  bms_external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
}

# Filter by network role using server-side tag filtering
data "nsxt_policy_baremetal_server_interfaces" "data_interfaces" {
  tag_scope = "network_role"
  tag       = "data"
}

# Filter by interface name using regex
data "nsxt_policy_baremetal_server_interfaces" "eth_interfaces" {
  display_name = "eth.*"
}

# Multiple filters
data "nsxt_policy_baremetal_server_interfaces" "prod_data_eth" {
  bms_external_id = data.nsxt_policy_baremetal_servers.prod.results[0].external_id
  display_name    = "eth.*"
  tag_scope       = "network_role"
  tag             = "data"
}

output "interface_ips" {
  value = [for i in data.nsxt_policy_baremetal_server_interfaces.data_interfaces.results : i.ip_addresses]
}
```

## Argument Reference

* `display_name` - (Optional) Regex pattern to filter bare metal server interfaces by display name. Default is empty (no filtering).
* `bms_external_id` - (Optional) Filter by parent bare metal server external ID. Default is empty (no filtering).
* `source_id` - (Optional) Filter by source (bare metal controller) ID. Default is empty (no filtering).
* `tag_scope` - (Optional) Filter by tag scope. This filter is applied server-side for better performance. Default is empty (no filtering).
* `tag` - (Optional) Filter by tag value. This filter is applied server-side for better performance. Default is empty (no filtering).

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the data source.
* `results` - List of bare metal server interfaces matching the criteria. Each interface contains:
    * `external_id` - External ID of the bare metal server interface.
    * `display_name` - Display name of the bare metal server interface.
    * `bms_external_id` - External ID of the parent bare metal server.
    * `source_id` - Source ID of the bare metal server interface.
    * `is_mgmt_interface` - Whether this is a management interface.
    * `ip_addresses` - List of IP addresses assigned to the interface.
    * `state` - State of the interface (UP, DOWN, etc.).
    * `resource_type` - Resource type.
    * `last_sync_time` - Last synchronization time.
    * `tags` - List of tags applied to the interface. Each tag contains:
        * `scope` - Tag scope.
        * `tag` - Tag value.
