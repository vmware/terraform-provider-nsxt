---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_policy"
description: A data source to retrieve an Intrusion Service (IDS) Policy.
---

# nsxt_policy_intrusion_service_policy

This data source provides information about an existing Intrusion Service Policy configured on NSX.
It can be useful for fetching the IDPS DFW Policy including all its embedded rules.

~> **NOTE:** This data source retrieves the policy including embedded rules, allowing you to refer IDPS DFW policy details and its rules in other resources. For different use cases, consider:

* `nsxt_policy_parent_intrusion_service_policy` - For **IDPS DFW Policy metadata only (no rules)**
* `nsxt_policy_intrusion_service_policy_rule` - For **individual standalone IDPS DFW rules**

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_intrusion_service_policy" "ids_policy" {
  display_name = "intrusion-service-policy"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_intrusion_service_policy" "ids_policy" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "intrusion-service-policy"
}
```

## Argument Reference

* `id` - (Optional) The ID of the policy to retrieve.
* `display_name` - (Optional) The display name of the policy to retrieve.
* `domain` - (Optional) The domain of the policy. Defaults to `default`.
* `category` - (Optional) Category of the policy.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `stateful` - When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed. Note: Intrusion Service Policies are always stateful.
* `sequence_number` - This field is used to resolve conflicts between multiple policies that have rules that match the same packet.
* `locked` - Indicates whether a security policy should be locked.
* `tag` - A list of scope + tag pairs to associate with this policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `comments` - Comments for security policy lock/unlock.
* `rule` - A list of rules in this policy. Each rule contains:
    * `display_name` - The display name of the rule.
    * `description` - The description of the rule.
    * `notes` - Additional notes for the rule.
    * `action` - Action for this rule (DETECT, DETECT_PREVENT, EXEMPT). Note: EXEMPT action requires NSX version 9.1.0 or higher.
    * `source_groups` - List of source groups.
    * `destination_groups` - List of destination groups.
    * `services` - List of services.
    * `scope` - List of objects where the rule is enforced.
    * `logged` - Flag to enable logging.
    * `log_label` - Additional information which will be propagated to the rule syslog.
    * `disabled` - Flag to disable the rule.
    * `direction` - Traffic direction.
    * `ip_version` - IP version.
    * `sources_excluded` - Flag to indicate whether sources are negated.
    * `destinations_excluded` - Flag to indicate whether destinations are negated.
    * `ids_profiles` - List of IDS profiles for this rule.
    * `oversubscription` - Indicates how rule performs when oversubscribed.
    * `sequence_number` - The sequence number of the rule.
    * `rule_id` - Unique positive number assigned by the system.
    * `revision` - Indicates current revision number of the rule.
    * `tag` - A list of scope + tag pairs associated with this rule.
