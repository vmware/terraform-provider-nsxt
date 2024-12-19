---
subcategory: "Realization"
layout: "nsxt"
page_title: "NSXT: transport_node_realization"
description: Transport node resource realization information.
---

# nsxt_transport_node_realization

This data source provides information about the realization of a transport node resource on NSX manager. This data source will wait until realization is determined as either success or error. It is recommended to use this data source if further configuration depends on transport node realization.

## Example Usage

```hcl
resource "nsxt_transport_node" "test" {
  description  = "Terraform-deployed edge node"
  display_name = "tf_edge_node"
  standard_host_switch {
    ip_assignment {
      static_ip_pool = data.nsxt_ip_pool.ipp1.id
    }
    transport_zone_endpoint {
      transport_zone          = data.nsxt_transport_zone.tz1.id
      transport_zone_profiles = ["52035bb3-ab02-4a08-9884-18631312e50a"]
    }
    host_switch_profile = [nsxt_policy_uplink_host_switch_profile.hsw_profile1.path]
    pnic {
      device_name = "fp-eth0"
      uplink_name = "uplink1"
    }
  }
  deployment_config {
    form_factor = "SMALL"
    node_user_settings {
      cli_password  = "some_cli_password"
      root_password = "some_other_password"
    }
    vm_deployment_config {
      management_network_id = data.vsphere_network.network1.id
      data_network_ids      = [data.vsphere_network.network1.id]
      compute_id            = data.vsphere_compute_cluster.compute_cluster1.id
      storage_id            = data.vsphere_datastore.datastore1.id
      vc_id                 = nsxt_compute_manager.vc1.id
      host_id               = data.vsphere_host.host1.id
    }
  }
  node_settings {
    hostname             = "tf-edge-node"
    allow_ssh_root_login = true
    enable_ssh           = true
  }
}

data "nsxt_transport_node_realization" "test" {
  id      = nsxt_transport_node.test.id
  timeout = 60
}
```

## Argument Reference

* `id` - (Required) ID of the resource.
* `delay` - (Optional) Delay (in seconds) before realization polling is started. Default is set to 1.
* `timeout` - (Optional) Timeout (in seconds) for realization polling. Default is set to 1200.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - The realization state of the resource. Transitional states are: "pending", "in_progress", "in_sync", "unknown". Target states are: "success", "failed", "partial_success", "orphaned", "error".
