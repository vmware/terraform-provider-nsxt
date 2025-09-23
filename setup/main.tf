terraform {
  required_providers {
    nsxt = {
      source  = "vmware/nsxtlocal"
      version = ">=9.0.0"
    }
  }
}

provider "nsxt" {
  host                 = "10.161.208.216"
  username             = "admin"
  password             = "D5-AcRWv#XD9BPVj"
  allow_unverified_ssl = true
  max_retries          = 2
}

# =============================================================================
# MULTITENANT ENTERPRISE NSX-T TOPOLOGY WITH TRANSIT GATEWAY
# =============================================================================
#
# ARCHITECTURE PURPOSE:
# Enterprise NSX-T VPC multi-tenant topology implementing:
# - Dual-project isolation (Development & Production environments)
# - Transit Gateway-based VPC interconnectivity with comprehensive NAT services
# - Multi-tier application architecture (Web & Database VPCs per environment)
# - Advanced security framework (VPC groups, security policies, connectivity policies)
# - External connectivity through service gateways and IP address allocations
# - Comprehensive networking services (DHCP configurations, NAT rules, static bindings)
# - Resource governance through IP block quotas and IPv6 readiness
#
# TOPOLOGY STRUCTURE:
# Infrastructure Foundation
# ├── Tier-0 Gateway (pepsi-tgw-enabled) with TGW transit subnets (169.254.0.0/28)
# ├── IP Blocks: Dev External(203.0.114.0/24), Prod External(203.0.113.0/24), VLAN Extension(192.168.100.0/24), IPv6(2001:db8::/64)
# ├── Private IP Blocks: Dev(10.0.0.0/12, 10.16.0.0/12), Prod(10.32.0.0/12, 10.48.0.0/12)
# ├── IP Block Quotas: Dev(10 IPs/5 subnets), Prod(20 IPs/10 subnets), IPv6(10 /64 subnets)
# └── Gateway Connections & TGW Attachments for external connectivity
# 
# Multi-Tenant Projects
# ├── Development Project
# │   ├── Transit Gateway (ACTIVE_STANDBY) with SNAT/DNAT NAT rules
# │   ├── VPCs: dev-web-vpc(192.168.10.0/24), dev-db-vpc(10.1.0.0/16), dev-dmz-vpc
# │   ├── Service Profile: Distributed DHCP, Google DNS(8.8.8.8), NTP(20.2.60.5)
# │   ├── Connectivity Profile with external IP blocks & service gateway
# │   ├── VPC Attachments for external connectivity
# │   ├── Subnets: Private(/32), Public(/32), Isolated(/28 with DHCP server & reserved ranges)
# │   ├── DHCP Static Bindings with reserved IP ranges for isolated subnets
# │   ├── VPC Groups: dev_apps (path-based criteria)
# │   ├── Connectivity Policies: devApps (ISOLATED scope)
# │   └── VPC NAT rules & IP address allocations
# └── Production Project
#     ├── Transit Gateway (ACTIVE_STANDBY) with SNAT/DNAT NAT rules
#     ├── VPCs: prod-web-vpc(10.30.0.0/16), prod-db-vpc(10.40.0.0/16), prod-dual-stack-vpc
#     ├── Service Profile: Centralized DHCP, Cloudflare DNS(1.1.1.1), NTP(20.2.60.3)
#     ├── Connectivity Profile with external IP blocks & service gateway
#     ├── VPC Attachments for external connectivity
#     ├── Subnets: Private(/128), Public(/64) with DHCP deactivated
#     ├── VPC Groups: prod_apps (path-based criteria)
#     ├── Connectivity Policies: prodApps (COMMUNITY scope)
#     └── DMZ Connectivity Profile with custom edge cluster configuration
#
# =============================================================================

# Data Sources for existing infrastructure
data "nsxt_policy_edge_cluster" "main_edge_cluster" {
  display_name = "EDGECLUSTER1"
}

data "nsxt_policy_edge_cluster" "main_edge_cluster2" {
  display_name = "EDGECLUSTER2"
}
# Create a new Tier-0 gateway with TGW transit subnets support
resource "nsxt_policy_tier0_gateway" "main_tier0" {
  display_name             = "pepsi-tgw-enabled"
  description              = "Enterprise Tier-0 gateway providing external connectivity and transit gateway integration for multi-tenant VPC environments"
  ha_mode                  = "ACTIVE_STANDBY"
  failover_mode            = "NON_PREEMPTIVE"
  enable_firewall          = true
  edge_cluster_path        = data.nsxt_policy_edge_cluster.main_edge_cluster.path
  
  # Required for TGW attachments in NSX 9.1.0+
  tgw_transit_subnets      = ["169.254.0.0/28"]
}

