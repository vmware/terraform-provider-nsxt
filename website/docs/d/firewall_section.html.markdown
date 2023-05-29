---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: firewall_section"
description: A firewall section data source.
---

# nsxt_firewall_section

This data source provides information about firewall section configured on NSX. It can be useful to enforce placement of newly created firewall sections.

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
