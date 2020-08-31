---
subcategory: "Policy - Firewall"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_fixed_gateway_policy"
description: A resource to Fixed Gateway Security Policies.
---

# nsxt_policy_fixed_gateway_policy

This resource provides a method to fine-tune a pre-created Gateway Policy and its rules.
There are two separate use cases for this resource:
* Modify certain settings of default Gateway Policy and its Default Rule.
* Modify fixed Gateway Policy that is not listed under Default category, and add rules to it.
  This use case is relevant for VMC.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage for NSX Policy Manager

```hcl
resource "nsxt_policy_fixed_gateway_policy" "test" {
  path = "/infra/domains/default/gateway-policies/Policy_Default_Infra"

  tag {
    scope = "color"
    tag   = "orange"
  }

  default_rule {
    scope     = nsxt_policy_tier0_gateway.main.path
    logged    = true
    log_label = "orange default"
    action  = "ALLOW"
  }

}
```

## Example Usage for VMC

```hcl
resource "nsxt_policy_fixed_gateway_policy" "test" {
  path = "/infra/domains/cgw/gateway-policies/default"

  rule {
    display_name  = "Allow ICMP"
    services      = [nsxt_policy_service.icmp.path]
    source_groups = [nsxt_policy_group.green.path]
    logged        = true
    log_label     = "icmp"
    action        = "ALLOW"
  }
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) Policy path for the fixed Gateway Policy to modify.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Gateway Policy.
* `rule` (Optional) A repeatable block to specify rules for the Gateway Policy. This setting is applicable to non-Default policies only. Each rule includes the following fields:
  * `display_name` - (Required) Display name of the resource.
  * `description` - (Optional) Description of the resource.
  * `destination_groups` - (Optional) A list of destination group paths to use for the policy.
  * `destinations_excluded` - (Optional) A boolean value indicating negation of destination groups.
  * `direction` - (Optional) The traffic direction for the policy. Must be one of: `IN`, `OUT` or `IN_OUT`. Defaults to `IN_OUT`.
  * `disabled` - (Optional) A boolean value to indicate the rule is disabled. Defaults to `false`.
  * `ip_version` - (Optional) The IP Protocol for the rule. Must be one of: `IPV4`, `IPV6` or `IPV4_IPV6`. Defaults to `IPV4_IPV6`.
  * `logged` - (Optional) A boolean flag to enable packet logging.
  * `notes` - (Optional) Text for additional notes on changes for the rule.
  * `profiles` - (Optional) A list of profiles for the rule.
  * `scope` - (Required) List of policy paths where the rule is applied.
  * `services` - (Optional) List of services to match.
  * `source_groups` - (Optional) A list of source group paths to use for the policy.
  * `source_excluded` - (Optional) A boolean value indicating negation of source groups.
  * `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this Rule.
  * `action` - (Optional) The action for the Rule. Must be one of: `ALLOW`, `DROP` or `REJECT`. Defaults to `ALLOW`.
* `defaultrule` (Optional) A repeatable block to modify default rules for the Gateway Policy in a `DEFAULT` category. Each rule includes the following fields:
  * `scope` - (Required) Scope for the default rule that should be modified. Only one default rule can be present for each scope.
  * `description` - (Optional) Description of the resource.
  * `logged` - (Optional) A boolean flag to enable packet logging.
  * `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this Rule.
  * `action` - (Optional) The action for the Rule. Must be one of: `ALLOW`, `DROP` or `REJECT`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Secuirty Policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `rule`:
  * `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
  * `path` - The NSX path of the policy resource.
  * `sequence_number` - Sequence number of the this rule, is defined by order of rules in the list.
  * `rule_id` - Unique positive number that is assigned by the system and is useful for debugging.

