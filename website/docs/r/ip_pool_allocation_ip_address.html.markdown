---
layout: "nsxt"
page_title: "NSXT: nsxt_ip_pool_allocation_ip_address"
sidebar_current: "docs-nsxt-resource-ip-pool-allocation-ip-address"
description: |-
  Provides a resource to allocate an IP address from an IP pool on NSX-T manager
---

# nsxt_ip_pool_allocation_ip_address

Provides a resource to allocate an IP address from an IP pool on NSX-T manager

## Example Usage

```hcl
data "nsxt_ip_pool" "ip_pool" {
  display_name = "DefaultIpPool"
}

resource "nsxt_ip_pool_allocation_ip_address" "pool_ip_address" {
  ip_pool_id = "${data.nsxt_ip_pool.ip_pool.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ip_pool_id` - (Required) Ip Pool ID from which the IP address will be allocated.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP pool allocation IP address (currently identical to `allocation_ip`).
* `allocation_ip` - Allocation IP address.
