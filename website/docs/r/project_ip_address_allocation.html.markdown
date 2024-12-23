---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_project_ip_address_allocation"
description: A resource to configure IP Address Allocation under Project.
---

# nsxt_policy_project_ip_address_allocation

This resource provides a method for allocating IP Address from IP block associated with VPC.

This resource is applicable to NSX Policy Manager.

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

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `allocation_size` - (Optional) The system will allocate IP addresses from unused IP addresses based on allocation size. Currently only size `1` is supported.
* `allocation_ips` - (Optional) If specified, IPs have to be within range of respective IP blocks.
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

The above command imports IP Address Allocation named `test` with the policy path `PATH`.
