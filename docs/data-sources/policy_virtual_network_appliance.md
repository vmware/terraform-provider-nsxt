---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_virtual_network_appliance"
description: A data source for a Virtual Network Appliance.
---

# nsxt_policy_virtual_network_appliance

Use this data source to get the ID and attributes of an existing Virtual Network Appliance (VNA) within a VNA Cluster.

This data source is supported with NSX 9.1.1 onwards.

## Example Usage

```hcl
data "nsxt_policy_virtual_network_appliance_cluster" "cluster" {
  display_name = "vna-cluster-1"
}

data "nsxt_policy_virtual_network_appliance" "vna" {
  display_name = "vna-1"
  cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.cluster.path
}
```

## Argument Reference

* `cluster_path` - (Required) Policy path of the parent `nsxt_policy_virtual_network_appliance_cluster`.
* `id` - (Optional) The NSX ID of the Virtual Network Appliance. Conflicts with `display_name`.
* `display_name` - (Optional) Display name of the Virtual Network Appliance. Prefix matching is supported when no exact match is found.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - Description of the Virtual Network Appliance.
* `path` - Policy path of the Virtual Network Appliance.
* `hostname` - Hostname or FQDN of the Virtual Network Appliance VM.
* `failure_domain_path` - Policy path of the failure domain.
