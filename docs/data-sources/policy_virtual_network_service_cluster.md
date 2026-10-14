---
subcategory: "Fabric"
page_title: "NSXT: policy_virtual_network_service_cluster"
description: Policy Virtual Network Service Cluster data source.
---

# nsxt_policy_virtual_network_service_cluster

This data source provides information about a Virtual Network Service (VNS) Cluster configured on NSX.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_virtual_network_service_cluster" "cluster" {
  display_name = "vns-cluster-1"
}
```

### Lookup by ID

```hcl
data "nsxt_policy_virtual_network_service_cluster" "cluster" {
  id = "abc123"
}
```

### Lookup under a specific site and enforcement point

```hcl
data "nsxt_policy_virtual_network_service_cluster" "cluster" {
  display_name      = "vns-cluster-1"
  site_path         = "/infra/sites/site1"
  enforcement_point = "default"
}
```

## Argument Reference

* `id` - (Optional) The ID of the Virtual Network Service Cluster to retrieve. If set, no additional argument is required.
* `display_name` - (Optional) The display name of the Virtual Network Service Cluster to retrieve. An exact match is preferred; otherwise a unique name prefix match is used.
* `site_path` - (Optional) Path to the site this cluster belongs to. Defaults to the default site path.
* `enforcement_point` - (Optional) ID of the enforcement point under the given `site_path`. Defaults to `default`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX policy path of the resource.
* `node_form_factor` - The form factor of the Virtual Network Service nodes (`SMALL`, `MEDIUM` or `LARGE`).
* `number_of_nodes` - Number of Virtual Network Service nodes in the cluster.
* `bfd_profile_path` - Policy path of the BFD profile applied to the cluster.
* `default_project_service_allocation` - Whether the cluster is available for default project service allocation.
* `vsphere_deployment_config` - vSphere deployment configurations of the cluster. Each entry contains:
    * `compute_manager_id` - The NSX identifier of the vCenter server.
    * `cluster_ids` - Set of vCenter clusters hosting the nodes.
    * `fault_domain_ids` - Set of fault domains or fault domain groups restricting the deployment.
