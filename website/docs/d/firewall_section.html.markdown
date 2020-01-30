---
layout: "nsxt"
page_title: "NSXT: firewall_section"
sidebar_current: "docs-nsxt-datasource-firewall-section"
description: A firewall section data source.
---

# nsxt_firewall_section

This data source provides information about firewall section in NSX.

## Example Usage

```hcl
data "nsxt_firewall_section" "block_all" {
  display_name = "block all"
}
```

## Argument Reference

* `id` - (Optional) The ID of resource to retrieve

* `display_name` - (Optional) The Display Name of resource to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of resource.
