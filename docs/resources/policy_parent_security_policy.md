---
subcategory: "Firewall"
page_title: "NSXT: nsxt_policy_parent_security_policy"
description: A resource to configure a Security Policy without rules.
---

# nsxt_policy_parent_security_policy

This resource provides a method for the management of Security Policy without rules.

Note: to avoid unexpected behavior, don't use this resource and resource `nsxt_policy_security_policy` to manage the same Security Policy at the same time.
To config rules under this resource, please use resource `nsxt_policy_security_policy_rule` to manage rules separately.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_parent_security_policy" "policy1" {
  display_name = "policy1"
  description  = "Terraform provisioned Security Policy"
  category     = "Application"
  locked       = false
  stateful     = true
  tcp_strict   = false
  scope        = [nsxt_policy_group.pets.path]

  lifecycle {
    create_before_destroy = true
  }
}

resource "nsxt_policy_security_policy_rule" "rule1" {
  display_name       = "rule1"
  description        = "Terraform provisioned Security Policy Rule"
  policy_path        = nsxt_policy_parent_security_policy.policy1.path
  sequence_number    = 1
  destination_groups = [nsxt_policy_group.cats.path, nsxt_policy_group.dogs.path]
  action             = "DROP"
  services           = [nsxt_policy_service.icmp.path]
  logged             = true
}
```

## Global Manager example

```hcl
data "nsxt_policy_site" "paris" {
  display_name = "Paris"
}
resource "nsxt_policy_parent_security_policy" "policy1" {
  display_name = "policy1"
  description  = "Terraform provisioned Security Policy"
  category     = "Application"
  locked       = false
  stateful     = true
  tcp_strict   = false
  scope        = [nsxt_policy_group.pets.path]
  domain       = data.nsxt_policy_site.paris.id

  lifecycle {
    create_before_destroy = true
  }
}

resource "nsxt_policy_security_policy_rule" "rule1" {
  display_name       = "rule1"
  description        = "Terraform provisioned Security Policy Rule"
  policy_path        = nsxt_policy_parent_security_policy.policy1.path
  sequence_number    = 1
  destination_groups = [nsxt_policy_group.cats.path, nsxt_policy_group.dogs.path]
  action             = "DROP"
  services           = [nsxt_policy_service.icmp.path]
  logged             = true
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_parent_security_policy" "policy1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "policy1"
  description  = "Terraform provisioned Security Policy"
  category     = "Application"
  locked       = false
  stateful     = true
  tcp_strict   = false
  scope        = [nsxt_policy_group.pets.path]

  lifecycle {
    create_before_destroy = true
  }
}

resource "nsxt_policy_security_policy_rule" "rule1" {
  display_name       = "rule1"
  description        = "Terraform provisioned Security Policy Rule"
  policy_path        = nsxt_policy_parent_security_policy.policy1.path
  sequence_number    = 1
  destination_groups = [nsxt_policy_group.cats.path, nsxt_policy_group.dogs.path]
  action             = "DROP"
  services           = [nsxt_policy_service.icmp.path]
  logged             = true
}
```

## Example Usage - Bare Metal Server Parent Policy

```hcl
# Create BMS groups for policy scope
resource "nsxt_policy_group" "production_bms" {
  display_name = "Production-BMS-Servers"
  description  = "Production bare metal servers"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "environment|production"
    }
  }
}

resource "nsxt_policy_group" "bms_data_interfaces" {
  display_name = "BMS-Data-Interfaces"
  description  = "BMS data plane interfaces"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "network-type|data-plane"
    }
  }
}

# Parent policy scoped to BMS groups
resource "nsxt_policy_parent_security_policy" "bms_parent_policy" {
  display_name = "BMS-Parent-Policy"
  description  = "Parent policy for bare metal server security"
  category     = "Application"
  locked       = false
  stateful     = true
  tcp_strict   = false
  scope = [
    nsxt_policy_group.production_bms.path,
    nsxt_policy_group.bms_data_interfaces.path
  ]

  lifecycle {
    create_before_destroy = true
  }
}

