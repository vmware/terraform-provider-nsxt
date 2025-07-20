---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_edge_cluster"
description: A resource to configure a PolicyEdgeCluster.
---

# nsxt_policy_edge_cluster

This resource provides a method for the management of a PolicyEdgeCluster.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.
## Example Usage

```hcl
resource "nsxt_policy_edge_cluster" "test" {
  display_name              = "test"
  description               = "Terraform provisioned PolicyEdgeCluster"
  edge_cluster_profile_path = nsxt_policy_edge_high_availability_profile.path
  password_managed_by_vcf   = true
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `site_path` - (Optional) The path of the site which the Host Transport Node belongs to. `path` field of the existing `nsxt_policy_site` can be used here. Defaults to default site path.
* `enforcement_point` - (Optional) The ID of enforcement point under given `site_path` to manage the Host Transport Node. Defaults to default enforcement point.
* `inter_site_forwarding_enabled` - (Optional) Inter site forwarding is enabled if true.
* `edge_cluster_profile_path` - (Required) Path of edge cluster high availability profile.
* `allocation_rule` - Allocation rules for auto placement.
    * `action_based_on_failure_domain_enabled` - (Optional) Auto place TIER1 logical routers, DHCP and MDProxy contexts on two edge nodes (active and standby) from different failure domains.
* `policy_edge_node` - (Optional) Policy Edge Cluster Member.
    * `edge_transport_node_path` - (Required) Edge Transport Node Path.
    * `id` - (Optional) ID of PolicyEdgeNode.
* `password_managed_by_vcf` - (Optional) Setting to true enables VCF password management for all edge nodes in the cluster. Default: false.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_edge_cluster.test PATH
```

The above command imports PolicyEdgeCluster named `test` with the NSX path `PATH`.
