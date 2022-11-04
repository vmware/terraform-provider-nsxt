---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: ns_groups"
description: A networking and security groups data source. This data source builds "display name to id" map representation of the whole table.
---

# nsxt_ns_groups

This data source builds a "name to uuid" map of the whole NS Group table. Such map can be referenced in configuration to obtain object uuids by display name at a cost of single roudtrip to NSX, which improves apply and refresh
time at scale, compared to multiple instances of `nsxt_ns_group` data source.

## Example Usage

```hcl
data "nsxt_ns_groups" "map" {
}

resource "nsxt_firewall_section" "s1" {
  display_name = "section1"

  applied_to {
    target_type = "NSGroup"
    target_id   = data.nsxt_ns_groups.map.items["group1"]
  }

  section_type = "LAYER3"
}

```

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `items` - Map of ns group uuids keyed by display name.
