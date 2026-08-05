---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_group_baremetal_server_interface_members"
description: A policy data source to list bare metal server interface members in a group.
---

# nsxt_policy_group_baremetal_server_interface_members

This data source provides a way to list effective bare metal server interface members that belong to a specific NSX-T policy group. This API is applicable only for groups containing BMSI member type. For groups containing other member types, it returns an empty list.

This data source is applicable to NSX Policy Manager and requires NSX-T version 9.0.0 or higher (Bare Metal Server support)

## Example Usage

```hcl
data "nsxt_policy_group_baremetal_server_interface_members" "example" {
  domain   = "default"
  group_id = "bmsi-group-id"
}

output "bmsi_members" {
  value = data.nsxt_policy_group_baremetal_server_interface_members.example.items
}
```

## Example Usage with Group Reference

```hcl
resource "nsxt_policy_group" "bmsi_group" {
  display_name = "bmsi-static-group"
  description  = "Group for bare metal server interfaces"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type  = "BareMetalServerInterface"
      external_ids = ["71be0142-2ed1-1d53-9c60-02005b4b7246"]
    }
  }
}

data "nsxt_policy_group_baremetal_server_interface_members" "bmsi_members" {
  domain   = "default"
  group_id = nsxt_policy_group.bmsi_group.id
}
```

## Argument Reference

* `domain` - (Optional) Domain ID for the group. Defaults to "default".
* `group_id` - (Required) ID of the group to read effective bare metal server interface members.
* `enforcement_point_path` - (Optional) The path of the enforcement point from which the list of members needs to be fetched.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the group.
* `items` - List of bare metal server interface members in the group.
    * `external_id` - External ID of the bare metal server interface.
    * `display_name` - Display name of the bare metal server interface.
    * `bms_external_id` - External ID of the parent bare metal server.
    * `source_id` - Source ID of the bare metal server interface.
    * `mac_address` - MAC address of the interface.
    * `is_mgmt_interface` - Whether this is a management interface.
    * `state` - State of the interface (e.g., "UP", "DOWN").
    * `ip_addresses` - List of IP addresses assigned to the interface.
    * `last_sync_time` - Last synchronization time.
