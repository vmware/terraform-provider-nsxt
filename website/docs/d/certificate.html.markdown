---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: certificate"
description: A certificate data source.
---

# nsxt_certificate

This data source provides information about various types of certificates imported into NSX trust management.

This data source is applicable to NSX-T Policy Manager.

## Example Usage

```hcl
data "nsxt_certificate" "CA" {
  display_name = "ca-cert"
}
```

## Argument Reference

* `id` - (Optional) The ID of Certificate to retrieve.
* `display_name` - (Optional) The Display Name of the Certificate to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the Certificate.
