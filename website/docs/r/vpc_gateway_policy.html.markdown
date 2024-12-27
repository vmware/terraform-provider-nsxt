---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_vpc_gateway_policy"
description: A resource for VPC Gateway security policies.
---

# nsxt_vpc_gateway_policy

This resource provides a method for the management of a VPC Gateway Policy and its Rules.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc" "demovpc" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "vpc1"
}

resource "nsxt_vpc_gateway_policy" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  display_name    = "tf-gw-policy"
  description     = "Terraform provisioned Gateway Policy"
  locked          = false
  sequence_number = 3
  stateful        = true
  tcp_strict      = false

  tag {
    scope = "color"
    tag   = "orange"
  }

  rule {
    display_name       = "rule1"
    destination_groups = [nsxt_vpc_group.group1.path, nsxt_vpc_group.group2.path]
    disabled           = true
    action             = "DROP"
    logged             = true
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

-> We recommend using `lifecycle` directive as in samples above, in order to avoid dependency issues when updating groups/services simultaneously with the rule.

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Gateway Policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the Gateway Policy resource.
* `context` - (Required) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
  * `vpc_id` - (Required) The ID of the VPC which the object belongs to
* `comments` - (Optional) Comments for this Gateway Policy including lock/unlock comments.
* `locked` - (Optional) A boolean value indicating if the policy is locked. If locked, no other users can update the resource.
* `sequence_number` - (Optional) An int value used to resolve conflicts between security policies
* `stateful` - (Optional) A boolean value to indicate if this Policy is stateful. When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed.
* `tcp_strict` - (Optional) A boolean value to enable/disable a 3 way TCP handshake is done before the data packets are sent.
* `rule` (Optional) A repeatable block to specify rules for the Gateway Policy. Each rule includes the following fields:
  * `display_name` - (Required) Display name of the resource.
  * `description` - (Optional) Description of the resource.
  * `destination_groups` - (Optional) Set of group paths that serve as the destination for this rule. IPs, IP ranges, or CIDRs. An empty set can be used to specify "Any".
  * `destinations_excluded` - (Optional) A boolean value indicating negation of destination groups.
  * `direction` - (Optional) The traffic direction for the policy. Must be one of: `IN`, `OUT` or `IN_OUT`. Defaults to `IN_OUT`.
  * `disabled` - (Optional) A boolean value to indicate the rule is disabled. Defaults to `false`.
  * `ip_version` - (Optional) The IP Protocol for the rule. Must be one of: `IPV4`, `IPV6` or `IPV4_IPV6`. Defaults to `IPV4_IPV6`.
  * `logged` - (Optional) A boolean flag to enable packet logging.
  * `notes` - (Optional) Text for additional notes on changes for the rule.
  * `profiles` - (Optional) A list of context profiles for the rule.
  * `services` - (Optional) List of services to match.
  * `source_groups` - (Optional) Set of group paths that serve as the source for this rule. IPs, IP ranges, or CIDRs. An empty set can be used to specify "Any".
  * `source_excluded` - (Optional) A boolean value indicating negation of source groups.
  * `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this Rule.
  * `action` - (Optional) The action for the Rule. Must be one of: `ALLOW`, `DROP` or `REJECT`. Defaults to `ALLOW`.
  * `sequence_number` - (Optional) It is recommended not to specify sequence number for rules, but rather rely on provider to auto-assign them. If you choose to specify sequence numbers, you must make sure the numbers are consistent with order of the rules in configuration. Please note that sequence numbers should start with 1, not 0. To avoid confusion, either specify sequence numbers in all rules, or none at all.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Security Policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `rule`:
  * `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
  * `path` - The NSX path of the policy resource.
  * `sequence_number` - Sequence number for the rule.
  * `rule_id` - Unique positive number that is assigned by the system and is useful for debugging.

## Importing

An existing VPC Gateway Policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_vpc_gateway_policy.policy1 PATH
```

The above command imports the VPC gateway policy named `policy1` with the NSX Policy path `PATH`.