# Keep the data source as backup reference
data "nsxt_policy_tier0_gateway" "original_tier0" {
  display_name = "pepsi"
}

# IP Blocks for different visibility types - Dev Project
resource "nsxt_policy_ip_block" "dev_external_block" {
  display_name = "dev-external-ip-block"
  description  = "Public IP address pool for development environment external services and NAT translations"
  cidr         = "203.0.114.0/24"
  visibility   = "EXTERNAL"
}

resource "nsxt_policy_ip_block" "dev_private_block" {
  context {
    project_id = nsxt_policy_project.dev_project.id
  }

  display_name = "dev-private-ip-block"
  description  = "Internal IP address space for development VPCs and private subnet allocations"
  cidrs        = ["10.0.0.0/12", "10.16.0.0/12"]
  visibility   = "PRIVATE"
}

# IP Blocks for different visibility types - Prod Project
resource "nsxt_policy_ip_block" "prod_external_block" {
  display_name = "prod-external-ip-block"
  description  = "Public IP address pool for production environment external services and customer-facing applications"
  cidr         = "203.0.113.0/24"
  visibility   = "EXTERNAL"
}

resource "nsxt_policy_ip_block" "prod_private_block" {
  context {
    project_id = nsxt_policy_project.prod_project.id
  }

  display_name = "prod-private-ip-block"
  description  = "Internal IP address space for production VPCs ensuring isolation from development networks"
  cidrs        = ["10.32.0.0/12", "10.48.0.0/12"]
  visibility   = "PRIVATE"
}


# IP Block Quotas for different projects
resource "nsxt_policy_ip_block_quota" "dev_quota" {
  display_name = "dev-project-quota"
  description  = "Resource limits controlling IP address consumption for development teams and applications"

  quota {
    ip_block_paths        = [nsxt_policy_ip_block.dev_external_block.path]
    ip_block_visibility   = "EXTERNAL"
    ip_block_address_type = "IPV4"
    single_ip_cidrs       = 10
    other_cidrs {
      total_count = 5
    }
  }
}

resource "nsxt_policy_ip_block_quota" "prod_quota" {
  display_name = "prod-project-quota"
  description  = "Resource limits ensuring controlled IP address allocation for production workloads and services"

  quota {
    ip_block_paths        = [nsxt_policy_ip_block.prod_external_block.path]
    ip_block_visibility   = "EXTERNAL"
    ip_block_address_type = "IPV4"
    single_ip_cidrs       = 20
    other_cidrs {
      total_count = 10
    }
  }
}

# Projects for multitenancy
resource "nsxt_policy_project" "dev_project" {
  display_name = "dev-project"
  description  = "Development project for multi-tenant isolation"
  tier0_gateway_paths = [nsxt_policy_tier0_gateway.main_tier0.path]
  external_ipv4_blocks = [nsxt_policy_ip_block.dev_external_block.path]
  tgw_external_connections = [nsxt_policy_gateway_connection.dev_gw_connection.path]

  site_info {
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.main_edge_cluster.path]
  }
}

resource "nsxt_policy_project" "prod_project" {
  display_name = "prod-project"
  description  = "Production project for multi-tenant isolation"
  tier0_gateway_paths = [nsxt_policy_tier0_gateway.main_tier0.path]
  external_ipv4_blocks = [nsxt_policy_ip_block.prod_external_block.path]
  tgw_external_connections = [nsxt_policy_gateway_connection.prod_gw_connection.path]

  site_info {
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.main_edge_cluster.path]
  }
}

# Data sources for projects to reference external IP blocks
data "nsxt_policy_project" "dev_project" {
  id = nsxt_policy_project.dev_project.id
}

data "nsxt_policy_project" "prod_project" {
  id = nsxt_policy_project.prod_project.id
}

# IP Address Allocations for NAT rules
resource "nsxt_policy_project_ip_address_allocation" "dev_nat_ip" {
  context {
    project_id = nsxt_policy_project.dev_project.id
  }
  
  display_name = "dev-nat-ip-allocation"
  description  = "IP allocation for development NAT rules"
  ip_block     = data.nsxt_policy_project.dev_project.external_ipv4_blocks[0]
  allocation_size = 1
}

resource "nsxt_policy_project_ip_address_allocation" "dev_web_dnat_ip" {
  context {
    project_id = nsxt_policy_project.dev_project.id
  }
  
  display_name = "dev-web-dnat-ip-allocation"
  description  = "IP allocation for development web DNAT rule"
  ip_block     = data.nsxt_policy_project.dev_project.external_ipv4_blocks[0]
  allocation_size = 1
}

