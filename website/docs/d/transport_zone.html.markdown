---
layout: "nsxt"
page_title: "NSXT: transport_zone"
sidebar_current: "docs-nsxt-datasource-transport-zone"
description: |-
  Provides Transport Zone data source.
---

# nsxt_transport_zone

Provides information about transport zones (TZ) configured on NSX-T manager.

## Example Usage

```
data "nsxt_transport_zone" "TZ1" {
  display_name = "1-transportzone-87"
}
```

## Argument Reference

* `id` - (Optional) The ID of Transport Zone to retrieve

* `display_name` - (Optional) Display Name prefix of the Transport Zone to retrieve

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - Description of the transport zone.

* `host_switch_name` - TName of the host switch on all transport nodes in this transport zone that will be used to run NSX network traffic.

* `transport_type` - The transport type of this transport zone (OVERLAY or VLAN).
