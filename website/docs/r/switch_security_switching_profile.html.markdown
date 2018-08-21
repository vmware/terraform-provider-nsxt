---
layout: "nsxt"
page_title: "NSXT: nsxt_switch_security_switching_profile"
sidebar_current: "docs-nsxt-switching-profile-resource"
description: |-
  Provides a resource to configure switch security switching profile on NSX-T manager
---

# nsxt_switch_security_switching_profile

Provides a resource to configure switch security switching profile on NSX-T manager

## Example Usage

```hcl
resource "nsxt_switch_security_switching_profile" "switch_security_switching_profile" {
  description      = "switch_security_switching_profile provisioned by Terraform"
  display_name     = "switch_security_switching_profile"
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

  tag = {
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
* `block_non_ip` - (Optional) Indicates whether blocking of all traffic except IP/(G)ARP/BPDU is enabled.
* `block_client_dhcp` - (Optional) Indicates whether DHCP client blocking is enabled
* `block_server_dhcp` - (Optional) Indicates whether DHCP server blocking is enabled
* `bpdu_filter_enabled` - (Optional) Indicates whether BPDU filter is enabled
* `bpdu_filter_whitelist` - (Optional) Set of allowed MAC addresses to be excluded from BPDU filtering, if enabled.
* `rate_limits` - (Optional) Rate limit definitions for broadcast and multicast traffic.
  * `enabled` - (Optional) Whether rate limitimg is enabled.
  * `rx_broadcast` - (Optional) Incoming broadcast traffic limit in packets per second.
  * `rx_multicast` - (Optional) Incoming multicast traffic limit in packets per second.
  * `tx_broadcast` - (Optional) Outgoing broadcast traffic limit in packets per second.
  * `tx_multicast` - (Optional) Outgoing multicast traffic limit in packets per second.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the QoS switching profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing switch security switching profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_switch_security_switching_profile.switch_security_switching_profile UUID
```

The above would import switching profile named `switch_security_switching_profile` with the nsx id `UUID`
