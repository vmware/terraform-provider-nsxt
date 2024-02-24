---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_edge_high_availability_profile"
description: A resource to configure an Edge High Availability Profile.
---

# nsxt_edge_high_availability_profile

This resource provides a method for the management of an Edge High Availability Profile.
This resource is supported with NSX 4.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_edge_high_availability_profile" "test" {
  description  = "Terraform provisioned Edge High Availability Profile"
  display_name = "test"
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `bfd_allowed_hops` - (Optional) BFD allowed hops.
* `bfd_declare_dead_multiple` - (Optional) Number of times a packet is missed before BFD declares the neighbor down.
* `bfd_probe_interval` - (Optional) the time interval (in millisecond) between probe packets for heartbeat purpose.
* `standby_relocation_threshold` - (Optional) Standby service context relocation wait time.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful
  for debugging.

## Importing

An existing Edge High Availability Profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_edge_high_availability_profile.test UUID
```

The above command imports Edge High Availability Profile named `test` with the NSX Edge High Availability Profile
ID `UUID`.
