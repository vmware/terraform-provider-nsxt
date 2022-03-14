---
subcategory: "Policy - Gateways and Routing"
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
    enabled                    = true
    locale_service             = "default"
    service_id                 = "default"
    tier0_id                   = "vmc"
    vpn_type                   = "RouteBasedIPSecVpnSession"
    authentication_mode        = "PSK"
    compliance_suite           = "NONE"
    subnets                    = ["169.254.152.2"]
    prefix_length              = 30
    peer_address               = "18.18.18.19"
    peer_id                    = "18.18.18.19"
    psk                        = "VMware123!"
    connection_initiation_mode = "INITIATOR"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `ike_profile_path` - (Optional) Policy path referencing IKE profile.
* `tunnel_profile_path` - (Optional) Policy path referencing Tunnel profile to be used. Default is set to system default profile.
* `enabled` - (Optional) Boolean. Enable/Disable IPsec VPN session. Default is "true" (session enabled).
* `locale_service` - (Optional) Unique identifier of the Locale Service resource. Default value is "default".
* `service_id` - (Optional) Unique identifier of the Service. Default value is "default".
* `tier0_id` - (Optional) Unique identifier of the T0 resource. Default value is "vmc".
* `dpd_profile_path` - (Optional) Policy path referencing Dead Peer Detection (DPD) profile. Default is set to system default profile.
* `vpn_type` - (Optional) "RouteBasedIPSecVpnSession" or "PolicyBasedIPSecVpnSession". Policy Based VPN requires to define protect rules that match local and peer subnets. IPSec security associations is negotiated for each pair of local and peer subnet. A Route Based VPN is more flexible, more powerful and recommended over policy based VPN. IP Tunnel port is created and all traffic routed via tunnel port is protected. Routes can be configured statically or can be learned through BGP. A route based VPN is must for establishing redundant VPN session to remote site.
* `compliance_suite` - (Optional) Connection initiation mode used by local endpoint to establish ike connection with peer site. INITIATOR - In this mode local endpoint initiates tunnel setup and will also respond to incoming tunnel setup requests from peer gateway. RESPOND_ONLY - In this mode, local endpoint shall only respond to incoming tunnel setup requests. It shall not initiate the tunnel setup. ON_DEMAND - In this mode local endpoint will initiate tunnel creation once first packet matching the policy rule is received and will also respond to incoming initiation request.
* `subnets` - (Optional) IP Tunnel interface (commonly referred as VTI) subnet.
* `prefix_length` - (Optional) Subnet Prefix Length.
* `peer_address` - (Optional) Public IPV4 address of the remote device terminating the VPN connection.
* `peer_id` - (Optional) Peer ID to uniquely identify the peer site. The peer ID is the public IP address of the remote device terminating the VPN tunnel. When NAT is configured for the peer, enter the private IP address of the peer.
* `local_endpoint_path` - (Optional) Policy path referencing Local endpoint. In VMC, Local Endpoints are pre-configured the user can refer to their path using `data nsxt_policy_ipsec_vpn_local_endpoint` and using the "Private IP1" or "Public IP1" values to refer to the private and public endpoints respectively.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_ipsec_vpn_session.test UUID
```

The above command imports IPSec VPN  session named `test` with the NSX IPSec VPN Ike session ID `UUID`.
