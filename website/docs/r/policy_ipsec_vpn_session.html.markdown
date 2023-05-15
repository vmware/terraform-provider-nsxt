---
subcategory: "VPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ipsec_vpn_session"
description: A resource to configure a IPSec VPN session.
---

# nsxt_policy_ipsec_vpn_session

This resource provides a method for the management of a IPSec VPN session.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_ipsec_vpn_session" "test" {
  display_name               = "Route-Based VPN Session"
  description                = "Terraform-provisioned IPsec Route-Based VPN"
  ike_profile_path           = nsxt_policy_ipsec_vpn_ike_profile.profile_ike_l3vpn.path
  tunnel_profile_path        = nsxt_policy_ipsec_vpn_tunnel_profile.profile_tunnel_l3vpn.path
  local_endpoint_path        = data.nsxt_policy_ipsec_vpn_local_endpoint.private_endpoint.path
  dpd_profile_path           = nsxt_policy_ipsec_vpn_dpd_profile.test.path
  enabled                    = true
  service_path               = nsxt_policy_ipsec_vpn_service.test.path
  vpn_type                   = "RouteBased"
  authentication_mode        = "PSK"
  compliance_suite           = "NONE"
  ip_addresses               = ["169.254.152.2"]
  prefix_length              = 30
  peer_address               = "18.18.18.19"
  peer_id                    = "18.18.18.19"
  psk                        = "BhvrlXXmH+TxXlFKNaF5mAXnnLja3lSQ"
  connection_initiation_mode = "INITIATOR"
}

resource "nsxt_policy_ipsec_vpn_session" "test2" {
  display_name               = "Policy-Based VPN Session"
  description                = "Terraform-provisioned IPsec Policy-Based VPN"
  ike_profile_path           = nsxt_policy_ipsec_vpn_ike_profile.profile_ike_l3vpn.path
  tunnel_profile_path        = nsxt_policy_ipsec_vpn_tunnel_profile.profile_tunnel_l3vpn.path
  local_endpoint_path        = data.nsxt_policy_ipsec_vpn_local_endpoint.private_endpoint.path
  dpd_profile_path           = nsxt_policy_ipsec_vpn_dpd_profile.test.path
  enabled                    = true
  service_path               = nsxt_policy_ipsec_vpn_service.test.path
  vpn_type                   = "PolicyBased"
  authentication_mode        = "PSK"
  compliance_suite           = "NONE"
  peer_address               = "18.18.18.19"
  peer_id                    = "18.18.18.19"
  psk                        = "BhvrlXXmH+TxXlFKNaF5mAXnnLja3lSQ"
  connection_initiation_mode = "INITIATOR"
  rule {
    sources      = ["192.168.10.0/24"]
    destinations = ["192.169.10.0/24"]
    action       = "BYPASS"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `ike_profile_path` - (Optional) Policy path referencing IKE profile. Note that if user wants to create session with `compliance_suite`, then this field should not be configured, the provider will use the default Profile for each compliance suite type.
* `tunnel_profile_path` - (Optional) Policy path referencing Tunnel profile to be used. Note that if user wants to create session with `compliance_suite`, then this field should not be configured, the provider will use the default Profile for each compliance suite type.
* `enabled` - (Optional) Boolean. Enable/Disable IPsec VPN session. Default is "true" (session enabled).
* `service_path` - (Required) The path of the IPSec VPN service for the VPN session.
* `dpd_profile_path` - (Optional) Policy path referencing Dead Peer Detection (DPD) profile. Default is set to system default profile.
* `vpn_type` - (Required) `RouteBased` or `PolicyBased`. Policy Based VPN requires to define protect rules that match local and peer subnets. IPSec security association is negotiated for each pair of local and peer subnet. For PolicyBased Session, `rule` must be specified with `sources`, `destination` and `action`. A Route Based VPN is more flexible, more powerful and recommended over policy based VPN. IP Tunnel port is created and all traffic routed via tunnel port is protected. Routes can be configured statically or can be learned through BGP. A route based VPN is a must for establishing redundant VPN session to remote site. For RouteBased VPN session, `ip_addresses` and `prefix_length` must be specified to create the tunnel interface and its subnet.
* `compliance_suite` -  (Optional) Compliance suite. Value is one of `CNSA`, `SUITE_B_GCM_128`, `SUITE_B_GCM_256`, `PRIME`, `FOUNDATION`, `FIPS`, `None`.
* `compliance_initiation_mode` - (Optional) Connection initiation mode used by local endpoint to establish ike connection with peer site. `INITIATOR` - In this mode local endpoint initiates tunnel setup and will also respond to incoming tunnel setup requests from peer gateway. `RESPOND_ONLY` - In this mode, local endpoint shall only respond to incoming tunnel setup requests. It shall not initiate the tunnel setup. `ON_DEMAND` - In this mode local endpoint will initiate tunnel creation once first packet matching the policy rule is received and will also respond to incoming initiation request.
* `authentication_mode` - (Optional) Peer authentication mode. `PSK` - In this mode a secret key shared between local and peer sites is to be used for authentication. The secret key can be a string with a maximum length of 128 characters. `CERTIFICATE` - In this mode a certificate defined at the global level is to be used for authentication. If user wants to configure compliance_suite, then the authentication_mode can only be `CERTIFICATE`.
* `ip_addresses` - (Optional) IP Tunnel interface (commonly referred as VTI) ip_addresses. Only applied for Route Based VPN Session. 
* `prefix_length` - (Optional) Subnet Prefix Length. Only applied for Route Based VPN Session. 
* `peer_address` - (Optional) Public IPV4 address of the remote device terminating the VPN connection.
* `peer_id` - (Optional) Peer ID to uniquely identify the peer site. The peer ID is the public IP address of the remote device terminating the VPN tunnel. When NAT is configured for the peer, enter the private IP address of the peer.
* `local_endpoint_path` - (Required) Policy path referencing Local endpoint. In VMC, Local Endpoints are pre-configured the user can refer to their path using `data nsxt_policy_ipsec_vpn_local_endpoint` and using the "Private IP1" or "Public IP1" values to refer to the private and public endpoints respectively. Note that if `authentication_mode` is `CERTIFICATE`, then the local_endpoint must be configured with `certificate_path` and `trust_ca_paths`.
* `rule` - (Optional) Bypass rules for this IPSec VPN Session. Only applicable to `PolicyBased` VPN Session. 
  * `sources` - (Optional) List of source subnets. Subnet format is ipv4 CIDR.
  * `destinations` - (Optional) List of distination subnets. Subnet format is ipv4 CIDR.
  * `action` - (Optional) `PROTECT` or `BYPASS`. Default is `PROTECT`.
* `direction` - (Optional) The traffic direction apply to the MSS clamping. Value is one of `NONE`, `INBOUND_CONNECTION`, `OUTBOUND_CONNECTION` AND `BOTH`.
* `max_segment_size` - (Optional) Maximum amount of data the host will accept in a TCP segment. Value is an int between `108` and `8860`. If not specified then the value would be the automatic calculated MSS value.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `rule` additional arguments:
  * `nsx_id` - The NSX ID of the bypass rule for this IPSec VPN Session.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ipsec_vpn_session.test PATH
```

The above command imports IPSec VPN session named `test` that corresponds to NSX IPSec VPN session with policy path `PATH`.
