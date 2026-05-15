---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_gateway_policy"
description: A data source to retrieve an Intrusion Service Gateway Policy.
---

# nsxt_policy_intrusion_service_gateway_policy

This data source provides information about an existing Intrusion Service Gateway Policy configured on NSX.
It can be useful for fetching the IDPS Gateway Policy including all its embedded rules.

~> **NOTE:** This data source retrieves the policy including embedded rules, allowing you to refer IDPS Gateway policy details and its rules in other resources. For different use cases, consider:

* `nsxt_policy_parent_intrusion_service_gateway_policy` - For **IDPS Gateway policy metadata only (no rules)**
* `nsxt_policy_intrusion_service_gateway_policy_rule` - For **individual standalone IDPS Gateway rules**

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_intrusion_service_gateway_policy" "idps_gateway_policy" {
  display_name = "intrusion-service-gateway-policy"
}
```

## Argument Reference

* `id` - (Optional) The ID of the policy to retrieve.
* `display_name` - (Optional) The display name of the policy to retrieve.
* `domain` - (Optional) The domain of the policy. Defaults to `default`.
* `category` - (Optional) Category of the policy.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `stateful` - When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed. Note: Intrusion Service Gateway Policies are always stateful.
* `sequence_number` - This field is used to resolve conflicts between multiple policies that have rules that match the same packet.
* `locked` - Indicates whether a security policy should be locked.
* `comments` - Comments for security policy lock/unlock.
* `tag` - A list of scope + tag pairs to associate with this policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `rule` - A list of rules in this policy. Each rule contains:
    * `display_name` - The display name of the rule.
    * `description` - The description of the rule.
    * `notes` - Additional notes for the rule.
    * `action` - Action for this rule (DETECT, DETECT_PREVENT).
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
    * `sequence_number` - The sequence number of the rule.
    * `rule_id` - Unique positive number assigned by the system.
    * `revision` - Indicates current revision number of the rule.
    * `tag` - A list of scope + tag pairs associated with this rule.
