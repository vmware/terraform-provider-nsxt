---
layout: "nsxt"
page_title: "VMware NSX-T Terraform Provider for VMConAWS"
description: |-
  The VMware NSX-T Terraform Provider for VMConAWS
---

# Terraform Provider for NSX-T support extended to VMConAWS

The NSX Terraform Resources and Data Sources applicable to NSX-T on VMConAWS are supported. More information on VMConAWS can be found on the [VMConAWS Product Page](https://cloud.vmware.com/vmc-aws)

Documentation on the NSX platform can be found on the [NSX Documentation Page](https://docs.vmware.com/en/VMware-NSX-T/index.html)

## Basic Authentication with VMConAWS


```hcl
provider "nsxt" {
  host                 = "nsx-54-203-97-32.rp.vmwarevmc.com/vmc/reverse-proxy/api/orgs/54f6629e-460b-4521-ac2b-ae3b5461916c/sddcs/1b3c6c7e-59ba-4df7-ab0e-af9b06369a83/sks-nsxt-manager"
  vmc_token            = var.vmc_token
  allow_unverified_ssl = true
  enforcement_point    = "vmc-enforcementpoint"
}

```

Note that `host` is the ProxyURL and `vmc_token` is the API Access Token

## Available Resources for use with VMConAWS

The following Resources are available to use with VMConAWS:

* Flexible Segment: [nsxt_policy_segment](https://www.terraform.io/docs/providers/nsxt/r/policy_segment)
* Fixed Segment: [nsxt_policy_fixed_segment](https://www.terraform.io/docs/providers/nsxt/r/policy_fixed_segment)
* Tier1 Gateway: [nsxt_policy_tier1_gateway](https://www.terraform.io/docs/providers/nsxt/r/policy_tier1_gateway)
* Distributed Firewall Section: [nsxt_policy_security_policy](https://www.terraform.io/docs/providers/nsxt/r/policy_security_policy)
* Predefined Distributed Firewall Section: [nsxt_policy_predefined_security_policy](https://www.terraform.io/docs/providers/nsxt/r/policy_predefined_security_policy)
* Predefined Gateway Firewall Section: [nsxt_policy_predefined_gateway_policy](https://www.terraform.io/docs/providers/nsxt/r/policy_predefined_gateway_policy)
* NAT: [nsxt_policy_nat_rule](https://www.terraform.io/docs/providers/nsxt/r/policy_nat_rule)
* Security Group: [nsxt_policy_group](https://www.terraform.io/docs/providers/nsxt/r/policy_group)
* Service: [nsxt_policy_service](https://www.terraform.io/docs/providers/nsxt/r/policy_service)
* DHCP Relay: [nsxt_policy_dhcp_relay](https://www.terraform.io/docs/providers/nsxt/r/policy_dhcp_relay)
* DHCP Service: [nsxt_policy_dhcp_service](https://www.terraform.io/docs/providers/nsxt/r/policy_dhcp_service)
* Virtual Machine Tags: [nsxt_policy_vm_tag](https://www.terraform.io/docs/providers/nsxt/r/policy_vm_tags)
* Context Profile: [nsxt_policy_context_profile](https://www.terraform.io/docs/providers/nsxt/r/policy_context_profile)
* Gateway QoS Profile: [nsxt_policy_gateway_qos_profile](https://www.terraform.io/docs/providers/nsxt/r/policy_gateway_qos_profile)
* IDS Policy: [nsxt_policy_intrusion_service_policy](https://www.terraform.io/docs/providers/nsxt/r/policy_intrusion_service_policy)
* IDS Profile: [nsxt_policy_intrusion_service_profile](https://www.terraform.io/docs/providers/nsxt/r/policy_intrusion_service_profile)
* DNS Forwarder Zone: [nsxt_policy_dns_porwarder_zone](https://www.terraform.io/docs/providers/nsxt/r/policy_dns_forwarder_zone)
* Gateway DNS Forwarder: [nsxt_policy_gateway_dns_forwarder](https://www.terraform.io/docs/providers/nsxt/r/policy_gateway_dns_forwarder)
* IPSec VPN DPD Profile: [nsxt_policy_ipsec_vpn_dpd_profile](https://www.terraform.io/docs/providers/nsxt/r/policy_ipsec_vpn_dpd_profile)
* IPSec VPN IKE Profile: [nsxt_policy_ipsec_vpn_ike_profile](https://www.terraform.io/docs/providers/nsxt/r/policy_ipsec_vpn_ike_profile)
* IPSec VPN Tunnel Profile: [nsxt_policy_ipsec_vpn_tunnel_profile](https://www.terraform.io/docs/providers/nsxt/r/policy_ipsec_vpn_tunnel_profile)
* IPSec VPN Local Endpoint: [nsxt_policy_ipsec_vpn_local_endpoint](https://www.terraform.io/docs/providers/nsxt/r/policy_ipsec_vpn_local_endpoint)
* IPSec VPN Service: [nsxt_policy_ipsec_vpn_service](https://www.terraform.io/docs/providers/nsxt/r/policy_ipsec_vpn_service)
* IPSec VPN Session: [nsxt_policy_ipsec_vpn_session](https://www.terraform.io/docs/providers/nsxt/r/policy_ipsec_vpn_session)
* L2 VPN Service: [nsxt_policy_l2_vpn_service](https://www.terraform.io/docs/providers/nsxt/r/policy_l2_vpn_service)
* L2 VPN Session: [nsxt_policy_l2_vpn_session](https://www.terraform.io/docs/providers/nsxt/r/policy_l2_vpn_session)

## Available Data Sources for use with VMConAWS

The following Data Sources are available to use with VMConAWS:

* Service: [nsxt_policy_service](https://www.terraform.io/docs/providers/nsxt/d/policy_service)
* Group: [nsxt_policy_group](https://www.terraform.io/docs/providers/nsxt/d/policy_group)
* Edge Cluster: [nsxt_policy_edge_cluster](https://www.terraform.io/docs/providers/nsxt/d/policy_edge_cluster)
* Transport Zone: [nsxt_policy_transport_zone](https://www.terraform.io/docs/providers/nsxt/d/policy_transport_zone)
* Tier-0 Gateway: [nsxt_policy_tier0_gateway](https://www.terraform.io/docs/providers/nsxt/d/policy_tier0_gateway)
* Tier-1 Gateway: [nsxt_policy_tier1_gateway](https://www.terraform.io/docs/providers/nsxt/d/policy_tier1_gateway)
* IPSec VPN Service: [nsxt_policy_ipsec_vpn_service](https://www.terraform.io/docs/providers/nsxt/d/policy_ipsec_vpn_service)
* IPSec VPN Local Endpoint: [nsxt_policy_ipsec_vpn_local_endpoint](https://www.terraform.io/docs/providers/nsxt/d/policy_ipsec_vpn_local_endpoint)
* L2 VPN Service: [nsxt_policy_l2_vpn_service](https://www.terraform.io/docs/providers/nsxt/d/policy_l2_vpn_service)
* Intrusion Service Profile: [nsxt_policy_intrusion_service_profile](https://www.terraform.io/docs/providers/nsxt/d/policy_intrusion_service_profile)
* DHCP Server: [nsxt_policy_dhcp_server](https://www.terraform.io/docs/providers/nsxt/d/policy_dhcp_server)
* IP Discovery Segment Profile: [nsxt_policy_ip_discovery_profile](https://www.terraform.io/docs/providers/nsxt/d/policy_ip_discovery_profile)
* MAC Discovery Segment Profile: [nsxt_policy_mac_discovery_profile](https://www.terraform.io/docs/providers/nsxt/d/policy_mac_discovery_profile)
* QoS Segment Profile: [nsxt_policy_qos_profile](https://www.terraform.io/docs/providers/nsxt/d/policy_qos_profile)
* Segment Security Profile: [nsxt_policy_segment_security_profile](https://www.terraform.io/docs/providers/nsxt/d/policy_segment_security_profile)
* Spoofguard Segment Profile: [nsxt_policy_spoof_guard_profile](https://www.terraform.io/docs/providers/nsxt/d/policy_spoof_guard_profile)
* Certificate: [nsxt_policy_certificate](https://www.terraform.io/docs/providers/nsxt/d/policy_certificate)
* Virtual Machine: [nsxt_policy_vm](https://www.terraform.io/docs/providers/nsxt/d/policy_vm)
* Virtual Machines with filters: [nsxt_policy_vms](https://www.terraform.io/docs/providers/nsxt/d/policy_vms)
* Security Policy: [nsxt_policy_security_policy](https://www.terraform.io/docs/providers/nsxt/d/policy_security_policy)
* Gateway Policy: [nsxt_policy_gateway_policy](https://www.terraform.io/docs/providers/nsxt/d/policy_gateway_policy)


Note that in relevant resources, a domain needs to be specified, since on VMC the domain is different from `default` domain. For example:

```
resource "nsxt_policy_group" "test" {
  display_name = "test"
  domain = "cgw"

  criteria {
    ipaddress_expression {
      ip_addresses = ["10.1.60.20", "10.1.60.21"]
    }
  }
}
```
