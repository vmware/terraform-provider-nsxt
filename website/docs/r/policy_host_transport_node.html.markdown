---
subcategory: "Beta"
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
  description        = "Terraform-deployed host transport node"
  display_name       = "tf_host_transport_node"
  discovered_node_id = data.nsxt_discovered_node.dn.id

  standard_host_switch {
    host_switch_id      = "50 0b 31 a4 b8 af 35 df-40 56 b6 f9 aa d3 ee 12"
    host_switch_profile = [data.nsxt_policy_uplink_host_switch_profile.uplink_host_switch_profile.path]

    ip_assignment {
      assigned_by_dhcp = true
    }

    transport_zone_endpoint {
      transport_zone = data.nsxt_policy_transport_zone.overlay_transport_zone.path
    }

    uplink {
      uplink_name     = "uplink-1"
      vds_uplink_name = "uplink1"
    }

    uplink {
      uplink_name     = "uplink-2"
      vds_uplink_name = "uplink2"
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
* `site_path` - (Optional) The path of the site which the Host Transport Node belongs to. `path` field of the existing `nsxt_policy_site` can be used here. Defaults to default site path.
* `enforcement_point` - (Optional) The ID of enforcement point under given `site_path` to manage the Host Transport Node. Defaults to default enforcement point.
* `discovered_node_id` - (Required)  Discovered node id to create Host Transport Node. Specify discovered node id to create Host Transport Node for Discovered Node. This field is required during Host Transport Node create from vCenter server managing the ESXi type HostNode.
* `standard_host_switch` - (Required) Standard host switch specification.
  * `cpu_config` - (Optional) Enhanced Networking Stack enabled HostSwitch CPU configuration.
    * `num_lcores` - (Required) Number of Logical cpu cores (Lcores) to be placed on a specified NUMA node.
    * `numa_node_index` - (Required) Unique index of the Non Uniform Memory Access (NUMA) node.
  * `host_switch_id` - (Optional) The host switch id. This ID will be used to reference a host switch.
  * `host_switch_name` - (Optional) Host switch name. This name will be used to reference a host switch.
  * `host_switch_mode` - (Optional) Operational mode of a HostSwitch. Accepted values - 'STANDARD', 'ENS', 'ENS_INTERRUPT' or 'LEGACY'.
  * `host_switch_profile` - (Optional) Policy path of host switch profiles to be associated with this host switch.
  * `ip_assignment` - (Required) - Specification for IPs to be used with host switch virtual tunnel endpoints. Should contain exatly one of the below:
    * `assigned_by_dhcp` - (Optional) Enables DHCP assignment.
    * `static_ip` - (Optional) IP assignment specification for Static IP List.
        * `ip_addresses` - (Required) List of IPs for transport node host switch virtual tunnel endpoints.
        * `subnet_mask` - (Required) Subnet mask.
        * `default_gateway` - (Required) Gateway IP.
    * `static_ip_pool` - (Optional) Policy path of Static IP Pool used for IP assignment specification.
  * `is_migrate_pnics` - (Optional) Migrate any pnics which are in use.
  * `pnic` - (Optional) Physical NICs connected to the host switch.
    * `device_name` - (Required) Device name or key.
    * `uplink_name` - (Required) Uplink name for this Pnic.
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
            * `static_ip_pool` - (Optional) IP assignment specification for Static IP Pool.
        * `uplink` - (Optional) Uplink/LAG of VMware vSphere Distributed Switch connected to the HostSwitch.
            * `uplink_name` - (Required) Uplink name from UplinkHostSwitch profile.
            * `vds_lag_name` - (Optional) Link Aggregation Group (LAG) name of Virtual Distributed Switch.
            * `vds_uplink_name` - (Optional) Uplink name of VMware vSphere Distributed Switch (VDS).
    * `name` - (Required) Name of the transport node profile config option.
  * `transport_zone_endpoint` - (Optional) Transport zone endpoints
    * `transport_zone` - (Required) Unique ID identifying the transport zone for this endpoint.
    * `transport_zone_profiles` - (Optional) Identifiers of the transport zone profiles associated with this transport zone endpoint on this transport node.
  * `uplink` - (Optional) Uplink/LAG of VMware vSphere Distributed Switch connected to the HostSwitch.
    * `uplink_name` - (Required) Uplink name from UplinkHostSwitch profile.
    * `vds_lag_name` - (Optional) Link Aggregation Group (LAG) name of Virtual Distributed Switch.
    * `vds_uplink_name` - (Optional) Uplink name of VMware vSphere Distributed Switch (VDS).
  * `vmk_install_migration` - (Optional) The vmknic and logical switch mappings.
    * `destination_network` - (Required) The network id to which the ESX vmk interface will be migrated.
    * `device_name` - (Required) ESX vmk interface name.
* `remove_nsx_on_destroy` - (Optional) Upon deletion, uninstall NSX from Transport Node. Default is true.

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
