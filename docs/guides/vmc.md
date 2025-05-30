---
page_title: "Support for VMware Cloud on AWS"
description: |-
  Support for VMware Cloud on AWS
---

# Support for VMware Cloud on AWS

The provider includes support for resources and datasources applicable on VMware Cloud on AWS.

For more information on VMware Cloud on AWS, refer to the [VMware Cloud on AWS product site](https://www.vmware.com/solutions/partner-managed-services/vmc-on-aws)

For more information on NSX, refer to the [NSX documentation](https://techdocs.broadcom.com/us/en/vmware-cis/nsx.html).

## Basic Authentication with VMware Cloud on AWS

```hcl
provider "nsxt" {
  host                 = "nsx-54-203-97-32.rp.vmwarevmc.com/vmc/reverse-proxy/api/orgs/54f6629e-460b-4521-ac2b-ae3b5461916c/sddcs/1b3c6c7e-59ba-4df7-ab0e-af9b06369a83/sks-nsxt-manager"
  vmc_token            = var.vmc_token
  allow_unverified_ssl = true
  enforcement_point    = "vmc-enforcementpoint"
}
```

Note that `host` is the ProxyURL and `vmc_token` is the API Access Token

## Available Resources for use with VMware Cloud on AWS

The following Resources are available to use with VMware Cloud on AWS:

* Flexible Segment: [nsxt_policy_segment](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_segment)
* Fixed Segment: [nsxt_policy_fixed_segment](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_fixed_segment)
* Tier1 Gateway: [nsxt_policy_tier1_gateway](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_tier1_gateway)
* Distributed Firewall Section: [nsxt_policy_security_policy](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_security_policy)
* Predefined Distributed Firewall Section: [nsxt_policy_predefined_security_policy](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_predefined_security_policy)
* Predefined Gateway Firewall Section: [nsxt_policy_predefined_gateway_policy](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_predefined_gateway_policy)
* NAT: [nsxt_policy_nat_rule](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_nat_rule)
* Security Group: [nsxt_policy_group](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_group)
* Service: [nsxt_policy_service](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_service)
* DHCP Relay: [nsxt_policy_dhcp_relay](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_dhcp_relay)
* DHCP Service: [nsxt_policy_dhcp_service](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_dhcp_service)
* Virtual Machine Tags: [nsxt_policy_vm_tag](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_vm_tags)
* Context Profile: [nsxt_policy_context_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_context_profile)
* Gateway QoS Profile: [nsxt_policy_gateway_qos_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_gateway_qos_profile)
* IDS Policy: [nsxt_policy_intrusion_service_policy](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_intrusion_service_policy)
* IDS Profile: [nsxt_policy_intrusion_service_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_intrusion_service_profile)
* DNS Forwarder Zone: [nsxt_policy_dns_porwarder_zone](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_dns_forwarder_zone)
* Gateway DNS Forwarder: [nsxt_policy_gateway_dns_forwarder](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_gateway_dns_forwarder)
* IPSec VPN DPD Profile: [nsxt_policy_ipsec_vpn_dpd_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_ipsec_vpn_dpd_profile)
* IPSec VPN IKE Profile: [nsxt_policy_ipsec_vpn_ike_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_ipsec_vpn_ike_profile)
* IPSec VPN Tunnel Profile: [nsxt_policy_ipsec_vpn_tunnel_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_ipsec_vpn_tunnel_profile)
* IPSec VPN Local Endpoint: [nsxt_policy_ipsec_vpn_local_endpoint](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_ipsec_vpn_local_endpoint)
* IPSec VPN Service: [nsxt_policy_ipsec_vpn_service](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_ipsec_vpn_service)
* IPSec VPN Session: [nsxt_policy_ipsec_vpn_session](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_ipsec_vpn_session)
* L2 VPN Service: [nsxt_policy_l2_vpn_service](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_l2_vpn_service)
* L2 VPN Session: [nsxt_policy_l2_vpn_session](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/resources/policy_l2_vpn_session)

## Available Data Sources for use with VMware Cloud on AWS

The following Data Sources are available to use with VMware Cloud on AWS:

* Service: [nsxt_policy_service](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_service)
* Group: [nsxt_policy_group](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_group)
* Edge Cluster: [nsxt_policy_edge_cluster](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_edge_cluster)
* Transport Zone: [nsxt_policy_transport_zone](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_transport_zone)
* Tier-0 Gateway: [nsxt_policy_tier0_gateway](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_tier0_gateway)
* Tier-1 Gateway: [nsxt_policy_tier1_gateway](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_tier1_gateway)
* IPSec VPN Service: [nsxt_policy_ipsec_vpn_service](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_ipsec_vpn_service)
* IPSec VPN Local Endpoint: [nsxt_policy_ipsec_vpn_local_endpoint](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_ipsec_vpn_local_endpoint)
* L2 VPN Service: [nsxt_policy_l2_vpn_service](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_l2_vpn_service)
* Intrusion Service Profile: [nsxt_policy_intrusion_service_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_intrusion_service_profile)
* DHCP Server: [nsxt_policy_dhcp_server](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_dhcp_server)
* IP Discovery Segment Profile: [nsxt_policy_ip_discovery_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_ip_discovery_profile)
* MAC Discovery Segment Profile: [nsxt_policy_mac_discovery_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_mac_discovery_profile)
* QoS Segment Profile: [nsxt_policy_qos_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_qos_profile)
* Segment Security Profile: [nsxt_policy_segment_security_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_segment_security_profile)
* Spoofguard Segment Profile: [nsxt_policy_spoof_guard_profile](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_spoof_guard_profile)
* Certificate: [nsxt_policy_certificate](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_certificate)
* Virtual Machine: [nsxt_policy_vm](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_vm)
* Virtual Machines with filters: [nsxt_policy_vms](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_vms)
* Security Policy: [nsxt_policy_security_policy](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_security_policy)
* Gateway Policy: [nsxt_policy_gateway_policy](https://registry.terraform.io/providers/vmware/nsxt/latest/docs/data-sources/policy_gateway_policy)

Note that in relevant resources, a domain needs to be specified, since on VMware Cloud on AWS the domain is different from `default` domain. For example:

```hcl
resource "nsxt_policy_group" "test" {
  display_name = "test"
  domain       = "cgw"

  criteria {
    ipaddress_expression {
      ip_addresses = ["10.1.60.20", "10.1.60.21"]
    }
  }
}
```
