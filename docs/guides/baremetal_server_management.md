---
page_title: "Bare Metal Server Management with NSX-T Terraform Provider"
description: |-
  This guide provides comprehensive examples and best practices for managing bare metal servers with the NSX-T Terraform provider.
---

# Bare Metal Server Management with NSX-T Terraform Provider

This guide demonstrates how to use the NSX-T Terraform provider to discover, tag, group, and apply security policies to bare metal servers managed by NSX-T.

## Prerequisites

- NSX-T Data Center 9.0.0 or higher
- Bare metal servers discovered and registered with NSX-T through a compute manager
- NSX-T Terraform Provider with BMS support
- Local Manager access (BMS features are not supported on Global Manager)

## Overview

The NSX-T Terraform provider supports the following BMS operations:

1. **Discovery**: Query BMS inventory (servers and interfaces)
2. **Tagging**: Apply tags to servers and interfaces for automation
3. **Grouping**: Create static and dynamic groups based on server attributes
4. **Security**: Apply firewall policies using BMS groups

## Basic Workflow

### Step 1: Discover Bare Metal Server Inventory

```hcl
# Discover all bare metal servers
data "nsxt_policy_baremetal_servers" "all_servers" {
}

# Discover all bare metal server interfaces
data "nsxt_policy_baremetal_server_interfaces" "all_interfaces" {
}

# Output inventory for review
output "bms_inventory_summary" {
  value = {
    total_servers    = length(data.nsxt_policy_baremetal_servers.all_servers.results)
    total_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all_interfaces.results)

    servers = [
      for server in data.nsxt_policy_baremetal_servers.all_servers.results : {
        id           = server.external_id
        display_name = server.display_name
        cpu_cores    = server.cpu_cores
        os_name      = server.os_name
        os_version   = server.os_version
      }
    ]
  }
}
```

### Step 2: Apply Tags for Automation

```hcl
# Tag servers based on naming conventions
resource "nsxt_policy_baremetal_server_tags" "production_servers" {
  for_each = {
    for idx, server in data.nsxt_policy_baremetal_servers.all_servers.results :
    server.external_id => server
    if length(regexall("prod|production", lower(server.display_name))) > 0
  }

  external_id = each.key

  tag {
    scope = "environment"
    tag   = "production"
  }

  tag {
    scope = "managed-by"
    tag   = "terraform"
  }

  tag {
    scope = "cpu-tier"
    tag   = each.value.cpu_cores >= 16 ? "high" : "standard"
  }
}

# Tag interfaces, excluding management interfaces
resource "nsxt_policy_baremetal_server_interface_tags" "data_interfaces" {
  for_each = {
    for idx, iface in data.nsxt_policy_baremetal_server_interfaces.all_interfaces.results :
    iface.external_id => iface
    if !iface.is_mgmt_interface && iface.state == "UP"
  }

  external_id = each.key

  tag {
    scope = "interface-role"
    tag   = "data-plane"
  }

  tag {
    scope = "speed"
    tag   = contains(each.value.display_name, "10G") ? "10Gbps" : "1Gbps"
  }
}
```

### Step 3: Create Groups for Policy Application

```hcl
# Static group for specific servers
resource "nsxt_policy_group" "critical_servers" {
  display_name = "Critical-BMS-Servers"
  description  = "Manually selected critical bare metal servers"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = [
        "server-001-uuid",
        "server-002-uuid",
        "server-003-uuid"
      ]
    }
  }
}

# Dynamic group based on tags
resource "nsxt_policy_group" "production_servers" {
  display_name = "Production-BMS-Servers"
  description  = "All production bare metal servers"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "environment|production"
    }
  }
}

# High-performance server group
resource "nsxt_policy_group" "high_performance_servers" {
  display_name = "High-Performance-BMS"
  description  = "High CPU bare metal servers"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "environment|production"
    }
  }

  conjunction {
    operator = "AND"
  }

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "cpu-tier|high"
    }
  }
}

# Interface groups for network policies
resource "nsxt_policy_group" "high_speed_interfaces" {
  display_name = "High-Speed-BMS-Interfaces"
  description  = "10Gbps+ bare metal server interfaces"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "speed|10Gbps"
    }
  }
}
```

### Step 4: Apply Security Policies

