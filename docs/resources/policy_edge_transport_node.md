---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_edge_transport_node"
description: A resource to configure a Policy Edge Transport Node.
---

# nsxt_policy_edge_transport_node

This resource provides a method for the management of a Policy Edge Transport Node.
This resource is supported with NSX 9.0.0 onwards.

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
data "nsxt_policy_realization_info" "info" {
  path        = nsxt_policy_edge_transport_node.test.path
  timeout     = 1200
  entity_type = "RealizedEdgeTransportNode"
}

```

**NOTE:** `data.vsphere_network`, `data.vsphere_compute_cluster`, `data.vsphere_datastore`, `data.vsphere_host` are
obtained using [hashicorp/vsphere](https://registry.terraform.io/providers/hashicorp/vsphere/) provider.

## Example Usage, with Edge Transport Node created externally and converted into a transport node using Terraform

```hcl
data "nsxt_policy_edge_node" "node1" {
  display_name = "tf_edge_node"
}

resource "nsxt_policy_edge_transport_node" "test" {
  node_id = data.nsxt_policy_edge_node.node1.id

  hostname = "test-edge-12"
  switch {
    overlay_transport_zone_path = data.nsxt_policy_transport_zone.overlay_tz.path
    pnic {
      device_name = "fp-eth0"
      uplink_name = "uplink1"
    }
    uplink_host_switch_profile_path = data.nsxt_policy_uplink_host_switch_profile.uplink_host_switch_profile.path
    tunnel_endpoint {
      ip_assignment {
        static_ipv4_pool = data.nsxt_policy_ip_pool.ip_pool1.path
      }
    }
    vlan_transport_zone_paths = [data.nsxt_policy_transport_zone.vlan_tz.path]
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `site_path` - (Optional) The path of the site which the Edge Transport Node belongs to. `path` field of the existing `nsxt_policy_site` can be used here. Defaults to default site path.
* `enforcement_point` - (Optional) The ID of enforcement point under given `site_path` to manage the Edge Transport Node. Defaults to default enforcement point.
* `node_id` - (Optional) The id of a pre-deployed Edge appliance to be converted into a Policy Edge transport node.
* `advanced_configuration` - (Optional) Array of additional specific properties for advanced or cloud-specific deployments in key-value format.
  * `key` - (Required) Key.
  * `value` - (Required) Value.
* `appliance_config` - (Optional) Applicance configuration.
  * `allow_ssh_root_login` (Optional) Allow root SSH logins.
  * `dns_servers` - (Optional) DNS servers.
  * `enable_ssh` - (Optional) Enable SSH.
  * `enable_upt_mode` - (Optional) Enable Uniform Passthrough mode.
  * `syslog_server` - (Optional) Syslog servers.
    * `log_level` - (Optional) Log level to be redirected. Valid values are `DEBUG`, `INFO`, `NOTICE`, `WARNING`, `ERROR`, `CRITICAL`, `ALERT`, `EMERGENCY`. Default: `INFO`.
    * `port` - (Optional) Syslog server port. Default: 514.
    * `protocol` - (Optional) Syslog protocol. Valid values: `TCP`, `UDP`, `TLS`, `LI`, `LI_TLS`. Default: `UDP`.
    * `server` - (Required) Server IP or fqdn.
* `credentials` - (Optional) Username and password settings for the node.
  * `audit_password` - (Optional) Node audit user password.
  * `audit_username` - (Optional) CLI "audit" username.
  * `cli_password` - (Required) Node cli password.
  * `cli_username` - (Optional) CLI "admin" username.
  * `root_password` - (Required) Node root user password.
* `failure_domain_path` - (Optional) Path of the failure domain.
* `form_factor` - (Optional) Appliance form factor. Valid values - 'SMALL', 'MEDIUM', 'LARGE', 'XLARGE'. The default value is 'MEDIUM'.
* `hostname` - (Required) Host name or FQDN for edge node.
* `management_interface` - (Optional) Applicable For LCM managed Node and contains the management interface info.
  * `ip_assignment` - (Required) IPv4 or Ipv6 Port subnets for management port. Can be one IPv4 and one IPv6 assignment but not requires both. Two assignments are supported here, one for IPv4 and one for IPv6.
    * `dhcp_v4` - (Optional) Enable DHCP based IPv4 assignment.
    * `static_ipv4` - (Optional) IP assignment specification for a Static IP.
      * `default_gateway` - (Required) Default IPv4 gateway for the node.
      * `management_port_subnet` - (Required) IPv4 Port subnet for management port.
        * `ip_addresses` - (Required) IPv4 Addresses.
        * `prefix_length` - (Required) Subnet Prefix Length.
    * `static_ipv6` - (Optional) IP assignment specification for a Static IP.
      * `default_gateway` - (Optional) Default IPv6 gateway for the node.
      * `management_port_subnet` - (Optional) Ipv6 Port subnet for management port.
        * `ip_addresses` - (Required) IPv6 Addresses list.
        * `prefix_length` - (Required) Subnet Prefix Length.
  * `network_id` - (Required) Portgroup, logical switch identifier or segment path for management network connectivity.
* `switch` - (Required) Edge Transport Node switches configuration.
  * `overlay_transport_zone_path` - (Optional) An overlay TransportZone path that is associated with the specified edge TN switch.
  * `pnic` - (Required) Physical NIC specification.
    * `datapath_network_id` - (Optional) A portgroup, logical switch identifier or segment path for datapath connectivity.
    * `device_name` - (Required) Device name or key e.g. fp-eth0, fp-eth1 etc.
    * `uplink_name` - (Required) Uplink name for this Pnic. This name will be used to reference this Pnic in other configurations.
  * `uplink_host_switch_profile_path` - (Optional) Uplink Host Switch Profile Path.
  * `lldp_host_switch_profile_path` - (Optional) LLDP Host Switch Profile Path.
  * `switch_name` - (Optional) Edge Tn switch name. This name will be used to reference an edge TN switch. Default: `nsxDefaultHostSwitch`.
  * `tunnel_endpoint` - (Optional) Tunnel Endpoint.
    * `ip_assignment` - (Required) IPv4 or Ipv6 Port subnets for management port.
      * `dhcp_v4` - (Optional) Enable DHCP based IPv4 assignment.
      * `dhcp_v6` - (Optional) Enable DHCP based IPv6 assignment.
      * `static_ipv4_list` - (Optional) IP assignment specification value for Static IPv4 List.
        * `default_gateway` - (Required) Gateway IP.
        * `ip_addresses` - (Required) List of IPV4 addresses for edge transport node host switch virtual tunnel endpoints.
        * `subnet_mask` - (Required) Subnet mask.
      * `static_ipv4_pool` - (Optional) IP assignment specification for Static IPv4 Pool. Input can be MP ip pool UUID or policy path of IP pool.
      * `static_ipv6_list` - (Optional) IP assignment specification value for Static IPv4 List.
        * `default_gateway` - (Required) Gateway IP.
        * `ip_addresses` - (Required) List of IPv6 IPs for edge transport node host switch virtual tunnel endpoints.
        * `prefix_length` - (Required) Prefix Length.
      * `static_ipv6_pool` - (Optional) IP assignment specification for Static IPv6 Pool. Input can be MP ip pool UUID or policy path of IP pool.
    * `vlan` - (Optional) VLAN ID for tunnel endpoint.
  * `vlan_transport_zone_paths` - (Optional) List of Vlan TransportZone paths that are to be associated with specified edge TN switch.
* `vm_deployment_config` - (Optional) VM deployment configuration for LCM nodes.
  * `compute_folder_id` - (Optional) Compute folder identifier in the specified vcenter server.
  * `compute_id` - (Required) Cluster identifier for specified vcenter server.
  * `edge_host_affinity_config` - (Optional) Edge VM to host affinity configuration.
    * `host_group_name` - (Required) Host group name.
  * `host_id` - (Optional) Host identifier in the specified vcenter server.
  * `reservation_info` - (Optional) Resource reservation settings.
    * `cpu_reservation_in_mhz` - (Optional) CPU reservation in MHz.
    * `cpu_reservation_in_shares` - (Optional) CPU reservation in shares. Accepted values - 'EXTRA_HIGH_PRIORITY', 'HIGH_PRIORITY', 'NORMAL_PRIORITY', 'LOW_PRIORITY'. The default value is 'HIGH_PRIORITY'.
    * `memory_reservation_percentage` - (Optional) Memory reservation percentage.
  * `storage_id` - (Required) Storage/datastore identifier in the specified vcenter server.
  * `compute_manager_id` - (Required) Vsphere compute identifier for identifying the vcenter server.

~> **NOTE:** `credentials` block and `allow_ssh_root_login` attribute cannot be modified after the resource is initially applied. To update credentials, use the Edge appliance CLI.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing Transport Node can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_edge_transport_node.test POLICY_PATH
```

The above command imports Edge Transport Node named `test` with the policy path `POLICY_PATH`.