resource "nsxt_policy_project_ip_address_allocation" "prod_nat_ip" {
  context {
    project_id = nsxt_policy_project.prod_project.id
  }
  
  display_name = "prod-nat-ip-allocation"
  description  = "IP allocation for production NAT rules"
  ip_block     = data.nsxt_policy_project.prod_project.external_ipv4_blocks[0]
  allocation_size = 1
}

resource "nsxt_policy_project_ip_address_allocation" "prod_web_dnat_ip" {
  context {
    project_id = nsxt_policy_project.prod_project.id
  }
  
  display_name = "prod-web-dnat-ip-allocation"
  description  = "IP allocation for production web DNAT rule"
  ip_block     = data.nsxt_policy_project.prod_project.external_ipv4_blocks[0]
  allocation_size = 1
}


# Transit Gateway NAT Rules
resource "nsxt_policy_transit_gateway_nat_rule" "dev_snat_rule" {


  display_name = "dev-snat-rule"
  description  = "SNAT rule for development outbound traffic"
  parent_path  = data.nsxt_policy_transit_gateway_nat.dev_nat.path
  action       = "SNAT"
  
  source_network      = "10.0.0.0/12"
  translated_network  = nsxt_policy_project_ip_address_allocation.dev_nat_ip.allocation_ips
  scope              = [nsxt_policy_transit_gateway_attachment.dev_tgw_attachment.path]
  logging            = true
  enabled            = true
}

resource "nsxt_policy_transit_gateway_nat_rule" "dev_web_dnat_rule" {


  display_name = "dev-web-dnat-rule"
  description  = "DNAT rule for development web services"
  parent_path  = data.nsxt_policy_transit_gateway_nat.dev_nat.path
  action       = "DNAT"
  
  destination_network = nsxt_policy_project_ip_address_allocation.dev_web_dnat_ip.allocation_ips
  translated_network  = "10.1.1.10"
  scope              = [nsxt_policy_transit_gateway_attachment.dev_tgw_attachment.path]
  logging            = true
  enabled            = true
}

resource "nsxt_policy_transit_gateway_nat_rule" "prod_snat_rule" {

  display_name = "prod-snat-rule"
  description  = "SNAT rule for production outbound traffic"
  parent_path  = data.nsxt_policy_transit_gateway_nat.prod_nat.path
  action       = "SNAT"
  
  source_network      = "10.32.0.0/12"
  translated_network  = nsxt_policy_project_ip_address_allocation.prod_nat_ip.allocation_ips
  scope              = [nsxt_policy_transit_gateway_attachment.prod_tgw_attachment.path]
  logging            = true
  enabled            = true
}

resource "nsxt_policy_transit_gateway_nat_rule" "prod_web_dnat_rule" {

  display_name = "prod-web-dnat-rule"
  description  = "DNAT rule for production web services"
  parent_path  = data.nsxt_policy_transit_gateway_nat.prod_nat.path
  action       = "DNAT"
  
  destination_network = nsxt_policy_project_ip_address_allocation.prod_web_dnat_ip.allocation_ips
  translated_network  = "10.33.1.10"
  scope              = [nsxt_policy_transit_gateway_attachment.prod_tgw_attachment.path]
  logging            = true
  enabled            = true
}

# Transit Gateways for each project
resource "nsxt_policy_transit_gateway" "dev_tgw" {
  context {
    project_id = nsxt_policy_project.dev_project.id
  }

  display_name = "dev-transit-gateway"
  description  = "Development transit gateway for VPC connectivity"

  high_availability_config {
    ha_mode            = "ACTIVE_STANDBY"
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.main_edge_cluster.path]
  }
}

resource "nsxt_policy_transit_gateway" "prod_tgw" {
  context {
    project_id = nsxt_policy_project.prod_project.id
  }

  display_name = "prod-transit-gateway"
  description  = "Production transit gateway for VPC connectivity"

  high_availability_config {
    ha_mode            = "ACTIVE_STANDBY"
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.main_edge_cluster.path]
  }
}

# Gateway Connections - moved to project context for proper scoping
resource "nsxt_policy_gateway_connection" "dev_gw_connection" {
  display_name = "dev-gateway-connection"
  description  = "Gateway connection for development project"
  tier0_path   = nsxt_policy_tier0_gateway.main_tier0.path
  advertise_outbound_networks {
    allow_external_blocks = [nsxt_policy_ip_block.dev_external_block.path]
  }
  depends_on = [ nsxt_policy_ip_block.dev_external_block,nsxt_policy_tier0_gateway.main_tier0]
}

