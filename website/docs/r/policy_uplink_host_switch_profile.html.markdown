---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_uplink_host_switch_profile"
description: A resource to configure a uplink host switch profile in NSX Policy.
---

# nsxt_policy_uplink_host_switch_profile

This resource provides a method for the management of a uplink host switch profile which can be used within NSX Policy.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_uplink_host_switch_profile" "uplink_host_switch_profile" {
  description  = "Uplink host switch profile provisioned by Terraform"
  display_name = "uplink_host_switch_profile"

  mtu            = 1500
  transport_vlan = 0
  overlay_encap  = "GENEVE"
  lag {
    name                   = "t_created_lag"
    load_balance_algorithm = "SRCDESTIPVLAN"
    mode                   = "ACTIVE"
    number_of_uplinks      = 2
    timeout_type           = "SLOW"
  }
  teaming {
    active {
      uplink_name = "t1"
      uplink_type = "PNIC"
    }
    standby {
      uplink_name = "t2"
      uplink_type = "LAG"
    }
    policy = "FAILOVER_ORDER"
  }
  named_teaming {
    active {
      uplink_name = "nt1"
      uplink_type = "PNIC"
    }
    standby {
      uplink_name = "nt2"
      uplink_type = "PNIC"
    }
    policy = "FAILOVER_ORDER"
    name   = "named_teaming1"
  }
  named_teaming {
    active {
      uplink_name = "nt3"
      uplink_type = "PNIC"
    }
    policy = "LOADBALANCE_SRCID"
    name   = "named_teaming2"
  }

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `mtu` - (Optional) Maximum Transmission Unit used for uplinks. Minimum: 1280.
* `transport_vlan` - (Optional) VLAN used for tagging Overlay traffic of associated HostSwitch. Default: 0.
* `overlay_encap` - (Optional) The protocol used to encapsulate overlay traffic. Possible values are: `VXLAN`, `GENEVE`. Default: `GENEVE`.
* `lag` - (Optional) List of LACP group.
  * `load_balance_algorithm` - (Required) LACP load balance Algorithm. Possible values are: `SRCMAC`, `DESTMAC`, `SRCDESTMAC`, `SRCDESTIPVLAN`, `SRCDESTMACIPPORT`.
  * `mode` - (Required) LACP group mode. Possible values are: `ACTIVE`, `PASSIVE`.
  * `name` - (Required) Lag name.
  * `number_of_uplinks` - (Required) Number of uplinks. Minimum: 2, maximum: 32.
  * `timeout_type` - (Optional) LACP timeout type. Possible values are: `SLOW`, `FAST`. Default: `SLOW`.
* `teaming` - (Required) Default TeamingPolicy associated with this UplinkProfile.
  * `active` - (Required) List of Uplinks used in active list.
    * `uplink_name` - (Required) Name of this uplink.
    * `uplink_type` - (Required) Type of the uplink. Possible values are: `PNIC`, `LAG`.
  * `policy` - (Required) Teaming policy. Possible values are: `FAILOVER_ORDER`, `LOADBALANCE_SRCID`, `LOADBALANCE_SRC_MAC`.
  * `standby` - (Optional) List of Uplinks used in standby list.
    * `uplink_name` - (Required) Name of this uplink.
    * `uplink_type` - (Required) Type of the uplink. Possible values are: `PNIC`, `LAG`.
* `named_teaming` - (Optional) List of named uplink teaming policies that can be used by logical switches.
  * `active` - (Required) List of Uplinks used in active list.
    * `uplink_name` - (Required) Name of this uplink.
    * `uplink_type` - (Required) Type of the uplink. Possible values are: `PNIC`, `LAG`.
  * `policy` - (Required) Teaming policy. Possible values are: `FAILOVER_ORDER`, `LOADBALANCE_SRCID`, `LOADBALANCE_SRC_MAC`.
  * `standby` - (Optional) List of Uplinks used in standby list.
    * `uplink_name` - (Required) Name of this uplink.
    * `uplink_type` - (Required) Type of the uplink. Possible values are: `PNIC`, `LAG`.
  * `name` - (Optional) An uplink teaming policy of a given name defined in UplinkHostSwitchProfile. The names of all NamedTeamingPolicies in an UplinkHostSwitchProfile must be different, but a name can be shared by different UplinkHostSwitchProfiles. Different TransportNodes can use different NamedTeamingPolicies having the same name in different UplinkHostSwitchProfiles to realize an uplink teaming policy on a logical switch. An uplink teaming policy on a logical switch can be any policy defined by a user; it does not have to be a single type of FAILOVER or LOADBALANCE. It can be a combination of types, for instance, a user can define a policy with name `MyHybridTeamingPolicy` as `FAILOVER on all ESX TransportNodes and LOADBALANCE on all KVM TransportNodes`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `realized_id` - Realized ID for the profile. For reference in fabric resources (such as `transport_node`), `realized_id` should be used rather than `id`.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_uplink_host_switch_profile.uplink_host_switch_profile UUID
```

The above command imports UplinkHostSwitchProfile named `uplink_host_switch_profile` with the NSX ID `UUID`.


```
terraform import nsxt_policy_uplink_host_switch_profile.uplink_host_switch_profile POLICY_PATH
```
The above command imports the uplink host switch profile named `uplink_host_switch_profile` with policy path `POLICY_PATH`.
