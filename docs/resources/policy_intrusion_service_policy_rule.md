---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_policy_rule"
description: A resource to configure a single rule in an Intrusion Service (IDS) Policy.
---

# nsxt_policy_intrusion_service_policy_rule

This resource provides a method for the management of a single rule in an Intrusion Service (IDS) Policy for East-West traffic inspection (Distributed Firewall context).

Note: to avoid unexpected behavior, don't use this resource and resource `nsxt_policy_intrusion_service_policy` to manage rules under an intrusion service policy at the same time.
Instead, please use this resource with resource `nsxt_policy_parent_intrusion_service_policy` to manage an intrusion service policy and its rules separately, and use `nsxt_policy_intrusion_service_policy` to manage a policy and its rules in one single resource.

This resource is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_intrusion_service_profile" "default" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_parent_intrusion_service_policy" "parent" {
  display_name    = "tf-intrusion-svc-parent-policy"
  description     = "Parent policy for standalone IDS rules"
  category        = "ThreatRules"
  locked          = false
  sequence_number = 3
  stateful        = true

  lifecycle {
    create_before_destroy = true
  }
}

resource "nsxt_policy_intrusion_service_policy_rule" "detect_rule" {
  display_name       = "detect-threats"
  description        = "Detect threats in East-West traffic"
  policy_path        = nsxt_policy_parent_intrusion_service_policy.parent.path
  action             = "DETECT"
  direction          = "IN_OUT"
  ip_version         = "IPV4_IPV6"
  oversubscription   = "INHERIT_GLOBAL"
  sequence_number    = 1
  source_groups      = [nsxt_policy_group.web_servers.path]
  destination_groups = [nsxt_policy_group.db_servers.path]
  services           = [nsxt_policy_service.http.path]
  ids_profiles       = [data.nsxt_policy_intrusion_service_profile.default.path]
  logged             = true
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_parent_intrusion_service_policy" "parent" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name    = "tf-intrusion-svc-parent-policy"
  locked          = false
  sequence_number = 3
  stateful        = true
}

resource "nsxt_policy_intrusion_service_policy_rule" "detect_rule" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name    = "detect-threats"
  policy_path     = nsxt_policy_parent_intrusion_service_policy.parent.path
  action          = "DETECT"
  direction       = "IN_OUT"
  ip_version      = "IPV4_IPV6"
  sequence_number = 1
  ids_profiles    = [data.nsxt_policy_intrusion_service_profile.default.path]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `policy_path` - (Required) Path of the Intrusion Service Policy this rule belongs to.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to. If not provided, it will be derived from `policy_path`.
    * `project_id` - (Required) The ID of the project which the object belongs to
* `destination_groups` - (Optional) Set of group paths that serve as the destination for this rule. IPs, IP ranges, or CIDRs may also be used starting in NSX-T 3.0. An empty set can be used to specify `ANY`. Default is `ANY`.
* `destinations_excluded` - (Optional) A boolean value indicating negation of destination groups. Default is `false`.
* `direction` - (Optional) The traffic direction for the rule. Must be one of: `IN`, `OUT` or `IN_OUT`. Default is `IN_OUT`.
* `disabled` - (Optional) A boolean value to indicate the rule is disabled. Default is `false`.
* `ip_version` - (Optional) The IP Protocol for the rule. Must be one of: `IPV4`, `IPV6` or `IPV4_IPV6`. Default is `IPV4_IPV6`.
* `logged` - (Optional) A boolean flag to enable packet logging. Default is `false`.
* `notes` - (Optional) Text for additional notes on changes for this rule.
* `scope` - (Optional) Set of policy object paths where the rule is applied for East-West traffic inspection.
* `services` - (Optional) Set of service paths to match for this rule. An empty set can be used to specify `ANY`. Default is `ANY`.
* `source_groups` - (Optional) Set of group paths that serve as the source for this rule. IPs, IP ranges, or CIDRs may also be used starting in NSX-T 3.0. An empty set can be used to specify `ANY`. Default is `ANY`.
* `sources_excluded` - (Optional) A boolean value indicating negation of source groups. Default is `false`.
* `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog for this rule.
* `tag` - (Optional) A list of scope + tag pairs to associate with this rule.
* `action` - (Optional) Rule action for intrusion detection/prevention. One of `DETECT`, `DETECT_PREVENT`, or `EXEMPT`. Default is `DETECT`. Note: `EXEMPT` action is only available from NSX version 9.1.0 onwards and allows traffic to bypass intrusion detection/prevention.
* `oversubscription` - (Optional) Action to take when IDPS engine is oversubscribed. One of `BYPASSED`, `DROPPED` or `INHERIT_GLOBAL`. Default is `INHERIT_GLOBAL`. `BYPASSED`: Traffic bypasses IDPS when oversubscribed. `DROPPED`: Traffic is dropped when oversubscribed. `INHERIT_GLOBAL`: Inherit the behavior from the global IDPS settings.
* `ids_profiles` - (Required) Set of IDS profile paths for this rule. These profiles define the intrusion detection signatures to be applied.
* `sequence_number` - (Required) Sequence number to determine the order of rule processing within the parent policy.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX path of the policy resource.
* `rule_id` - Unique positive number that is assigned by the system and is useful for debugging.

## Importing

An existing Intrusion Service Policy Rule can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_intrusion_service_policy_rule.detect_rule RULE_PATH
```

Example: `/infra/domains/default/intrusion-service-policies/ids-policy/rules/detect-rule`.
