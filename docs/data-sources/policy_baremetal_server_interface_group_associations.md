---
subcategory: "Beta"
page_title: "NSXT: policy_baremetal_server_interface_group_associations"
description: Groups a bare metal server interface is a member of.
---

# nsxt_policy_baremetal_server_interface_group_associations

This data source provides information about policy groups for which the given bare metal server interface is a member. This is useful for discovering group membership and understanding network segmentation.

This data source is applicable to NSX Policy Manager and requires NSX-T version 9.0.0 or higher (Bare Metal Server support)

## Example Usage

```hcl
# Get groups that a specific bare metal server interface belongs to
data "nsxt_policy_baremetal_server_interface_group_associations" "if1_groups" {
  external_id = "71be0142-2ed1-1d53-9c60-02005b4b7246"
}

# Define the data-network group referenced in the security policy
resource "nsxt_policy_group" "data-network" {
  display_name = "Data-Network-Group"
  description  = "Group containing data network interfaces"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "network_zone|data"
    }
  }
}

# Use the groups in security policies
resource "nsxt_policy_security_policy" "interface_policy" {
  display_name = "BMS Interface Security Policy"
  category     = "Infrastructure"

  rule {
    display_name       = "Allow Interface Traffic"
    source_groups      = data.nsxt_policy_baremetal_server_interface_group_associations.if1_groups.groups[*].path
    destination_groups = [nsxt_policy_group.data-network.path]
    action             = "ALLOW"
    services           = [] # Empty list means "ANY"
  }
}

output "interface_groups" {
  value = {
    for group in data.nsxt_policy_baremetal_server_interface_group_associations.if1_groups.groups :
    group.display_name => group.path
  }
}
```

## Argument Reference

* `external_id` - (Required) External ID of the bare metal server interface.
* `enforcement_point_path` - (Optional) Path of the enforcement point.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the data source.
* `groups` - List of groups this bare metal server interface is a member of. Each group contains:
    * `path` - Policy path of the group.
    * `display_name` - Display name of the group.
    * `target_type` - Type of the target resource.
    * `is_valid` - Indicates if the referenced NSX resource is valid.
