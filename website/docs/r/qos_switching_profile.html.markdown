---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_qos_switching_profile"
description: |-
  Provides a resource to configure QoS switching profile on NSX-T manager
---

# nsxt_qos_switching_profile

Provides a resource to configure Qos switching profile on NSX-T manager

## Example Usage

```hcl
resource "nsxt_qos_switching_profile" "qos_switching_profile" {
  description      = "qos_switching_profile provisioned by Terraform"
  display_name     = "qos_switching_profile"
  class_of_service = "5"
  dscp_trusted     = "true"
  dscp_priority    = "53"

  ingress_rate_shaper {
    enabled         = "true"
    peak_bw_mbps    = "800"
    burst_size      = "200"
    average_bw_mbps = "100"
  }

  egress_rate_shaper {
    enabled         = "true"
    peak_bw_mbps    = "800"
    burst_size      = "200"
    average_bw_mbps = "100"
  }

  ingress_broadcast_rate_shaper {
    enabled         = "true"
    average_bw_kbps = "111"
    burst_size      = "222"
    peak_bw_kbps    = "500"
  }

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this qos switching profile.
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

* `id` - ID of the QoS switching profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing qos switching profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_qos_switching_profile.qos_switching_profile UUID
```

The above would import the Qos switching profile named `qos_switching_profile` with the nsx id `UUID`
