---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_predefined_gateway_policy"
description: A resource to update Predefined Gateway Security Policies.
---

# nsxt_policy_predefined_gateway_policy

This resource provides a method to fine-tune a pre-created Gateway Policy and its rules.
There are two separate use cases for this resource:
* Modify certain settings of default Gateway Policy and its Default Rule.
* Modify predefined Gateway Policy that is not listed under Default category, and add rules to it.
  This use case is relevant for VMC.

~> **NOTE:** Importing the resource first is recommended if you with to reconfigure the policy from scratch. Terraform state does is not aware of attributes/rules that are already configured on your NSX!

~> **NOTE:** An absolute path can be provided for this resource (this approach will work slightly faster, as the roundtrip for data source retrieval will be spared). In one of the examples below a data source is used in order to pull the predefined policy, while the other uses absolute path.

~> **NOTE:** Default gateway policy generation behavior have changed in NSX 3.1.0. Below this version, there is a single default policy, while default rules are created under it per Gateway. Above NSX 3.1.0, a default policy is generated per Gateway. The first example provided here is limited to NSX 3.0.0 and below.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage for NSX Policy Manager

```hcl
data "nsxt_policy_gateway_policy" "default" {
  category = "Default"
}

resource "nsxt_policy_predefined_gateway_policy" "test" {
  path = data.nsxt_policy_gateway_policy.default.path

  tag {
    scope = "color"
    tag   = "orange"
  }

  default_rule {
    scope     = nsxt_policy_tier0_gateway.main.path
    logged    = true
    log_label = "orange default"
    action    = "ALLOW"
  }

}
```

## Example Usage for VMC
~> **NOTE:** VMC may auto-generate rules under the default gateway policy. If you want these rules to remain post-apply, please copy them into your terraform configuration, including matching `nsx_id`. Sample below includes example of such auto-generated rule (`default-vti-rule`). If you prefer to re-configure your policy from scratch and dispose of generated rules, it is recommended to import the policy first (otherwise you may need to apply twice in order to sync backend state with terraform configuration).

```hcl
resource "nsxt_policy_predefined_gateway_policy" "test" {
  path = "/infra/domains/cgw/gateway-policies/default"

  rule {
    display_name = "Default VTI Rule"
    nsx_id       = "default-vti-rule"
    action       = "DROP"
    scope        = ["/infra/labels/cgw-vpn"]
  }

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

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_gateway_policy" "default" {
  category = "Default"
}

resource "nsxt_policy_predefined_gateway_policy" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  path = data.nsxt_policy_gateway_policy.default.path

  tag {
    scope = "color"
    tag   = "orange"
  }

  default_rule {
    scope     = nsxt_policy_tier0_gateway.main.path
    logged    = true
    log_label = "orange default"
    action    = "ALLOW"
  }
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) Policy path for the predefined Gateway Policy to modify.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Gateway Policy.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `rule` (Optional) A repeatable block to specify rules for the Gateway Policy. This setting is not applicable to policy belonging to `DEFAULT` category. Each rule includes the following fields:
  * `display_name` - (Required) Display name of the resource.
  * `description` - (Optional) Description of the resource.
  * `destination_groups` - (Optional) Set of group paths that serve as the destination for this rule. IPs, IP ranges, or CIDRs may also be used starting in NSX-T 3.0. An empty set can be used to specify "Any".
  * `destinations_excluded` - (Optional) A boolean value indicating negation of destination groups.
  * `direction` - (Optional) The traffic direction for the policy. Must be one of: `IN`, `OUT` or `IN_OUT`. Defaults to `IN_OUT`.
  * `disabled` - (Optional) A boolean value to indicate the rule is disabled. Defaults to `false`.
  * `ip_version` - (Optional) The IP Protocol for the rule. Must be one of: `IPV4`, `IPV6` or `IPV4_IPV6`. Defaults to `IPV4_IPV6`.
  * `logged` - (Optional) A boolean flag to enable packet logging.
  * `notes` - (Optional) Text for additional notes on changes for the rule.
  * `profiles` - (Optional) A list of context profiles for the rule. Note: due to platform issue, this setting is only supported with NSX 3.2 onwards.
  * `scope` - (Required) List of policy paths where the rule is applied.
  * `services` - (Optional) List of services to match.
  * `source_groups` - (Optional) Set of group paths that serve as the source for this rule. IPs, IP ranges, or CIDRs may also be used starting in NSX-T 3.0. An empty set can be used to specify "Any".
  * `source_excluded` - (Optional) A boolean value indicating negation of source groups.
  * `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this Rule.
  * `action` - (Optional) The action for the Rule. Must be one of: `ALLOW`, `DROP` or `REJECT`. Defaults to `ALLOW`.
* `default_rule` (Optional) A repeatable block to modify default rules for the Gateway Policy in a `DEFAULT` category. Each rule includes the following fields:
  * `scope` - (Required) Scope for the default rule that should be modified. Only one default rule can be present for each scope.
  * `description` - (Optional) Description of the resource.
  * `logged` - (Optional) A boolean flag to enable packet logging.
  * `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this Rule.
  * `action` - (Optional) The action for the Rule. Must be one of: `ALLOW`, `DROP` or `REJECT`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `rule`, `default_rule`:
  * `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
  * `path` - The NSX path of the policy resource.
  * `sequence_number` - Sequence number of the this rule, is defined by order of rules in the list.
  * `rule_id` - Unique positive number that is assigned by the system and is useful for debugging.

## Importing

An existing Gateway Policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_gateway_policy.default POLICY_PATH
```
The above command imports the policy Gateway Policy named `default` with policy path `POLICY_PATH`.
The import command is recommended in case the NSX policy in question already has rules configured, and you wish to reconfigure the policy from scratch. If your terraform configuration copies existing rules, like in VMC example above, import step can be skipped.
