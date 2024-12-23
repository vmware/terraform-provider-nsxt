---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_vpc_ip_address_allocation"
description: A resource to configure IP Address Allocation from IP block associated with VPC.
---

# nsxt_vpc_ip_address_allocation

This resource provides a method for allocating IP Address from IP block associated with VPC.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_vpc_ip_address_allocation" "nat" {
  context {
    project_id = data.nsxt_policy_project.dev.id
    vpc_id     = nsxt_vpc.test.id
  }
  display_name    = "nat"
  allocation_size = 1
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `allocation_size` - (Optional) The system will allocate IP addresses from unused IP addresses based on allocation size.
* `allocation_ips` - (Optional) If specified, IPs have to be within range of respectiveIP blocks.
* `ip_address_block_visibility` - (Optional) Represents visibility of IP address block. This field is not applicable if IPAddressType at VPC is IPV6. Valid values are `EXTERNAL`, `PRIVATE`, `PRIVATE_TGW`, default is `EXTERNAL`.
* `ip_address_type` - (Optional) This defines the type of IP address block that will be used to allocate IP. This field is applicable only if IP address type at VPC is `DUAL`. Valid values are `IPV4` and `IPV6`, default is `IPV4`.
* `ip_block` - (Optional) Policy path for IP Block for the allocation.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ip_address_allocation.test PATH
```

The above command imports VpcIpAddressAllocation named `test` with the policy path `PATH`.
