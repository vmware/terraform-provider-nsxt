---
subcategory: "Beta"
page_title: "NSXT: policy_baremetal_server_interface_tags"
description: A Bare Metal Server Interface Tags data source.
---

# nsxt_policy_baremetal_server_interface_tags

This data source provides information about tags applied to a Bare Metal Server Interface in NSX-T Policy.

This data source is applicable to NSX Policy Manager and requires NSX-T version 9.0.0 or higher (Bare Metal Server support)

## Example Usage

```hcl
data "nsxt_policy_baremetal_server_interface_tags" "test" {
  external_id = "71be0142-2ed1-1d53-9c60-02005b4b7246"
}
```

## Example Usage - Using with BMS Interface Discovery

```hcl
# First discover BMS interfaces
data "nsxt_policy_baremetal_server_interfaces" "data_interfaces" {
  bms_external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
}

# Then get tags for a specific interface
data "nsxt_policy_baremetal_server_interface_tags" "interface_tags" {
  external_id = data.nsxt_policy_baremetal_server_interfaces.data_interfaces.results[0].external_id
}

# Use tags to determine network type
locals {
  network_type = [
    for tag in data.nsxt_policy_baremetal_server_interface_tags.interface_tags.tag :
    tag.tag if tag.scope == "network-type"
  ][0]
}

# Conditional resource creation based on network type
resource "nsxt_policy_group" "data_plane_interfaces" {
  count = local.network_type == "data-plane" ? 1 : 0

  display_name = "Data-Plane-Interfaces"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "network-type|data-plane"
    }
  }
}
```

## Argument Reference

* `external_id` - (Required) External ID of the Bare Metal Server Interface to retrieve tags for.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the data source (same as `external_id`).
* `tag` - Set of tags applied to the Bare Metal Server Interface. Each tag has the following attributes:
    * `scope` - Tag scope.
    * `tag` - Tag value.

## Common Tag Scopes

Common tag scopes used for BMS interfaces include:

* `network-type` - Values like `data-plane`, `management`, `storage`
* `vlan` - VLAN ID for the interface
* `network-zone` - Network zone designation (e.g., `dmz`, `internal`)
* `environment` - Environment designation (e.g., `production`, `staging`)

## Notes

* This data source retrieves tags that have been applied to BMS interfaces using the `nsxt_policy_baremetal_server_interface_tags` resource or other NSX-T mechanisms.
* Interface tags are commonly used for network segmentation and policy application.
* The `external_id` must be the external ID of a Bare Metal Server Interface that exists in NSX-T inventory.
