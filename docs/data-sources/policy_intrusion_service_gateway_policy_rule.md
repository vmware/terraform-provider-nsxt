---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_gateway_policy_rule"
description: A data source to retrieve an Intrusion Service Gateway Policy Rule.
---

# nsxt_policy_intrusion_service_gateway_policy_rule

This data source provides information about an existing Intrusion Service Gateway Policy Rule configured on NSX.
It can be useful to retrieve individual IDPS Gateway rules that are managed separately from their parent policy.

~> **NOTE:** This data source retrieves **standalone rules** that are managed separately from their parent policy, allowing you to refer specific IDPS Gateway rule in other resources. For different use cases, consider:

* `nsxt_policy_intrusion_service_gateway_policy` - For **IDPS Gateway policy with embedded rules**
* `nsxt_policy_parent_intrusion_service_gateway_policy` - For **parent IDPS Gateway policy metadata only**

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
# Get parent Gateway policy for its path
data "nsxt_policy_parent_intrusion_service_gateway_policy" "ids_gw_policy" {
  display_name = "my-ids-gateway-policy"
}

# Get individual Gateway rule from that policy
data "nsxt_policy_intrusion_service_gateway_policy_rule" "ids_gw_rule" {
  display_name = "detect-north-south-threats"
  policy_path  = data.nsxt_policy_parent_intrusion_service_gateway_policy.ids_gw_policy.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the rule to retrieve.
* `display_name` - (Optional) The display name of the rule to retrieve.
* `policy_path` - (Required) The path of the parent gateway policy containing this rule.
* `domain` - (Optional) The domain of the policy containing this rule. Defaults to `default`.

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
* `rule_id` - Unique positive number that is assigned by the system and is useful for debugging.
* `tag` - A list of scope + tag pairs to associate with this rule.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
