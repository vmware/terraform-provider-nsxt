---
subcategory: "Manager"
page_title: "NSXT: provider_info"
description: A NSX-T provider information data source.
---

# nsxt_provider_info

This data source provides build and version information about the installed NSX-T Terraform provider.

## Example Usage

```hcl
data "nsxt_provider_info" "info" {}
```

## Attributes Reference

* `commit` - Latest commit hash of the provider build.
* `date` - Date and time when the provider was compiled.