# Standalone rules under the BMS parent policy
resource "nsxt_policy_security_policy_rule" "allow_bms_internal" {
  display_name       = "allow-bms-internal"
  description        = "Allow internal BMS communication"
  policy_path        = nsxt_policy_parent_security_policy.bms_parent_policy.path
  sequence_number    = 100
  source_groups      = [nsxt_policy_group.production_bms.path]
  destination_groups = [nsxt_policy_group.production_bms.path]
  action             = "ALLOW"
  services           = ["SSH", "HTTP", "HTTPS"]
  logged             = true
}

resource "nsxt_policy_security_policy_rule" "block_mgmt_on_data" {
  display_name    = "block-mgmt-on-data"
  description     = "Block management traffic on data plane interfaces"
  policy_path     = nsxt_policy_parent_security_policy.bms_parent_policy.path
  sequence_number = 200
  scope           = [nsxt_policy_group.bms_data_interfaces.path]
  action          = "DROP"
  services        = ["SSH", "SNMP", "TELNET"]
  logged          = true
}
```

## Example Usage - Mixed BMS and VM Parent Policy

```hcl
# Create mixed BMS and VM groups
resource "nsxt_policy_group" "all_app_servers" {
  display_name = "All-Application-Servers"
  description  = "Application servers across BMS and VMs"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "tier|application"
    }
  }

  conjunction {
    operator = "OR"
  }

  criteria {
    condition {
      key         = "Tag"
      member_type = "VirtualMachine"
      operator    = "EQUALS"
      value       = "tier|application"
    }
  }
}

resource "nsxt_policy_group" "web_tier_bms" {
  display_name = "Web-Tier-BMS"
  description  = "BMS servers in web tier (static group)"

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

# Parent policy with mixed scope
resource "nsxt_policy_parent_security_policy" "mixed_environment_policy" {
  display_name = "Mixed-Environment-Policy"
  description  = "Parent policy for mixed BMS and VM environment"
  category     = "Application"
  stateful     = true
  scope = [
    nsxt_policy_group.all_app_servers.path,
    nsxt_policy_group.web_tier_bms.path
  ]

  lifecycle {
    create_before_destroy = true
  }
}

# Rules under the mixed environment policy
resource "nsxt_policy_security_policy_rule" "allow_web_traffic" {
  display_name       = "allow-web-traffic"
  description        = "Allow web traffic to all web servers"
  policy_path        = nsxt_policy_parent_security_policy.mixed_environment_policy.path
  sequence_number    = 100
  destination_groups = [nsxt_policy_group.web_tier_bms.path]
  action             = "ALLOW"
  services           = ["HTTP", "HTTPS"]
  logged             = true
}

resource "nsxt_policy_security_policy_rule" "app_tier_communication" {
  display_name       = "app-tier-communication"
  description        = "Allow communication within application tier"
  policy_path        = nsxt_policy_parent_security_policy.mixed_environment_policy.path
  sequence_number    = 200
  source_groups      = [nsxt_policy_group.all_app_servers.path]
  destination_groups = [nsxt_policy_group.all_app_servers.path]
  action             = "ALLOW"
  services           = ["HTTP", "HTTPS", "MYSQL", "POSTGRESQL"]
  logged             = true
}
```

-> We recommend using `lifecycle` directive as in samples above, in order to avoid dependency issues when updating groups/services simultaneously with the rule.

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `domain` - (Optional) The domain to use for the resource. This domain must already exist. For VMware Cloud on AWS use `cgw`. For Global Manager, please use site id for this field. If not specified, this field is default to `default`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `category` - (Required) Category of this policy. For local manager must be one of `Ethernet`, `Emergency`, `Infrastructure`, `Environment`, `Application`. For global manager must be one of: `Infrastructure`, `Environment`, `Application`.
* `comments` - (Optional) Comments for security policy lock/unlock.
* `locked` - (Optional) Indicates whether a security policy should be locked. If locked by a user, no other user would be able to modify this policy.
* `scope` - (Optional) The list of policy object paths where the rules in this policy will get applied.
* `sequence_number` - (Optional) This field is used to resolve conflicts between security policies across domains.
* `stateful` - (Optional) If true, state of the network connects are tracked and a stateful packet inspection is performed. Default is true.
* `tcp_strict` - (Optional) Ensures that a 3 way TCP handshake is done before the data packets are sent. Default is false.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Security Policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing security policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_parent_security_policy.policy1 domain/ID
```

The above command imports the security policy named `policy1` under NSX domain `domain` with the NSX Policy ID `ID`.
