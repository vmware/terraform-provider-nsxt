---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ip_address_allocation"
sidebar_current: "docs-nsxt-resource-policy-ip-address-allocation"
description: A resource to configure a IP Address allocations.
---

# nsxt_policy_ip_address_allocation

This resource provides a method for the management of a IP Address Allocations.
Note that IP Address Allocations cannot be updated once created and changing any attributes of
an existing allocation will re-create it.
 
## Example Usage

```hcl
resource "nsxt_policy_ip_address_allocation" "test" {
  display_name  = "test"
  description   = "Terraform provisioned IpAddressAllocation"
  pool_path     = nsxt_policy_ip_pool.pool1.path
  allocation_ip = "12.12.12.12"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `allocation_ip` - (Optional) The IP Address to allocate. If unspecified any free IP in the pool will be allocated.
* `pool_path` - (Required) The policy path to the IP Pool for this Allocation.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Allocation.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `allocation_ip` - If the `allocation_ip` is not specified in the resource, any free IP is allocated and its value is exported on this attribute.

## Importing

An existing IP Allocation can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_ip_address_allocation.test POOL-ID/ID
```

The above command imports IpAddressAllocation named `test` with the NSX IpAddressAllocation ID `ID` in IP Pool `POOL-ID`.
