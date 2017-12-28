---
layout: "nsxt"
page_title: "NSXT: transport_zone"
sidebar_current: "docs-nsxt-datasource-transport-zone"
description: |-
  Provides Transport Zone data source.
---

# nsxt_transport_zone

Provides infromation about transport zones (TZ) configured on NSX-T manager.

## Example Usage

```
data "nsxt_transport_zone" "TZ1" {
  display_name = "1-transportzone-87"
}
```

## Argument Reference

* `id` - (Optional) The ID of Transport Zone to retrieve

* `display_name` - (Optional) Display Name of the Transport Zone to retrieve

## Attributes Reference

Same as Argument Reference
