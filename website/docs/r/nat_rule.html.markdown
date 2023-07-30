---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_nat_rule"
description: A resource to configure a NAT rule in NSX.
---

# nsxt_nat_rule

This resource provides a means to configure a NAT rule in NSX. NAT provides network address translation between one IP address space and another IP address space. NAT rules can be destination NAT or source NAT rules.

## Example Usage

```hcl
resource "nsxt_nat_rule" "rule1" {
  logical_router_id         = nsxt_logical_tier1_router.rtr1.id
  description               = "NR provisioned by Terraform"
  display_name              = "NR"
  action                    = "SNAT"
  enabled                   = true
  logging                   = true
  nat_pass                  = false
  translated_network        = "4.4.0.0/24"
  match_destination_network = "3.3.3.0/24"
  match_source_network      = "5.5.5.0/24"

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `logical_router_id` - (Required) ID of the logical router.
* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this NAT rule.
* `action` - (Required) NAT rule action type. Valid actions are: SNAT, DNAT, NO_NAT and REFLEXIVE. All rules in a logical router are either stateless or stateful. Mix is not supported. SNAT and DNAT are stateful, and can NOT be supported when the logical router is running at active-active HA mode. The REFLEXIVE action is stateless. The NO_NAT action has no translated_fields, only match fields.
* `enabled` - (Optional) enable/disable the rule.
* `logging` - (Optional) enable/disable the logging of rule.
* `match_destination_network` - (Required for action=DNAT, not allowed for action=REFLEXIVE) IP Address | CIDR. Omitting this field implies Any.
* `match_source_network` - (Required for action=NO_NAT or REFLEXIVE, Optional for the other actions) IP Address | CIDR. Omitting this field implies Any.
* `nat_pass` - (Optional) Enable/disable to bypass following firewall stage. The default is true, meaning that the following firewall stage will be skipped. Please note, if action is NO_NAT, then nat_pass must be set to true or omitted.
* `translated_network` - (Required for action=DNAT or SNAT) IP Address | IP Range | CIDR.
* `translated_ports` - (Optional) port number or port range. Allowed only when action=DNAT.
* `rule_priority` - (Optional) The priority of the rule which is ascending, valid range [0-2147483647]. If multiple rules have the same priority, evaluation sequence is undefined.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the NAT rule.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing NAT rule can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_nat_rule.rule1 logical-router-uuid/nat-rule-num
```

The above command imports the NAT rule named `rule1` with the number id `nat-rule-num` that belongs to the tier 1 logical router with the NSX id `logical-router-uuid`.
