---
subcategory: "Segments"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_segment_security_profile"
description: A resource to configure a Segment Security Profile.
---

# nsxt_policy_segment_security_profile

This resource provides a method for the management of a Segment Security Profile.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_segment_security_profile" "test" {
  display_name                 = "test"
  description                  = "Terraform provisioned Segment Security Profile"
  bpdu_filter_allow            = ["01:80:c2:00:00:05"]
  bpdu_filter_enable           = true
  dhcp_client_block_enabled    = true
  dhcp_client_block_v6_enabled = true
  dhcp_server_block_enabled    = false
  dhcp_server_block_v6_enabled = true
  non_ip_traffic_block_enabled = true
  ra_guard_enabled             = true
  rate_limits_enabled          = true

  rate_limit {
    rx_broadcast = 1800
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_segment_security_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name                 = "test"
  description                  = "Terraform provisioned Segment Security Profile"
  bpdu_filter_allow            = ["01:80:c2:00:00:05"]
  bpdu_filter_enable           = true
  dhcp_client_block_enabled    = true
  dhcp_client_block_v6_enabled = true
  dhcp_server_block_enabled    = false
  dhcp_server_block_v6_enabled = true
  non_ip_traffic_block_enabled = true
  ra_guard_enabled             = true
  rate_limits_enabled          = true

  rate_limit {
    rx_broadcast = 1800
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `bpdu_filter_allow` - (Optional) List of allowed MAC addresses to be excluded from BPDU filtering. List of allowed MACs - `01:80:c2:00:00:00`, `01:80:c2:00:00:01`, `01:80:c2:00:00:02`, `01:80:c2:00:00:03`, `01:80:c2:00:00:04`, `01:80:c2:00:00:05`, `01:80:c2:00:00:06`, `01:80:c2:00:00:07`, `01:80:c2:00:00:08`, `01:80:c2:00:00:09`, `01:80:c2:00:00:0a`, `01:80:c2:00:00:0b`, `01:80:c2:00:00:0c`, `01:80:c2:00:00:0d`, `01:80:c2:00:00:0e`, `01:80:c2:00:00:0f`, `00:e0:2b:00:00:00`, `00:e0:2b:00:00:04`, `00:e0:2b:00:00:06`, `01:00:0c:00:00:00`, `01:00:0c:cc:cc:cc`, `01:00:0c:cc:cc:cd`, `01:00:0c:cd:cd:cd`, `01:00:0c:cc:cc:c0`, `01:00:0c:cc:cc:c1`, `01:00:0c:cc:cc:c2`, `01:00:0c:cc:cc:c3`, `01:00:0c:cc:cc:c4`, `01:00:0c:cc:cc:c5`, `01:00:0c:cc:cc:c6`, `01:00:0c:cc:cc:c7`.
* `bpdu_filter_enable` - (Optional) Indicates whether BPDU filter is enabled. Default is `True`.
* `dhcp_client_block_enabled` - (Optional) Filters DHCP server and/or client traffic. Default is `False`.
* `dhcp_client_block_v6_enabled` - (Optional) Filters DHCP server and/or client IPv6 traffic. Default is `False`.
* `dhcp_server_block_enabled` - (Optional) Filters DHCP server and/or client traffic. Default is `True`.
* `dhcp_server_block_v6_enabled` - (Optional) Filters DHCP server and/or client IPv6 traffic. Default is `True`.
* `non_ip_traffic_block_enabled` - (Optional) A flag to block all traffic except IP/(G)ARP/BPDU. Default is `False`.
* `ra_guard_enabled` - (Optional) Enable or disable Router Advertisement Guard. Default is `False`
* `rate_limit` - (Optional) Rate limits.
  * `rx_broadcast` - (Optional) Incoming broadcast traffic limit in packets per second.
  * `rx_multicast` - (Optional) Incoming multicast traffic limit in packets per second.
  * `tx_broadcast` - (Optional) Outgoing broadcast traffic limit in packets per second.
  * `tx_multicast` - (Optional) Outgoing multicast traffic limit in packets per second.
* `rate_limits_enabled` - (Optional) Enable or disable Rate Limits. Default is `False`.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_segment_security_profile.test UUID
```

The above command imports Segment Security Profile named `test` with the NSX ID `UUID`.

```
terraform import nsxt_policy_segment_security_profile.test POLICY_PATH
```

The above command imports Segment Security Profile named `test` with the policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
