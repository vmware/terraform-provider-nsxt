---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_gateway_policy"
sidebar_current: "docs-nsxt-resource-policy-gateway-policy"
description: A resource to Gateway security policies.
---

# nsxt_policy_gateway_policy

This resource provides a method for the management of a Gateway Policy and its Rules.

## Example Usage

```hcl
resource "nsxt_policy_gateway_policy" "test" {
  display_name    = "tf-gw-policy"
  description     = "Terraform provisioned Gateway Policy"
  category        = "LocalGatewayRules"
  locked          = false
  sequence_number = 3
  stateful        = true
  tcp_strict      = false

  tag {
    scope = "color"
    tag   = "orange"
  }

  rule {
    display_name          = "rule1"
    destination_groups    = [nsxt_policy_group.group1.path, nsxt_policy_group.group2.path]
    disabled              = true
    action                = "DROP"
    logged                = true
    scope                 = [nsxt_policy_tier1_gateway.policygateway.path]
  }

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `category` - (Required) The category to use for priority of this Gateway Policy. Must be one of: `Emergency`, `SystemRules`, `SharedPreRules`, `LocalGatewayRules`, `AutoServiceRules` and `Default`.
* `description` - (Optional) Description of the resource.
* `domain` - (Optional) The domain to use for the Gateway Policy. This domain must already exist. For VMware Cloud on AWS use `cgw`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Gateway Policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the Gateway Policy resource.
* `comments` - (Optional) Comments for this Gateway Policy including lock/unlock comments.
* `locked` - (Optional) A boolean value indicating if the policy is locked. If locked, no other users can update the resource.
* `sequence_number` - (Optional) An int value used to resolve conflicts between security policies across domains
* `stateful` - (Optional) A boolean value to indicate if this Policy is stateful. When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed.
* `tcp_strict` - (Optional) A boolean value to enable/disable a 3 way TCP handshake is done before the data packets are sent.
* `rule` (Optional) A repeatable block to specify rules for the Gateway Policy. Each rule includes the following fields:
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

## Importing

An existing Gateway Policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_gateway_policy.gwpolicy1 ID
```

The above command imports the policy Gateway Policy named `gwpolicy1` with the NSX Policy id `ID`.

If the Policy to import isn't in the `default` domain, the domain name can be added to the `ID` before a slash.

For example to import a Group with `ID` in the `MyDomain` domain:

```
terraform import nsxt_policy_gateway_policy.gwpolicy1 MyDomain/ID
```
