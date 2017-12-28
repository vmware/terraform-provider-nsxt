---
layout: "nsxt"
page_title: "NSXT: switching_profile"
sidebar_current: "docs-nsxt-datasource-switching-profile"
description: |-
  Provides Switching Profile data source.
---

# nsxt_switching_profile

Provides infromation about switching profiles configured on NSX-T manager.

## Example Usage

```
data "nsxt_switching_profile" "profile1" {
  display_name = "qos-profile"
}
```

## Argument Reference

* `id` - (Optional) The ID of Switching Profile to retrieve

* `display_name` - (Optional) Display Name of the Switching Profile to retrieve

## Attributes Reference

Same as Argument Reference
