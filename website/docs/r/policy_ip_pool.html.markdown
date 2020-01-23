---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ip_pool"
sidebar_current: "docs-nsxt-resource-policy-ip-pool"
description: A resource to configure IP Pools in NSX Policy.
---

# nsxt_policy_ip_pool

  This resource provides a means to configure IP Pools in NSX Policy.

## Example Usage

```hcl
resource "nsxt_policy_ip_pool" "pool1" {
  display_name = "ip-pool1"

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

* `display_name` - (Required) The display name for the IP Pool.
* `description` - (Optional) Description of the resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this IP Pool.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP Pool.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the resource.

## Importing

An existing IP Pool can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_ip_pool.pool1 ID
```

The above would import NSX IP Pool as a resource named `pool1` with the NSX ID `ID`, where `ID` is NSX ID of the IP Pool.
