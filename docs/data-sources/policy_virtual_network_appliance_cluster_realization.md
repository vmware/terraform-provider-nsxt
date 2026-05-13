---
subcategory: "Beta"
page_title: "NSXT: policy_virtual_network_appliance_cluster_realization"
description: Policy Virtual Network Appliance Cluster realization information.
---

# nsxt_policy_virtual_network_appliance_cluster_realization

This data source provides realization state information for a Policy Virtual Network Appliance Cluster resource on NSX Manager. It waits until realization is determined as either success or error. Use this data source when downstream configuration depends on the VNA cluster being fully realized (i.e., all VNA instances provisioned on their respective edge transport nodes).

This data source is available since NSX 9.1.1.

## Example Usage

```hcl
resource "nsxt_policy_virtual_network_appliance_cluster" "example" {
  display_name          = "example-vna-cluster"
  appliance_form_factor = "MEDIUM"
  service_type          = "VPC_SERVICES"

  member {
    edge_transport_node_path = data.nsxt_policy_edge_transport_node.edge1.path
  }

  advanced_configuration {
    overlay_transport_zone_path = data.nsxt_policy_transport_zone.overlay_tz.path
  }
}

data "nsxt_policy_virtual_network_appliance_cluster_realization" "wait" {
  path    = nsxt_policy_virtual_network_appliance_cluster.example.path
  timeout = 1800
}
```

## Argument Reference

* `path` - (Required) The policy path of the `nsxt_policy_virtual_network_appliance_cluster` resource.
* `delay` - (Optional) Delay in seconds before realization polling begins. Default is `1`.
* `timeout` - (Optional) Timeout in seconds for realization polling. Default is `1800`.

## Attributes Reference

* `state` - The current realization state of the cluster. Possible values are `SUCCESS`, `IN_PROGRESS`, `ERROR`, `DOWN`, `DEGRAGED`, `DISABLED`, `UNINITIALIZED`, `SANDBOXED_REALIZATION_PENDING`, and `UNKNOWN`.
* `vna_paths` - List of policy paths for the individual Virtual Network Appliance nodes deployed within the cluster. Use `tolist(data.<name>.vna_paths)[0]` to reference the first appliance node path (e.g. as `virtual_network_appliance_path` in `nsxt_policy_route_controller_interface`).
