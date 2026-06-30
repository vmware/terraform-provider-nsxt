---
subcategory: "Beta"
page_title: "NSXT: policy_virtual_network_appliance_realization"
description: Policy Virtual Network Appliance realization information.
---

# nsxt_policy_virtual_network_appliance_realization

This data source provides realization state information for a Policy Virtual Network Appliance resource on NSX Manager. It polls the per-appliance state endpoint until realization is determined as either success or a terminal error state. Use this data source when downstream configuration depends on a specific VNA node being fully realized.

This data source is available since NSX 9.1.1.

## Example Usage

```hcl
resource "nsxt_policy_virtual_network_appliance" "example" {
  display_name = "example-vna"
  cluster_path = nsxt_policy_virtual_network_appliance_cluster.example.path
  hostname     = "example-vna.example.com"

  management_interface {
    network_id = var.portgroup_id

    ip_assignment {
      dhcp_v4 = true
    }
  }

  vm_deployment_config {
    compute_manager_id          = data.nsxt_compute_manager.example.id
    cluster_or_resource_pool_id = data.nsxt_compute_collection.example.cm_local_id
    datastore_id                = var.datastore_id
  }
}

data "nsxt_policy_virtual_network_appliance_realization" "wait" {
  path    = nsxt_policy_virtual_network_appliance.example.path
  timeout = 1800
}
```

## Argument Reference

* `path` - (Required) The policy path of the `nsxt_policy_virtual_network_appliance` resource.
* `delay` - (Optional) Delay in seconds before realization polling begins. Default is `1`.
* `timeout` - (Optional) Timeout in seconds for realization polling. Default is `1800`.

## Attributes Reference

* `state` - The current realization state of the appliance. Possible values are `SUCCESS`, `IN_PROGRESS`, `ERROR`, `DOWN`, `DEGRAGED`, `DISABLED`, `UNINITIALIZED`, `SANDBOXED_REALIZATION_PENDING`, and `UNKNOWN`.
