---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ip_block"
sidebar_current: "docs-nsxt-resource-policy-ip-block"
description: A resource to configure IP Blocks in NSX Policy.
---

# nsxt_policy_ip_block

  This resource provides a means to configure IP Blocks in NSX Policy.

## Example Usage

```hcl
resource "nsxt_policy_ip_block" "block1" {
  display_name = "ip-block1"
  cidr         = "192.168.1.0/24"

  tag {
    scope = "color"
    tag   = "blue"
  }

  tag {
    scope = "env"
    tag   = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The display name for the IP Block.
* `description` - (Optional) Description of the resource.
* `cidr` - (Required) Network address and the prefix length which will be associated with a layer-2 broadcast domain.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this IP Block.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP Block.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the resource.

## Importing

An existing IP Block can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_ip_block.block1 ID
```

The above would import NSX IP Block as a resource named `block1` with the NSX id `ID`, where `ID` is NSX ID of the IP Block.
