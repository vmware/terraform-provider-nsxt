---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: failure_domain"
description: Failure Domain data source.
---

# nsxt_failure_domain

This data source provides information about Failure Domain configured in NSX.

## Example Usage

```hcl
data "nsxt_failure_domain" "failure_domain" {
  display_name = "failuredomain1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Failure Domain to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Failure Domain to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
