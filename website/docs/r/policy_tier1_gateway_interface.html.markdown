---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_tier1_gateway_interface"
description: A resource to configure an Interface on Tier-1 gateway on NSX Policy manager.
---

# nsxt_policy_tier1_gateway_interface

This resource provides a method for the management of a Tier-1 gateway Interface. Note that edge cluster must be configured on Tier-1 Gateway in order to configure interfaces on it.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

# Example Usage

```hcl
data "nsxt_policy_tier1_gateway" "gw1" {
  display_name = "gw1"
}

data "nsxt_policy_ipv6_ndra_profile" "slaac" {
  display_name = "slaac"
}

resource "nsxt_policy_vlan_segment" "segment1" {
  display_name = "segment1"
  vlan_ids     = [12]
}

resource "nsxt_policy_tier1_gateway_interface" "if1" {
  display_name           = "segment1_interface"
  description            = "connection to segment1"
  gateway_path           = data.nsxt_policy_tier1_gateway.gw1.path
  segment_path           = nsxt_policy_vlan_segment.segment1.path
  subnets                = ["12.12.2.13/24"]
  mtu                    = 1500
  ipv6_ndra_profile_path = data.nsxt_policy_ipv6_ndra_profile.slaac.path
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_tier1_gateway_interface" "if1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name           = "segment1_interface"
  description            = "connection to segment1"
  gateway_path           = data.nsxt_policy_tier1_gateway.gw1.path
  segment_path           = nsxt_policy_vlan_segment.segment1.path
  subnets                = ["12.12.2.13/24"]
  mtu                    = 1500
  ipv6_ndra_profile_path = data.nsxt_policy_ipv6_ndra_profile.slaac.path
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `gateway_path` - (Required) Policy path for the Tier-1 Gateway.
* `segment_path` - (Required) Policy path for segment to be connected with this Tier1 Gateway.
* `subnets` - (Required) list of Ip Addresses/Prefixes in CIDR format, to be associated with this interface.
* `mtu` - (Optional) Maximum Transmission Unit for this interface.
* `ipv6_ndra_profile_path` - (Optional) IPv6 NDRA profile to be associated with this interface.
* `dhcp_relay_path` - (Optional) DHCP relay path to be associated with this interface.
* `urpf_mode` - (Optional) Unicast Reverse Path Forwarding mode, one of `NONE`, `STRICT`. Default is `STRICT`. This attribute is supported with NSX 3.0.0 onwards.
* `site_path` - (Required for global manager only) Path of the site the Tier1 edge cluster belongs to. This configuration is required for global manager only. `path` field of the existing `nsxt_policy_site` can be used here.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing policy Tier-1 Gateway Interface can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_tier1_gateway_interface.interface1 GW-ID/LOCALE-SERVICE-ID/ID
```
The above command imports the policy Tier-1 gateway interface named `interface1` with the NSX Policy ID `ID` on Tier1 Gateway `GW-ID`, under locale service `LOCALE-SERVICE-ID`.

```
terraform import nsxt_policy_tier1_gateway_interface.interface1 POLICY_PATH
```
The above command imports the policy Tier-1 gateway interface named `interface1` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.