resource "nsxt_policy_gateway_connection" "prod_gw_connection" {
  display_name = "prod-gateway-connection"
  description  = "Gateway connection for production project"
  tier0_path   = nsxt_policy_tier0_gateway.main_tier0.path
advertise_outbound_networks {
    allow_external_blocks = [nsxt_policy_ip_block.prod_external_block.path]
  }
}

# Transit Gateway Attachments
resource "nsxt_policy_transit_gateway_attachment" "dev_tgw_attachment" {
  display_name    = "dev-tgw-attachment"
  description     = "Transit gateway attachment for development"
  parent_path     = nsxt_policy_transit_gateway.dev_tgw.path
  connection_path = nsxt_policy_gateway_connection.dev_gw_connection.path
}

resource "nsxt_policy_transit_gateway_attachment" "prod_tgw_attachment" {
  display_name    = "prod-tgw-attachment"
  description     = "Transit gateway attachment for production"
  parent_path     = nsxt_policy_transit_gateway.prod_tgw.path
  connection_path = nsxt_policy_gateway_connection.prod_gw_connection.path
}

# VPC Service Profiles for different environments
resource "nsxt_vpc_service_profile" "dev_service_profile" {
  context {
    project_id = nsxt_policy_project.dev_project.id
  }

  display_name = "dev-vpc-service-profile"
  description  = "Service profile for development VPCs"

  dhcp_config {
    dhcp_server_config {
      ntp_servers = ["20.2.60.5"]
      lease_time  = 86400

      dns_client_config {
        dns_server_ips = ["8.8.8.8", "8.8.4.4"]
      }

      advanced_config {
        is_distributed_dhcp = true
      }
    }
  }
}

resource "nsxt_vpc_service_profile" "prod_service_profile" {
  context {
    project_id = nsxt_policy_project.prod_project.id
  }

  display_name = "prod-vpc-service-profile"
  description  = "Service profile for production VPCs"

  dhcp_config {
    dhcp_server_config {
      ntp_servers = ["20.2.60.3"]
      lease_time  = 43200

      dns_client_config {
        dns_server_ips = ["1.1.1.1", "1.0.0.1"]
      }

      advanced_config {
        is_distributed_dhcp = false
      }
    }
  }
}

# VPC Connectivity Profiles
resource "nsxt_vpc_connectivity_profile" "dev_connectivity_profile" {
  context {
    project_id = nsxt_policy_project.dev_project.id
  }

  display_name         = "dev-connectivity-profile"
  description          = "Connectivity profile for development VPCs"
  transit_gateway_path = nsxt_policy_transit_gateway.dev_tgw.path
  external_ip_blocks = [nsxt_policy_ip_block.dev_external_block.path]

  service_gateway {
    enable             = true
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.main_edge_cluster.path]
    nat_config {
      enable_default_snat = true
    }
  }
}

resource "nsxt_vpc_connectivity_profile" "prod_connectivity_profile" {
  context {
    project_id = nsxt_policy_project.prod_project.id
  }

  display_name         = "prod-connectivity-profile"
  description          = "Connectivity profile for production VPCs"
  transit_gateway_path = nsxt_policy_transit_gateway.prod_tgw.path
  external_ip_blocks = [nsxt_policy_ip_block.prod_external_block.path]

  service_gateway {
    enable             = true
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.main_edge_cluster.path]
    nat_config {
      enable_default_snat = true
    }
  }
}

# VPCs for different applications
resource "nsxt_vpc" "dev_web_vpc" {
  context {
    project_id = nsxt_policy_project.dev_project.id
  }

  display_name           = "dev-web-vpc"
  description            = "Development VPC for web applications"
  private_ips           = ["192.168.10.0/24"]
  short_id              = "dev-web"
  vpc_service_profile   = nsxt_vpc_service_profile.dev_service_profile.path

  load_balancer_vpc_endpoint {
    enabled = false
  }

  depends_on = [
    nsxt_policy_project.dev_project,
    nsxt_vpc_service_profile.dev_service_profile,
    nsxt_policy_transit_gateway.dev_tgw
  ]
}


# VPC Subnets using VPC's private IP space
resource "nsxt_vpc_subnet" "dev_web_private" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_web_vpc.id
  }

  display_name     = "dev-web-private-subnet"
  description      = "Private subnet for dev web VPC"
  ipv4_subnet_size = 32
  access_mode      = "Private"
}

