---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_vpc_nat_rule"
description: A resource to configure a NAT rule under VPC.
---

# nsxt_vpc_nat_rule

This resource provides a method for the management of a VPC NAT Rule.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_vpc_ip_address_allocation" "nat" {
  context {
    project_id = data.nsxt_policy_project.dev.id
    vpc_id     = nsxt_vpc.vpc1.id
  }
  display_name    = "for-nat"
  allocation_size = 1
}

data "nsxt_vpc_nat" "test" {
  context {
    project_id = data.nsxt_policy_project.dev.id
    vpc_id     = nsxt_vpc.vpc1.id
  }
  nat_type = "USER"
}

resource "nsxt_vpc_nat_rule" "test" {
  display_name        = "test"
  description         = "terraform provisioned nat rule for vpc"
  parent_path         = data.nsxt_vpc_nat.test.path
  destination_network = nsxt_vpc_ip_address_allocation.nat.allocation_ips
  action              = "DNAT"
  source_network      = "10.205.1.13"
  translated_network  = "2.2.2.13"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `parent_path` - (Required) Policy path of parent NAT object, typically reference to `path` in `nsxt_vpc_nat` data source.
* `translated_network` - (Optional) For `SNAT`, `DNAT` and `REFLEXIVE` rules, this ia a required field, which
represents the translated network address. In case of `SNAT` and `REFLEXIVE` rule, translated network address should be single IPv4 address allocated from External Block associated with VPC.
* `logging` - (Optional) Boolean flag to indicate whether logging is enabled. The default is `false`.
* `destination_network` - (Optional) For `DNAT` rules, this is a required field, and represents the destination network for the incoming packets. For other type of rules, it may contain destination network of outgoing packets. In case of `DNAT` rule, destination network address should be IPv4 address allocated from External Block associated with VPC.
* `action` - (Required) NAT action, one of `SNAT` (translates a source IP address into an outbound packet so that
the packet appears to originate from a different network), `DNAT` (translates the destination IP address of inbound packets so that packets are delivered to a target address into another network), and `REFLEXIVE` (one-to-one mapping of source and destination IP addresses).
* `firewall_match` - (Optional) Indicates how the firewall matches the address after NATing if firewall
stage is not skipped, one of `MATCH_EXTERNAL_ADDRESS`, `MATCH_INTERNAL_ADDRESS` or `BYPASS`. Default is `MATCH_INTERNAL_ADDRESS`.
* `source_network` - (Optional) Source network. For `SNAT` and `REFLEXIVE` rules, this is a required field. For `DNAT` rules, it may contain source network for incoming packets.
* `enabled` - (Optional) Flag for enabling the NAT rule, default is `true`.
* `sequence_number` - (Optional) The sequence_number decides the rule_priority of a NAT rule.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_vpc_nat_rule.test PATH
```

The above command imports PolicyVpcNatRule named `test` with the policy path `PATH`.
