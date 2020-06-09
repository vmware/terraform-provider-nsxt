---
layout: "nsxt"
page_title: "Provider: NSXT"
sidebar_current: "docs-nsxt-guides-vmc-index"
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

* Flexible Segments: [nsxt_policy_segment](https://www.terraform.io/docs/providers/nsxt/r/policy_segment.html)
* Distributed Firewall Section: [nsxt_policy_security_policy](https://www.terraform.io/docs/providers/nsxt/r/policy_security_policy.html)
* NAT: [nsxt_policy_nat_rule](https://www.terraform.io/docs/providers/nsxt/r/policy_nat_rule.html)
* Security Groups: [nsxt_policy_group](https://www.terraform.io/docs/providers/nsxt/r/policy_group.html)
* Services: [nsxt_policy_service](https://www.terraform.io/docs/providers/nsxt/r/policy_service.html)
* DHCP Relay: [nsxt_policy_dhcp_relay](https://www.terraform.io/docs/providers/nsxt/r/policy_dhcp_relay.html)
* Virtual Machine Tags: [nsxt_policy_vm_tag](https://www.terraform.io/docs/providers/nsxt/r/policy_vm_tags.html)

## Available Data Sources for use with VMConAWS

The following Data Sources are available to use with VMConAWS:

* Services: [nsxt_policy_service](https://www.terraform.io/docs/providers/nsxt/d/policy_service.html)
* Edge Cluster: [nsxt_policy_edge_cluster](https://www.terraform.io/docs/providers/nsxt/d/policy_edge_cluster.html)
* Transport Zone: [nsxt_policy_transport_zone](https://www.terraform.io/docs/providers/nsxt/d/policy_transport_zone.html)
* Tier-0 Gateway: [nsxt_policy_tier0_gateway](https://www.terraform.io/docs/providers/nsxt/d/policy_tier0_gateway.html)
* Tier-1 Gateway: [nsxt_policy_tier1_gateway](https://www.terraform.io/docs/providers/nsxt/d/policy_tier1_gateway.html)
