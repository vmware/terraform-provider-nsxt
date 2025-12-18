---
subcategory: "Segments"
page_title: "NSXT: nsxt_policy_segment_port_binding"
description: A resource to configure profile bindings for an existing Segment Port.
---

# nsxt_policy_segment_port_binding

This resource provides a method for managing profile bindings (discovery, QoS, and security profiles) for an existing Segment Port. This resource is useful when you need to modify the profiles of a segment port that was created outside of Terraform or by another resource.

~> **NOTE:** This resource modifies an existing segment port's profile bindings. It does not create a new segment port. The segment port must already exist.

## Example Usage

```hcl
data "nsxt_policy_segment" "segment1" {
  display_name = "existing-segment"
}

data "nsxt_policy_segment_port" "existing_port" {
  display_name = "existing-port"
  segment_path = data.nsxt_policy_segment.segment1.path
}

data "nsxt_policy_ip_discovery_profile" "ip_profile" {
  display_name = "default-ip-discovery-profile"
}

data "nsxt_policy_mac_discovery_profile" "mac_profile" {
  display_name = "default-mac-discovery-profile"
}

resource "nsxt_policy_segment_port_binding" "existing_port_binding" {
  segment_port_id = data.nsxt_policy_segment_port.existing_port.id
  segment_path    = data.nsxt_policy_segment.segment1.path

  discovery_profile {
    ip_discovery_profile_path  = data.nsxt_policy_ip_discovery_profile.ip_profile.path
    mac_discovery_profile_path = data.nsxt_policy_mac_discovery_profile.mac_profile.path
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_segment_port_binding" "port1_binding" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }

  segment_port_id = "port1"
  segment_path    = data.nsxt_policy_segment.segment1.path

  discovery_profile {
    ip_discovery_profile_path  = data.nsxt_policy_ip_discovery_profile.profile1.path
    mac_discovery_profile_path = data.nsxt_policy_mac_discovery_profile.profile1.path
  }

  security_profile {
    spoofguard_profile_path = data.nsxt_policy_spoofguard_profile.profile1.path
    security_profile_path   = data.nsxt_policy_segment_security_profile.profile1.path
  }
}
```

## Argument Reference

The following arguments are supported:

* `segment_port_id` - (Required) The ID of the existing segment port to bind profiles to. Changing this forces a new resource to be created.
* `segment_path` - (Required) The policy path of the segment that contains the port. Changing this forces a new resource to be created.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `discovery_profile` - (Optional) IP and MAC discovery profiles for this segment port.
    * `ip_discovery_profile_path` - (Optional) Policy path of the IP Discovery Profile to bind.
    * `mac_discovery_profile_path` - (Optional) Policy path of the MAC Discovery Profile to bind.
* `qos_profile` - (Optional) QoS profile for this segment port.
    * `qos_profile_path` - (Required) Policy path of the QoS Profile to bind.
* `security_profile` - (Optional) Security profiles for this segment port.
    * `spoofguard_profile_path` - (Optional) Policy path of the Spoofguard Profile to bind.
    * `security_profile_path` - (Optional) Policy path of the Segment Security Profile to bind.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Segment Port (same as `segment_port_id`).
* `discovery_profile`:
    * `binding_map_path` - Policy path of the discovery profile binding map.
    * `revision` - Revision number of the binding map.
* `qos_profile`:
    * `binding_map_path` - Policy path of the QoS profile binding map.
    * `revision` - Revision number of the binding map.
* `security_profile`:
    * `binding_map_path` - Policy path of the security profile binding map.
    * `revision` - Revision number of the binding map.

## Resource Lifecycle

~> **NOTE:** This resource manages profile bindings for an existing segment port. When the resource is destroyed (via `terraform destroy` or resource removal), the profile bindings are **not** removed from the segment port. This is by design to prevent unintended changes to existing infrastructure. If you need to remove profile bindings, you must do so manually or by updating the segment port directly.