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
- Local Manager access — BMS features are not supported on Global Manager

## Overview

The NSX-T Terraform provider supports the following BMS operations:

1. **Discovery**: Query BMS inventory (servers and interfaces)
2. **Tagging**: Apply tags to servers and interfaces for automation
3. **Grouping**: Create static and dynamic groups based on server attributes
4. **Security**: Apply firewall policies using BMS groups

## Resources and Data Sources Reference

| Resource / Data Source | Type | Purpose |
|---|---|---|
| `nsxt_policy_baremetal_server_tags` | Resource | Apply tags to a bare metal server |
| `nsxt_policy_baremetal_server_interface_tags` | Resource | Apply tags to a bare metal server interface |
| `nsxt_policy_group` | Resource | Create BMS groups with `group_type = "BareMetalServer"` and `member_type = "BareMetalServer"` or `"BareMetalServerInterface"`; supports dynamic conditions, static external ID expressions, and nested path expressions |
| `nsxt_policy_security_policy` | Resource | Apply DFW policies scoped to BMS groups with inline rules |
| `nsxt_policy_parent_security_policy` | Resource | Define a BMS-scoped DFW policy without inline rules (use with `nsxt_policy_security_policy_rule`) |
| `nsxt_policy_security_policy_rule` | Resource | Create standalone DFW rules that target BMS groups as source, destination, or scope |

| `nsxt_policy_baremetal_servers` | Data Source | Query BMS inventory with filtering (bulk) |
| `nsxt_policy_baremetal_server` | Data Source | Look up a single BMS by external ID or display name |
| `nsxt_policy_baremetal_server_interfaces` | Data Source | Query BMS interface inventory with filtering |
| `nsxt_policy_baremetal_server_interface` | Data Source | Look up a single BMS interface by external ID or display name |
| `nsxt_policy_baremetal_server_tags` | Data Source | Read current tags on a BMS server |
| `nsxt_policy_baremetal_server_interface_tags` | Data Source | Read current tags on a BMS interface |
| `nsxt_policy_baremetal_server_group_associations` | Data Source | List all groups a BMS server belongs to; group paths can be used directly in policy rules |
| `nsxt_policy_baremetal_server_interface_group_associations` | Data Source | List all groups a BMS interface belongs to |
| `nsxt_policy_group_baremetal_server_members` | Data Source | List all BMS server members of a group |
| `nsxt_policy_group_baremetal_server_interface_members` | Data Source | List all BMS interface members of a group |
| `nsxt_policy_group` | Data Source | Look up an existing BMS group by display name to obtain its path for use in policies |
| `nsxt_policy_security_policy` | Data Source | Look up an existing security policy by display name to obtain its path |

## Known Limitations

### Platform Limitations

- **Multitenancy not supported** — BMS features are not available in multi-tenant (VPC/Project) deployments.
- **Federation / Global Manager not supported** — BMS management is available on Local Manager only.

### Firewall Limitations

- **Only DFW (Distributed Firewall) is supported** — BMS groups cannot be used in Gateway Firewall (`nsxt_policy_gateway_policy`) or IDPS policies (`nsxt_policy_intrusion_service_policy`).
- **Ethernet category not supported** — The `Ethernet` DFW category operates at Layer 2 only and does not support BMS groups. Use `Emergency`, `Infrastructure`, `Environment`, or `Application` categories.
- **Layer 3 / Layer 4 only** — Application-layer (L7) rules are not supported for BMS targets.
- **No L7 / FQDN context profiles** — The `context_profiles` field in security policy rules is not supported when BMS groups are in the policy scope, source, or destination.
- **No time-based rules** — Rule schedulers are silently ignored by NSX for BMS rules.
- **DFW Exclusion List does not apply to BMS** — Adding a BMS server to the DFW Exclusion List has no effect.

### Group Constraints

- **`group_type = "BareMetalServer"` is always required** when a group contains `BareMetalServer` or `BareMetalServerInterface` member types.
- **BMS groups may only contain BMS or BMSI members** — mixing `BareMetalServer` / `BareMetalServerInterface` with `VirtualMachine` or any other member type in the same group is not supported.
- Separate `BareMetalServer`-typed groups from VM groups. If a policy must cover both, reference them as separate groups in the policy's `scope`, `source_groups`, or `destination_groups`.