```hcl
# Basic security policy for production servers
resource "nsxt_policy_security_policy" "production_bms_security" {
  display_name = "Production-BMS-Security"
  description  = "Security policy for production bare metal servers"
  category     = "Application"

  rule {
    display_name       = "Allow-Internal-Communication"
    source_groups      = [nsxt_policy_group.production_servers.path]
    destination_groups = [nsxt_policy_group.production_servers.path]
    action             = "ALLOW"
    services           = ["SSH", "HTTP", "HTTPS"]
    logged             = true
  }

  rule {
    display_name       = "Allow-Management-Access"
    source_groups      = ["management-network-group"] # Pre-existing group
    destination_groups = [nsxt_policy_group.production_servers.path]
    action             = "ALLOW"
    services           = ["SSH", "SNMP"]
    logged             = true
  }

  rule {
    display_name       = "Drop-Unauthorized"
    source_groups      = ["ANY"]
    destination_groups = [nsxt_policy_group.production_servers.path]
    action             = "DROP"
    logged             = true
  }
}

# High-security policy for critical servers
resource "nsxt_policy_security_policy" "critical_bms_security" {
  display_name = "Critical-BMS-Security"
  description  = "High security policy for critical servers"
  category     = "Application"

  rule {
    display_name       = "Allow-Specific-Management"
    source_groups      = ["critical-management-group"] # Pre-existing group
    destination_groups = [nsxt_policy_group.critical_servers.path]
    action             = "ALLOW"
    services           = ["SSH"]
    logged             = true
  }

  rule {
    display_name       = "Deny-All-Others"
    source_groups      = ["ANY"]
    destination_groups = [nsxt_policy_group.critical_servers.path]
    action             = "DROP"
    logged             = true
  }
}
```

## Advanced Use Cases

### Multi-Tier Application Architecture

```hcl
# Separate groups for different application tiers
locals {
  app_tiers = ["web", "app", "db"]
}

resource "nsxt_policy_baremetal_server_tags" "app_tier_tags" {
  for_each = {
    for server in data.nsxt_policy_baremetal_servers.all_servers.results :
    server.external_id => {
      id   = server.external_id
      tier = length([for tier in local.app_tiers : tier if contains(lower(server.display_name), tier)]) > 0 ? [for tier in local.app_tiers : tier if contains(lower(server.display_name), tier)][0] : "general"
    }
  }

  external_id = each.key

  tag {
    scope = "app-tier"
    tag   = each.value.tier
  }

  tag {
    scope = "managed-by"
    tag   = "terraform"
  }
}

# Create tier-specific groups
resource "nsxt_policy_group" "app_tier_groups" {
  for_each = toset(local.app_tiers)

  display_name = "BMS-${title(each.key)}-Tier"
  description  = "${title(each.key)} tier bare metal servers"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "app-tier|${each.key}"
    }
  }
}

# Tier-specific security policies
resource "nsxt_policy_security_policy" "tier_communication" {
  display_name = "Multi-Tier-BMS-Communication"
  description  = "Inter-tier communication rules"
  category     = "Application"

  # Web tier can access app tier
  rule {
    display_name       = "Web-to-App"
    source_groups      = [nsxt_policy_group.app_tier_groups["web"].path]
    destination_groups = [nsxt_policy_group.app_tier_groups["app"].path]
    action             = "ALLOW"
    services           = ["HTTP", "HTTPS"]
    logged             = true
  }

  # App tier can access database tier
  rule {
    display_name       = "App-to-DB"
    source_groups      = [nsxt_policy_group.app_tier_groups["app"].path]
    destination_groups = [nsxt_policy_group.app_tier_groups["db"].path]
    action             = "ALLOW"
    services           = ["MySQL", "PostgreSQL"]
    logged             = true
  }

  # Deny direct web to database access
  rule {
    display_name       = "Block-Web-to-DB"
    source_groups      = [nsxt_policy_group.app_tier_groups["web"].path]
    destination_groups = [nsxt_policy_group.app_tier_groups["db"].path]
    action             = "DROP"
    logged             = true
  }
}
```

### Network Segmentation by VLAN

```hcl
# Tag interfaces by VLAN information
resource "nsxt_policy_baremetal_server_interface_tags" "vlan_segmentation" {
  for_each = {
    for iface in data.nsxt_policy_baremetal_server_interfaces.all_interfaces.results :
    iface.external_id => {
      id   = iface.external_id
      vlan = can(regex("vlan[_-]?(\\d+)", lower(iface.display_name))) ? regex("vlan[_-]?(\\d+)", lower(iface.display_name))[0] : "unknown"
    }
    if !iface.is_mgmt_interface
  }

  external_id = each.key

  tag {
    scope = "vlan"
    tag   = each.value.vlan
  }

  tag {
    scope = "segment-type"
    tag   = each.value.vlan == "100" ? "dmz" : each.value.vlan == "200" ? "internal" : "general"
  }
}

# Create VLAN-specific groups
resource "nsxt_policy_group" "dmz_interfaces" {
  display_name = "DMZ-BMS-Interfaces"
  description  = "DMZ VLAN interfaces"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "segment-type|dmz"
    }
  }
}

resource "nsxt_policy_group" "internal_interfaces" {
  display_name = "Internal-BMS-Interfaces"
  description  = "Internal VLAN interfaces"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "segment-type|internal"
    }
  }
}
```

