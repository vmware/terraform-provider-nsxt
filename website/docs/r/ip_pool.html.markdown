---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_ip_pool"
description: |-
  Provides a resource to configure IP pool on NSX-T manager
---

# nsxt_ip_pool

Provides a resource to configure IP pool on NSX-T manager

## Example Usage

```hcl
resource "nsxt_ip_pool" "ip_pool" {
  description  = "ip_pool provisioned by Terraform"
  display_name = "ip_pool"

  tag {
    scope = "color"
    tag   = "red"
  }

  subnet {
    allocation_ranges = ["2.1.1.1-2.1.1.11", "2.1.1.21-2.1.1.100"]
    cidr              = "2.1.1.0/24"
    gateway_ip        = "2.1.1.12"
    dns_suffix        = "abc"
    dns_nameservers   = ["33.33.33.33"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this IP pool.
* `subnet` - (Optional) Subnets can be IPv4 or IPv6 and they should not overlap. The maximum number will not exceed 5 subnets. Each subnet has the following arguments:
  * `allocation_ranges` - (Required) A collection of IPv4 Pool Ranges
  * `cidr` - (Required) Network address and the prefix length which will be associated with a layer-2 broadcast domainIPv4 Pool Ranges
  * `dns_nameservers` - (Optional) A collection of up to 3 DNS servers for the subnet
  * `dns_suffix` - (Optional) The DNS suffix for the DNS server
  * `gateway_ip` - (Optional) The default gateway address on a layer-3 router

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP pool.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing IP pool can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_ip_pool.ip_pool UUID
```

The above would import the IP pool named `ip_pool` with the nsx id `UUID`
