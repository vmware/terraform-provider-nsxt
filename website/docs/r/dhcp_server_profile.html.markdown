---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_dhcp_server_profile"
description: |-
  Provides a resource to configure DHCP server profile on NSX-T manager
---

# nsxt_dhcp_server_profile

Provides a resource to configure DHCP server profile on NSX-T manager

## Example Usage

```hcl
data "nsxt_edge_cluster" "edge_cluster1" {
  display_name = "edgecluster"
}

resource "nsxt_dhcp_server_profile" "dhcp_profile" {
  description                 = "dhcp_profile provisioned by Terraform"
  display_name                = "dhcp_profile"
  edge_cluster_id             = data.nsxt_edge_cluster.edge_cluster1.id
  edge_cluster_member_indexes = [0, 1]

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
* `edge_cluster_id` - (Required) Edge cluster uuid.
* `edge_cluster_member_indexes` - (Optional) Up to 2 edge nodes from the given cluster. If none is provided, the NSX will auto-select two edge-nodes from the given edge cluster. If user provides only one edge node, there will be no HA support.
* `tag` - (Optional) A list of scope + tag pairs to associate with this DHCP profile.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the DHCP server profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing DHCP profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_dhcp_server_profile.dhcp_profile UUID
```

The above would import the DHCP server profile named `dhcp_profile` with the nsx id `UUID`
