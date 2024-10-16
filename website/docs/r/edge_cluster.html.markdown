---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_edge_cluster"
description: A resource to configure an Edge Cluster.
---

# nsxt_edge_cluster

This resource provides a method for the management of an Edge Cluster.
This resource is supported with NSX 4.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_edge_cluster" "test" {
  description  = "Terraform provisioned Edge Cluster"
  display_name = "test"
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `edge_ha_profile_id` - (Optional) Edge high availability cluster profile ID.
* `failure_domain_allocation` - (Optional) Flag to enable failure domain based allocation. Enable placement algorithm to consider failure domain of edge transport nodes and place active and standby contexts in different failure domains. Supported values are `enable` and `disable`. 
* `member` - (Optional) Edge cluster members
  * `description` - (Optional) Description of this resource.
  * `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
  * `member_index` - (Optional) System generated index for cluster member.
  * `transport_node_id` - (Optional) UUID of edge transport node.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `deployment_type` - Edge cluster deployment type.
* `member_node_type` - Node type of the cluster members.
* `node_rtep_ips` - Remote tunnel endpoint ip address.
  * `member_index` - System generated index for cluster member
  * `rtep_ips` - Remote tunnel endpoint ip address
  * `transport_node_id` - UUID of edge transport node

## Importing

An existing Edge Cluster can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_edge_cluster.test UUID
```
The above command imports Edge Cluster named `test` with the NSX Edge Cluster ID `UUID`.
