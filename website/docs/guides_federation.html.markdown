---
layout: "nsxt"
page_title: "Provider: NSXT"
sidebar_current: "docs-nsxt-guides-vmc-index"
description: |-
  The VMware NSX-T Terraform Provider support for Federation 
---

# Terraform Provider for NSX-T extension for Federation

The NSX Terraform Resources and Data Sources are now extended to support Federation. More information on Fedeartion that was introduced in NSX-T 3.0 can be found on the [NSX-T Product Page for Fedeartion](https://docs.vmware.com/en/VMware-NSX-T-Data-Center/3.0/administration/GUID-D5B6DC79-6733-44A7-8072-50221CF2122A.html)

Documentation on the NSX platform can be found on the [NSX Documentation Page](https://docs.vmware.com/en/VMware-NSX-T/index.html)

## Basic Authentication with Federation


```hcl
provider "nsxt" {
  host           = "192.168.110.41"
  username       = "admin"
  password       = "default"
  global_manager = true
}
```

Note that authentication with the Global Manager uses the same NSX-T Terraform provider but uses a `global_manager = true` flag for identifiation.

## Using Resources and Data Sources
Just like authentication, the resources and data soruces are common between the Local Manager (LM) and Global Manager (GM). However, there could be small differences in the attributes for each of the objects. These differences, when applicable, is called out in the documentation.

For example, consider the resource nsxt_policy_tier0_gateway. While the same resource works for both GM and LM, when using with GM, the attribute `locale_service` is required and is optional when connecting to LM.

Remember to check out the documentation for the resource you are interested in for such differences.

Another item to note is that only Policy resources are avaiable to use with Federation.

## Available Resources for use with NSX-T Federation

The following Resources are available to use with Federation:

 * Tier-0 Gateway [nsxt_policy_tier0_gateway](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_tier0_gateway.html)
 * Tier-0 Interface [nsxt_policy_tier0_interface](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_tier0_interface.html)
 * Tier-1 Gateway [nsxt_policy_tier1_gateway](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_tier1_gateway.html)
 * Tier-1 Interface [nsxt_policy_tier1_interface](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_tier1_interface.html)
 * Segment [nsxt_policy_segment](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_segment.html)
 * VLAN Segment [nsxt_policy_vlan_segment](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_vlan_segment.html)
 * Groups [nsxt_policy_group](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_group.html)
 * Services [nsxt_policy_service](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_service.html)
 * DFW Security Policy [nsxt_policy_security_policy](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_security_policy.html)
 * Gateway Policy [nsxt_policy_gateway_policy](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_gateway_policy.html)
 * NAT Rules [nsxt_policy_nat_rules](https://www.terraform.io/docs/providers/nsxt/r/nsxt_policy_nat_rules.html)

## Available Data Sources for use with NSX-T Federation

The following Data Sources are available to use with Federation:

 * Content Profile: [nsxt_policy_context_profile](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_contnext_profile.html)
 * Services: [nsxt_policy_service](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_service.html)
 * IP Discovery Profile: [nsxt_policy_ip_discovery_profile](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_ip_discovery_profile.html)
 * QOS Profile: [nsxt_policy_qos_profile](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_qos_profile.html)
 * Segment Security Profile: [nsxt_policy_segment_security_profile](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_segment_security_profile.html)
 * MAC Discovery Profile [nsxt_policy_mac_discovery_profile](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_mac_discovery_profile.html)
 * Federation Site [nsxt_policy_site](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_site.html)
 * Transport Zone [nsxt_policy_transport_zone](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_transport_zone.html)
 * Edge Cluster [nsxt_policy_edge_cluster](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_edge_cluster.html)
 * Gateway QoS Profile [nsxt_policy_gateway_qos_profile](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_gateway_qos_profile.html)
 * Edge Node [nsxt_policy_edge_node](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_edge_node.html)
 * Tier-0 Gateway [nsxt_policy_tier0_gateway](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_tier0_gateway.html)
 * Realization Info [nsxt_policy_realization_info](https://www.terraform.io/docs/providers/nsxt/d/nsxt_policy_realization_info.html)