## Workflow Overview

```
  Discover BMS Inventory
  (nsxt_policy_baremetal_servers / nsxt_policy_baremetal_server)
              |
              v
  Apply Tags for Classification
  (nsxt_policy_baremetal_server_tags / nsxt_policy_baremetal_server_interface_tags)
              |
              v
  Create Groups based on Attributes / Tags
  (nsxt_policy_group with group_type = "BareMetalServer")
              |
              v
  Apply DFW Security Policies           <----+
  (nsxt_policy_security_policy /              |
   nsxt_policy_parent_security_policy +        |  group paths used
   nsxt_policy_security_policy_rule)           |  as source_groups /
              |                               |  destination_groups
              v                               |
  Query Group Memberships & Associations -----+
  (nsxt_policy_baremetal_server_group_associations —
   discover which groups a server belongs to, then
   feed those group paths directly into a policy rule)
```

`nsxt_policy_baremetal_server_group_associations` returns the list of group paths a given BMS server (or interface) is already a member of. You can pass those paths directly to `source_groups` or `destination_groups` in a security policy rule — see the [example in Step 6](#step-6-audit-group-membership-and-associations) below.

---

## BMS Group Condition Reference

The following keys and operators are supported for BMS member types in `condition` blocks within `nsxt_policy_group` (requires NSX-T 9.0.0 or higher):

**`member_type = "BareMetalServer"`**

| Key | Supported Operators | Notes |
|---|---|---|
| `Tag` | `EQUALS`, `CONTAINS`, `ENDSWITH`, `STARTSWITH` | Use `scope\|value` notation for scoped tags |
| `Name` | `EQUALS`, `CONTAINS`, `ENDSWITH`, `STARTSWITH`, `NOTEQUALS` | |
| `OSName` | `EQUALS`, `CONTAINS`, `ENDSWITH`, `STARTSWITH`, `NOTEQUALS` | |
| `OSVersion` | `EQUALS` only | Full version string, e.g. `"8.5.2111"`, `"8.10"` |
| `ALL` | `EQUALS` only | Value must be `"ALL_BMS"` |

**`member_type = "BareMetalServerInterface"`**

| Key | Supported Operators | Notes |
|---|---|---|
| `Tag` | `EQUALS`, `NOTEQUALS`, `NOTIN` | Use `scope\|value` notation for scoped tags |
| `ManagementInterface` | `EQUALS` only | Values: `"TRUE"` or `"FALSE"` |

---

## Step 1: Discover Bare Metal Server Inventory

```hcl
# List all bare metal servers
data "nsxt_policy_baremetal_servers" "all_servers" {}

# List all bare metal server interfaces
data "nsxt_policy_baremetal_server_interfaces" "all_interfaces" {}

# Filter by production environment tag
data "nsxt_policy_baremetal_servers" "prod_servers" {
  tag_scope = "environment"
  tag       = "production"
}

# Filter by OS name substring
data "nsxt_policy_baremetal_servers" "linux_servers" {
  os_name = "redhat"
}

# Filter by display name using a regex pattern
data "nsxt_policy_baremetal_servers" "web_servers" {
  display_name = "web-.*"
}

# Filter interfaces by parent server
data "nsxt_policy_baremetal_server_interfaces" "server_interfaces" {
  bms_external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
}

# Look up a single server by external ID
data "nsxt_policy_baremetal_server" "my_server" {
  external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
}

# Look up a single server by display name (must return exactly one result)
data "nsxt_policy_baremetal_server" "web_server_01" {
  display_name = "web-server-01"
}

output "bms_inventory" {
  value = {
    total_servers    = length(data.nsxt_policy_baremetal_servers.all_servers.results)
    total_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all_interfaces.results)
    servers = [
      for s in data.nsxt_policy_baremetal_servers.all_servers.results : {
        id           = s.external_id
        display_name = s.display_name
        os           = "${s.os_name} ${s.os_version}"
        cpu_cores    = s.cpu_cores
      }
    ]
  }
}
```

---

## Step 2: Apply Tags for Classification

Tags are the primary mechanism for dynamic group membership. Apply them with `nsxt_policy_baremetal_server_tags` (for servers) and `nsxt_policy_baremetal_server_interface_tags` (for interfaces).

```hcl
# Tag a single server with a known external ID
resource "nsxt_policy_baremetal_server_tags" "prod_web_01" {
  external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"

  tag {
    scope = "environment"
    tag   = "production"
  }

  tag {
    scope = "application"
    tag   = "web-server"
  }

  tag {
    scope = "managed-by"
    tag   = "terraform"
  }
}

# Tag all servers matching a naming pattern using for_each
resource "nsxt_policy_baremetal_server_tags" "production_servers" {
  for_each = {
    for s in data.nsxt_policy_baremetal_servers.all_servers.results :
    s.external_id => s
    if length(regexall("prod", lower(s.display_name))) > 0
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
}

# Tag data-plane interfaces (skip management interfaces)
resource "nsxt_policy_baremetal_server_interface_tags" "data_interfaces" {
  for_each = {
    for iface in data.nsxt_policy_baremetal_server_interfaces.all_interfaces.results :
    iface.external_id => iface
    if !iface.is_mgmt_interface && iface.state == "UP"
  }

  external_id = each.key

  tag {
    scope = "network-type"
    tag   = "data-plane"
  }

  tag {
    scope = "managed-by"
    tag   = "terraform"
  }
}

# Tag management interfaces
resource "nsxt_policy_baremetal_server_interface_tags" "mgmt_interfaces" {
  for_each = {
    for iface in data.nsxt_policy_baremetal_server_interfaces.all_interfaces.results :
    iface.external_id => iface
    if iface.is_mgmt_interface
  }

  external_id = each.key

  tag {
    scope = "network-type"
    tag   = "management"
  }
}
```

---

## Step 3: Read Back Existing Tags

Use the tag data sources to read the current tag state on servers and interfaces — useful for drift detection or conditional logic.

```hcl
# Read tags on a specific server
data "nsxt_policy_baremetal_server_tags" "server_tags" {
  external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
}

# Extract a tag value by scope
locals {
  server_environment = try(
    [for t in data.nsxt_policy_baremetal_server_tags.server_tags.tag : t.tag if t.scope == "environment"][0],
    "unknown"
  )
}

# Read tags on a specific interface
data "nsxt_policy_baremetal_server_interface_tags" "if_tags" {
  external_id = "71be0142-2ed1-1d53-9c60-02005b4b7246"
}

output "interface_tags" {
  value = data.nsxt_policy_baremetal_server_interface_tags.if_tags.tag
}
```

---

## Step 4: Create Groups

All groups containing BMS or BMSI members **must** set `group_type = "BareMetalServer"`. Groups may contain either BMS or BMSI members (or both), but cannot mix them with `VirtualMachine` or other member types.


```hcl
# Dynamic group — all production BMS servers using Tags
resource "nsxt_policy_group" "production_bms" {
  display_name = "Production-BMS-Servers"
  description  = "All production bare metal servers"
  group_type   = "BareMetalServer"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "environment|production"
    }
  }
}

# Dynamic group — BMS servers AND/OR interfaces combined using OR conjunction
resource "nsxt_policy_group" "web_tier_bms" {
  display_name = "Web-Tier-BMS"
  description  = "Web tier BMS servers and their data-plane interfaces"
  group_type   = "BareMetalServer"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "application|web-server"
    }
  }

  conjunction {
    operator = "OR"
  }

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "network-type|data-plane"
    }
  }
}

# Dynamic group — servers by OS version (EQUALS only)
resource "nsxt_policy_group" "redhat_810_servers" {
  display_name = "RedHat-8-10-Servers"
  description  = "BMS servers running RedHat 8.10"
  group_type   = "BareMetalServer"

  criteria {
    condition {
      key         = "OSVersion"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "8.10"
    }
  }
}

# Dynamic group — management interfaces only
resource "nsxt_policy_group" "mgmt_interfaces" {
  display_name = "BMS-Management-Interfaces"
  description  = "All BMS management interfaces"
  group_type   = "BareMetalServer"

  criteria {
    condition {
      key         = "ManagementInterface"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "TRUE"
    }
  }
}

# Dynamic group — Linux servers using OSName
resource "nsxt_policy_group" "linux_bms" {
  display_name = "Linux-BMS-Servers"
  description  = "All Linux bare metal servers"
  group_type   = "BareMetalServer"

  criteria {
    condition {
      key         = "OSName"
      member_type = "BareMetalServer"
      operator    = "CONTAINS"
      value       = "Linux"
    }
  }
}
```

### Static Groups (External ID-based)

```hcl
# Static group — specific servers by external ID
resource "nsxt_policy_group" "critical_servers" {
  display_name = "Critical-BMS-Servers"
  description  = "Manually selected critical bare metal servers"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = [
        "71be0142-2ed1-1d53-9c60-5564cf4b7e2e",
        "81be0142-2ed1-1d53-9c60-5564cf4b7e2f"
      ]
    }
  }
}

# Static group — specific interfaces by external ID
resource "nsxt_policy_group" "critical_interfaces" {
  display_name = "Critical-BMS-Interfaces"
  description  = "Manually selected critical BMS interfaces"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = [
        "71be0142-2ed1-1d53-9c60-02005b4b7246"
      ]
    }
  }
}

# Static group — mix of servers and interfaces
resource "nsxt_policy_group" "mixed_static_bms" {
  display_name = "Mixed-Static-BMS"
  description  = "Static group with both BMS servers and interfaces"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type  = "BareMetalServer"
      external_ids = ["71be0142-2ed1-1d53-9c60-5564cf4b7e2e"]
    }
  }

  conjunction {
    operator = "OR"
  }

  criteria {
    external_id_expression {
      member_type  = "BareMetalServerInterface"
      external_ids = ["71be0142-2ed1-1d53-9c60-02005b4b7246"]
    }
  }
}

# Nested group — reference another BMS group by path
resource "nsxt_policy_group" "nested_bms" {
  display_name = "Nested-BMS-Group"
  description  = "Group that includes another BMS group by reference"
  group_type   = "BareMetalServer"

  criteria {
    path_expression {
      member_paths = [nsxt_policy_group.production_bms.path]
    }
  }
}
```

### Tag-based Group from Discovery

```hcl
# Discover and tag, then group — all in one workflow
data "nsxt_policy_baremetal_servers" "all" {}

resource "nsxt_policy_baremetal_server_tags" "env_tags" {
  for_each = {
    for s in data.nsxt_policy_baremetal_servers.all.results :
    s.external_id => s
  }

  external_id = each.key

  tag {
    scope = "environment"
    tag   = length(regexall("prod", lower(each.value.display_name))) > 0 ? "production" : "non-production"
  }
}

resource "nsxt_policy_group" "prod_bms_group" {
  display_name = "All-Production-BMS"
  group_type   = "BareMetalServer"

  depends_on = [nsxt_policy_baremetal_server_tags.env_tags]

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "environment|production"
    }
  }
}
```

-> For the full list of supported keys and operators per member type, see the [BMS Group Condition Reference](#bms-group-condition-reference) section above.

---

## Step 5: Apply DFW Security Policies

```hcl
resource "nsxt_policy_service" "ssh" {
  display_name = "SSH"
  l4_port_set_entry {
    protocol          = "TCP"
    destination_ports = ["22"]
  }
}

resource "nsxt_policy_service" "http" {
  display_name = "HTTP"
  l4_port_set_entry {
    protocol          = "TCP"
    destination_ports = ["80"]
  }
}

resource "nsxt_policy_service" "https" {
  display_name = "HTTPS"
  l4_port_set_entry {
    protocol          = "TCP"
    destination_ports = ["443"]
  }
}

# Basic security policy scoped to a BMS group
resource "nsxt_policy_security_policy" "prod_bms_policy" {
  display_name = "Production-BMS-Policy"
  description  = "Security policy for production bare metal servers"
  category     = "Application"
  stateful     = true
  scope        = [nsxt_policy_group.production_bms.path]

  rule {
    display_name       = "allow_internal_communication"
    description        = "Allow traffic between production BMS servers"
    source_groups      = [nsxt_policy_group.production_bms.path]
    destination_groups = [nsxt_policy_group.production_bms.path]
    action             = "ALLOW"
    services           = [nsxt_policy_service.http.path, nsxt_policy_service.https.path]
    logged             = true
    sequence_number    = 100
  }

  rule {
    display_name       = "allow_ssh_management"
    description        = "Allow SSH from management network"
    destination_groups = [nsxt_policy_group.production_bms.path]
    action             = "ALLOW"
    services           = [nsxt_policy_service.ssh.path]
    logged             = true
    sequence_number    = 200
  }

  rule {
    display_name    = "default_deny"
    description     = "Deny all other traffic"
    action          = "DROP"
    logged          = true
    sequence_number = 1000
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Policy using standalone rules with nsxt_policy_parent_security_policy
resource "nsxt_policy_parent_security_policy" "bms_parent" {
  display_name = "BMS-Parent-Policy"
  category     = "Application"
  stateful     = true
  scope = [
    nsxt_policy_group.production_bms.path,
    nsxt_policy_group.web_tier_bms.path
  ]

  lifecycle {
    create_before_destroy = true
  }
}

resource "nsxt_policy_security_policy_rule" "allow_web_traffic" {
  display_name       = "allow-web-traffic"
  description        = "Allow web traffic to BMS web tier"
  policy_path        = nsxt_policy_parent_security_policy.bms_parent.path
  sequence_number    = 100
  destination_groups = [nsxt_policy_group.web_tier_bms.path]
  action             = "ALLOW"
  services           = [nsxt_policy_service.http.path, nsxt_policy_service.https.path]
  logged             = true
}

resource "nsxt_policy_security_policy_rule" "allow_bms_internal" {
  display_name       = "allow-bms-internal"
  description        = "Allow BMS-to-BMS communication"
  policy_path        = nsxt_policy_parent_security_policy.bms_parent.path
  sequence_number    = 200
  source_groups      = [nsxt_policy_group.production_bms.path]
  destination_groups = [nsxt_policy_group.production_bms.path]
  action             = "ALLOW"
  services           = [nsxt_policy_service.ssh.path]
  logged             = true
}

resource "nsxt_policy_security_policy_rule" "default_deny" {
  display_name    = "default-deny"
  policy_path     = nsxt_policy_parent_security_policy.bms_parent.path
  sequence_number = 1000
  action          = "DROP"
  logged          = true
}
```

### Example — BMS Interface Security Policy

~> **NOTE:** `group_type = "BareMetalServer"` is **required** whenever a group contains `BareMetalServer` or `BareMetalServerInterface` member types. BMS groups may only contain BMS or BMSI members — mixing with `VirtualMachine` or other member types is not supported. Create separate policies for BMS and VM workloads.

```hcl
# BMS interface group — data-plane interfaces only
resource "nsxt_policy_group" "bms_data_only" {
  display_name = "BMS-Data-Interfaces"
  description  = "Data-plane interfaces for all BMS servers"
  group_type   = "BareMetalServer"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "network-type|data-plane"
    }
  }
}

resource "nsxt_policy_group" "trusted_management_hosts" {
  display_name = "Trusted-Management-Hosts"
  description  = "IP addresses of trusted management hosts"

  criteria {
    ipaddress_expression {
      ip_addresses = ["10.0.0.0/24"]
    }
  }
}

# Policy applied specifically to BMS data-plane interfaces
resource "nsxt_policy_security_policy" "bms_interface_policy" {
  display_name = "BMS-Interface-Policy"
  description  = "Security policy for BMS data-plane interfaces"
  category     = "Application"
  stateful     = true
  scope        = [nsxt_policy_group.bms_data_only.path]

  rule {
    display_name       = "allow_management_ssh"
    description        = "Allow SSH from trusted hosts to BMS servers"
    source_groups      = [nsxt_policy_group.trusted_management_hosts.path]
    destination_groups = [nsxt_policy_group.bms_data_only.path]
    action             = "ALLOW"
    services           = [nsxt_policy_service.ssh.path]
    logged             = true
  }

  rule {
    display_name    = "default_deny"
    description     = "Default deny on data-plane interfaces"
    action          = "DROP"
    logged          = true
    sequence_number = 1000
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

---

## Step 6: Audit Group Membership and Associations

Use the association and member-listing data sources to verify that servers and interfaces are in the right groups.

```hcl
# Which groups does a specific server belong to?
data "nsxt_policy_baremetal_server_group_associations" "server_memberships" {
  external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
}

output "server_group_paths" {
  value = [for g in data.nsxt_policy_baremetal_server_group_associations.server_memberships.groups : g.path]
}

# Which interfaces belong to a specific group?
data "nsxt_policy_group_baremetal_server_interface_members" "group_interfaces" {
  domain   = "default"
  group_id = nsxt_policy_group.web_tier_bms.id
}

output "group_interface_external_ids" {
  value = [for m in data.nsxt_policy_group_baremetal_server_interface_members.group_interfaces.items : m.external_id]
}

# Which servers belong to a specific group?
data "nsxt_policy_group_baremetal_server_members" "group_servers" {
  domain   = "default"
  group_id = nsxt_policy_group.production_bms.id
}

output "group_server_names" {
  value = [for m in data.nsxt_policy_group_baremetal_server_members.group_servers.items : m.display_name]
}

# Which groups does a specific interface belong to?
data "nsxt_policy_baremetal_server_interface_group_associations" "if_memberships" {
  external_id = "71be0142-2ed1-1d53-9c60-02005b4b7246"
}

output "interface_group_paths" {
  value = [for g in data.nsxt_policy_baremetal_server_interface_group_associations.if_memberships.groups : g.path]
}
```

### Using Group Associations as Policy Input

`nsxt_policy_baremetal_server_group_associations` does not only serve as an audit tool — you can feed the discovered group paths directly into a security policy rule's `source_groups` or `destination_groups`. This is useful when you want to allow or restrict traffic based on the groups a specific server already belongs to, without hard-coding group paths.

```hcl
# Discover all groups the server is currently a member of
data "nsxt_policy_baremetal_server_group_associations" "bms_groups" {
  external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
}

# Target destination group — e.g. a monitoring server group
resource "nsxt_policy_group" "monitoring_servers" {
  display_name = "Monitoring-Servers"
  description  = "Monitoring server IP addresses"

  criteria {
    ipaddress_expression {
      ip_addresses = ["10.1.0.0/24"]
    }
  }
}

resource "nsxt_policy_service" "monitoring" {
  display_name = "Monitoring-Port"
  l4_port_set_entry {
    protocol          = "TCP"
    destination_ports = ["9090", "9091"]
  }
}

# Allow the server's existing groups to reach the monitoring servers
resource "nsxt_policy_security_policy" "bms_monitoring_policy" {
  display_name = "BMS-Monitoring-Access"
  description  = "Allow discovered BMS groups to reach monitoring infrastructure"
  category     = "Application"

  rule {
    display_name       = "allow_bms_to_monitoring"
    description        = "Allow BMS group traffic to monitoring servers"
    source_groups      = data.nsxt_policy_baremetal_server_group_associations.bms_groups.groups[*].path
    destination_groups = [nsxt_policy_group.monitoring_servers.path]
    action             = "ALLOW"
    services           = [nsxt_policy_service.monitoring.path]
    logged             = true
    sequence_number    = 100
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Similarly, use interface group associations in rules
data "nsxt_policy_baremetal_server_interface_group_associations" "bms_if_groups" {
  external_id = "71be0142-2ed1-1d53-9c60-02005b4b7246"
}

resource "nsxt_policy_security_policy" "bms_interface_access_policy" {
  display_name = "BMS-Interface-Access"
  description  = "Policy using interface group membership as source"
  category     = "Application"

  rule {
    display_name       = "allow_interface_groups"
    source_groups      = data.nsxt_policy_baremetal_server_interface_group_associations.bms_if_groups.groups[*].path
    destination_groups = [nsxt_policy_group.monitoring_servers.path]
    action             = "ALLOW"
    services           = [nsxt_policy_service.monitoring.path]
    logged             = true
    sequence_number    = 100
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

---

## Advanced Use Cases

### Multi-Tier Application

```hcl
locals {
  tiers = ["web", "app", "db"]
}

# Tag servers by application tier
resource "nsxt_policy_baremetal_server_tags" "tier_tags" {
  for_each = {
    for s in data.nsxt_policy_baremetal_servers.all.results :
    s.external_id => (
      contains(local.tiers, split("-", lower(s.display_name))[0])
      ? split("-", lower(s.display_name))[0]
      : "general"
    )
  }

  external_id = each.key

  tag {
    scope = "app-tier"
    tag   = each.value
  }
}

# One group per tier
resource "nsxt_policy_group" "tier_groups" {
  for_each = toset(local.tiers)

  display_name = "BMS-${title(each.key)}-Tier"
  group_type   = "BareMetalServer"

  depends_on = [nsxt_policy_baremetal_server_tags.tier_tags]

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "app-tier|${each.key}"
    }
  }
}

# Tier-to-tier communication rules
resource "nsxt_policy_security_policy" "tier_communication" {
  display_name = "BMS-Tier-Communication"
  category     = "Application"

  rule {
    display_name       = "web-to-app"
    source_groups      = [nsxt_policy_group.tier_groups["web"].path]
    destination_groups = [nsxt_policy_group.tier_groups["app"].path]
    action             = "ALLOW"
    services           = [nsxt_policy_service.http.path, nsxt_policy_service.https.path]
    logged             = true
    sequence_number    = 100
  }

  rule {
    display_name       = "app-to-db"
    source_groups      = [nsxt_policy_group.tier_groups["app"].path]
    destination_groups = [nsxt_policy_group.tier_groups["db"].path]
    action             = "ALLOW"
    logged             = true
    sequence_number    = 200
  }

  rule {
    display_name       = "block-web-to-db"
    source_groups      = [nsxt_policy_group.tier_groups["web"].path]
    destination_groups = [nsxt_policy_group.tier_groups["db"].path]
    action             = "DROP"
    logged             = true
    sequence_number    = 300
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

### Interface Segmentation by Network Zone

```hcl
# Tag interfaces by network zone derived from naming
resource "nsxt_policy_baremetal_server_interface_tags" "zone_tags" {
  for_each = {
    for iface in data.nsxt_policy_baremetal_server_interfaces.all_interfaces.results :
    iface.external_id => (
      length(regexall("dmz", lower(iface.display_name))) > 0 ? "dmz" :
      length(regexall("internal", lower(iface.display_name))) > 0 ? "internal" : "general"
    )
    if !iface.is_mgmt_interface
  }

  external_id = each.key

  tag {
    scope = "network-zone"
    tag   = each.value
  }
}

resource "nsxt_policy_group" "dmz_interfaces" {
  display_name = "DMZ-BMS-Interfaces"
  group_type   = "BareMetalServer"

  depends_on = [nsxt_policy_baremetal_server_interface_tags.zone_tags]

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "network-zone|dmz"
    }
  }
}

resource "nsxt_policy_group" "internal_interfaces" {
  display_name = "Internal-BMS-Interfaces"
  group_type   = "BareMetalServer"

  depends_on = [nsxt_policy_baremetal_server_interface_tags.zone_tags]

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "network-zone|internal"
    }
  }
}
```

---

## Importing Existing Resources

Both BMS tag resources support import using the server or interface external ID.

```shell
# Import server tags by external ID
terraform import nsxt_policy_baremetal_server_tags.my_server 71be0142-2ed1-1d53-9c60-5564cf4b7e2e

# Import interface tags by external ID
terraform import nsxt_policy_baremetal_server_interface_tags.my_interface 71be0142-2ed1-1d53-9c60-02005b4b7246
```

After import, run `terraform plan` to verify the state matches expectations. If tags exist on the server that are not in your Terraform configuration, they will be removed on the next `terraform apply` (the resource manages all tags on the server).

---

## Drift Detection

Use the tag data sources to detect configuration drift — tags applied outside Terraform.

```hcl
data "nsxt_policy_baremetal_server_tags" "live_tags" {
  external_id = nsxt_policy_baremetal_server_tags.prod_web_01.id
}

resource "nsxt_policy_baremetal_server_tags" "prod_web_01" {
  external_id = "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"

  tag {
    scope = "environment"
    tag   = "production"
  }
}

# Compare live tags with expected — drift detected if lengths differ
output "drift_check" {
  value = {
    expected_count = 1
    actual_count   = length(data.nsxt_policy_baremetal_server_tags.live_tags.tag)
    is_drifted     = length(data.nsxt_policy_baremetal_server_tags.live_tags.tag) != 1
  }
}
```

---

## Best Practices

### Consistent Tagging Strategy

Define tag scopes as locals and reuse them across resources:

```hcl
locals {
  bms_tag_scopes = {
    environment = "environment"
    application = "application"
    tier        = "app-tier"
    managed_by  = "managed-by"
    owner       = "owner"
  }

  env_map = {
    prod = "production"
    dev  = "development"
    stg  = "staging"
    test = "testing"
  }
}

resource "nsxt_policy_baremetal_server_tags" "standard" {
  for_each = {
    for s in data.nsxt_policy_baremetal_servers.all.results :
    s.external_id => s
  }

  external_id = each.key

  tag {
    scope = local.bms_tag_scopes.managed_by
    tag   = "terraform"
  }

  tag {
    scope = local.bms_tag_scopes.environment
    tag   = lookup(local.env_map, split("-", lower(each.value.display_name))[0], "unknown")
  }
}
```

### Guard Against Missing Inventory

```hcl
resource "nsxt_policy_security_policy" "conditional_policy" {
  count = length(data.nsxt_policy_baremetal_servers.prod_servers.results) > 0 ? 1 : 0

  display_name = "Conditional-BMS-Policy"
  category     = "Application"

  rule {
    display_name       = "allow_internal"
    source_groups      = [nsxt_policy_group.production_bms.path]
    destination_groups = [nsxt_policy_group.production_bms.path]
    action             = "ALLOW"
    logged             = true
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

### Module Pattern

```hcl
# modules/bms-tier/variables.tf
variable "tier_name" { type = string }
variable "server_external_ids" { type = list(string) }

# modules/bms-tier/main.tf
resource "nsxt_policy_baremetal_server_tags" "tags" {
  for_each    = toset(var.server_external_ids)
  external_id = each.key

  tag { scope = "app-tier"; tag = var.tier_name }
}

resource "nsxt_policy_group" "group" {
  display_name = "BMS-${var.tier_name}"
  group_type   = "BareMetalServer"

  depends_on = [nsxt_policy_baremetal_server_tags.tags]

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "app-tier|${var.tier_name}"
    }
  }
}

# modules/bms-tier/outputs.tf
output "group_path" { value = nsxt_policy_group.group.path }

# root/main.tf
module "web_tier" {
  source              = "./modules/bms-tier"
  tier_name           = "web"
  server_external_ids = ["uuid-1", "uuid-2"]
}

module "app_tier" {
  source              = "./modules/bms-tier"
  tier_name           = "app"
  server_external_ids = ["uuid-3", "uuid-4"]
}
```

---

## Troubleshooting

### No BMS Servers Discovered

- Verify that a compute manager is configured in NSX-T and bare metal servers are registered.
- Check that the NSX version is 9.0.0 or higher.
- Run `data "nsxt_policy_baremetal_servers" "all" {}` with no filters to confirm the inventory is non-empty.

### Group Membership Not Reflecting Expected Servers

- Ensure tag scope and value strings match exactly (case-sensitive).
- Tags must be applied before the group is evaluated. Add a `depends_on` from the group to the tag resource.
- Use `nsxt_policy_group_baremetal_server_members` data source to inspect actual effective membership.

### Policy Apply Fails with BMS Group

- Confirm `group_type = "BareMetalServer"` is set on every group containing BMS/BMSI members.
- Confirm the policy `category` is not `Ethernet` (unsupported for BMS).
- Confirm no `context_profiles` are referenced in rules that target BMS groups.

### Import Differences After `terraform plan`

- If tags exist on the server outside of Terraform, they appear as planned deletions. Either add them to the Terraform configuration or accept removal on next apply.
- Check that the `external_id` used for import matches the one in your configuration.

### Debugging Inventory and Tags

```hcl
# Dump full server inventory
output "debug_servers" {
  value = {
    for s in data.nsxt_policy_baremetal_servers.all_servers.results :
    s.external_id => {
      display_name = s.display_name
      os           = "${s.os_name} ${s.os_version}"
      cpu_cores    = s.cpu_cores
      tags         = s.tags
    }
  }
}

# Dump full interface inventory
output "debug_interfaces" {
  value = {
    for iface in data.nsxt_policy_baremetal_server_interfaces.all_interfaces.results :
    iface.external_id => {
      display_name      = iface.display_name
      bms_external_id   = iface.bms_external_id
      is_mgmt_interface = iface.is_mgmt_interface
      state             = iface.state
      ip_addresses      = iface.ip_addresses
    }
  }
}
```
