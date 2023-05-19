---
subcategory: "Grouping and Tagging"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_domain"
description: A resource to configure a Global manager domain.
---

# nsxt_policy_domain

This resource provides a method for the management of a global manager domain (a.k.a Region) and its locations.

This resource is applicable to NSX Global Manager.

## Example Usage

```hcl
resource "nsxt_policy_domain" "domain1" {
  display_name = "Domain 1"
  sites        = ["here", "there"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `sites` - (Required) A list of sites (a.k.a locations) for this domain.

* `description` - (Optional) Description of the resource.

* `tag` - (Optional) A list of scope + tag pairs to associate with this Domain.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the Domain resource.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Domain.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing policy Domain can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_domain.domain1 ID
```

The above command imports the policy Domain named `domain` with the NSX Policy ID `ID`.
