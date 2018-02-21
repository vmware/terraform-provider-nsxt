---
layout: "nsxt"
page_title: "NSXT: switching_profile"
sidebar_current: "docs-nsxt-datasource-switching-profile"
description: |-
  Provides Switching Profile data source.
---

# nsxt_switching_profile

Provides information about switching profiles configured on NSX-T manager.

## Example Usage

```
data "nsxt_switching_profile" "qos_profile" {
  display_name = "qos-profile"
}
```

## Argument Reference

* `id` - (Optional) The ID of Switching Profile to retrieve

* `display_name` - (Optional) Display Name of the Switching Profile to retrieve

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `resource_type` - The resource type representing the specific type of this profile.

* `system_owned` - A boolean that indicates whether this resource is system-owned and thus read-only.

* `description` - Description of the edge cluster.
