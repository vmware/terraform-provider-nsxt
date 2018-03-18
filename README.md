

# terraform-provider-nsxt
This repository contains the NSX Terraform provider for VMware NSX-T.

NSX is the VMware network and security virtualization platform. More more infomration on the NSX product please visit the [NSX Product Page](https://www.vmware.com/products/nsx.html)

Documenation on the NSX plaform can be found on the [NSX Documenatation page](https://docs.vmware.com/en/VMware-NSX-T/index.html)

# Overview

Supported data sources:

* edge_cluster
* logical_tier0_router
* switching_profile
* transport_zone

Supported resources:

* alg_type_ns_service
* dhcp_relay_profile
* dhcp_relay_service
* ether_type_ns_service
* firewall_section
* icmp_type_ns_service
* igmp_type_ns_service
* ip_protocol_ns_service
* ip_set
* logical_switch
* logical_port
* logical_tier1_router
* logical_router_downlink_port
* logical_router_link_port_on_tier0
* logical_router_link_port_on_tier1
* l4_port_set_ns_service
* ns_group
* nat_rule
* static_route

# Interoperability

The following versions of NSX-T are supported:

- NSX-T 2.1.*

# Prerequisites

The following are prerequisistes that need to be installed in order to run the NSX Terraform provider:

- Go 1.9.x onwards - [Go installation instructions](https://golang.org/doc/install)

- Terraform 0.10.x - [Terraform installation instructions](https://www.terraform.io/intro/getting-started/install.html)

Make sure that both terraform and go are in your path so they can be executed. To check versions of each you can run the following:

    go version

    terraform version  

Make sure you have the following directory:

    ~/.terraform.d/plugins/

# Installation

These commands will allow you to install the NSX Terraform provider:

    go get github.com/vmware/terraform-provider-nsxt

    cd $GOROOT/src/github.com/vmware/terraform-provider-nsxt

    go build -o ~/.terraform.d/plugins/terraform-provider-nsxt

# Contributing

The terraform-provider-nsxt project team welcomes contributions from the community. If you wish to contribute code and you have not signed our contributor license agreement (CLA), our bot will update the issue when you open a Pull Request. For any questions about the CLA process, please refer to our [FAQ](https://cla.vmware.com/faq). For more detailed information, refer to [CONTRIBUTING](https://github.com/vmware/terraform-provider-nsxt/blob/master/CONTRIBUTING.md).

# Support

The NSX Terraform provider is community supported. For bugs and feature requests please open a Github Issue and label it approprately. As this is a community supported solution there is no SLA for resolutions.

# License

Copyright © 2015-2018 VMware, Inc. All Rights Reserved.

The NSX Terraform provider is available under [MPL2.0 license](https://github.com/vmware/terraform-provider-nsxt/blob/master/LICENSE.txt).
