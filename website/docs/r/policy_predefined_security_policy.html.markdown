---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_predefined_security_policy"
description: A resource to update Predefined (Default) Security Security Policies.
---

# nsxt_policy_predefined_security_policy

This resource provides a method to modify default Security Policy and its rules.
This can be default layer2 policy or default layer3 policy. Maximum one resource
for each type should exist in your configuration.

~> **NOTE:** An absolute path, such as `/infra/domains/default/security-policies/default-layer3-section`, can be provided for this resource (this approach will work slightly faster, as the roundtrip for data source retrieval will be spared) In the example below a data source is used in order to pull the predefined policy.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_security_policy" "default_l3" {
  is_default = true
  category   = "Application"
}

resource "nsxt_policy_predefined_security_policy" "test" {
  path = data.nsxt_policy_security_policy.default_l3.path

  tag {
    scope = "color"
    tag   = "orange"
  }

  rule {
    display_name       = "allow_icmp"
    destination_groups = [nsxt_policy_group.cats.path, nsxt_policy_group.dogs.path]
    action             = "ALLOW"
    services           = [nsxt_policy_service.icmp.path]
    logged             = true
  }

  rule {
    display_name     = "allow_udp"
    source_groups    = [nsxt_policy_group.fish.path]
    sources_excluded = true
    scope            = [nsxt_policy_group.aquarium.path]
    action           = "ALLOW"
    services         = [nsxt_policy_service.udp.path]
    logged           = true
    disabled         = true
  }

  default_rule {
    action = "DROP"
  }

}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_security_policy" "default_l3" {
  is_default = true
  category   = "Application"
}

resource "nsxt_policy_predefined_security_policy" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  path = data.nsxt_policy_security_policy.default_l3.path

  tag {
    scope = "color"
    tag   = "orange"
  }

  rule {
    display_name       = "allow_icmp"
    destination_groups = [nsxt_policy_group.cats.path, nsxt_policy_group.dogs.path]
    action             = "ALLOW"
    services           = [nsxt_policy_service.icmp.path]
    logged             = true
  }

  rule {
    display_name     = "allow_udp"
    source_groups    = [nsxt_policy_group.fish.path]
    sources_excluded = true
    scope            = [nsxt_policy_group.aquarium.path]
    action           = "ALLOW"
    services         = [nsxt_policy_service.udp.path]
    logged           = true
    disabled         = true
  }

  default_rule {
    action = "DROP"
  }
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) Policy path for the predefined Security Policy to modify.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Security Policy.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `rule` (Optional) A repeatable block to specify rules for the Security Policy. This setting is applicable to non-Default policies only. Each rule includes the following fields:
  * `display_name` - (Required) Display name of the resource.
  * `description` - (Optional) Description of the resource.
  * `destination_groups` - (Optional) Set of group paths that serve as the destination for this rule. IPs, IP ranges, or CIDRs may also be used starting in NSX-T 3.0. An empty set can be used to specify "Any".
  * `destinations_excluded` - (Optional) A boolean value indicating negation of destination groups.
  * `direction` - (Optional) The traffic direction for the policy. Must be one of: `IN`, `OUT` or `IN_OUT`. Defaults to `IN_OUT`.
  * `disabled` - (Optional) A boolean value to indicate the rule is disabled. Defaults to `false`.
  * `ip_version` - (Optional) The IP Protocol for the rule. Must be one of: `IPV4`, `IPV6` or `IPV4_IPV6`. Defaults to `IPV4_IPV6`.
  * `logged` - (Optional) A boolean flag to enable packet logging.
  * `notes` - (Optional) Text for additional notes on changes for the rule.
  * `profiles` - (Optional) A list of profiles for the rule.
  * `scope` - (Required) List of policy paths where the rule is applied.
  * `services` - (Optional) List of services to match.
  * `source_groups` - (Optional) Set of group paths that serve as the source for this rule. IPs, IP ranges, or CIDRs may also be used starting in NSX-T 3.0. An empty set can be used to specify "Any".
  * `source_excluded` - (Optional) A boolean value indicating negation of source groups.
  * `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this Rule.
  * `action` - (Optional) The action for the Rule. Must be one of: `ALLOW`, `DROP` or `REJECT`. Defaults to `ALLOW`. Note that `REJECT` action is not applicable for L2 policy.

* `default_rule` (Optional) A repeatable block to modify default rules for the Security Policy in a `DEFAULT` category. Each rule includes the following fields:
  * `description` - (Optional) Description of the resource.
  * `logged` - (Optional) A boolean flag to enable packet logging.
  * `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this Rule.
  * `action` - (Optional) The action for the Rule. Must be one of: `ALLOW`, `DROP` or `REJECT`. Note that `REJECT` action is not applicable for L2 policy.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `rule`, `default_rule`:
  * `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
  * `path` - The NSX path of the policy resource.
  * `sequence_number` - Sequence number of the this rule, is defined by order of rules in the list.
  * `rule_id` - Unique positive number that is assigned by the system and is useful for debugging.

