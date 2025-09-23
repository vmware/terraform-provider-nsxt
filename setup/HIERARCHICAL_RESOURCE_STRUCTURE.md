# NSX-T Terraform Resource Hierarchical Structure

## Overview
This document represents the hierarchical structure and dependencies of all resources in the main.tf configuration, showing the creation order and relationships.

## Resource Dependency Hierarchy

### Level 0: Foundation (External Dependencies)
```
External Infrastructure (Pre-existing)
├── data.nsxt_policy_edge_cluster.main_edge_cluster (EDGECLUSTER1)
├── data.nsxt_policy_edge_cluster.main_edge_cluster2 (EDGECLUSTER2)
└── data.nsxt_policy_tier0_gateway.original_tier0 (pepsi)
```

### Level 1: Core Infrastructure
```
Core Network Infrastructure
├── nsxt_policy_tier0_gateway.main_tier0 (pepsi-tgw-enabled)
│   ├── HA Mode: ACTIVE_STANDBY
│   ├── TGW Transit Subnets: 169.254.0.0/28
│   └── Edge Cluster: main_edge_cluster
├── nsxt_policy_ip_block.dev_external_block (203.0.114.0/24)
├── nsxt_policy_ip_block.prod_external_block (203.0.113.0/24)
└── nsxt_policy_ip_block.vlan_extension_block (192.168.100.0/24)
    ├── Visibility: EXTERNAL
    └── Subnet Exclusive: true
```

### Level 2: Resource Governance
```
IP Block Quotas
├── nsxt_policy_ip_block_quota.dev_quota
│   └── Depends on: dev_external_block
├── nsxt_policy_ip_block_quota.prod_quota
    └── Depends on: prod_external_block
```

### Level 3: Gateway Connections & VLAN Extensions
```
Gateway Connections
├── nsxt_policy_gateway_connection.dev_gw_connection
│   └── Depends on: main_tier0, dev_external_block
├── nsxt_policy_gateway_connection.prod_gw_connection
│   └── Depends on: main_tier0, prod_external_block
└── nsxt_policy_distributed_vlan_connection.legacy_vlan_connection
    ├── Gateway Addresses: 192.168.100.254/24
    ├── VLAN ID: 100
    ├── IP Block: vlan_extension_block
    ├── VPC Gateway Connection: Enabled
    └── Depends on: vlan_extension_block
```

### Level 4: Projects (Multi-tenancy)
```
Projects
├── nsxt_policy_project.dev_project
│   ├── Depends on: main_tier0, dev_external_block, dev_gw_connection
│   ├── TGW External Connections: dev_gw_connection, legacy_vlan_connection
│   └── Contains: dev private IP blocks
└── nsxt_policy_project.prod_project
    ├── Depends on: main_tier0, prod_external_block, prod_gw_connection
    ├── TGW External Connections: prod_gw_connection
    └── Contains: prod private IP blocks

Private IP Blocks (Project-scoped)
├── nsxt_policy_ip_block.dev_private_block
│   └── Context: dev_project
└── nsxt_policy_ip_block.prod_private_block
    └── Context: prod_project
```

### Level 5: Transit Gateways
```
Transit Gateways
├── nsxt_policy_transit_gateway.dev_tgw
│   ├── Context: dev_project
│   └── Edge Cluster: main_edge_cluster
└── nsxt_policy_transit_gateway.prod_tgw
    ├── Context: prod_project
    └── Edge Cluster: main_edge_cluster

Transit Gateway Attachments
├── nsxt_policy_transit_gateway_attachment.dev_tgw_attachment
│   ├── Parent: dev_tgw
│   └── Connection: dev_gw_connection
└── nsxt_policy_transit_gateway_attachment.prod_tgw_attachment
    ├── Parent: prod_tgw
    └── Connection: prod_gw_connection
```

