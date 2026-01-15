---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_attachment"
description: A resource to configure VPC attachment.
---

# nsxt_vpc_attachment

This resource provides a method for the management of VPC Attachment.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

## Example Usage

```hcl
resource "nsxt_vpc_attachment" "test" {
  display_name             = "testVpcAttachment"
  description              = "terraform provisioned vpc attachment"
  parent_path              = nsxt_vpc.test.path
  vpc_connectivity_profile = nsxt_vpc_connectivity_profile.test.path

  depends_on = [
    nsxt_vpc.test,
    nsxt_vpc_connectivity_profile.test,
    nsxt_policy_transit_gateway.test
  ]
}

```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `parent_path` - (Required) Policy path of parent VPC object.
* `vpc_connectivity_profile` - (Required) Path of VPC connectivity profile to attach to the VPC.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_attachment.test PATH
```

The above command imports VPC Attachment named `test` with the policy path `PATH`.
