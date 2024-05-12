---
layout: "nsxt"
page_title: "VMware NSX Terraform Provider support for Multi-tenancy"
description: |-
  The VMware NSX Terraform Provider support for Multi-tenancy
---

# NSX Provider for NSX Multi-tenancy

NSX Terraform Provider offers support for NSX [multi-tenancy feature](https://docs.vmware.com/en/VMware-NSX/4.1/administration/GUID-52180BC5-A1AB-4BC2-B1CE-666292505317.html).

## NSX Multi-Tenancy
NSX introduced a new construct called Project in order to offer tenancy by isolating security and networking objects across tenants in a single NSX deployment.

Project allows multiple users to work on the platform in parallel, isolating configurations and defining scope for security (distributed security configured in a Project only applies to VMs connected to networks from that Project)

NSX also introduced NSX VPCs, offering a second level of tenancy and cloud consumption. Currently the Terraform Provider does not yet support NSX VPCs.

## NSX Project Creation and Reference

NSX Project objects can be created and referenced with the `nsxt_policy_project` [resource](../resources/policy_project.html.markdown) and [data source](../data-sources/policy_project.html.markdown).

For example, an NSX multi-tenancy Project could be created as below:

```hcl
resource "nsxt_policy_project" "test" {
  display_name        = "test"
  description         = "Terraform provisioned Project"
  short_id            = "test"
  tier0_gateway_paths = ["/infra/tier-0s/test"]
}
```

To support multi-tenancy, a context object has been introduced to the supported data sources and resources: context may refer to the Project ID in order to associate the object with a Project.
For example, to create a Project-associated IP Discovery Profile object:

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_ip_discovery_profile" "ip_discovery_profile" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  description  = "ip discovery profile for demo project"
  display_name = "ip_discovery_profile1"

  arp_nd_binding_timeout         = 20
  duplicate_ip_detection_enabled = false

  arp_binding_limit     = 140
  arp_snooping_enabled  = false
  dhcp_snooping_enabled = false
  vmtools_enabled       = false

  dhcp_snooping_v6_enabled = false
  nd_snooping_enabled      = false
  nd_snooping_limit        = 12
  vmtools_v6_enabled       = false
  tofu_enabled             = false

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

# Importing a Project resource

To import a resource which is associated with a Project, use the complete object policy path tp identify the imported object.
For example, to import a Tier1 gateway which is associated with Project demoproj, use:

```
terraform import nsxt_policy_tier1_gateway.tier1_gw /orgs/default/projects/demoproj/infra/tier-1s/t1gwdemo
```

# Supported data sources

* [nsxt_policy_gateway_locale_service](../data-sources/policy_gateway_locale_service.html.markdown)
* [nsxt_policy_ip_pool](../data-sources/policy_ip_pool.html.markdown)
* [nsxt_policy_segment](../data-sources/policy_segment.html.markdown)
* [nsxt_policy_context_profile](../data-sources/policy_context_profile.html.markdown)
* [nsxt_policy_ip_discovery_profile](../data-sources/policy_ip_discovery_profile.html.markdown)
* [nsxt_policy_ipv6_dad_profile](../data-sources/policy_ipv6_dad_profile.html.markdown)
* [nsxt_policy_tier1_gateway](../data-sources/policy_tier1_gateway.html.markdown)
* [nsxt_policy_spoofguard_profile](../data-sources/policy_spoofguard_profile.html.markdown)
* [nsxt_policy_group](../data-sources/policy_group.html.markdown)
* [nsxt_policy_mac_discovery_profile](../data-sources/policy_mac_discovery_profile.html.markdown)
* [nsxt_policy_gateway_qos_profile](../data-sources/policy_gateway_qos_profile.html.markdown)
* [nsxt_policy_ip_block](../data-sources/policy_ip_block.html.markdown)
* [nsxt_policy_gateway_policy](../data-sources/policy_gateway_policy.html.markdown)
* [nsxt_policy_dhcp_server](../data-sources/policy_dhcp_server.html.markdown)
* [nsxt_policy_realization_info](../data-sources/policy_realization_info.html.markdown)
* [nsxt_policy_segment_security_profile](../data-sources/policy_segment_security_profile.html.markdown)
* [nsxt_policy_segment_realization](../data-sources/policy_segment_realization.html.markdown)
* [nsxt_policy_qos_profile](../data-sources/policy_qos_profile.html.markdown)
* [nsxt_policy_ipv6_ndra_profile](../data-sources/policy_ipv6_ndra_profile.html.markdown)
* [nsxt_policy_service](../data-sources/policy_service.html.markdown)
* [nsxt_policy_security_policy](../data-sources/policy_security_policy.html.markdown)

# Supported resources

* [nsxt_policy_nat_rule](../resources/policy_nat_rule.html.markdown)
* [nsxt_policy_dns_forwarder_zone](../resources/policy_dns_forwarder_zone.html.markdown)
* [nsxt_policy_predefined_security_policy](../resources/policy_predefined_security_policy.html.markdown)
* [nsxt_policy_static_route](../resources/policy_static_route.html.markdown)
* [nsxt_policy_ip_pool](../resources/policy_ip_pool.html.markdown)
* [nsxt_policy_segment](../resources/policy_segment.html.markdown)
* [nsxt_policy_context_profile](../resources/policy_context_profile.html.markdown)
* [nsxt_policy_ip_discovery_profile](../resources/policy_ip_discovery_profile.html.markdown)
* [nsxt_policy_ip_pool_static_subnet](../resources/policy_ip_pool_static_subnet.html.markdown)
* [nsxt_policy_predefined_gateway_policy](../resources/policy_predefined_gateway_policy.html.markdown)
* [nsxt_policy_fixed_segment](../resources/policy_fixed_segment.html.markdown)
* [nsxt_policy_tier1_gateway](../resources/policy_tier1_gateway.html.markdown)
* [nsxt_policy_tier1_gateway_interface](../resources/policy_tier1_gateway_interface.html.markdown)
* [nsxt_policy_group](../resources/policy_group.html.markdown)
* [nsxt_policy_mac_discovery_profile](../resources/policy_mac_discovery_profile.html.markdown)
* [nsxt_policy_ip_block](../resources/policy_ip_block.html.markdown)
* [nsxt_policy_gateway_dns_forwarder](../resources/policy_gateway_dns_forwarder.html.markdown)
* [nsxt_policy_gateway_policy](../resources/policy_gateway_policy.html.markdown)
* [nsxt_policy_dhcp_server](../resources/policy_dhcp_server.html.markdown)
* [nsxt_policy_ip_address_allocation](../resources/policy_ip_address_allocation.html.markdown)
* [nsxt_policy_segment_security_profile](../resources/policy_segment_security_profile.html.markdown)
* [nsxt_policy_context_profile_custom_attribute](../resources/policy_context_profile_custom_attribute.html.markdown)
* [nsxt_policy_vm_tags](../resources/policy_vm_tags.html.markdown)
* [nsxt_policy_dhcp_relay](../resources/policy_dhcp_relay.html.markdown)
* [nsxt_policy_spoof_guard_profile](../resources/policy_spoof_guard_profile.html.markdown)
* [nsxt_policy_qos_profile](../resources/policy_qos_profile.html.markdown)
* [nsxt_policy_service](../resources/policy_service.html.markdown)
* [nsxt_policy_ip_pool_block_subnet](../resources/policy_ip_pool_block_subnet.html.markdown)
* [nsxt_policy_security_policy](../resources/policy_security_policy.html.markdown)