### Level 6: IP Address Allocations
```
Project IP Allocations
├── nsxt_policy_project_ip_address_allocation.dev_nat_ip
│   └── Context: dev_project
├── nsxt_policy_project_ip_address_allocation.dev_web_dnat_ip
│   └── Context: dev_project
├── nsxt_policy_project_ip_address_allocation.prod_nat_ip
│   └── Context: prod_project
└── nsxt_policy_project_ip_address_allocation.prod_web_dnat_ip
    └── Context: prod_project
```

### Level 7: NAT Data Sources & Rules
```
Transit Gateway NAT Data Sources
├── data.nsxt_policy_transit_gateway_nat.dev_nat
│   └── Transit Gateway: dev_tgw
└── data.nsxt_policy_transit_gateway_nat.prod_nat
    └── Transit Gateway: prod_tgw

Transit Gateway NAT Rules
├── nsxt_policy_transit_gateway_nat_rule.dev_snat_rule
│   ├── Parent: dev_nat
│   ├── Source: 10.0.0.0/12
│   └── Translated: dev_nat_ip
├── nsxt_policy_transit_gateway_nat_rule.dev_web_dnat_rule
│   ├── Parent: dev_nat
│   ├── Destination: dev_web_dnat_ip
│   └── Translated: 10.1.1.10
├── nsxt_policy_transit_gateway_nat_rule.prod_snat_rule
│   ├── Parent: prod_nat
│   ├── Source: 10.32.0.0/12
│   └── Translated: prod_nat_ip
└── nsxt_policy_transit_gateway_nat_rule.prod_web_dnat_rule
    ├── Parent: prod_nat
    ├── Destination: prod_web_dnat_ip
    └── Translated: 10.33.1.10
```

### Level 8: Service Profiles
```
VPC Service Profiles
├── nsxt_vpc_service_profile.dev_service_profile
│   ├── Context: dev_project
│   ├── DHCP: Distributed, DNS: 8.8.8.8
│   └── NTP: 20.2.60.5
└── nsxt_vpc_service_profile.prod_service_profile
    ├── Context: prod_project
    ├── DHCP: Centralized, DNS: 1.1.1.1
    └── NTP: 20.2.60.3
```

### Level 9: Connectivity Profiles
```
VPC Connectivity Profiles
├── nsxt_vpc_connectivity_profile.dev_connectivity_profile
│   ├── Context: dev_project
│   ├── Transit Gateway: dev_tgw
│   └── External IP Blocks: dev_external_block
├── nsxt_vpc_connectivity_profile.prod_connectivity_profile
│   ├── Context: prod_project
│   ├── Transit Gateway: prod_tgw
│   └── External IP Blocks: prod_external_block
└── #nsxt_vpc_connectivity_profile.dmz_connectivity_profile
    ├── Context: dev_project
    └── Transit Gateway: dev_tgw
```

### Level 10: VPCs
```
Development VPCs
├── nsxt_vpc.dev_web_vpc
│   ├── Context: dev_project
│   ├── Private IPs: 192.168.10.0/24
│   ├── Service Profile: dev_service_profile
│   └── Depends on: dev_project, dev_service_profile, dev_tgw
├── nsxt_vpc.dev_db_vpc
│   ├── Context: dev_project
│   ├── Private IPs: 10.1.0.0/16
│   ├── Service Profile: dev_service_profile
│   └── Depends on: dev_project, dev_service_profile, dev_tgw
└── #nsxt_vpc.dev_dmz_vpc
    ├── Context: dev_project
    └── IP Address Type: IPv4

Production VPCs
├── nsxt_vpc.prod_web_vpc
│   ├── Context: prod_project
│   ├── Private IPs: 10.30.0.0/16
│   ├── Service Profile: prod_service_profile
│   └── Depends on: prod_project, prod_service_profile, prod_tgw
├── nsxt_vpc.prod_db_vpc
│   ├── Context: prod_project
│   ├── Private IPs: 10.40.0.0/16
│   ├── Service Profile: prod_service_profile
│   └── Depends on: prod_project, prod_service_profile, prod_tgw
```

