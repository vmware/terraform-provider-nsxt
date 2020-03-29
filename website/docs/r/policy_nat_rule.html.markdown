---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_nat_rule"
sidebar_current: "docs-nsxt-resource-policy-nat-rule"
description: A resource to configure NAT Ruels in NSX Policy manager.
---

# nsxt_policy_nat_rule

This resource provides a method for the management of a NAT Rule.

## Example Usage

```hcl
resource "nsxt_policy_nat_rule" "dnat1" {
  display_name         = "dnat_rule1"
  action               = "DNAT"
  source_networks      = ["9.1.1.1", "9.2.1.1"]
  destination_networks = ["11.1.1.1"]
  translated_networks  = ["10.1.1.1"]
  gateway_path         = nsxt_policy_tier1_gateway.t1gateway.path
  logging              = false
  firewall_match       = "MATCH_INTERNAL_ADDRESS"

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this NAT Rule.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `gateway_path` - (Required) The NSX Policy path to the Tier0 or Tier1 Gateway for this NAT Rule.
* `action` - (Required) The action for the NAT Rule. One of `SNAT`, `DNAT`, `REFLEXIVE`, `NO_SNAT`, `NO_DNAT`, `NAT64`.
* `destination_networks` - (Optional) A list of destination network IP addresses or CIDR.
* `enabled` - (Optional) Enable/disable the Rule. Defaults to `true`.
* `firewall_match` - (Optional) Firewall match flag. One of `MATCH_EXTERNAL_ADDRESS`, `MATCH_INTERNAL_ADDRESS`, `BYPASS`.
* `logging` - (Optional) Enable/disable rule logging. Defaults to `false`.
* `rule_priority` - (Optional) The priority of the rule. Valid values between 0 to 2147483647. Defaults to `100`.
* `service` - (Optional) Policy path of Service on which the NAT rule will be applied.
* `source_networks` - (Optional) A list of source network IP addresses or CIDR.
* `translated_networks` - (Optional) A list of translated network IP addresses or CIDR.
* `translated_ports` - (Optional) Port number or port range. For use with `DNAT` action only.
* `scope` - (Optional) A list of paths to interfaces and/or labels where the NAT Rule is enforced.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing policy NAT Rule can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_nat_rule.rule1 GWID/ID
```

The above command imports the policy NAT Rule named `rule1` for the NSX Tier0 or Tier1 Gateway `GWID` with the NSX Policy ID `ID`.
