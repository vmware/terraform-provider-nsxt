---
subcategory: "Beta"
page_title: "NSXT: policy_group_baremetal_server_members"
description: A policy data source to list bare metal server members in a group.
---

# nsxt_policy_group_baremetal_server_members

This data source provides a way to list effective bare metal server members that belong to a specific NSX-T policy group. This API is applicable only for groups containing BMS member type. For groups containing other member types, it returns an empty list.

This data source is applicable to NSX Policy Manager and requires NSX-T version 9.0.0 or higher (Bare Metal Server support)

## Example Usage

```hcl
data "nsxt_policy_group_baremetal_server_members" "example" {
  domain   = "default"
  group_id = "bms-group-id"
}

output "bms_members" {
  value = data.nsxt_policy_group_baremetal_server_members.example.items
}
```

## Example Usage with Group Reference

```hcl
resource "nsxt_policy_group" "bms_group" {
  display_name = "bms-static-group"
  description  = "Group for bare metal servers"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type  = "BareMetalServer"
      external_ids = ["71be0142-2ed1-1d53-9c60-5564cf4b7e2e"]
    }
  }
}

data "nsxt_policy_group_baremetal_server_members" "bms_members" {
  domain   = "default"
  group_id = nsxt_policy_group.bms_group.id
}
```

## Argument Reference

* `domain` - (Optional) Domain ID for the group. Defaults to "default".
* `group_id` - (Required) ID of the group to read effective bare metal server members.
* `enforcement_point_path` - (Optional) The path of the enforcement point from which the list of members needs to be fetched.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the group.
* `items` - List of bare metal server members in the group.
    * `external_id` - External ID of the bare metal server.
    * `display_name` - Display name of the bare metal server.
    * `source_id` - Source ID of the bare metal server.
    * `cpu_cores` - Number of CPU cores on the bare metal server.
    * `os_name` - Operating system name.
    * `os_version` - Operating system version.
    * `last_sync_time` - Last synchronization time.
