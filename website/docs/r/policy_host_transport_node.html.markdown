---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_host_transport_node"
description: A resource to configure a HostTransportNode.
---

# nsxt_policy_host_transport_node

This resource provides a method for the management of a HostTransportNode.
This resource is supported with NSX 4.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_policy_host_transport_node" "test" {
  description       = "Terraform-deployed host transport node"
  display_name      = "tf_host_transport_node"
  site_path         = "/infra/sites/default"
  enforcement_point = "default"

  node_deployment_info {
    ip_addresses = ["10.168.186.150"]

    os_type    = "ESXI"
    os_version = "7.0.3"

    host_credential {
      username   = "user1"
      password   = "password1"
      thumbprint = "thumbprint1"
    }
  }

  standard_host_switch {
    host_switch_mode    = "STANDARD"
    host_switch_type    = "NVDS"
    host_switch_profile = [data.nsxt_policy_uplink_host_switch_profile.hsw_profile1.path]

    ip_assignment {
      assigned_by_dhcp = true
    }

    transport_zone_endpoint {
      transport_zone = data.nsxt_transport_zone.tz1.path
    }

    pnic {
      device_name = "fp-eth0"
      uplink_name = "uplink1"
    }
  }

  tag {
    scope = "app"
    tag   = "web"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `site_path` - (Optional) The path of the site which the Host Transport Node belongs to. `path` field of the existing `nsxt_policy_site` can be used here.
* `enforcement_point` - (Optional) The ID of enforcement point under given `site_path` to manage the Host Transport Node.
* `discovered_node_id` - (Optional)  Discovered node id to create Host Transport Node. Specify discovered node id to create Host Transport Node for Discovered Node. This field is required during Host Transport Node create from vCenter server managing the ESXi type HostNode.
* `node_deployment_info` - (Optional)
  * `fqdn` - (Optional) Fully qualified domain name of the fabric node.
  * `ip_addresses` - (Required) IP Addresses of the Node, version 4 or 6.
  * `host_credential` - (Optional) Host login credentials.
    * `password` - (Required) The authentication password of the host node.
    * `thumbprint` - (Optional) ESXi thumbprint or SSH key fingerprint of the host node.
    * `username` - (Required) The username of the account on the host node.
  * `os_type` - (Required) Hypervisor OS type. Accepted values - 'ESXI', 'RHELSERVER', 'WINDOWSSERVER', 'RHELCONTAINER', 'UBUNTUSERVER', 'HYPERV', 'CENTOSSERVER', 'CENTOSCONTAINER', 'SLESSERVER' or 'OELSERVER'.
  * `os_version` - (Optional) Hypervisor OS version.
  * `windows_install_location` - (Optional) Install location of Windows Server on baremetal being managed by NSX. Defaults to 'C:\Program Files\VMware\NSX\'.
* `standard_host_switch` - (Optional) Standard host switch specification.
  * `cpu_config` - (Optional) Enhanced Networking Stack enabled HostSwitch CPU configuration.
    * `num_lcores` - (Required) Number of Logical cpu cores (Lcores) to be placed on a specified NUMA node.
    * `numa_node_index` - (Required) Unique index of the Non Uniform Memory Access (NUMA) node.
  * `host_switch_id` - (Optional) The host switch id. This ID will be used to reference a host switch.
  * `host_switch_mode` - (Optional) Operational mode of a HostSwitch. Accepted values - 'STANDARD', 'ENS', 'ENS_INTERRUPT' or 'LEGACY'. The default value is 'STANDARD'.
  * `host_switch_profile` - (Optional) Policy path of host switch profiles to be associated with this host switch.
  * `host_switch_type` - (Optional) Type of HostSwitch. Accepted values - 'NVDS' or 'VDS'. The default value is 'NVDS'.
  * `ip_assignment` - (Required) - Specification for IPs to be used with host switch virtual tunnel endpoints. Should contain exatly one of the below:
    * `assigned_by_dhcp` - (Optional) Enables DHCP assignment.
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
        * `static_ip_pool` - (Optional) Policy path of Static IP Pool used for IP assignment specification.
  * `is_migrate_pnics` - (Optional) Migrate any pnics which are in use.
  * `pnic` - (Optional) Physical NICs connected to the host switch.
    * `device_name` - (Required) Device name or key.
    * `uplink_name` - (Required) Uplink name for this Pnic.
  * `portgroup_transport_zone` - (Optional) Transport Zone ID representing the DVS used in NSX on DVPG.
  * `transport_node_profile_sub_config` - (Optional) Transport Node Profile sub-configuration Options.
    * `host_switch_config_option` - (Required) Subset of the host switch configuration.
      * `host_switch_id` - (Optional) The host switch id. This ID will be used to reference a host switch.
      * `host_switch_profile` - (Optional) Identifiers of host switch profiles to be associated with this host switch.
      * `ip_assignment` - (Required) - Specification for IPs to be used with host switch virtual tunnel endpoints. Should contain exatly one of the below:
        * `assigned_by_dhcp` - (Optional) Enables DHCP assignment.
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
            * `static_ip_pool` - (Optional) IP assignment specification for Static IP Pool.
      * `uplink` - (Optional) Uplink/LAG of VMware vSphere Distributed Switch connected to the HostSwitch.
        * `uplink_name` - (Required) Uplink name from UplinkHostSwitch profile.
        * `vds_lag_name` - (Optional) Link Aggregation Group (LAG) name of Virtual Distributed Switch.
        * `vds_uplink_name` - (Optional) Uplink name of VMware vSphere Distributed Switch (VDS).
    * `name` - (Required) Name of the transport node profile config option.
  * `transport_zone_endpoint` - (Optional) Transport zone endpoints
    * `transport_zone` - (Required) Unique ID identifying the transport zone for this endpoint.
    * `transport_zone_profile` - (Optional) Identifiers of the transport zone profiles associated with this transport zone endpoint on this transport node.
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
    * `transport_zone` - (Required) Unique ID identifying the transport zone for this endpoint.
    * `transport_zone_profile` - (Optional) Identifiers of the transport zone profiles associated with this transport zone endpoint on this transport node.

**NOTE:** Resource should contain either `standard_host_switch` or `preconfigured_host_switch`

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing Transport Node can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_host_transport_node.test POLICY_PATH
```
The above command imports Host Transport Node named `test` with the policy path `POLICY_PATH`.
