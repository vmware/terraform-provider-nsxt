## 1.1.2 (Unreleased)

FEATURES:
* **New Data Source**: `nsxt_ip_pool` ([#109](https://github.com/terraform-providers/terraform-provider-nsxt/pull/190))
* **New Resource**: `nsxt_ip_pool_allocation_ip_address` ([#109](https://github.com/terraform-providers/terraform-provider-nsxt/pull/190))

IMPROVEMENTS:
* `resource/nsxt_ns_group`: Support IPSet type in membership criteria ([#195](https://github.com/terraform-providers/terraform-provider-nsxt/pull/195))

BUG FIXES:
* Fix refresh failures for most of resources. When resource was deleted on backend, the provider is expected to refresh state, discover resource absence and re-create it on next apply. Instead, the provider errored out ([#195]https://github.com/terraform-providers/terraform-provider-nsxt/pull/191))
* `resource/nsxt_ip_set`: Allow force-deletion of IPSet even if its referenced in ns groups.
* `resource/nsxt_logical_router_downlink_port`: Fix crash that happened during import with specific configuration ([#193](https://github.com/terraform-providers/terraform-provider-nsxt/pull/193))
* `resource/nsxt_logical_router_link_port_on_tier1`: Fix crash that happened during import with specific configuration ([#193](https://github.com/terraform-providers/terraform-provider-nsxt/pull/193))

## 1.1.1 (August 05, 2019)

NOTES:

* The provider is now aligned with Terraform 0.12 SDK which is required for Terraform 0.12 support. This version of terraform is more strict with syntax enforcement. If you old configuration errors out post upgrade, please verify syntax against the updated provider documentation.

IMPROVEMENTS:

* `resource/nsxt_vm_tag`: Support tagging of logical port for the VM ([#171](https://github.com/terraform-providers/terraform-provider-nsxt/pull/171))
* `resource/nsxt_firewall_section`: Add ability to control order of FW sections ([#150](https://github.com/terraform-providers/terraform-provider-nsxt/pull/150))
* `resource/nsxt_firewall_section`: Add support for LogicalRouter and LogicalRouterPort in as applied_to type ([#157](https://github.com/terraform-providers/terraform-provider-nsxt/pull/157))
* Introduce flag to tolerate partial_success realization state. This can be controlled by tolerate_partial_success provider attribute or NSXT_TOLERATE_PARTIAL_SUCCESS environment variable. The default is False ([#181](https://github.com/terraform-providers/terraform-provider-nsxt/pull/181))
* Add Go Modules support ([#155](https://github.com/terraform-providers/terraform-provider-nsxt/pull/155))
* Fix syntax in documentation and tests according to terraform 0.12 requirements ([#178](https://github.com/terraform-providers/terraform-provider-nsxt/pull/178))
* Verify interoperability with NSX 2.5
* Improve documentation and test coverage


BUG FIXES:

* `resource/nsxt_nat_rule`: Fix deletion of NAT rule that was due to a platform bug in versions 2.4 and below ([#166](https://github.com/terraform-providers/terraform-provider-nsxt/pull/166)).
* `resource/nsxt_firewall_section`: Do not enforce order of services in rules. This fixes the bug of non-empty plan when services were registered on backend in order different that defined in terraform ([#156](https://github.com/terraform-providers/terraform-provider-nsxt/pull/156))
* `resource/nsxt_firewall_section`: Prevent re-creation of rules by retaining rule ids ([#154](https://github.com/terraform-providers/terraform-provider-nsxt/pull/154))
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
