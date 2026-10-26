---
subcategory: "Fabric"
page_title: "NSXT: nsxt_policy_virtual_network_service_cluster"
description: A resource to configure a Virtual Network Service Cluster.
---

# nsxt_policy_virtual_network_service_cluster

This resource provides a method for the management of a Virtual Network Service (VNS) Cluster under an NSX enforcement point.

A VNS Cluster is a logical group of Virtual Network Service nodes that share the same form factor, BFD profile and vSphere deployment configuration.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_compute_manager" "vc" {
  display_name = "vc1"
}

data "nsxt_compute_collection" "cluster" {
  display_name = "compute-cluster-1"
}

resource "nsxt_policy_virtual_network_service_cluster" "example" {
  display_name                       = "vns-cluster-1"
  description                        = "VNS cluster for VPC services"
  node_form_factor                   = "MEDIUM"
  number_of_nodes                    = 2
  default_project_service_allocation = true

  vsphere_deployment_config {
    compute_manager_id = data.nsxt_compute_manager.vc.id
    cluster_ids        = [data.nsxt_compute_collection.cluster.cm_local_id]
  }

  tag {
    scope = "environment"
    tag   = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `site_path` - (Optional) The path of the site this cluster belongs to. `path` field of the existing `nsxt_policy_site` can be used here. Defaults to default site path. Changing this forces a new resource.
* `enforcement_point` - (Optional) The ID of enforcement point under given `site_path` to manage the cluster. Defaults to `default`. Changing this forces a new resource.
* `node_form_factor` - (Optional, Computed) The form factor of the Virtual Network Service nodes in the cluster. Supported values: `SMALL`, `MEDIUM`, `LARGE`.
* `number_of_nodes` - (Optional, Computed) Number of Virtual Network Service nodes to deploy in the cluster. Must be at least 1.
* `bfd_profile_path` - (Optional, Computed) Policy path of the BFD profile applied to all Service Hosts in this cluster. If not set, NSX applies a default profile.
* `default_project_service_allocation` - (Optional, Computed) When set to `true`, the cluster is available for default project service allocation.
* `vsphere_deployment_config` - (Required) Exactly one vSphere deployment configuration must be supplied, matching the NSX API requirement. Keep this block in configuration for existing and imported clusters. Each block supports:
    * `compute_manager_id` - (Required) The NSX assigned identifier of the vCenter server registered with NSX (for example `id` of the `nsxt_compute_manager` data source).
    * `cluster_ids` - (Required) Set of vCenter clusters on which the Virtual Network Service nodes are deployed. All hosts in the specified clusters must have NSX fabric prepared.
    * `fault_domain_ids` - (Optional) Set of infrastructure fault domains or fault domain groups restricting the deployment within the clusters.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `path` - The NSX policy path of the resource.
* `revision` - Indicates current revision number of the object as seen by the NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing Virtual Network Service Cluster can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_virtual_network_service_cluster.example POLICY_PATH
```

The above command imports the Virtual Network Service Cluster named `example` with the NSX policy path `POLICY_PATH`.

Note: The policy path format is `/infra/sites/<site-id>/enforcement-points/<ep-id>/virtual-network-service-clusters/<cluster-id>`.
