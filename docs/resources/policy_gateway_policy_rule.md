---
subcategory: "Firewall"
page_title: "NSXT: nsxt_policy_gateway_policy_rule"
description: A resource to add Gateway security policy rules.
---

# nsxt_policy_gateway_policy_rule

This resource provides a method for adding rules for  Gateway Policy.

This resource is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_gateway_policy" "predefined" {
  display_name = "tf-gw-policy"
}

data "nsxt_policy_tier0_gateway" "t0_pepsi" {
  display_name = "pepsi"
}

resource "nsxt_policy_gateway_policy_rule" "rule1" {
  display_name       = "rule2"
  description        = "Terraform provisioned gateway Policy Rule"
  policy_path        = data.nsxt_policy_gateway_policy.predefined.path
  sequence_number    = 1
  destination_groups = [nsxt_policy_group.cats.path, nsxt_policy_group.dogs.path]
  action             = "DROP"
  logged             = true
  scope              = [data.nsxt_policy_tier0_gateway.t0_pepsi.path]
}
```

## Argument Reference

The following arguments are supported:
    
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
* `service_entries` - (Optional) Set of explicit protocol/port service definition
    * `icmp_entry` - (Optional) Set of ICMP type service entries
        * `display_name` - (Optional) Display name of the service entry
        * `protocol` - (Required) Version of ICMP protocol: `ICMPv4` or `ICMPv6`
        * `icmp_code` - (Optional) ICMP message code
        * `icmp_type` - (Optional) ICMP message type
    * `l4_port_set_entry` - (Optional) Set of L4 ports set service entries
        * `display_name` - (Optional) Display name of the service entry
        * `protocol` - (Required) L4 protocol: `TCP` or `UDP`
        * `destination_ports` - (Optional) Set of destination ports
        * `source_ports` - (Optional) Set of source ports
    * `igmp_entry` - (Optional) Set of IGMP type service entries
        * `display_name` - (Optional) Display name of the service entry
    * `ether_type_entry` - (Optional) Set of Ether type service entries
        * `display_name` - (Optional) Display name of the service entry
        * `ether_type` - (Required) Type of the encapsulated protocol
    * `ip_protocol_entry` - (Optional) Set of IP Protocol type service entries
        * `display_name` - (Optional) Display name of the service entry
        * `protocol` - (Required) IP protocol number
    * `algorithm_entry` - (Optional) Set of Algorithm type service entries
        * `display_name` - (Optional) Display name of the service entry
        * `destination_port` - (Required) a single destination port
        * `source_ports` - (Optional) Set of source ports/ranges
        * `algorithm` - (Required) Algorithm: `FTP` or `TFTP`
* `source_groups` - (Optional) Set of group paths that serve as the source for this rule. IPs, IP ranges, or CIDRs may also be used starting in NSX-T 3.0. An empty set can be used to specify "Any".
* `source_excluded` - (Optional) A boolean value indicating negation of source groups.
* `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Rule.
* `action` - (Optional) The action for the Rule. Must be one of: `ALLOW`, `DROP` or `REJECT`. Defaults to `ALLOW`.
* `sequence_number` - (Optional) It is recommended not to specify sequence number for rules, but rather rely on provider to auto-assign them. If you choose to specify sequence numbers, you must make sure the numbers are consistent with order of the rules in configuration. Please note that sequence numbers should start with 1, not 0. To avoid confusion, either specify sequence numbers in all rules, or none at all.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `sequence_number` - Sequence number for the rule.
* `rule_id` - Unique positive number that is assigned by the system and is useful for debugging.

