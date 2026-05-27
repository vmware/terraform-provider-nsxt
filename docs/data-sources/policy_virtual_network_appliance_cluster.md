---
subcategory: "Beta"
page_title: "NSXT: policy_virtual_network_appliance_cluster"
description: Policy Virtual Network Appliance Cluster data source.
---

# nsxt_policy_virtual_network_appliance_cluster

This data source provides information about a Virtual Network Appliance (VNA) Cluster configured on NSX.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.1.1 onwards.

## Example Usage

```hcl
data "nsxt_policy_virtual_network_appliance_cluster" "cluster" {
  display_name = "vna-cluster-1"
}
```

### Lookup by ID

```hcl
data "nsxt_policy_virtual_network_appliance_cluster" "cluster" {
  id = "abc123"
}
```

### Lookup under a specific site and enforcement point

```hcl
data "nsxt_policy_virtual_network_appliance_cluster" "cluster" {
  display_name      = "vna-cluster-1"
  site_path         = "/infra/sites/site1"
  enforcement_point = "default"
}
```

## Argument Reference

* `id` - (Optional) The ID of the Virtual Network Appliance Cluster to retrieve. If set, no additional argument is required.
* `display_name` - (Optional) The display name prefix of the Virtual Network Appliance Cluster to retrieve.
* `site_path` - (Optional) Path to the site this cluster belongs to. Defaults to the default site path.
* `enforcement_point` - (Optional) ID of the enforcement point under the given `site_path`. Defaults to `default`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX policy path of the resource.
* `appliance_form_factor` - The form factor of the virtual network appliances in the cluster (`SMALL`, `MEDIUM`, `LARGE`, `XLARGE`).
* `service_type` - The service type of the cluster (`VPC_SERVICES` or `ROUTE_CONTROLLER`).
