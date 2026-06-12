---
subcategory: "Beta"
page_title: "NSXT: policy_baremetal_server_interface_tags"
description: A resource for applying tags to Bare Metal Server Interfaces.
---

# nsxt_policy_baremetal_server_interface_tags

This resource provides a way to configure tags on Bare Metal Server Interfaces discovered and managed by NSX-T.
Tags applied to bare metal server interfaces can be used in dynamic group membership criteria and security policies.

## Example Usage

```hcl
# Discover bare metal server interfaces
data "nsxt_policy_baremetal_server_interfaces" "server_interfaces" {
  bms_external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
}

# Apply tags to a specific interface
resource "nsxt_policy_baremetal_server_interface_tags" "eth0_tags" {
  external_id = data.nsxt_policy_baremetal_server_interfaces.server_interfaces.results[0].external_id

  tag {
    scope = "network-type"
    tag   = "data-plane"
  }

  tag {
    scope = "vlan"
    tag   = "100"
  }

  tag {
    scope = "speed"
    tag   = "10Gbps"
  }
}

# Use tagged interfaces in a dynamic group
resource "nsxt_policy_group" "data_plane_interfaces" {
  display_name = "Data-Plane-Interfaces"
  description  = "Dynamic group of data plane interfaces"

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

## Example with Interface Filtering

```hcl
# Tag only non-management interfaces
resource "nsxt_policy_baremetal_server_interface_tags" "data_interfaces" {
  for_each = {
    for idx, iface in data.nsxt_policy_baremetal_server_interfaces.all.results :
    idx => iface.external_id
    if !iface.is_mgmt_interface && iface.state == "UP"
  }

  external_id = each.value

  tag {
    scope = "interface-role"
    tag   = "data-plane"
  }

  tag {
    scope = "managed-by"
    tag   = "terraform"
  }
}

# Create group for high-speed interfaces
resource "nsxt_policy_group" "high_speed_interfaces" {
  display_name = "High-Speed-Interfaces"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "speed|10Gbps"
    }
  }

  conjunction {
    operator = "OR"
  }

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "speed|25Gbps"
    }
  }
}
```

## Example with VLAN-based Segmentation

```hcl
# Tag interfaces by VLAN for network segmentation
resource "nsxt_policy_baremetal_server_interface_tags" "vlan_100_interfaces" {
  count = length([
    for iface in data.nsxt_policy_baremetal_server_interfaces.all.results :
    iface if length(regexall("vlan-100", lower(iface.display_name))) > 0
  ])

  external_id = [
    for iface in data.nsxt_policy_baremetal_server_interfaces.all.results :
    iface.external_id if length(regexall("vlan-100", lower(iface.display_name))) > 0
  ][count.index]

  tag {
    scope = "vlan"
    tag   = "100"
  }

  tag {
    scope = "network-segment"
    tag   = "dmz"
  }
}

# Create VLAN-specific group
resource "nsxt_policy_group" "vlan_100_group" {
  display_name = "VLAN-100-Interfaces"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "vlan|100"
    }
  }
}
```

## Argument Reference

* `external_id` - (Required) External ID of the bare metal server interface to tag. This is typically obtained from the `nsxt_policy_baremetal_server_interfaces` data source.
* `tag` - (Optional) A list of scope + tag pairs to associate with this bare metal server interface. Default is empty (no tags applied).

The `tag` block supports:

* `scope` - (Required) Tag scope.
* `tag` - (Required) Tag value.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the bare metal server interface being tagged.

## Import

An NSX Policy Bare Metal Server Interface Tags resource can be imported using the interface external ID:

```shell
terraform import nsxt_policy_baremetal_server_interface_tags.eth0_tags 71be0142-2ed1-1d53-9c60-02005b4b7246
```

The above command imports the interface tags for the interface with external ID `71be0142-2ed1-1d53-9c60-02005b4b7246`.

## Notes

* This resource requires NSX-T version 9.0.0 or higher (Bare Metal Server support)
* The bare metal server interface must be discovered and registered with NSX-T through a compute manager
* Tags applied through this resource will be visible in the NSX-T UI and available for use in policy groups
* Use the `is_mgmt_interface` attribute from the data source to distinguish between management and data plane interfaces
* Only tag operations are supported; other interface properties cannot be modified through this resource
* This resource uses NSX Policy APIs and requires local manager access (not supported on Global Manager)
* When the resource is destroyed, all tags managed by this resource are removed from the interface
* Consider using interface naming conventions or existing attributes to automate tag assignment
* Tags can be used to implement network segmentation and micro-segmentation policies
