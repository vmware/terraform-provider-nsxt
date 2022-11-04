---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_ip_pool_allocation_ip_address"
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
  ip_pool_id = data.nsxt_ip_pool.ip_pool.id
}
```

## Argument Reference

The following arguments are supported:

* `ip_pool_id` - (Required) Ip Pool ID from which the IP address will be allocated.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP pool allocation IP address (currently identical to `allocation_ip`).
* `allocation_ip` - Allocation IP address.

## Importing

An existing IP pool allocation address can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_ip_pool_allocation_ip_address.ip1 POOL-UUID/UUID
```

The above would import the IP pool allocation address named `ip_pool` with the nsx ID `UUID`, from IP Pool with nsx ID `POOL-UUID`.