resource "nsxt_vpc" "dev_db_vpc" {
  context {
    project_id = nsxt_policy_project.dev_project.id
  }

  display_name           = "dev-web-vpc"
  description            = "Development VPC for web applications"
  private_ips           = ["10.1.0.0/16"]
  short_id              = "dev-db"
  vpc_service_profile   = nsxt_vpc_service_profile.dev_service_profile.path

  load_balancer_vpc_endpoint {
    enabled = false
  }

  depends_on = [
    nsxt_policy_project.dev_project,
    nsxt_vpc_service_profile.dev_service_profile,
    nsxt_policy_transit_gateway.dev_tgw
  ]
}

resource "nsxt_vpc_attachment" "dev" {
  display_name             = "App1Attachment"
  description              = "terraform provisioned dev vpc attachment"
  parent_path              = nsxt_vpc.dev_web_vpc.path
  vpc_connectivity_profile = nsxt_vpc_connectivity_profile.dev_connectivity_profile.path
}

resource "nsxt_vpc_attachment" "dev_db" {
  display_name             = "devDBAttachment"
  description              = "terraform provisioned dev vpc attachment"
  parent_path              = nsxt_vpc.dev_db_vpc.path
  vpc_connectivity_profile = nsxt_vpc_connectivity_profile.dev_connectivity_profile.path
}

resource "nsxt_policy_group" "dev" {
  context {
    project_id = data.nsxt_policy_project.dev_project.id
  }

  display_name = "dev_apps"
  description  = "Terraform provisioned Group"

  criteria {
    path_expression {
      member_paths = [nsxt_vpc.dev_web_vpc.path, nsxt_vpc.dev_db_vpc.path]
    }
  }
  depends_on = [ nsxt_vpc.dev_db_vpc,nsxt_vpc.dev_web_vpc ]
}

resource "nsxt_policy_connectivity_policy" "devapps" {
  display_name = "devApps"
  group_path = nsxt_policy_group.dev.path
  parent_path = nsxt_policy_transit_gateway.dev_tgw.path
  connectivity_scope = "ISOLATED"

  depends_on = [
    nsxt_vpc_attachment.dev,
    nsxt_vpc_attachment.dev_db,
    nsxt_policy_group.dev
  ]
}

resource "nsxt_vpc_ip_address_allocation" "dev_external_nat_pool" {
  context {
    project_id = data.nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_web_vpc.id
  }
  
  display_name                = "dev_external-nat-pool"
  description                 = "IP allocation for NAT rules and external connectivity"
  allocation_size             = 1
  ip_address_block_visibility = "EXTERNAL"
  ip_address_type            = "IPV4"

  depends_on = [
    nsxt_vpc.dev_web_vpc,
    nsxt_vpc_connectivity_profile.dev_connectivity_profile,
    nsxt_vpc_attachment.dev,
    nsxt_policy_ip_block.dev_external_block
  ]
}



data "nsxt_vpc_nat" "app_nat" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_web_vpc.id
  }
  nat_type = "DEFAULT"
}

resource "nsxt_vpc_nat_rule" "outbound_snat" {
  display_name        = "outbound-snat-rule"
  description         = "SNAT rule for outbound internet access"
  parent_path         = data.nsxt_vpc_nat.app_nat.path
  action              = "SNAT"
  source_network      = "192.168.10.0/24"
  translated_network  = nsxt_vpc_ip_address_allocation.dev_external_nat_pool.allocation_ips
  logging             = true
  enabled             = true
  sequence_number     = 100
}

# Note: VPC External Address resources removed due to provider limitations
# These would require subnet port paths which are not available in this configuration

# Static binding for web server
resource "nsxt_vpc_dhcp_v4_static_binding" "dev_web_server_binding" {
  parent_path     = nsxt_vpc_subnet.dev_db_isolated_subnet.path
  display_name    = "dev-web-server-binding"
  description     = "Static DHCP binding for development web server"
  
  mac_address     = "00:50:56:11:22:33"
  ip_address      = "192.168.10.71"  # Within reserved range
  gateway_address = "192.168.10.72"
  host_name       = "dev-web01.dev.example.com"
  lease_time      = 86400

  options {
    option121 {
      static_route {
        network  = "0.0.0.0/0"
        next_hop = "192.168.10.73"
      }
    }
    
    other {
      code   = 6   # DNS servers
      values = ["8.8.8.8", "8.8.4.4"]
    }
  }
}



data "nsxt_policy_vm" "web_server" {
  context {
    project_id = nsxt_policy_project.dev_project.id
  }
  display_name = "Web-VM-2"
}

