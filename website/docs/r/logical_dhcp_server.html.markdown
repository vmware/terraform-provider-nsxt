---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_logical_dhcp_server"
description: |-
  Provides a resource to configure logical DHCP server on NSX-T manager
---

# nsxt_logical_dhcp_server

Provides a resource to configure logical DHCP server on NSX-T manager

## Example Usage

```hcl
data "nsxt_edge_cluster" "edgecluster" {
  display_name = "edgecluster1"
}

resource "nsxt_dhcp_server_profile" "serverprofile" {
  edge_cluster_id = data.nsxt_edge_cluster.edgecluster.id
}

resource "nsxt_logical_dhcp_server" "logical_dhcp_server" {
  display_name     = "logical_dhcp_server"
  description      = "logical_dhcp_server provisioned by Terraform"
  dhcp_profile_id  = nsxt_dhcp_server_profile.PRF.id
  dhcp_server_ip   = "1.1.1.10/24"
  gateway_ip       = "1.1.1.20"
  domain_name      = "abc.com"
  dns_name_servers = ["5.5.5.5"]

  dhcp_option_121 {
    network  = "6.6.6.0/24"
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
* `dhcp_profile_id` - (Required) DHCP profile uuid.
* `dhcp_server_ip` - (Required) DHCP server IP in cidr format.
* `gateway_ip` - (Optional) Gateway IP.
* `domain_name` - (Optional) Domain name.
* `dns_name_servers` - (Optional) DNS IPs.
* `dhcp_option_121` - (Optional) DHCP classless static routes.
  * `network` - (Required) Destination in cidr format.
  * `next_hop` - (Required) IP address of next hop.
* `dhcp_generic_option` - (Optional) Generic DHCP options.
  * `code` - (Required) DHCP option code. Valid values are from 0 to 255.
  * `values` - (Required) List of DHCP option values.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical DHCP server.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical DHCP server.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `attached_logical_port_id` - ID of the attached logical port.


## Importing

An existing logical DHCP server can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_logical_dhcp_server.logical_dhcp_server UUID
```

The above would import the logical DHCP server named `logical_dhcp_server` with the nsx id `UUID`
