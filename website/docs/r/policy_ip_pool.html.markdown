---
subcategory: "IPAM"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ip_pool"
description: A resource to configure IP Pools in NSX Policy.
---

# nsxt_policy_ip_pool

This resource provides a means to configure IP Pools in NSX Policy.

This resource is applicable to NSX Policy Manager.

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

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_ip_pool" "pool1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
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
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP Pool.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the resource.
* `realized_id` - The id of realized pool object. This id should be used in `nsxt_transport_node` resource.

## Importing

An existing IP Pool can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ip_pool.pool1 ID
```
The above would import NSX IP Pool as a resource named `pool1` with the NSX ID `ID`, where `ID` is NSX ID of the IP Pool.

```
terraform import nsxt_policy_ip_pool.pool1 POLICY_PATH
```
The above would import NSX IP Pool as a resource named `pool1` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
