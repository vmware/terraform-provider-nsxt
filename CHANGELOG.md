## 1.1.0 (Unreleased)

NEW DATA SOURCES:

* `nsxt_mac_pool`
* `nsxt_ns_group`
* `nsxt_ns_service`
* `nsxt_certificate`

NEW RESOURCES:

* `nsxt_dhcp_relay_profile`
* `nsxt_dhcp_relay_service`
* `nsxt_dhcp_server_profile`
* `nsxt_logical_dhcp_server`
* `nsxt_dhcp_server_ip_pool`
* `nsxt_vlan_logical_switch`
* `nsxt_logical_dhcp_port`
* `nsxt_logical_tier0_router`
* `nsxt_logical_router_centralized_service_port`
* `nsxt_ip_block`
* `nsxt_ip_block_subnet`
* `nsxt_ip_pool`
* `nsxt_ip_set`
* `nsxt_lb_icmp_monitor`
* `nsxt_lb_tcp_monitor`
* `nsxt_lb_udp_monitor`
* `nsxt_lb_http_monitor`
* `nsxt_lb_https_monitor`
* `nsxt_lb_passive_monitor`
* `nsxt_lb_pool`
* `nsxt_lb_tcp_virtual_server`
* `nsxt_lb_udp_virtual_server`
* `nsxt_lb_http_forwarding_rule`
* `nsxt_lb_http_request_rewrite_rule`
* `nsxt_lb_http_response_rewrite_rule`
* `nsxt_lb_cookie_persistence_profile`
* `nsxt_lb_source_ip_persistence_profile`
* `nsxt_lb_client_ssl_profile`
* `nsxt_lb_server_ssl_profile`
* `nsxt_lb_service`
* `nsxt_lb_fast_tcp_application_profile`
* `nsxt_lb_fast_udp_application_profile`
* `nsxt_lb_http_application_profile`

DEPRECATED ATTRIBUTES:

* Attribute `vlan` of `nsxt_logical_switch` is deprecated. Please use new resource `nsxt_vlan_logical_switch` to manage vlan based logical switches.

## 1.0.0 (April 09, 2018)

Initial release.
