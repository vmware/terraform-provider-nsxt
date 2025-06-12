---
subcategory: "Realization"
page_title: "NSXT: policy_edge_transport_node_realization"
description: Policy Edge transport node resource realization information.
---

# nsxt_policy_edge_transport_node_realization

This data source provides information about the realization of a Policy Edge transport node resource on NSX manager. This data source will wait until realization is determined as either success or error. It is recommended to use this data source if further configuration depends on Policy Edge transport node realization.
The data source will fail the plan when the Edge transport node deployment fails.

## Example Usage

```hcl
resource "nsxt_policy_edge_transport_node" "test" {
  display_name = "test-edge-appliance"
  description  = "Terraform deployed Edge appliance"

  appliance_config {
    allow_ssh_root_login = true
    enable_ssh           = true
  }
  credentials {
    cli_password  = var.CLI_PASSWORD
    root_password = var.ROOT_PASSWORD
  }
  form_factor = "SMALL"
  hostname    = "test-edge"
  management_interface {
    ip_assignment {
      dhcp_v4 = true
    }
    ip_assignment {
      static_ipv6 {
        default_gateway = "2001:0000:130F:0000:0000:09C0:876A:1"
        management_port_subnet {
          ip_addresses  = ["2001:0000:130F:0000:0000:09C0:876A:130B"]
          prefix_length = 64
        }
      }
    }
    network_id = data.vsphere_network.mgt_network.id
  }
  switch {
    overlay_transport_zone_path = data.nsxt_policy_transport_zone.overlay_tz.path
    pnic {
      device_name = "fp-eth0"
      uplink_name = "uplink1"
    }
    uplink_host_switch_profile_path = data.nsxt_policy_uplink_host_switch_profile.uplink_host_switch_profile.path
    tunnel_endpoint {
      ip_assignment {
        dhcp_v4 = true
      }
    }
    vlan_transport_zone_paths = [data.nsxt_policy_transport_zone.vlan_tz.path]
  }
  vm_deployment_config {
    compute_id         = data.vsphere_compute_cluster.edge_cluster.id
    storage_id         = data.vsphere_datastore.datastore.id
    compute_manager_id = data.nsxt_compute_manager.vc.id
    host_id            = data.vsphere_host.host.id
  }
}

# Wait for realization of the Edge transport node
data "nsxt_policy_edge_transport_node_realization" "info" {
  path    = nsxt_policy_edge_transport_node.test.path
  timeout = 1200
}
```

## Argument Reference

* `path` - (Required) The path for the Policy Edge transport node.
* `delay` - (Optional) Delay (in seconds) before realization polling is started. Default is set to 1.
* `timeout` - (Optional) Timeout (in seconds) for realization polling. Default is set to 1200.
