---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_baremetal_server_interface"
description: A bare metal server interface inventory data source.
---

# nsxt_policy_baremetal_server_interface

This data source provides information about a single bare metal server interface from NSX inventory. Use either `external_id` (exact match) or `display_name` (exact match returning exactly one result).

This data source is applicable to NSX Policy Manager and requires NSX-T version 9.0.0 or higher (Bare Metal Server support)

## Example Usage

```hcl
# Find interface by external ID
data "nsxt_policy_baremetal_server_interface" "baremetalinterface" {
  external_id = "71be0142-2ed1-1d53-9c60-02005b4b7246"
}

# Find interface by display name (must return exactly one result)
data "nsxt_policy_baremetal_server_interface" "eth0" {
  display_name = "eth0"
}

# Use interface information in other resources
resource "nsxt_policy_baremetal_server_interface_tags" "interface_tags" {
  external_id = data.nsxt_policy_baremetal_server_interface.baremetalinterface.external_id

  tag {
    scope = "network_role"
    tag   = "data_plane"
  }
}

output "interface_ips" {
  value = data.nsxt_policy_baremetal_server_interface.baremetalinterface.ip_addresses
}
```

## Argument Reference

* `external_id` - (Optional) External ID of the bare metal server interface for exact match lookup. Default is empty (must specify either external_id or display_name).
* `display_name` - (Optional) Display name of the bare metal server interface for exact match lookup. Default is empty (must specify either external_id or display_name).

**Note**: Either `external_id` or `display_name` must be specified. If using `display_name`, exactly one interface must match.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the data source.
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
