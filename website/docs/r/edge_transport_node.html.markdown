---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_edge_transport_node"
description: A resource to configure an Edge Transport Node.
---

# nsxt_edge_transport_node

This resource provides a method for the management of an Edge Transport Node.
This resource is supported with NSX 4.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_edge_transport_node" "test" {
  description  = "Terraform-deployed edge node"
  display_name = "tf_edge_node"
  standard_host_switch {
    ip_assignment {
      assigned_by_dhcp = true
    }
    transport_zone_endpoint {
      transport_zone          = data.nsxt_policy_transport_zone.tz1.id
      transport_zone_profiles = ["52035bb3-ab02-4a08-9884-18631312e50a"]
    }
    host_switch_profile = [nsxt_policy_uplink_host_switch_profile.hsw_profile1.id]
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
```

**NOTE:** `data.vsphere_network`, `data.vsphere_compute_cluster`, `data.vsphere_datastore`, `data.vsphere_host` are 
obtained using [hashicorp/vsphere](https://registry.terraform.io/providers/hashicorp/vsphere/) provider.

**NOTE** while policy path values are acceptable for `transport_zone`,  `host_switch_profile` attributes, it is better to use id values to avoid state issues.

## Example Usage, with Edge Transport Node created externally and converted into a transport node using Terraform
```hcl
data "nsxt_transport_node" "test_node" {
  display_name = "tf_edge_node"
}

resource "nsxt_edge_transport_node" "test_node" {
  node_id      = data.nsxt_transport_node.test_node.id
  description  = "Terraform-deployed edge node"
  display_name = "tf_edge_node"
  standard_host_switch {
    ip_assignment {
      static_ip_pool = data.nsxt_policy_ip_pool.vtep_ip_pool.realized_id
    }
    transport_zone_endpoint {
      transport_zone = data.nsxt_policy_transport_zone.overlay_tz.id
    }
    transport_zone_endpoint {
      transport_zone = data.nsxt_policy_transport_zone.vlan_tz.id
    }
    host_switch_profile = [data.nsxt_policy_uplink_host_switch_profile.edge_uplink_profile.id]
    pnic {
      device_name = "fp-eth0"
      uplink_name = "uplink1"
    }
  }
  node_settings {
    hostname = "tf-edge-node"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `node_id` - (Optional) The id of a pre-deployed Edge appliance to be converted into a transport node. Note that `node_id` attribute conflicts with `external_id`, `fqdn`, `ip_addresses` `deployment_config` and `node_settings` and those will be ignored while specifying `node_id`.   
* `failure_domain` - (Optional)  Id of the failure domain.
* `standard_host_switch` - (Required) Standard host switch specification.
  * `host_switch_id` - (Optional) The host switch id. This ID will be used to reference a host switch.
  * `host_switch_name` - (Optional) Host switch name. This name will be used to reference a host switch.
  * `host_switch_profile` - (Optional) Identifiers of host switch profiles to be associated with this host switch.
  * `ip_assignment` - (Required) - Specification for IPs to be used with host switch virtual tunnel endpoints. Should contain exatly one of the below:
    * `assigned_by_dhcp` - (Optional) Enables DHCP assignment.
    * `static_ip` - (Optional) IP assignment specification for Static IP List.
      * `ip_addresses` - (Required) List of IPs for transport node host switch virtual tunnel endpoints.
      * `subnet_mask` - (Required) Subnet mask.
      * `default_gateway` - (Required) Gateway IP.
    * `static_ip_pool` - (Optional) IP assignment specification for Static IP Pool. Pool ID is expected here (if `nsxt_policy_ip_pool` object is referenced here, please use `realized_id` attribute)
  * `pnic` - (Optional) Physical NICs connected to the host switch.
    * `device_name` - (Required) Device name or key.
    * `uplink_name` - (Required) Uplink name for this Pnic.
  * `transport_zone_endpoint` - (Optional) Transport zone endpoints
    * `transport_zone` - (Required) Unique ID identifying the transport zone for this endpoint.
    * `transport_zone_profiles` - (Optional) Identifiers of the transport zone profiles associated with this transport zone endpoint on this transport node.
* `external_id` - (Optional) ID of the Node.
* `fqdn` - (Optional) Fully qualified domain name of the fabric node.
* `id` - (Optional) Unique identifier of this resource.
* `ip_addresses` - (Optional) IP Addresses of the Node, version 4 or 6.
* `deployment_config` - (Optional) Config for automatic deployment of edge node virtual machine.
  * `form_factor` - (Optional) Accepted values - 'SMALL', 'MEDIUM', 'LARGE', 'XLARGE'. The default value is 'MEDIUM'.
  * `node_user_settings` - (Required) Node user settings.
    * `audit_password` - (Optional) Node audit user password.
    * `audit_username` - (Optional) CLI "audit" username.
    * `cli_password` - (Required) Node cli password.
    * `cli_username` - (Optional) CLI "admin" username. Defaults to "admin".
    * `root_password` - (Required) Node root user password.
  * `vm_deployment_config` - (Required) The vSphere deployment configuration determines where to deploy the edge node.
    * `compute_folder_id` - (Optional) Compute folder identifier in the specified vcenter server.
    * `compute_id` - (Required) Cluster identifier or resourcepool identifier for specified vcenter server.
    * `data_network_ids` - (Required) List of portgroups, logical switch identifiers or segment paths for datapath connectivity.
    * `default_gateway_address` - (Optional) Default gateway for the node.
    * `host_id` - (Optional) Host identifier in the specified vcenter server.
    * `management_network_id` - (Required) Portgroup, logical switch identifier or segment path for management network connectivity.
    * `management_port_subnet` - (Optional) Port subnets for management port. IPv4, IPv6 and Dual Stack Address is supported.
      * `ip_addresses` - (Required) List of IP addresses.
      * `prefix_length` - (Required) Subnet Prefix Length.
    * `reservation_info` - (Optional) Resource reservation settings.
      * `cpu_reservation_in_mhz` - (Optional) CPU reservation in MHz.
      * `cpu_reservation_in_shares` - (Optional) CPU reservation in shares. Accepted values - 'EXTRA_HIGH_PRIORITY', 'HIGH_PRIORITY', 'NORMAL_PRIORITY', 'LOW_PRIORITY'. The default value is 'HIGH_PRIORITY'.
      * `memory_reservation_percentage` - (Optional) Memory reservation percentage.
    * `storage_id` - (Required) Storage/datastore identifier in the specified vcenter server
    * `vc_id` - (Required) Vsphere compute identifier for identifying the vcenter server.
* `node_settings` - (Required) Current configuration on edge node.
  * `advanced_configuration` - (Optional) Advanced configuration.
    * `key` - (Required)
    * `value` - (Required)
  * `allow_ssh_root_login` - (Optional) Allow root SSH logins. Defaults to false.
  * `dns_servers` - (Optional) List of DNS servers.
  * `enable_ssh` - (Optional) Enable SSH. Defaults to false.
  * `enable_upt_mode` - (Optional) Enable Uniform Passthrough mode. Defaults to false.
  * `hostname` - (Required) Host name or FQDN for edge node.
  * `ntp_servers` - (Optional) List of NTP servers.
  * `search_domains` - (Optional) List of Search domain names.
  * `syslog_server` - (Optional) List of Syslog servers.
    * `log_level` - (Optional) Log level to be redirected. Accepted values - 'EMERGENCY', 'ALERT', 'CRITICAL', 'ERROR', 'WARNING', 'NOTICE', 'INFO' or 'DEBUG'. The default value is 'INFO'.
    * `name` - (Optional) Display name of the syslog server.
    * `port` - (Optional) Syslog server port. Defaults to 514.
    * `protocol` - (Optional) Syslog protocol. Accepted values - 'TCP', 'UDP', 'TLS', 'LI', 'LI_TLS'. The default value is 'UDP'.
    * `server` - (Required) Server IP or fqdn.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing Edge Transport Node can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_edge_transport_node.test UUID
```
The above command imports Edge Transport Node named `test` with the NSX Transport Node ID `UUID`.