# data "nsxt_vpc_subnet_port" "dev_web_server_port" {
#   subnet_path = nsxt_vpc_subnet.dev_web_public_subnet.path
#   vm_id       = data.nsxt_policy_vm.web_server.instance_id # or NSXT_TEST_VPC_VM_ID
# }

# resource "nsxt_vpc_external_address" "web_server_public_ip" {
#   parent_path                = data.nsxt_vpc_subnet_port.dev_web_server_port.path
#   allocated_external_ip_path = nsxt_vpc_ip_address_allocation.dev_external_nat_pool.path
  
#   depends_on = [
#     nsxt_vpc_ip_address_allocation.dev_external_nat_pool,
#     data.nsxt_vpc_subnet_port.dev_web_server_port
#   ]
# }


resource "nsxt_vpc" "prod_web_vpc" {
  context {
    project_id = nsxt_policy_project.prod_project.id
  }

  display_name           = "prod-web-vpc"
  description            = "Production VPC for web applications"
  private_ips           = ["10.30.0.0/16"]
  short_id              = "prod-web"
  vpc_service_profile   = nsxt_vpc_service_profile.prod_service_profile.path

  load_balancer_vpc_endpoint {
    enabled = false
  }

  depends_on = [
    nsxt_policy_project.prod_project,
    nsxt_vpc_service_profile.prod_service_profile,
    nsxt_policy_transit_gateway.prod_tgw
  ]
}

resource "nsxt_vpc" "prod_db_vpc" {
  context {
    project_id = nsxt_policy_project.prod_project.id
  }

  display_name           = "prod-db-vpc"
  description            = "Production VPC for database services"
  private_ips           = ["10.40.0.0/16"]
  short_id              = "prod-db"
  vpc_service_profile   = nsxt_vpc_service_profile.prod_service_profile.path

  load_balancer_vpc_endpoint {
    enabled = false
  }

  depends_on = [
    nsxt_policy_project.prod_project,
    nsxt_vpc_service_profile.prod_service_profile,
    nsxt_policy_transit_gateway.prod_tgw
  ]
}

resource "nsxt_vpc_attachment" "prod_db" {
  display_name             = "prodDBAttachment"
  description              = "terraform provisioned dev vpc attachment"
  parent_path              = nsxt_vpc.prod_db_vpc.path
  vpc_connectivity_profile = nsxt_vpc_connectivity_profile.prod_connectivity_profile.path
}

resource "nsxt_vpc_attachment" "prod_web" {
  display_name             = "prodDBAttachment"
  description              = "terraform provisioned dev vpc attachment"
  parent_path              = nsxt_vpc.prod_web_vpc.path
  vpc_connectivity_profile = nsxt_vpc_connectivity_profile.prod_connectivity_profile.path
}


resource "nsxt_policy_group" "prod" {
  context {
    project_id = data.nsxt_policy_project.prod_project.id
  }

  display_name = "prod_apps"
  description  = "Terraform provisioned Group"

  criteria {
    path_expression {
      member_paths = [nsxt_vpc.prod_web_vpc.path, nsxt_vpc.prod_db_vpc.path]
    }
  }
  depends_on = [ nsxt_vpc.prod_db_vpc,nsxt_vpc.prod_web_vpc ]
}

resource "nsxt_policy_connectivity_policy" "prodapps" {
  display_name = "prodApps"
  group_path = nsxt_policy_group.prod.path
  parent_path = nsxt_policy_transit_gateway.prod_tgw.path
  connectivity_scope = "COMMUNITY"

  depends_on = [
    nsxt_vpc_attachment.prod_web,
    nsxt_vpc_attachment.prod_db,
    nsxt_policy_group.prod
  ]
}

resource "nsxt_vpc_subnet" "dev_web_private_subnet" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_web_vpc.id
  }

  display_name     = "dev-web-private-subnet"
  description      = "Private subnet for development web tier"
  ipv4_subnet_size = 64
  access_mode      = "Private"

  dhcp_config {
    mode = "DHCP_SERVER"
    dhcp_server_additional_config {
      options {
        option121 {
          static_route {
            network  = "0.0.0.0/0"
            next_hop = "192.168.10.1"
          }
        }
      }
    }
  }

  depends_on = [nsxt_vpc.dev_web_vpc]
}

resource "nsxt_vpc_subnet" "dev_web_public_subnet" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_web_vpc.id
  }

  display_name      = "dev-web-private-subnet-2"
  description       = "public subnet for development web tier"
  ipv4_subnet_size  = 32
  access_mode       = "Public"

  dhcp_config {
    mode = "DHCP_SERVER" 
    }
  depends_on = [
    nsxt_vpc.dev_web_vpc,
    nsxt_vpc_attachment.dev  
  ]
}

