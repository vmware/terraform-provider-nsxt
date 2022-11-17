---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_dhcp_server_ip_pool"
description: |-
  Provides a resource to configure IP Pool for logical DHCP server on NSX-T manager
---

# nsxt_dhcp_server_ip_pool

Provides a resource to configure IP Pool for logical DHCP server on NSX-T manager

## Example Usage

```hcl
data "nsxt_edge_cluster" "edgecluster" {
  display_name = "edgecluster1"
}

resource "nsxt_dhcp_server_profile" "serverprofile" {
  edge_cluster_id = data.nsxt_edge_cluster.edgecluster.id
}

resource "nsxt_logical_dhcp_server" "logical_dhcp_server" {
  display_name    = "logical_dhcp_server"
  dhcp_profile_id = nsxt_dhcp_server_profile.PRF.id
  dhcp_server_ip  = "1.1.1.10/24"
  gateway_ip      = "1.1.1.20"
}

resource "nsxt_dhcp_server_ip_pool" "dhcp_ip_pool" {
  display_name           = "ip pool"
  description            = "ip pool"
  logical_dhcp_server_id = nsxt_logical_dhcp_server.logical_dhcp_server.id
  gateway_ip             = "1.1.1.21"
  lease_time             = 1296000
  error_threshold        = 98
  warning_threshold      = 70

  ip_range {
    start = "1.1.1.40"
    end   = "1.1.1.60"
  }

  dhcp_option_121 {
    network  = "5.5.5.0/24"
    next_hop = "1.1.1.21"
  }

  dhcp_generic_option {
    code   = "119"
    values = ["abc"]
  }

  tag {
    scope = "color"
    tag   = "red"
  }

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `description` - (Optional) Description of this resource.
* `logical_dhcp_server_id` - (Required) DHCP server uuid. Changing this would force new pool to be created.
* `gateway_ip` - (Optional) Gateway IP.
* `ip_range` - (Required) IP Ranges to be used within this pool.
  * `start` - (Required) IP address that indicates range start.
  * `end` - (Required) IP address that indicates range end.
* `lease_time` - (Optional) Lease time in seconds. Default is 86400.
* `error_threshold` - (Optional) Error threshold in percent. Valid values are from 80 to 100, default is 100.
* `warning_threshold` - (Optional) Warning threshold in percent. Valid values are from 50 to 80, default is 80.
* `dhcp_option_121` - (Optional) DHCP classless static routes. If specified, overrides DHCP server settings.
  * `network` - (Required) Destination in cidr format.
  * `next_hop` - (Required) IP address of next hop.
* `dhcp_generic_option` - (Optional) Generic DHCP options. If specified, overrides DHCP server settings.
  * `code` - (Required) DHCP option code. Valid values are from 0 to 255.
  * `values` - (Required) List of DHCP option values.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical DHCP server.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the DHCP server IP pool.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing DHCP server IP Pool can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_dhcp_server_ip_pool.ip_pool DHCP_SERVER_UUID POOL_UUID
```

The above would import the IP pool named `ip pool` for dhcp server with nsx ID `DHCP_SERVER_UUID` and pool nsx id `POOL_UUID`
