---
layout: "nsxt"
page_title: "NSXT: nsxt_license"
sidebar_current: "docs-nsxt-resource-license"
description: |-
  Provides a resource to configure license on NSX-T manager
---

# nsxt_license

Provides a resource to configure license on NSX-T manager

## Example Usage

```hcl
resource "nsxt_license" "license" {
  license_key = "0000-00000-00000-00000-00000"
}
```

## Argument Reference

The following arguments are supported:

* `license_key` - (Required) NSX-T license key.
