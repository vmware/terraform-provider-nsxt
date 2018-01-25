---
layout: "nsxt"
page_title: "NSXT: nsxt_firewall_section"
sidebar_current: "docs-nsxt-resource-firewall-section"
description: |-
  Provides a resource to configure firewall section on NSX-T manager
---

# nsxt_firewall_section

Provides a resource to configure firewall section on NSX-T manager

## Example Usage

```hcl
resource "nsxt_firewall_section" "FS" {
  description = "FS provisioned by Terraform"
  display_name = "FS"
  tags = [{ scope = "color"
            tag = "red" }
  ]
  applied_tos = ...
  rules = ...
  section_type = "LAYER3"
  stateful = true
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) Defaults to ID if not set.
* `tags` - (Optional) A list of scope + tag pairs to associate with this firewall_section.
* `applied_tos` - (Optional) List of objects where the rules in this section will be enforced. This will take precedence over rule level appliedTo.
* `section_type` - (Required) Type of the rules which a section can contain. Either LAYER2 or LAYER3. Only homogeneous sections are supported.
* `stateful` - (Required) Stateful or Stateless nature of firewall section is enforced on all rules inside the section. Layer3 sections can be stateful or stateless. Layer2 sections can only be stateless.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the firewall_section.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `system_owned` - A boolean that indicates whether this resource is system-owned and thus read-only.
* `is_default` - It is a boolean flag which reflects whether a firewall section is default section or not. Each Layer 3 and Layer 2 section will have at least and at most one default section.
* `rule_count` - (Optional) Number of rules in this section.
