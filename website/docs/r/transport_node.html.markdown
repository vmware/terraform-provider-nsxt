---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_transport_node"
description: A resource to configure a Transport Node.
---

# nsxt_transport_node

This resource provides a method for the management of a Transport Node.
This resource is supported with NSX 4.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_transport_node" "test" {
  description  = "Terraform-deployed edge node"
  display_name = "tf_edge_node"
  standard_host_switch {
    host_switch_mode = "STANDARD"
    host_switch_type = "NVDS"
    ip_assignment {
      assigned_by_dhcp = true
    }
    transport_zone_endpoint {
      transport_zone_id         = data.nsxt_transport_zone.tz1.id
      transport_zone_profile_id = ["52035bb3-ab02-4a08-9884-18631312e50a"]
    }
    host_switch_profile_id = [nsxt_uplink_host_switch_profile.hsw_profile1.id]
    is_migrate_pnics       = false
    pnic {
      device_name = "fp-eth0"
      uplink_name = "uplink1"
    }
  }
  edge_node {
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
      hostname             = "tf_edge_node"
      allow_ssh_root_login = true
      enable_ssh           = true
    }
  }
}
```

**NOTE:** `data.vsphere_network`, `data.vsphere_compute_cluster`, `data.vsphere_datastore`, `data.vsphere_host` are 
obtained using [hashicorp/vsphere](https://registry.terraform.io/providers/hashicorp/vsphere/) provider.

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `failure_domain` - (Optional)  Id of the failure domain.
* `standard_host_switch` - (Optional) Standard host switch specification.
  * `cpu_config` - (Optional) Enhanced Networking Stack enabled HostSwitch CPU configuration.
    * `num_lcores` - (Required) Number of Logical cpu cores (Lcores) to be placed on a specified NUMA node.
    * `numa_node_index` - (Required) Unique index of the Non Uniform Memory Access (NUMA) node.
  * `host_switch_id` - (Optional) The host switch id. This ID will be used to reference a host switch.
  * `host_switch_mode` - (Optional) Operational mode of a HostSwitch. Accepted values - 'STANDARD', 'ENS', 'ENS_INTERRUPT' or 'LEGACY'. The default value is 'STANDARD'.
  * `host_switch_profile_id` - (Optional) Identifiers of host switch profiles to be associated with this host switch.
  * `host_switch_type` - (Optional) Type of HostSwitch. Accepted values - 'NVDS' or 'VDS'. The default value is 'NVDS'.
  * `ip_assignment` - (Required) - Specification for IPs to be used with host switch virtual tunnel endpoints. Should contain exatly one of the below:
    * `assigned_by_dhcp` - (Optional) Enables DHCP assignment. Should be set to true.
    * `static_ip` - (Optional) IP assignment specification for Static IP List.
      * `ip_addresses` - (Required) List of IPs for transport node host switch virtual tunnel endpoints.
      * `subnet_mask` - (Required) Subnet mask.
      * `default_gateway` - (Required) Gateway IP.
    * `static_ip_mac` - (Optional) IP and MAC assignment specification for Static IP List.
      * `default_gateway` - (Required) Gateway IP.
      * `subnet_mask` - (Required) Subnet mask.
      * `ip_mac_pair` - (Required) List of IPs and MACs for transport node host switch virtual tunnel endpoints.
        * `ip` - (Required) IP address.
        * `mac` - (Required) MAC address.
        * `static_ip_pool_id` - (Optional) IP assignment specification for Static IP Pool.
  * `is_migrate_pnics` - (Optional) Migrate any pnics which are in use.
  * `pnic` - (Optional) Physical NICs connected to the host switch.
    * `device_name` - (Required) Device name or key.
    * `uplink_name` - (Required) Uplink name for this Pnic.
  * `portgroup_transport_zone_id` - (Optional) Transport Zone ID representing the DVS used in NSX on DVPG.
  * `transport_node_profile_sub_config` - (Optional) Transport Node Profile sub-configuration Options.
    * `host_switch_config_option` - (Required) Subset of the host switch configuration.
      * `host_switch_id` - (Optional) The host switch id. This ID will be used to reference a host switch.
      * `host_switch_profile_id` - (Optional) Identifiers of host switch profiles to be associated with this host switch.
      * `ip_assignment` - (Required) - Specification for IPs to be used with host switch virtual tunnel endpoints. Should contain exatly one of the below:
        * `assigned_by_dhcp` - (Optional) Enables DHCP assignment. Should be set to true.
        * `static_ip` - (Optional) IP assignment specification for Static IP List.
          * `ip_addresses` - (Required) List of IPs for transport node host switch virtual tunnel endpoints.
          * `subnet_mask` - (Required) Subnet mask.
          * `default_gateway` - (Required) Gateway IP.
        * `static_ip_mac` - (Optional) IP and MAC assignment specification for Static IP List.
          * `default_gateway` - (Required) Gateway IP.
          * `subnet_mask` - (Required) Subnet mask.
          * `ip_mac_pair` - (Required) List of IPs and MACs for transport node host switch virtual tunnel endpoints.
            * `ip` - (Required) IP address.
            * `mac` - (Required) MAC address.
            * `static_ip_pool_id` - (Optional) IP assignment specification for Static IP Pool.
      * `uplink` - (Optional) Uplink/LAG of VMware vSphere Distributed Switch connected to the HostSwitch.
        * `uplink_name` - (Required) Uplink name from UplinkHostSwitch profile.
        * `vds_lag_name` - (Optional) Link Aggregation Group (LAG) name of Virtual Distributed Switch.
        * `vds_uplink_name` - (Optional) Uplink name of VMware vSphere Distributed Switch (VDS).
    * `name` - (Required) Name of the transport node profile config option.
  * `transport_zone_endpoint` - (Optional) Transport zone endpoints
    * `transport_zone_id` - (Required) Unique ID identifying the transport zone for this endpoint.
    * `transport_zone_profile_id` - (Optional) Identifiers of the transport zone profiles associated with this transport zone endpoint on this transport node.
  * `uplink` - (Optional) Uplink/LAG of VMware vSphere Distributed Switch connected to the HostSwitch.
    * `uplink_name` - (Required) Uplink name from UplinkHostSwitch profile.
    * `vds_lag_name` - (Optional) Link Aggregation Group (LAG) name of Virtual Distributed Switch.
    * `vds_uplink_name` - (Optional) Uplink name of VMware vSphere Distributed Switch (VDS).
  * `vmk_install_migration` - (Optional) The vmknic and logical switch mappings.
    * `destination_network` - (Required) The network id to which the ESX vmk interface will be migrated.
    * `device_name` - (Required) ESX vmk interface name.
* `preconfigured_host_switch` - (Optional) Preconfigured host switch.
  * `endpoint` - (Optional) Name of the virtual tunnel endpoint which is preconfigured on this host switch.
  * `host_switch_id` - (Required) External Id of the preconfigured host switch.
  * `transport_zone_endpoint` - (Optional) Transport zone endpoints
    * `transport_zone_id` - (Required) Unique ID identifying the transport zone for this endpoint.
    * `transport_zone_profile_id` - (Optional) Identifiers of the transport zone profiles associated with this transport zone endpoint on this transport node.
* `node` - (Optional)
  * `external_id` - (Optional) ID of the Node.
  * `fqdn` - (Optional) Fully qualified domain name of the fabric node.
  * `id` - (Optional) Unique identifier of this resource.
  * `ip_addresses` - (Optional) IP Addresses of the Node, version 4 or 6.
* `host_node` - (Optional)
  * `external_id` - (Optional) ID of the Node.
  * `fqdn` - (Optional) Fully qualified domain name of the fabric node.
  * `id` - (Optional) Unique identifier of this resource.
  * `ip_addresses` - (Optional) IP Addresses of the Node, version 4 or 6.
  * `host_credential` - (Optional) Host login credentials.
    * `password` - (Required) The authentication password of the host node.
    * `thumbprint` - (Optional) ESXi thumbprint or SSH key fingerprint of the host node.
    * `username` - (Required) The username of the account on the host node.
  * `os_type` - (Required) Hypervisor OS type. Accepted values - 'ESXI', 'RHELSERVER', 'WINDOWSSERVER', 'RHELCONTAINER', 'UBUNTUSERVER', 'HYPERV', 'CENTOSSERVER', 'CENTOSCONTAINER', 'SLESSERVER' or 'OELSERVER'.
  * `os_version` - (Optional) Hypervisor OS version.
  * `windows_install_location` - (Optional) Install location of Windows Server on baremetal being managed by NSX. Defaults to 'C:\Program Files\VMware\NSX\'.
* `edge_node` - (Optional)
  * `deployment_config` - (Optional) Config for automatic deployment of edge node virtual machine.
    * `form_factor` - (Optional) Accepted values - 'SMALL', 'MEDIUM', 'LARGE', 'XLARGE'. The default value is 'MEDIUM'.
    * `node_user_settings` - (Required) Node user settings.
      * `audit_password` - (Optional) Node audit user password.
      * `audit_username` - (Optional) CLI "audit" username.
      * `cli_password` - (Required) Node cli password.
      * `cli_username` - (Optional) CLI "admin" username. Defaults to "admin".
      * `root_password` - (Required) Node root user password.
    * `vm_deployment_config` - (Required) The vSphere deployment configuration determines where to deploy the edge node.
      * `compute_folder_id` - (Optional) Cluster identifier or resourcepool identifier for specified vcenter server.
      * `data_network_ids` - (Required) List of portgroups, logical switch identifiers or segment paths for datapath connectivity.
      * `default_gateway_address` - (Optional) Default gateway for the node.
      * `host_id` - (Optional) Host identifier in the specified vcenter server.
      * `ipv4_assignment_enabled` - (Optional) This flag represents whether IPv4 configuration is enabled or not. Defaults to true.
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
* `public_cloud_gateway_node` - (Optional) 
  * `deployment_config` - (Optional) Config for automatic deployment of edge node virtual machine.
    * `form_factor` - (Optional) Accepted values - 'SMALL', 'MEDIUM', 'LARGE', 'XLARGE'. The default value is 'MEDIUM'.
    * `node_user_settings` - (Required) Node user settings.
      * `audit_password` - (Optional) Node audit user password.
      * `audit_username` - (Optional) CLI "audit" username.
      * `cli_password` - (Required) Node cli password.
      * `cli_username` - (Optional) CLI "admin" username. Defaults to "admin".
      * `root_password` - (Required) Node root user password.
    * `vm_deployment_config` - (Required) The vSphere deployment configuration determines where to deploy the edge node.
      * `compute_folder_id` - (Optional) Cluster identifier or resourcepool identifier for specified vcenter server.
      * `data_network_ids` - (Required) List of portgroups, logical switch identifiers or segment paths for datapath connectivity.
      * `default_gateway_address` - (Optional) Default gateway for the node.
      * `host_id` - (Optional) Host identifier in the specified vcenter server.
      * `ipv4_assignment_enabled` - (Optional) This flag represents whether IPv4 configuration is enabled or not. Defaults to true.
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
* `remote_tunnel_endpoint` - (Optional) Configuration for a remote tunnel endpoint.
  * `host_switch_name` - (Required) The host switch name to be used for the remote tunnel endpoint.
  * `ip_assignment` - (Required) - Specification for IPs to be used with host switch virtual tunnel endpoints. Should contain exatly one of the below:
    * `assigned_by_dhcp` - (Optional) Enables DHCP assignment. Should be set to true.
    * `static_ip` - (Optional) IP assignment specification for Static IP List.
      * `ip_addresses` - (Required) List of IPs for transport node host switch virtual tunnel endpoints.
      * `subnet_mask` - (Required) Subnet mask.
      * `default_gateway` - (Required) Gateway IP.
    * `static_ip_mac` - (Optional) IP and MAC assignment specification for Static IP List.
      * `default_gateway` - (Required) Gateway IP.
      * `subnet_mask` - (Required) Subnet mask.
      * `ip_mac_pair` - (Required) List of IPs and MACs for transport node host switch virtual tunnel endpoints.
        * `ip` - (Required) IP address.
        * `mac` - (Required) MAC address.
        * `static_ip_pool_id` - (Optional) IP assignment specification for Static IP Pool.
  * `named_teaming_policy` - (Optional) The named teaming policy to be used by the remote tunnel endpoint.
  * `rtep_vlan` - (Required) VLAN id for remote tunnel endpoint.



**NOTE:** Resource should contain either `standard_host_switch` or `preconfigured_host_switch`
**NOTE:** Resource should contain one of: `node`, `edge_node`, `host_node` or `public_cloud_gateway_node`

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing Transport Node can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_transport_node.test UUID
```
The above command imports Transport Node named `test` with the NSX Transport Node ID `UUID`.