resource "nsxt_vpc_subnet" "dev_db_isolated_subnet" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_db_vpc.id
  }

  display_name = "dev-db-isolated-subnet"
  description  = "Isolated subnet for development database"
  ip_addresses = ["192.168.10.64/28"] 
  access_mode  = "Isolated"

  advanced_config {
    connectivity_state = "DISCONNECTED"
    static_ip_allocation {
      enabled = false
    }
  }

  dhcp_config {
    mode = "DHCP_SERVER"
        dhcp_server_additional_config {
         # Reserve IP ranges for static bindings and infrastructure
    reserved_ip_ranges = [
        "192.168.10.68",  
        "192.168.10.69-192.168.10.78"   
      ]
    }
  }

  depends_on = [nsxt_vpc.dev_db_vpc]
}


resource "nsxt_vpc_subnet" "prod_web_private_subnet" {
  context {
    project_id = nsxt_policy_project.prod_project.id
    vpc_id     = nsxt_vpc.prod_web_vpc.id
  }

  display_name     = "prod-web-private-subnet"
  description      = "Private subnet for production web tier"
  ipv4_subnet_size = 32
  access_mode      = "Private"

  dhcp_config {
    mode = "DHCP_DEACTIVATED"
  }

  depends_on = [nsxt_vpc.prod_web_vpc]
}

resource "nsxt_vpc_subnet" "prod_web_public_subnet" {
  context {
    project_id = nsxt_policy_project.prod_project.id
    vpc_id     = nsxt_vpc.prod_web_vpc.id
  }

  display_name     = "prod-web-private-subnet-2"
  description      = "Public subnet for production web tier"
  ipv4_subnet_size = 64
  access_mode      = "Public"

  dhcp_config {
    mode = "DHCP_DEACTIVATED"
  }

  depends_on = [nsxt_vpc.prod_web_vpc, nsxt_vpc_attachment.prod_web]
}

# VPC Groups for security policies
resource "nsxt_vpc_group" "dev_web_servers" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_web_vpc.id
  }

  display_name = "dev-web-servers"
  description  = "Development web servers group"

  criteria {
    condition {
      key         = "Name"
      member_type = "VirtualMachine"
      operator    = "STARTSWITH"
      value       = "dev-web"
    }
  }
}

resource "nsxt_vpc_group" "dev_db_servers" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_db_vpc.id
  }

  display_name = "dev-db-servers"
  description  = "Development database servers group"

  criteria {
    condition {
      key         = "Name"
      member_type = "VirtualMachine"
      operator    = "STARTSWITH"
      value       = "dev-db"
    }
  }
}

resource "nsxt_vpc_group" "prod_web_servers" {
  context {
    project_id = nsxt_policy_project.prod_project.id
    vpc_id     = nsxt_vpc.prod_web_vpc.id
  }

  display_name = "prod-web-servers"
  description  = "Production web servers group"

  criteria {
    condition {
      key         = "Name"
      member_type = "VirtualMachine"
      operator    = "STARTSWITH"
      value       = "prod-web"
    }
  }
}

resource "nsxt_vpc_group" "prod_db_servers" {
  context {
    project_id = nsxt_policy_project.prod_project.id
    vpc_id     = nsxt_vpc.prod_db_vpc.id
  }

  display_name = "prod-db-servers"
  description  = "Production database servers group"

  criteria {
    condition {
      key         = "Name"
      member_type = "VirtualMachine"
      operator    = "STARTSWITH"
      value       = "prod-db"
    }
  }
}

# VPC Security Policies
resource "nsxt_vpc_security_policy" "dev_web_security_policy" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_web_vpc.id
  }

  display_name = "dev-web-security-policy"
  description  = "Security policy for development web tier"
  locked       = false
  stateful     = true
  tcp_strict   = false
  scope        = [nsxt_vpc_group.dev_web_servers.path]

  rule {
    display_name       = "allow_http_https"
    description        = "Allow HTTP and HTTPS traffic"
    action             = "ALLOW"
    #source_groups      = []
    destination_groups = [nsxt_vpc_group.dev_web_servers.path]
   # services           = []
    logged             = true
  }

  rule {
    display_name       = "allow_ssh"
    description        = "Allow SSH access for management"
    action             = "ALLOW"
    source_groups      = []
    destination_groups = [nsxt_vpc_group.dev_web_servers.path]
    services           = []
    logged             = true
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "nsxt_vpc_security_policy" "dev_db_security_policy" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_db_vpc.id
  }

  display_name = "dev-db-security-policy"
  description  = "Security policy for development database tier"
  locked       = false
  stateful     = true
  tcp_strict   = true
  scope        = [nsxt_vpc_group.dev_db_servers.path]

  rule {
    display_name       = "allow_db_access"
    description        = "Allow database access within VPC"
    action             = "ALLOW"
    source_groups      = []
    destination_groups = [nsxt_vpc_group.dev_db_servers.path]
    services           = []
    logged             = true
  }

  rule {
    display_name       = "deny_all_other"
    description        = "Deny all other traffic"
    action             = "DROP"
    source_groups      = []
    destination_groups = [nsxt_vpc_group.dev_db_servers.path]
    services           = []
    logged             = true
  }

  lifecycle {
    create_before_destroy = true
  }
}

