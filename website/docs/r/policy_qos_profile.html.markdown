---
subcategory: "Segments"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_qos_profile"
description: A resource to configure a QoS profile.
---

# nsxt_policy_qos_profile

This resource provides a method for the management of Qos profiles.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_qos_profile" "qos_profile" {
  description      = "qos profile provisioned by Terraform"
  display_name     = "qos_profile1"
  class_of_service = 5
  dscp_trusted     = true
  dscp_priority    = 53

  ingress_rate_shaper {
    enabled         = true
    peak_bw_mbps    = 800
    burst_size      = 200
    average_bw_mbps = 100
  }

  egress_rate_shaper {
    enabled         = true
    peak_bw_mbps    = 800
    burst_size      = 200
    average_bw_mbps = 100
  }

  ingress_broadcast_rate_shaper {
    enabled         = true
    average_bw_kbps = 111
    burst_size      = 222
    peak_bw_kbps    = 500
  }

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_qos_profile" "qos_profile" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  description      = "qos profile provisioned by Terraform"
  display_name     = "qos_profile1"
  class_of_service = 5
  dscp_trusted     = true
  dscp_priority    = 53

  ingress_rate_shaper {
    enabled         = true
    peak_bw_mbps    = 800
    burst_size      = 200
    average_bw_mbps = 100
  }

  egress_rate_shaper {
    enabled         = true
    peak_bw_mbps    = 800
    burst_size      = 200
    average_bw_mbps = 100
  }

  ingress_broadcast_rate_shaper {
    enabled         = true
    average_bw_kbps = 111
    burst_size      = 222
    peak_bw_kbps    = 500
  }

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `class_of_service` - (Optional) Class of service.
* `dscp_trusted` - (Optional) Trust mode for DSCP (False by default)
* `dscp_priority` - (Optional) DSCP Priority (0-63)
* `ingress_rate_shaper` - (Optional) Ingress rate shaper configuration:
  * `enabled` - (Optional) Whether this rate shaper is enabled.
  * `average_bw_mbps` - (Optional) Average Bandwidth in MBPS.
  * `peak_bw_mbps` - (Optional) Peak Bandwidth in MBPS.
  * `burst_size` - (Optional) Burst size in bytes.
* `egress_rate_shaper` - (Optional) Egress rate shaper configuration:
  * `enabled` - (Optional) Whether this rate shaper is enabled.
  * `average_bw_mbps` - (Optional) Average Bandwidth in MBPS.
  * `peak_bw_mbps` - (Optional) Peak Bandwidth in MBPS.
  * `burst_size` - (Optional) Burst size in bytes.
* `ingress_broadcast_rate_shaper` - (Optional) Ingress rate shaper configuration:
  * `enabled` - (Optional) Whether this rate shaper is enabled.
  * `average_bw_kbps` - (Optional) Average Bandwidth in KBPS.
  * `peak_bw_kbps` - (Optional) Peak Bandwidth in KBPS.
  * `burst_size` - (Optional) Burst size in bytes.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_qos_profile.qos_profile ID
```
The above command imports the qos profile named `qos_profile` with the NSX ID `ID`.

```
terraform import nsxt_policy_qos_profile.qos_profile POLICY_PATH
```
The above command imports the qos profile named `qos_profile` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
