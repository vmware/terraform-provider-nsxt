---
layout: "nsxt"
page_title: "NSXT: nsxt_nat_rule"
sidebar_current: "docs-nsxt-resource-nat-rule"
description: |-
  Provides a resource to configure nat rule on NSX-T manager
---

# nsxt_nat_rule

Provides a resource to configure nat rule on NSX-T manager

## Example Usage

```hcl
resource "nsxt_nat_rule" "NR" {
    logical_router_id = "${nsxt_logical_tier1_router.RTR1.id}"
    description = "NR provisioned by Terraform"
    display_name = "NR"
    action = "SNAT"
    enabled = true
    logging = true
    nat_pass = false
    translated_network = "4.4.0.0/24"
    match_destination_network = "3.3.3.0/24"
    match_source_network = "5.5.5.0/24"}
    tag {
        scope = "color"
        tag = "blue"
    }
```

## Argument Reference

The following arguments are supported:

* `logical_router_id` - (Required) ID of the logical router.
* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this nat_rule.
* `action` - (Required) valid actions: SNAT, DNAT, NO_NAT, REFLEXIVE. All rules in a logical router are either stateless or stateful. Mix is not supported. SNAT and DNAT are stateful, can NOT be supported when the logical router is running at active-active HA mode; REFLEXIVE is stateless. NO_NAT has no translated_fields, only match fields.
* `enabled` - (Optional) enable/disable the rule.
* `logging` - (Optional) enable/disable the logging of rule.
* `match_destination_network` - (Optional) IP Address | CIDR | (null implies Any).
* `match_source_network` - (Optional) IP Address | CIDR | (null implies Any).
* `nat_pass` - (Optional) Default is true. If the nat_pass is set to true, the following firewall stage will be skipped. Please note, if action is NO_NAT, then nat_pass must be set to true or omitted.
* `translated_network` - (Optional) IP Address | IP Range | CIDR.
* `translated_ports` - (Optional) port number or port range. DNAT only.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the nat_rule.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `system_owned` - A boolean that indicates whether this resource is system-owned and thus read-only.
* `rule_priority` - Ascending, valid range [0-2147483647]. If multiple rules have the same priority, evaluation sequence is undefined.
