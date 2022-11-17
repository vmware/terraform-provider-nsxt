---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_firewall_section"
description: A resource that can be used to configure a firewall section on NSX.
---

# nsxt_firewall_section

This resource provides a way to configure a firewall section on the NSX manager. A firewall section is a collection of firewall rules that are grouped together.
Order of firewall sections can be controlled with 'insert_before' attribute.

## Example Usage

```hcl
data "nsxt_firewall_section" "block_all" {
  display_name = "BLOCK"
}

resource "nsxt_firewall_section" "firewall_sect" {
  description  = "FS provisioned by Terraform"
  display_name = "FS"

  tag {
    scope = "color"
    tag   = "blue"
  }

  applied_to {
    target_type = "NSGroup"
    target_id   = nsxt_ns_group.group1.id
  }

  section_type  = "LAYER3"
  stateful      = true
  insert_before = data.nsxt_firewall_section.block_all.id

  rule {
    display_name          = "out_rule"
    description           = "Out going rule"
    action                = "ALLOW"
    logged                = true
    ip_protocol           = "IPV4"
    direction             = "OUT"
    destinations_excluded = "false"
    sources_excluded      = "true"

    source {
      target_type = "LogicalSwitch"
      target_id   = nsxt_logical_switch.switch1.id
    }

    destination {
      target_type = "LogicalSwitch"
      target_id   = nsxt_logical_switch.switch2.id
    }
  }

  rule {
    display_name = "in_rule"
    description  = "In going rule"
    action       = "DROP"
    logged       = true
    ip_protocol  = "IPV4"
    direction    = "IN"

    service {
      target_type = "NSService"
      target_id   = "e8d59e13-484b-4825-ae3b-4c11f83249d9"
    }

    service {
      target_type = "NSService"
      target_id   = nsxt_l4_port_set_ns_service.http.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) The display name of this firewall section. Defaults to ID if not set.
* `description` - (Optional) Description of this firewall section.
* `tag` - (Optional) A list of scope + tag pairs to associate with this firewall section.
* `applied_to` - (Optional) List of objects where the rules in this section will be enforced. This will take precedence over rule level applied_to. [Supported target types: "LogicalPort", "LogicalSwitch", "NSGroup", "LogicalRouter"]
* `section_type` - (Required) Type of the rules which a section can contain. Either LAYER2 or LAYER3. Only homogeneous sections are supported.
* `stateful` - (Required) Stateful or Stateless nature of firewall section is enforced on all rules inside the section. Layer3 sections can be stateful or stateless. Layer2 sections can only be stateless.
* `insert_before` - (Optional) Firewall section id that should come immediately after this one. It is user responsibility to use this attribute in consistent manner (for example, if same value would be set in two separate sections, the outcome would depend on order of creation). Changing this attribute would force recreation of the firewall section.
* `rule` - (Optional) A list of rules to be applied in this section. each rule has the following arguments:
  * `display_name` - (Optional) The display name of this rule. Defaults to ID if not set.
  * `description` - (Optional) Description of this rule.
  * `action` - (Required) Action enforced on the packets which matches the firewall rule. [Allowed values: "ALLOW", "DROP", "REJECT"]
  * `applied_to` - (Optional) List of objects where rule will be enforced. The section level field overrides this one. Null will be treated as any. [Supported target types: "LogicalPort", "LogicalSwitch", "NSGroup", "LogicalRouterPort"]
  * `destination` - (Optional) List of the destinations. Null will be treated as any. [Allowed target types: "IPSet", "LogicalPort", "LogicalSwitch", "NSGroup", "MACSet" (depending on the section type)]
  * `destinations_excluded` - (Optional) When this boolean flag is set to true, the rule destinations will be negated.
  * `direction` - (Optional) Rule direction in case of stateless firewall rules. This will only considered if section level parameter is set to stateless. Default to IN_OUT if not specified. [Allowed values: "IN", "OUT", "IN_OUT"]
  * `disabled` - (Optional) Flag to disable rule. Disabled will only be persisted but never provisioned/realized.
  * `ip_protocol` - (Optional) Type of IP packet that should be matched while enforcing the rule. [allowed values: "IPV4", "IPV6", "IPV4_IPV6"]
  * `logged` - (Optional) Flag to enable packet logging. Default is disabled.
  * `notes` - (Optional) User notes specific to the rule.
  * `rule_tag` - (Optional) User level field which will be printed in CLI and packet logs.
  * `service` - (Optional) List of the services. Null will be treated as any. [Allowed target types: "NSService", "NSServiceGroup"]
  * `source` - (Optional) List of sources. Null will be treated as any. [Allowed target types: "IPSet", "LogicalPort", "LogicalSwitch", "NSGroup", "MACSet" (depending on the section type)]
  * `sources_excluded` - (Optional) When this boolean flag is set to true, the rule sources will be negated.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the firewall section.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `is_default` - A boolean flag which reflects whether a firewall section is default section or not. Each Layer 3 and Layer 2 section will have at least and at most one default section.

## Importing

An existing Firewall section can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_firewall_section.firewall_sect UUID
```

The above command imports the firewall section named `firewall_sect` with the NSX id `UUID`.