### Level 11: VPC Attachments
```
VPC Attachments
├── nsxt_vpc_attachment.dev
│   ├── Parent: dev_web_vpc
│   └── Connectivity Profile: dev_connectivity_profile
├── nsxt_vpc_attachment.dev_db
│   ├── Parent: dev_db_vpc
│   └── Connectivity Profile: dev_connectivity_profile
├── nsxt_vpc_attachment.prod_web
│   ├── Parent: prod_web_vpc
│   └── Connectivity Profile: prod_connectivity_profile
└── nsxt_vpc_attachment.prod_db
    ├── Parent: prod_db_vpc
    └── Connectivity Profile: prod_connectivity_profile
```

### Level 12: VPC Subnets
```
Development Subnets
├── nsxt_vpc_subnet.dev_web_private
│   ├── Context: dev_project, dev_web_vpc
│   ├── Subnet Size: /32
│   └── Access Mode: Private
├── nsxt_vpc_subnet.dev_web_private_subnet
│   ├── Context: dev_project, dev_web_vpc
│   ├── Subnet Size: /64
│   ├── Access Mode: Private
│   └── DHCP: DHCP_SERVER with Option121
├── nsxt_vpc_subnet.dev_web_public_subnet
│   ├── Context: dev_project, dev_web_vpc
│   ├── Subnet Size: /32
│   ├── Access Mode: Public
│   └── Depends on: dev VPC attachment
└── nsxt_vpc_subnet.dev_db_isolated_subnet
    ├── Context: dev_project, dev_db_vpc
    ├── IP Addresses: 192.168.10.64/28
    ├── Access Mode: Isolated
    ├── DHCP: DHCP_SERVER with reserved ranges
    └── Connectivity: DISCONNECTED

Production Subnets
├── nsxt_vpc_subnet.prod_web_private_subnet
│   ├── Context: prod_project, prod_web_vpc
│   ├── Subnet Size: /128
│   ├── Access Mode: Private
│   └── DHCP: DHCP_DEACTIVATED
└── nsxt_vpc_subnet.prod_web_public_subnet
    ├── Context: prod_project, prod_web_vpc
    ├── Subnet Size: /64
    ├── Access Mode: Public
    ├── DHCP: DHCP_DEACTIVATED
    └── Depends on: prod_web VPC attachment
```

### Level 13: Security Groups
```
VPC Groups (Security Groups)
├── nsxt_policy_group.dev
│   ├── Context: dev_project
│   ├── Members: dev_web_vpc, dev_db_vpc
│   └── Depends on: dev VPCs
├── nsxt_policy_group.prod
│   ├── Context: prod_project
│   ├── Members: prod_web_vpc, prod_db_vpc
│   └── Depends on: prod VPCs
├── nsxt_vpc_group.dev_web_servers
│   ├── Context: dev_project, dev_web_vpc
│   └── Criteria: VMs starting with "dev-web"
├── nsxt_vpc_group.dev_db_servers
│   ├── Context: dev_project, dev_db_vpc
│   └── Criteria: VMs starting with "dev-db"
├── nsxt_vpc_group.prod_web_servers
│   ├── Context: prod_project, prod_web_vpc
│   └── Criteria: VMs starting with "prod-web"
└── nsxt_vpc_group.prod_db_servers
    ├── Context: prod_project, prod_db_vpc
    └── Criteria: VMs starting with "prod-db"
```

### Level 14: Security & Gateway Policies
```
Connectivity Policies
├── nsxt_policy_connectivity_policy.devapps
│   ├── Group: dev group
│   ├── Parent: dev_tgw
│   └── Scope: ISOLATED
└── nsxt_policy_connectivity_policy.prodapps
    ├── Group: prod group
    ├── Parent: prod_tgw
    └── Scope: COMMUNITY

VPC Security Policies
├── nsxt_vpc_security_policy.dev_web_security_policy
│   ├── Context: dev_project, dev_web_vpc
│   ├── Scope: dev_web_servers
│   └── Rules: Allow HTTP/HTTPS, Allow SSH
└── nsxt_vpc_security_policy.dev_db_security_policy
    ├── Context: dev_project, dev_db_vpc
    ├── Scope: dev_db_servers
    └── Rules: Allow DB access, Deny all other

VPC Gateway Policies
├── nsxt_vpc_gateway_policy.dev_web_gateway_policy
│   ├── Context: dev_project, dev_web_vpc
│   └── Rules: Allow inbound web, Allow outbound all
└── nsxt_vpc_gateway_policy.prod_web_gateway_policy
    ├── Context: prod_project, prod_web_vpc
    └── Rules: Allow inbound web, Allow controlled outbound
```