## Best Practices

### 1. Consistent Tagging Strategy

```hcl
# Define standard tags as locals
locals {
  standard_server_tags = {
    managed-by          = "terraform"
    provisioned         = formatdate("YYYY-MM-DD", timestamp())
    terraform-workspace = terraform.workspace
  }

  environment_mapping = {
    prod = "production"
    dev  = "development"
    test = "testing"
    stg  = "staging"
  }
}

# Apply standard tags to all servers
resource "nsxt_policy_baremetal_server_tags" "standard_tags" {
  for_each = {
    for server in data.nsxt_policy_baremetal_servers.all_servers.results :
    server.external_id => server
  }

  external_id = each.key

  dynamic "tag" {
    for_each = local.standard_server_tags
    content {
      scope = tag.key
      tag   = tag.value
    }
  }

  # Environment-specific tag based on naming convention
  tag {
    scope = "environment"
    tag = lookup(
      local.environment_mapping,
      split("-", each.value.display_name)[0],
      "unknown"
    )
  }
}
```

### 2. Conditional Resource Creation

```hcl
# Only create security policies if servers exist
resource "nsxt_policy_security_policy" "conditional_bms_policy" {
  count = length(data.nsxt_policy_baremetal_servers.all_servers.results) > 0 ? 1 : 0

  display_name = "BMS-Conditional-Policy"
  description  = "Policy created only when BMS servers exist"
  category     = "Application"

  rule {
    display_name       = "Default-Allow-Internal"
    source_groups      = [nsxt_policy_group.production_servers.path]
    destination_groups = [nsxt_policy_group.production_servers.path]
    action             = "ALLOW"
    logged             = true
  }
}
```

### 3. Error Handling and Validation

```hcl
# Validate that required servers exist before proceeding
locals {
  required_servers = ["critical-server-001", "critical-server-002"]

  discovered_critical_servers = [
    for server in data.nsxt_policy_baremetal_servers.all_servers.results :
    server.display_name
    if contains(local.required_servers, server.display_name)
  ]
}

# Use validation to ensure critical servers are discovered
resource "null_resource" "validate_critical_servers" {
  count = length(local.discovered_critical_servers) == length(local.required_servers) ? 0 : 1

  provisioner "local-exec" {
    command = "echo 'ERROR: Not all required servers discovered' && exit 1"
  }
}
```

### 4. Modular Organization

```hcl
# variables.tf
variable "bms_tagging_strategy" {
  description = "BMS tagging strategy configuration"
  type = object({
    environment_tag_scope = string
    application_tag_scope = string
    owner_tag_scope       = string
  })
  default = {
    environment_tag_scope = "environment"
    application_tag_scope = "application"
    owner_tag_scope       = "owner"
  }
}

# modules/bms-management/main.tf
resource "nsxt_policy_baremetal_server_tags" "servers" {
  for_each = var.server_configurations

  external_id = each.key

  dynamic "tag" {
    for_each = each.value.tags
    content {
      scope = tag.key
      tag   = tag.value
    }
  }
}

# Call module from main configuration
module "production_bms" {
  source = "./modules/bms-management"

  server_configurations = {
    for server in data.nsxt_policy_baremetal_servers.production.results :
    server.external_id => {
      tags = {
        environment = "production"
        application = "web-app"
        owner       = "platform-team"
      }
    }
  }
}
```

## Troubleshooting

### Common Issues and Solutions

1. **No BMS servers discovered**: Ensure compute manager is configured and servers are registered
2. **Tagging failures**: Verify server external IDs are correct and accessible
3. **Group membership issues**: Check tag scopes and values match exactly
4. **Policy application failures**: Ensure groups exist before referencing in policies

### Debugging Data Sources

```hcl
# Debug server discovery
output "debug_servers" {
  value = {
    for idx, server in data.nsxt_policy_baremetal_servers.all_servers.results :
    idx => {
      external_id  = server.external_id
      display_name = server.display_name
      cpu_cores    = server.cpu_cores
      os_name      = server.os_name
      os_version   = server.os_version
    }
  }
}

# Debug interface discovery
output "debug_interfaces" {
  value = {
    for idx, iface in data.nsxt_policy_baremetal_server_interfaces.all_interfaces.results :
    idx => {
      external_id       = iface.external_id
      display_name      = iface.display_name
      bms_external_id   = iface.bms_external_id
      is_mgmt_interface = iface.is_mgmt_interface
      state             = iface.state
      ip_addresses      = iface.ip_addresses
    }
  }
}
```

This guide provides a foundation for managing bare metal servers with the NSX-T Terraform provider. Adapt the examples to your specific environment and requirements.
