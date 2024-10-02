---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_host_transport_node_profile"
description: A resource to configure a Policy Host Transport Node Profile.
---

# nsxt_policy_host_transport_node_profile

This resource provides a method for the management of a Policy Host Transport Node Profile.
This resource is supported with NSX 4.1.0 onwards.

## Example Usage
```hcl
resource "nsxt_policy_host_transport_node_profile" "test" {
  display_name = "test_policy_host_tnp"

  standard_host_switch {
    ip_assignment {
      assigned_by_dhcp = true
    }
    transport_zone_endpoint {
      transport_zone = data.nsxt_policy_transport_zone.tz1.path
    }
    host_switch_profile = [nsxt_policy_uplink_host_switch_profile.hsw_profile1.path]
    is_migrate_pnics    = false
    pnic {
      device_name = "fp-eth0"
      uplink_name = "uplink1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `ignore_overridden_hosts` - (Optional) Determines if cluster-level configuration should be applied on overridden hosts
* `standard_host_switch` - (Required) Standard host switch specification.
    * `host_switch_id` - (Optional) The host switch id. This ID will be used to reference a host switch.
    * `host_switch_name` - (Optional) Host switch name. This name will be used to reference a host switch.
    * `host_switch_mode` - (Optional) Operational mode of a HostSwitch. Accepted values - 'STANDARD', 'ENS', 'ENS_INTERRUPT' or 'LEGACY'.
    * `host_switch_profile` - (Optional) Policy paths of host switch profiles to be associated with this host switch.
    * `ip_assignment` - (Required) - Specification for IPs to be used with host switch virtual tunnel endpoints. Should contain exactly one of the below:
        * `assigned_by_dhcp` - (Optional) Enables DHCP assignment.
        * `static_ip` - (Optional) IP assignment specification for Static IP List.
            * `ip_addresses` - (Required) List of IPs for transport node host switch virtual tunnel endpoints.
            * `subnet_mask` - (Required) Subnet mask.
            * `default_gateway` - (Required) Gateway IP.
        * `static_ip_pool` - (Optional) IP assignment specification for Static IP Pool.
    * `is_migrate_pnics` - (Optional) Migrate any pnics which are in use.
    * `pnic` - (Optional) Physical NICs connected to the host switch.
        * `device_name` - (Required) Device name or key.
        * `uplink_name` - (Required) Uplink name for this Pnic.
    * `transport_node_profile_sub_config` - (Optional) Transport Node Profile sub-configuration Options.
        * `host_switch_config_option` - (Required) Subset of the host switch configuration.
            * `host_switch_id` - (Optional) The host switch id. This ID will be used to reference a host switch.
            * `host_switch_profile` - (Optional) Policy paths of host switch profiles to be associated with this host switch.
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
        * `transport_zone` - (Required) Policy path of the transport zone for this endpoint.
        * `transport_zone_profiles` - (Optional) Policy paths of the transport zone profiles associated with this transport zone endpoint on this transport node.
    * `uplink` - (Optional) Uplink/LAG of VMware vSphere Distributed Switch connected to the HostSwitch.
        * `uplink_name` - (Required) Uplink name from UplinkHostSwitch profile.
        * `vds_lag_name` - (Optional) Link Aggregation Group (LAG) name of Virtual Distributed Switch.
        * `vds_uplink_name` - (Optional) Uplink name of VMware vSphere Distributed Switch (VDS).
    * `vmk_install_migration` - (Optional) The vmknic and logical switch mappings.
        * `destination_network` - (Required) The network id to which the ESX vmk interface will be migrated.
        * `device_name` - (Required) ESX vmk interface name.

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing policy Host Transport Node profile can be [imported][docs-import] into this resource, via the following command:

```
terraform import nsxt_policy_host_transport_node_profile.test POLICY_PATH
```
The above command imports Policy Host Transport Node Profile named `test` with NSX policy path `POLICY_PATH`.
