---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_project_ip_address_allocation"
description: A resource to configure IP Address Allocation under Project.
---

# nsxt_policy_project_ip_address_allocation

This resource provides a method for allocating IP Address from IP block associated with VPC.

This resource is applicable to NSX Policy Manager. IPv6-related arguments (`ip_address_type`, `ipv6_allocation_prefix_length`) are only available in v9.2.0 and above.

## Example Usage

```hcl
resource "nsxt_policy_project_ip_address_allocation" "nat" {
  context {
    project_id = data.nsxt_policy_project.dev.id
  }
  display_name = "nat"
  ip_block     = data.nsxt_policy_project.dev.external_ipv4_blocks[0]
}
```

## Example Usage - IPv6 Allocation (NSX 9.2.0+)

```hcl
resource "nsxt_policy_project_ip_address_allocation" "ipv6" {
  context {
    project_id = data.nsxt_policy_project.dev.id
  }
  display_name                  = "ipv6-alloc"
  ip_address_type               = "IPV6"
  ipv6_allocation_prefix_length = 64
  ip_block                      = data.nsxt_policy_project.dev.ipv6_blocks[0]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `allocation_size` - (Optional) The system will allocate IP addresses from unused IP addresses based on allocation size. Currently only size `1` is supported. Conflicts with `ipv6_allocation_prefix_length`.
* `allocation_ips` - (Optional) If specified, IPs have to be within range of respective IP blocks.
* `ip_address_type` - (Optional) Type of IP address to allocate. Allowed values are `IPV4` and `IPV6`. Defaults to `IPV4`. Immutable after creation (forces new resource). This attribute is only available in v9.2.0 and above.
* `ipv6_allocation_prefix_length` - (Optional) Prefix length of the allocated IPv6 subnet. Allowed values are between `64` and `128`. Defaults to `64`. Conflicts with `allocation_size`. Immutable after creation (forces new resource). This attribute is only available in v9.2.0 and above.
* `ip_block` - (Optional) Policy path for IP Block for the allocation.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_ip_address_allocation.test PATH
```

The above command imports IP Address Allocation named `test` with the policy path `PATH`.
