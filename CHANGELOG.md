## 3.1.1 (January 4, 2021)

FEATURES:

* **New Resource**: `nsxt_policy_dhcp_server`.
* **New Resource**: `nsxt_policy_domain` (Global Manager only).
* **New Resource**: `nsxt_policy_dhcp_v4_static_binding`.
* **New Resource**: `nsxt_policy_dhcp_v6_static_binding`.

EXPERIMENTAL FEATURES:

* **New Data Source**: `nsxt_policy_bfd_profile`.

* **New Resource**: `nsxt_policy_dns_forwarder_zone`.
* **New Resource**: `nsxt_policy_gateway_dns_forwarder`.
* **New Resource**: `nsxt_policy_intrusion_service_policy`.
* **New Resource**: `nsxt_policy_gateway_community_list`.
* **New Resource**: `nsxt_policy_fixed_segment` (VMC only).

IMPROVEMENTS:

* New provider attributes `client_auth_cert`, `client_auth_key` to allow passing these values as string rather than a file ([#524](https://github.com/vmware/terraform-provider-nsxt/pull/524))
* Allow Bearer token authorization type for VMC deployments (Experimental). This behavior is configured by setting new provider attribute `vmc_auth_mode` to `Bearer` ([#539](https://github.com/vmware/terraform-provider-nsxt/pull/539))
* Complete Global Manager support for data sources (T1 Gateway, IPv6 Profiles, Ceritificate)
* `resource/nsxt_policy_tier1_gateway`: Enhance T0 Gateway resource with `rd_admin_address` attribute ([#503](https://github.com/vmware/terraform-provider-nsxt/pull/503))
* `resource/nsxt_policy_predefined_gateway_policy`: Add Importer for this resource to match user expectations in case predefined rules exist. Documentation was also extended to cover import and no-import usage ([#527](https://github.com/vmware/terraform-provider-nsxt/pull/527))


BUG FIXES:
* Allow maximum subnet length in Gateway Interface validation ([#528](https://github.com/vmware/terraform-provider-nsxt/pull/528))
* Make sure policy data sources ignore deleted objects ([#516](https://github.com/vmware/terraform-provider-nsxt/pull/516))
* `resource/nsxt_policy_segment`: Allow configuration of segment on Global Manager without transport zone ([#513](https://github.com/vmware/terraform-provider-nsxt/pull/513)).
* Determine major NSX version behind VMC deployment, thus making 3.0.0 features (such as segment DHCP) available for VMC. This requires a more robust solution in futire ([#531](https://github.com/vmware/terraform-provider-nsxt/pull/531)).

## 3.1.0 (October 20, 2020)

FEATURES:

* **New Data Source**: `nsxt_policy_security_policy`.
* **New Data Source**: `nsxt_policy_gateway_policy`.
* **New Data Source**: `nsxt_policy_group`.
* **New Data Source**: `nsxt_policy_context_profile` (official support).

* **New Resource**: `nsxt_policy_context_profile` (official support).
* **New Resource**: `nsxt_policy_tier0_gateway_ha_vip_config` (official support).
* **New Resource**: `nsxt_policy_gateway_prefix_list` (official support).

EXPERIMENTAL FEATURES:

* **New Data Source**: `nsxt_management_cluster`.

* **New Resource**: `nsxt_policy_predefined_security_policy`. This resource allows users to modify default security policy. Please refer to docs for more details.
* **New Resource**: `nsxt_policy_predefined_gateway_policy`. This resource enables gateway policy configuration for VMC. Please refer to docs for more details.

IMPROVEMENTS:

* Allow specifying `vlan_ids` for overlay segments ([#462](https://github.com/vmware/terraform-provider-nsxt/pull/462))
* Allow specifying NSX license via provider attribute. Note: the lisence is not considered part of configuration, and is applied at plan time! ([#423](https://github.com/vmware/terraform-provider-nsxt/pull/423))

BUG FIXES:

* `resource/nsxt_policy_tier0_gateway`: Fix non-empty state issue for VRF use case ([#478](https://github.com/vmware/terraform-provider-nsxt/pull/478))
* `resource/nsxt_policy_segment`: Fix a bug with `excluded_range` assignment ([#473](https://github.com/vmware/terraform-provider-nsxt/pull/473))
* `resource/nsxt_policy_lb_pool`: Fix read function for `member_group` attribute ([#473](https://github.com/vmware/terraform-provider-nsxt/pull/473))
* `resource/nsxt_policy_ip_address_allocation`: Fix address allocation with older NSX versions ([#468](https://github.com/vmware/terraform-provider-nsxt/pull/468))
* `data/nsxt_policy_realization_info`: Fix realization polling with older NSX versions ([#468](https://github.com/vmware/terraform-provider-nsxt/pull/468))
* `data/nsxt_ns_group`: Add pagination support to fix group retrieval with many group objects defined ([#440](https://github.com/vmware/terraform-provider-nsxt/pull/440))
* `resource/nsxt_policy_lb_virtual_server`: Preserve existing rules that are defined outside terraform ([#482](https://github.com/vmware/terraform-provider-nsxt/pull/482))

## 3.0.0 (August 24, 2020)

* The provider is extended to support NSXT Global Manager. Only a subset of objects is supported, check the documentation for more details.

* **New Data Source**: `nsxt_policy_site`. Applicable for NSX Global Manager only.

* **New Resource**: `nsxt_policy_bgp_config`. Applicable for NSX Global Manager only.

EXPERIMENTAL FEATURES:

* **New Data Source**: `nsxt_policy_context_profile`.

* **New Resource**: `nsxt_policy_context_profile`.
* **New Resource**: `nsxt_policy_tier0_gateway_ha_vip_config`.
* **New Resource**: `nsxt_policy_gateway_prefix_list`.

IMPROVEMENTS:

* Improve error handling for policy resources. This fixes some scenarios (mostly relevant for VMConAWS) where error was swallowed by the provider ([#428] (https://github.com/vmware/terraform-provider-nsxt/pull/428))
* Improve provider host validation and allow schema to be specified ([#413](https://github.com/vmware/terraform-provider-nsxt/pull/413))
* `resource/nsxt_policy_vm_tags`: Support tagging specific logical port on the VM, based on segment path ([#406](https://github.com/vmware/terraform-provider-nsxt/pull/406))
* `resource/nsxt_policy_group`: Support MAC address criteria ([#388](https://github.com/vmware/terraform-provider-nsxt/pull/388))
* `resource/nsxt_policy_segment`, `resource/nsxt_policy_vlan_segment`: Support assigning custome segment profiles ([#384](https://github.com/vmware/terraform-provider-nsxt/pull/384))
* `resource/nsxt_policy_segment`, `resource/nsxt_policy_vlan_segment`: Wait for VM ports to be deleted before proceeding with segment delete. This avoids potential dependency error on deletion ([#311](https://github.com/vmware/terraform-provider-nsxt/pull/311))
* `resource/nsxt_policy_vlan_segment`: Allow specifying vlan range ([#342](https://github.com/vmware/terraform-provider-nsxt/pull/342))
* `resource/nsxt_policy_tier0_gateway`: Support assigning custom segment profiles ([#363](https://github.com/vmware/terraform-provider-nsxt/pull/363))

BUG FIXES:

* Fix to bypass certificate validation against cert request ([#381](https://github.com/vmware/terraform-provider-nsxt/pull/381))
* Fix potential crashes in some policy resources ([#305](https://github.com/vmware/terraform-provider-nsxt/pull/305))
* `resource/nsxt_policy_segment`: Fix error reporting on segment deletion ([#321](https://github.com/vmware/terraform-provider-nsxt/pull/321))
* `resource/nsxt_policy_vlan_segment`: Allow to specify zero as vlan id ([#304](https://github.com/vmware/terraform-provider-nsxt/pull/304))
* `resource/nsxt_policy_bgp_neighbor`: Fix route filters configuration ([#387](https://github.com/vmware/terraform-provider-nsxt/pull/387))
* `resource/nsxt_ip_pool_allocation_ip_address`: Fix import ([#319](https://github.com/vmware/terraform-provider-nsxt/pull/319))

## 2.1.0 (June 09, 2020)

* The provider is extended to support NSXT on VMConAWS. Only a subset of objects is supported, check the documentation for more details.

BUG FIXES:

* Fix remote authentication(vIDM) for policy objects. This fix is relevant for NSX version below 3.0.0. ([#302](https://github.com/vmware/terraform-provider-nsxt/pull/302))
* Fix client certificate authentication for policy objects ([#292](https://github.com/vmware/terraform-provider-nsxt/pull/292))
* Fix an issue related to non-admin NSX credentials ([#293](https://github.com/vmware/terraform-provider-nsxt/pull/293))
* `resource/nsxt_policy_vlan_segment`: Allow to specify vlan range ([#342](https://github.com/vmware/terraform-provider-nsxt/pull/342))
* `resource/nsxt_policy_segment`: Fix handling of segment deletion error ([#321](https://github.com/vmware/terraform-provider-nsxt/pull/321))
* `resource/nsxt_policy_segment`: Wait for potential VMs to free segment port before deleting the segment. ([#311](https://github.com/vmware/terraform-provider-nsxt/pull/311))
* `resource/nsxt_policy_vlan_segment`: Allow zero vlan ID ([#297](https://github.com/vmware/terraform-provider-nsxt/pull/297))
* `resource/nsxt_policy_tierX_gateway_interface`: Fix a use case of preconfigured locale service on gateway ([#300](https://github.com/vmware/terraform-provider-nsxt/pull/300))
* `resource/nsxt_policy_security_policy`: Fix import crash ([#299](https://github.com/vmware/terraform-provider-nsxt/pull/299))
* `resource/nsxt_policy_security_policy`: Expose `log_label` argument ([#298](https://github.com/vmware/terraform-provider-nsxt/pull/298))
* `resource/nsxt_policy_group`: Fix issues with group subresource import ([#288](https://github.com/vmware/terraform-provider-nsxt/pull/288))
* `resource/nsxt_policy_nat_rule`: Make `source_networks` argument optional ([#294](https://github.com/vmware/terraform-provider-nsxt/pull/294))
* `resource/nsxt_ip_pool_allocation_ip_address`: Fix resource import ([#319](https://github.com/vmware/terraform-provider-nsxt/pull/319))
* `data/nsxt_policy_segment_realization`: Expose computed attribute network_name. This attribute can be used as network name in vsphere provider, which forms the necessary dependency ([#308](https://github.com/vmware/terraform-provider-nsxt/pull/308))

## 2.0.0 (March 30, 2020)

NOTES:

* The provider is extended to support NSX-T policy API. Policy API is intended to be primary consumtion for NSX-T logical constructs, thus users are encouraged to use new data sources/resources, with _policy_ in the name.

FEATURES:
* **New Data Source**: `nsxt_policy_certificate`
* **New Data Source**: `nsxt_policy_edge_cluster`
* **New Data Source**: `nsxt_policy_edge_node`
* **New Data Source**: `nsxt_policy_tier0_gateway`
* **New Data Source**: `nsxt_policy_tier1_gateway`
* **New Data Source**: `nsxt_policy_segment`
* **New Data Source**: `nsxt_policy_vlan_segment`
* **New Data Source**: `nsxt_policy_service`
* **New Data Source**: `nsxt_policy_ip_discovery_profile`
* **New Data Source**: `nsxt_policy_spoofguard_profile`
* **New Data Source**: `nsxt_policy_qos_profile`
* **New Data Source**: `nsxt_policy_segment_security_profile`
* **New Data Source**: `nsxt_policy_mac_discovery_profile`
* **New Data Source**: `nsxt_policy_ipv6_ndra_profile`
* **New Data Source**: `nsxt_policy_ipv6_dad_profile`
* **New Data Source**: `nsxt_policy_vm`
* **New Data Source**: `nsxt_policy_lb_app_profile`
* **New Data Source**: `nsxt_policy_lb_client_ssl_profile`
* **New Data Source**: `nsxt_policy_lb_server_ssl_profile`
* **New Data Source**: `nsxt_policy_lb_monitor`
* **New Data Source**: `nsxt_policy_lb_persistence_profile`
* **New Data Source**: `nsxt_policy_vni_pool`
* **New Data Source**: `nsxt_policy_realization_info`
* **New Data Source**: `nsxt_policy_segment_realization`
* **New Data Source**: `nsxt_firewall_section`

* **New Resource**: `nsxt_policy_tier0_gateway`
* **New Resource**: `nsxt_policy_tier1_gateway`
* **New Resource**: `nsxt_policy_tier0_gateway_interface`
* **New Resource**: `nsxt_policy_tier1_gateway_interface`
* **New Resource**: `nsxt_policy_group`
* **New Resource**: `nsxt_policy_service`
* **New Resource**: `nsxt_policy_security_policy`
* **New Resource**: `nsxt_policy_gateway_policy`
* **New Resource**: `nsxt_policy_segment`
* **New Resource**: `nsxt_policy_vlan_segment`
* **New Resource**: `nsxt_policy_static_route`
* **New Resource**: `nsxt_policy_nat_rule`
* **New Resource**: `nsxt_policy_vm_tags`
* **New Resource**: `nsxt_policy_ip_block`
* **New Resource**: `nsxt_policy_ip_pool`
* **New Resource**: `nsxt_policy_ip_pool_block_subnet`
* **New Resource**: `nsxt_policy_ip_pool_static_subnet`
* **New Resource**: `nsxt_policy_ip_address_allocation`
* **New Resource**: `nsxt_policy_lb_pool`
* **New Resource**: `nsxt_policy_lb_service`
* **New Resource**: `nsxt_policy_lb_virtual_server`
* **New Resource**: `nsxt_policy_bgp_neighbor`
* **New Resource**: `nsxt_policy_dhcp_relay`
* **New Resource**: `nsxt_policy_dhcp_server`

IMPROVEMENTS:
* Migrate to Terraform Plugin SDK ([#210](https://github.com/vmware/terraform-provider-nsxt/pull/210))
* `resource/nsxt_vm_tags`: Avoid backend calls if no change required in corresponding tags ([#261](https://github.com/vmware/terraform-provider-nsxt/pull/261))

BUG FIXES:
* Fix client authentication error that used to occur when client certificate is not self signed ([#207](https://github.com/vmware/terraform-provider-nsxt/pull/207))
* Allow IPv6 in IP addresses and CIDR validations ([#204](https://github.com/vmware/terraform-provider-nsxt/pull/204))
* `resource/nsxt_vm_tags`: Fix tag removal ([#240](https://github.com/vmware/terraform-provider-nsxt/pull/240))
* `resource/nsxt_vm_tags`: Apply tags to all logical ports on given vm ([#235](https://github.com/vmware/terraform-provider-nsxt/pull/235))
* `resource/nsxt_logical_dhcp_server`: Mark gateway_ip as optional rather than required ([#245](https://github.com/vmware/terraform-provider-nsxt/pull/245))


## 1.1.2 (November 18, 2019)

FEATURES:
* **New Data Source**: `nsxt_ip_pool` ([#190](https://github.com/vmware/terraform-provider-nsxt/pull/190))
* **New Resource**: `nsxt_ip_pool_allocation_ip_address` ([#190](https://github.com/vmware/terraform-provider-nsxt/pull/190))

IMPROVEMENTS:
* `resource/nsxt_ns_group`: Support IPSet type in membership criteria ([#195](https://github.com/vmware/terraform-provider-nsxt/pull/195))

BUG FIXES:
* Fix refresh failures for most of resources. When resource was deleted on backend, the provider is expected to refresh state, discover resource absence and re-create it on next apply. Instead, the provider errored out ([#195]https://github.com/vmware/terraform-provider-nsxt/pull/191))
* `resource/nsxt_ip_set`: Allow force-deletion of IPSet even if its referenced in ns groups.
* `resource/nsxt_logical_router_downlink_port`: Fix crash that happened during import with specific configuration ([#193](https://github.com/vmware/terraform-provider-nsxt/pull/193))
* `resource/nsxt_logical_router_link_port_on_tier1`: Fix crash that happened during import with specific configuration ([#193](https://github.com/vmware/terraform-provider-nsxt/pull/193))
* `resource/nsxt_*_switching_profile`: Fix update error that occured in some cases due to omitted revision ([#201](https://github.com/vmware/terraform-provider-nsxt/pull/201))
* `resource/nsxt_logical_switch`: On delete operation, detach logical switch in order to avoid possible dependency errors ([#202](https://github.com/vmware/terraform-provider-nsxt/pull/202))

## 1.1.1 (August 05, 2019)

NOTES:

* The provider is now aligned with Terraform 0.12 SDK which is required for Terraform 0.12 support. This version of terraform is more strict with syntax enforcement. If you old configuration errors out post upgrade, please verify syntax against the updated provider documentation.

IMPROVEMENTS:

* `resource/nsxt_vm_tag`: Support tagging of logical port for the VM ([#171](https://github.com/vmware/terraform-provider-nsxt/pull/171))
* `resource/nsxt_firewall_section`: Add ability to control order of FW sections ([#150](https://github.com/vmware/terraform-provider-nsxt/pull/150))
* `resource/nsxt_firewall_section`: Add support for LogicalRouter and LogicalRouterPort in as applied_to type ([#157](https://github.com/vmware/terraform-provider-nsxt/pull/157))
* Introduce flag to tolerate partial_success realization state. This can be controlled by tolerate_partial_success provider attribute or NSXT_TOLERATE_PARTIAL_SUCCESS environment variable. The default is False ([#181](https://github.com/vmware/terraform-provider-nsxt/pull/181))
* Add Go Modules support ([#155](https://github.com/vmware/terraform-provider-nsxt/pull/155))
* Fix syntax in documentation and tests according to terraform 0.12 requirements ([#178](https://github.com/vmware/terraform-provider-nsxt/pull/178))
* Verify interoperability with NSX 2.5
* Improve documentation and test coverage


BUG FIXES:

* `resource/nsxt_nat_rule`: Fix deletion of NAT rule that was due to a platform bug in versions 2.4 and below ([#166](https://github.com/vmware/terraform-provider-nsxt/pull/166)).
* `resource/nsxt_firewall_section`: Do not enforce order of services in rules. This fixes the bug of non-empty plan when services were registered on backend in order different that defined in terraform ([#156](https://github.com/vmware/terraform-provider-nsxt/pull/156))
* `resource/nsxt_firewall_section`: Prevent re-creation of rules by retaining rule ids ([#154](https://github.com/vmware/terraform-provider-nsxt/pull/154))
* `resource/nsxt_nat_rule`: Allow setting rule_priority ([#182](https://github.com/te    rraform-providers/terraform-provider-nsxt/pull/182))

## 1.1.0 (February 22, 2019)

NOTES:

* resource/nsxt_logical_switch: Attribute `vlan` is deprecated. Please use new resource `nsxt_vlan_logical_switch` to manage vlan based logical switches.

FEATURES:

* **New Data Source**: `nsxt_mac_pool`
* **New Data Source**: `nsxt_ns_group`
* **New Data Source**: `nsxt_ns_service`
* **New Data Source**: `nsxt_certificate`
* **New Resource**: `nsxt_dhcp_relay_profile`
* **New Resource**: `nsxt_dhcp_relay_service`
* **New Resource**: `nsxt_dhcp_server_profile`
* **New Resource**: `nsxt_logical_dhcp_server`
* **New Resource**: `nsxt_dhcp_server_ip_pool`
* **New Resource**: `nsxt_vlan_logical_switch`
* **New Resource**: `nsxt_logical_dhcp_port`
* **New Resource**: `nsxt_logical_tier0_router`
* **New Resource**: `nsxt_logical_router_centralized_service_port`
* **New Resource**: `nsxt_ip_block`
* **New Resource**: `nsxt_ip_block_subnet`
* **New Resource**: `nsxt_ip_pool`
* **New Resource**: `nsxt_ip_set`
* **New Resource**: `nsxt_lb_icmp_monitor`
* **New Resource**: `nsxt_lb_tcp_monitor`
* **New Resource**: `nsxt_lb_udp_monitor`
* **New Resource**: `nsxt_lb_http_monitor`
* **New Resource**: `nsxt_lb_https_monitor`
* **New Resource**: `nsxt_lb_passive_monitor`
* **New Resource**: `nsxt_lb_pool`
* **New Resource**: `nsxt_lb_tcp_virtual_server`
* **New Resource**: `nsxt_lb_udp_virtual_server`
* **New Resource**: `nsxt_lb_http_forwarding_rule`
* **New Resource**: `nsxt_lb_http_request_rewrite_rule`
* **New Resource**: `nsxt_lb_http_response_rewrite_rule`
* **New Resource**: `nsxt_lb_cookie_persistence_profile`
* **New Resource**: `nsxt_lb_source_ip_persistence_profile`
* **New Resource**: `nsxt_lb_client_ssl_profile`
* **New Resource**: `nsxt_lb_server_ssl_profile`
* **New Resource**: `nsxt_lb_service`
* **New Resource**: `nsxt_lb_fast_tcp_application_profile`
* **New Resource**: `nsxt_lb_fast_udp_application_profile`
* **New Resource**: `nsxt_lb_http_application_profile`

## 1.0.0 (April 09, 2018)

Initial release.
