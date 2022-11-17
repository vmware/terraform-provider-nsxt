---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_switch_security_switching_profile"
description: |-
  Provides a resource to configure switch security switching profile on NSX-T manager
---

# nsxt_switch_security_switching_profile

Provides a resource to configure switch security switching profile on NSX-T manager

## Example Usage

```hcl
resource "nsxt_switch_security_switching_profile" "switch_security_switching_profile" {
  description           = "switch_security_switching_profile provisioned by Terraform"
  display_name          = "switch_security_switching_profile"
  block_non_ip          = true
  block_client_dhcp     = false
  block_server_dhcp     = false
  bpdu_filter_enabled   = true
  bpdu_filter_whitelist = ["01:80:c2:00:00:01"]

  rate_limits {
    enabled      = true
    rx_broadcast = 32
    rx_multicast = 32
    tx_broadcast = 32
    tx_multicast = 32
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

* `id` - ID of the switch security switching profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing switch security switching profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_switch_security_switching_profile.switch_security_switching_profile UUID
```

The above would import switching profile named `switch_security_switching_profile` with the nsx id `UUID`
