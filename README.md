

# terraform-provider-nsxt
This is terraform provider for vmware NSX-T

## Overview

Supported data sources:

* transport_zone
* switching_profile
* logical_tier0_router
* edge_cluster

Supported resources:

* logical_switch
* logical_port
* logical_tier1_router
* logical_router_downlink_port
* logical_router_link_port_on_tier0
* logical_router_link_port_on_tier1
* l4_port_set_ns_service
* icmp_type_ns_service
* igmp_type_ns_service
* ether_type_ns_service
* alg_type_ns_service
* ip_protocol_ns_service
* ns_group
* firewall_section
* nat_rule
* ip_set
* dhcp_relay_profile
* dhcp_relay_service
* static_route


## Try it out

### Prerequisites

* Go 1.9.x onwards
* Terraform 0.10.x

### Build & Run

1. go get github.com/vmware/terraform-provider-nsxt
2. cd $GOROOT/src/github.com/vmware/terraform-provider-nsxt
3. make

## Contributing

The terraform-provider-nsxt project team welcomes contributions from the community. If you wish to contribute code and you have not signed our contributor license agreement (CLA), our bot will update the issue when you open a Pull Request. For any questions about the CLA process, please refer to our [FAQ](https://cla.vmware.com/faq). For more detailed information, refer to [CONTRIBUTING.md(CONTRIBUTING.md).

## License

This terraform provider is available under MPL2.0 license.
