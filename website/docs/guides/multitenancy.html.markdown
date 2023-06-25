---
layout: "nsxt"
page_title: "VMware NSX-T Terraform Provider support for Multi-tenancy"
description: |-
  The VMware NSX-T Terraform Provider support for Multi-tenancy
---

NSX-T Terraform Provider offers support for NSX-T [multi-tenancy feature](https://docs.vmware.com/en/VMware-NSX/4.1/administration/GUID-52180BC5-A1AB-4BC2-B1CE-666292505317.html).

# NSX-T Project Creation and Reference

NSX-T Project objects can be created and referenced with the `nsxt_policy_project` [resource](../r/policy_project.html.markdown) and [data source](../d/policy_project.html.markdown).

For example, an NSX-T multi-tenancy Project could be created as below:

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
  description  = "ip discovery profile provisioned by Terraform"
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

* [nsxt_policy_gateway_locale_service](../d/policy_gateway_locale_service.html.markdown)
* [nsxt_policy_ip_pool](../d/policy_ip_pool.html.markdown)
* [nsxt_policy_segment](../d/policy_segment.html.markdown)
* [nsxt_policy_context_profile](../d/policy_context_profile.html.markdown)
* [nsxt_policy_ip_discovery_profile](../d/policy_ip_discovery_profile.html.markdown)
* [nsxt_policy_ipv6_dad_profile](../d/policy_ipv6_dad_profile.html.markdown)
* [nsxt_policy_tier1_gateway](../d/policy_tier1_gateway.html.markdown)
* [nsxt_policy_spoofguard_profile](../d/policy_spoofguard_profile.html.markdown)
* [nsxt_policy_group](../d/policy_group.html.markdown)
* [nsxt_policy_mac_discovery_profile](../d/policy_mac_discovery_profile.html.markdown)
* [nsxt_policy_gateway_qos_profile](../d/policy_gateway_qos_profile.html.markdown)
* [nsxt_policy_ip_block](../d/policy_ip_block.html.markdown)
* [nsxt_policy_gateway_policy](../d/policy_gateway_policy.html.markdown)
* [nsxt_policy_dhcp_server](../d/policy_dhcp_server.html.markdown)
* [nsxt_policy_realization_info](../d/policy_realization_info.html.markdown)
* [nsxt_policy_segment_security_profile](../d/policy_segment_security_profile.html.markdown)
* [nsxt_policy_segment_realization](../d/policy_segment_realization.html.markdown)
* [nsxt_policy_qos_profile](../d/policy_qos_profile.html.markdown)
* [nsxt_policy_ipv6_ndra_profile](../d/policy_ipv6_ndra_profile.html.markdown)
* [nsxt_policy_service](../d/policy_service.html.markdown)
* [nsxt_policy_security_policy](../d/policy_security_policy.html.markdown)

# Supported resources

* [nsxt_policy_nat_rule](../r/policy_nat_rule.html.markdown)
* [nsxt_policy_dns_forwarder_zone](../r/policy_dns_forwarder_zone.html.markdown)
* [nsxt_policy_predefined_security_policy](../r/policy_predefined_security_policy.html.markdown)
* [nsxt_policy_static_route](../r/policy_static_route.html.markdown)
* [nsxt_policy_ip_pool](../r/policy_ip_pool.html.markdown)
* [nsxt_policy_segment](../r/policy_segment.html.markdown)
* [nsxt_policy_context_profile](../r/policy_context_profile.html.markdown)
* [nsxt_policy_ip_discovery_profile](../r/policy_ip_discovery_profile.html.markdown)
* [nsxt_policy_ip_pool_static_subnet](../r/policy_ip_pool_static_subnet.html.markdown)
* [nsxt_policy_predefined_gateway_policy](../r/policy_predefined_gateway_policy.html.markdown)
* [nsxt_policy_fixed_segment](../r/policy_fixed_segment.html.markdown)
* [nsxt_policy_tier1_gateway](../r/policy_tier1_gateway.html.markdown)
* [nsxt_policy_tier1_gateway_interface](../r/policy_tier1_gateway_interface.html.markdown)
* [nsxt_policy_group](../r/policy_group.html.markdown)
* [nsxt_policy_mac_discovery_profile](../r/policy_mac_discovery_profile.html.markdown)
* [nsxt_policy_ip_block](../r/policy_ip_block.html.markdown)
* [nsxt_policy_gateway_dns_forwarder](../r/policy_gateway_dns_forwarder.html.markdown)
* [nsxt_policy_gateway_policy](../r/policy_gateway_policy.html.markdown)
* [nsxt_policy_dhcp_server](../r/policy_dhcp_server.html.markdown)
* [nsxt_policy_ip_address_allocation](../r/policy_ip_address_allocation.html.markdown)
* [nsxt_policy_segment_security_profile](../r/policy_segment_security_profile.html.markdown)
* [nsxt_policy_context_profile_custom_attribute](../r/policy_context_profile_custom_attribute.html.markdown)
* [nsxt_policy_vm_tags](../r/policy_vm_tags.html.markdown)
* [nsxt_policy_dhcp_relay](../r/policy_dhcp_relay.html.markdown)
* [nsxt_policy_spoof_guard_profile](../r/policy_spoof_guard_profile.html.markdown)
* [nsxt_policy_qos_profile](../r/policy_qos_profile.html.markdown)
* [nsxt_policy_service](../r/policy_service.html.markdown)
* [nsxt_policy_ip_pool_block_subnet](../r/policy_ip_pool_block_subnet.html.markdown)
* [nsxt_policy_security_policy](../r/policy_security_policy.html.markdown)
