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

* Flexible Segments: [nsxt_policy_segment](https://www.terraform.io/docs/providers/nsxt/r/policy_segment)
* Distributed Firewall Section: [nsxt_policy_security_policy](https://www.terraform.io/docs/providers/nsxt/r/policy_security_policy)
* NAT: [nsxt_policy_nat_rule](https://www.terraform.io/docs/providers/nsxt/r/policy_nat_rule)
* Security Groups: [nsxt_policy_group](https://www.terraform.io/docs/providers/nsxt/r/policy_group)
* Services: [nsxt_policy_service](https://www.terraform.io/docs/providers/nsxt/r/policy_service)
* DHCP Relay: [nsxt_policy_dhcp_relay](https://www.terraform.io/docs/providers/nsxt/r/policy_dhcp_relay)
* Virtual Machine Tags: [nsxt_policy_vm_tag](https://www.terraform.io/docs/providers/nsxt/r/policy_vm_tags)

## Available Data Sources for use with VMConAWS

The following Data Sources are available to use with VMConAWS:

* Services: [nsxt_policy_service](https://www.terraform.io/docs/providers/nsxt/d/policy_service)
* Edge Cluster: [nsxt_policy_edge_cluster](https://www.terraform.io/docs/providers/nsxt/d/policy_edge_cluster)
* Transport Zone: [nsxt_policy_transport_zone](https://www.terraform.io/docs/providers/nsxt/d/policy_transport_zone)
* Tier-0 Gateway: [nsxt_policy_tier0_gateway](https://www.terraform.io/docs/providers/nsxt/d/policy_tier0_gateway)
* Tier-1 Gateway: [nsxt_policy_tier1_gateway](https://www.terraform.io/docs/providers/nsxt/d/policy_tier1_gateway)


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
