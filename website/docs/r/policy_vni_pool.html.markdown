---
subcategory: "EVPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_vni_pool"
description: A resource to configure VNI Pool on NSX Policy manager.
---

# nsxt_policy_vni_pool

This resource provides a method for the management of VNI Pools.

This resource is applicable to NSX Policy Manager.

This resource is supported with NSX 3.0.0 onwards.

## Example Usage

```hcl
resource "nsxt_policy_vni_pool" "test_vni_pool" {
  description  = "vnipool1"
  display_name = "terraform provisioned VNI pool"

  start = 80000
  end   = 90000

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `start` - (Required) Start value of VNI Pool range. Minimum: 75001, Maximum: 16777215.
* `end` - (Required) End value of VNI Pool range. Minimum: 75001, Maximum: 16777215.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing VNI Pool can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_vni_pool.vnipool1 ID
```

The above command imports the VNI Pool named `vnipool1` with the NSX Policy ID `ID`.

```
terraform import nsxt_policy_vni_pool.vnipool1 POLICY_PATH
```
The above command imports the VNI pool named `vnipool1` with the policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