# VPC Gateway Policies
resource "nsxt_vpc_gateway_policy" "dev_web_gateway_policy" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_web_vpc.id
  }

  display_name    = "dev-web-gateway-policy"
  description     = "Gateway policy for development web VPC"
  locked          = false
  sequence_number = 1
  stateful        = true
  tcp_strict      = false

  rule {
    display_name       = "allow_inbound_web"
    description        = "Allow inbound web traffic"
    action             = "ALLOW"
    direction          = "IN"
    destination_groups = [nsxt_vpc_group.dev_web_servers.path]
    logged             = true
  }

  rule {
    display_name = "allow_outbound_all"
    description  = "Allow all outbound traffic"
    action       = "ALLOW"
    direction    = "OUT"
    logged       = false
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "nsxt_vpc_gateway_policy" "prod_web_gateway_policy" {
  context {
    project_id = nsxt_policy_project.prod_project.id
    vpc_id     = nsxt_vpc.prod_web_vpc.id
  }

  display_name    = "prod-web-gateway-policy"
  description     = "Gateway policy for production web VPC"
  locked          = false
  sequence_number = 1
  stateful        = true
  tcp_strict      = true

  rule {
    display_name       = "allow_inbound_web"
    description        = "Allow inbound web traffic"
    action             = "ALLOW"
    direction          = "IN"
    destination_groups = [nsxt_vpc_group.prod_web_servers.path]
    logged             = true
  }

  rule {
    display_name = "allow_outbound_controlled"
    description  = "Allow controlled outbound traffic"
    action       = "ALLOW"
    direction    = "OUT"
    logged       = true
  }

  lifecycle {
    create_before_destroy = true
  }
}


# VPC Static Routes
resource "nsxt_vpc_static_route" "dev_web_static_route" {
  context {
    project_id = nsxt_policy_project.dev_project.id
    vpc_id     = nsxt_vpc.dev_web_vpc.id
  }

  display_name = "dev-web-static-route"
  description  = "Static route for development web VPC"

  network = "10.10.0.0/16"
  next_hop {
    ip_address = "10.10.1.1"
    admin_distance = 1
  }
}

resource "nsxt_vpc_static_route" "prod_web_static_route" {
  context {
    project_id = nsxt_policy_project.prod_project.id
    vpc_id     = nsxt_vpc.prod_web_vpc.id
  }

  display_name = "prod-web-static-route"
  description  = "Static route for production web VPC"

  network = "10.30.0.0/16"
  next_hop {
    ip_address = "10.30.1.1"
    admin_distance = 1
  }
}

data "nsxt_policy_transit_gateway_nat" "prod_nat" {
  transit_gateway_path = nsxt_policy_transit_gateway.prod_tgw.path
}

data "nsxt_policy_transit_gateway_nat" "dev_nat" {
  transit_gateway_path = nsxt_policy_transit_gateway.dev_tgw.path
}

# Transit Gateway Static Routes
resource "nsxt_policy_transit_gateway_static_route" "dev_tgw_static_route" {
  display_name = "dev-tgw-static-route"
  description  = "Static route for development transit gateway"
  parent_path  = nsxt_policy_transit_gateway.dev_tgw.path
  network = "192.168.0.0/16"
  next_hop {
    scope = [nsxt_policy_transit_gateway_attachment.dev_tgw_attachment.path]
    admin_distance = 1
  }
}

resource "nsxt_policy_transit_gateway_static_route" "prod_tgw_static_route" {
  display_name = "prod-tgw-static-route"
  description  = "Static route for production transit gateway"
  parent_path  = nsxt_policy_transit_gateway.prod_tgw.path
  network = "192.168.0.0/16"
  next_hop {
    scope = [nsxt_policy_transit_gateway_attachment.prod_tgw_attachment.path]
    admin_distance = 1
  }
}