### Level 15: Network Services
```
VPC IP Allocations
└── nsxt_vpc_ip_address_allocation.dev_external_nat_pool
    ├── Context: dev_project, dev_web_vpc
    └── Allocation Size: 1

VPC NAT Data Sources & Rules
├── data.nsxt_vpc_nat.app_nat
│   └── Context: dev_project, dev_web_vpc
└── nsxt_vpc_nat_rule.outbound_snat
    ├── Parent: app_nat
    ├── Source: 192.168.10.0/24
    └── Translated: dev_external_nat_pool

Static Routes
├── nsxt_vpc_static_route.dev_web_static_route
│   ├── Context: dev_project, dev_web_vpc
│   └── Network: 10.10.0.0/16 → 10.10.1.1
├── nsxt_vpc_static_route.prod_web_static_route
│   ├── Context: prod_project, prod_web_vpc
│   └── Network: 10.30.0.0/16 → 10.30.1.1
├── nsxt_policy_transit_gateway_static_route.dev_tgw_static_route
│   ├── Parent: dev_tgw
│   └── Network: 192.168.0.0/16 → dev_tgw_attachment
└── nsxt_policy_transit_gateway_static_route.prod_tgw_static_route
    ├── Parent: prod_tgw
    └── Network: 192.168.0.0/16 → prod_tgw_attachment

DHCP Static Bindings
└── nsxt_vpc_dhcp_v4_static_binding.dev_web_server_binding
    ├── Parent: dev_db_isolated_subnet
    ├── MAC: 00:50:56:11:22:33
    ├── IP: 192.168.10.71
    └── Options: DNS, routing
```

### Level 16: VM Data Sources
```
VM References
└── data.nsxt_policy_vm.web_server
    ├── Context: dev_project
    └── Display Name: Web-VM-2
```

## Resource Creation Order Summary

1. **External Dependencies** (Pre-existing)
2. **Core Infrastructure** (Tier-0, IP blocks)
3. **Resource Governance** (IP quotas)
4. **Gateway Connections**
5. **Projects** (Multi-tenancy)
6. **Private IP Blocks** (Project-scoped)
7. **Transit Gateways & Attachments**
8. **IP Address Allocations**
9. **NAT Data Sources & Rules**
10. **Service Profiles**
11. **Connectivity Profiles**
12. **VPCs**
13. **VPC Attachments**
14. **VPC Subnets**
15. **Security Groups**
16. **Security & Gateway Policies**
17. **Network Services** (NAT, Routes, DHCP)
18. **VM References**

## Key Dependencies

- **Projects** require Tier-0 gateway, IP blocks, and gateway connections
- **Distributed VLAN connections** require dedicated subnet-exclusive IP blocks
- **VLAN connections** can only be assigned to one project at a time
- **VPCs** require projects, service profiles, and transit gateways
- **Subnets** require VPCs and (for public subnets) VPC attachments
- **Security policies** require VPC groups and VPCs
- **NAT rules** require IP allocations and NAT data sources
- **Static routes** require parent resources (VPCs or transit gateways)

Total Resources: **62+** across 16 hierarchical levels

## Legacy System Integration

The distributed VLAN connection provides hybrid connectivity for legacy system integration:
- **VLAN Extension**: Enables L2 connectivity to existing infrastructure
- **Project Exclusivity**: Each VLAN connection can only be used by one project
- **Subnet Exclusive IP Block**: Requires dedicated IP block with exclusive subnet configuration
- **VPC Gateway Connection**: Enables VPC-to-VLAN connectivity when needed
