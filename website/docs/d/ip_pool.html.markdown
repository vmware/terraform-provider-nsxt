---
layout: "nsxt"
page_title: "NSXT: ip_pool"
sidebar_current: "docs-nsxt-datasource-ip-pool"
description: A IP pool data source.
---

# nsxt_ip_pool

This data source provides information about a IP pool configured in NSX.

## Example Usage

```hcl
data "nsxt_ip_pool" "ip_pool" {
  display_name = "DefaultIpPool"
}
```

## Argument Reference

* `id` - (Optional) The ID of IP pool to retrieve

* `display_name` - (Optional) The Display Name of the IP pool to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the IP pool.
