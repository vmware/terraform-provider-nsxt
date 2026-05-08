---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_policy_rule"
description: A data source to retrieve an Intrusion Service (IDS) Policy Rule.
---

# nsxt_policy_intrusion_service_policy_rule

This data source provides information about an existing Intrusion Service Policy Rule configured on NSX.
It can be useful to retrieve individual IDPS DFW rules that are managed separately from their parent policy.

~> **NOTE:** This data source retrieves **standalone rules** that are managed separately from their parent policy, allowing you to refer specific IDPS DFW rule in other resources. For different use cases, consider:

* `nsxt_policy_intrusion_service_policy` - For **IDPS DFW Policy with embedded rules**
* `nsxt_policy_parent_intrusion_service_policy` - For **parent IDPS DFW Policy metadata only**

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
# Get parent policy for its path
data "nsxt_policy_parent_intrusion_service_policy" "ids_policy" {
  display_name = "my-ids-policy"
}

# Get individual rule from that policy
data "nsxt_policy_intrusion_service_policy_rule" "ids_rule" {
  display_name = "detect-threats"
  policy_path  = data.nsxt_policy_parent_intrusion_service_policy.ids_policy.path
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

# Get parent policy for its path
data "nsxt_policy_parent_intrusion_service_policy" "ids_policy" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "my-ids-policy"
}

# Get individual rule from that policy
data "nsxt_policy_intrusion_service_policy_rule" "ids_rule" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "detect-threats"
  policy_path  = data.nsxt_policy_parent_intrusion_service_policy.ids_policy.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the rule to retrieve.
* `display_name` - (Optional) The display name of the rule to retrieve.
* `policy_path` - (Required) The path of the parent policy containing this rule.
* `domain` - (Optional) The domain of the policy containing this rule. Defaults to `default`.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the rule resource.
* `sequence_number` - The sequence number of the rule.
* `source_groups` - List of source groups.
* `destination_groups` - List of destination groups.
* `services` - List of services.
* `scope` - List of policy objects where the rule is enforced.
* `action` - Action for this rule.
* `direction` - Traffic direction.
* `disabled` - Flag to disable the rule.
* `ip_version` - IP version.
* `logged` - Flag to enable logging.
* `log_label` - Additional information which will be propagated to the rule syslog.
* `notes` - Text for additional notes on changes for the rule.
* `sources_excluded` - Flag to indicate whether sources are negated.
* `destinations_excluded` - Flag to indicate whether destinations are negated.
* `ids_profiles` - List of IDS profiles for this rule.
* `oversubscription` - Indicates how rule performs when oversubscribed.
* `rule_id` - Unique positive number that is assigned by the system and is useful for debugging.
* `tag` - A list of scope + tag pairs to associate with this rule.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
