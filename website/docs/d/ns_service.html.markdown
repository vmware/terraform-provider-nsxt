---
layout: "nsxt"
page_title: "NSXT: ns_service"
sidebar_current: "docs-nsxt-datasource-ns-service"
description: |-
  Provides NS service data source.
---

# nsxt_ns_service

Provides information about NS service configured on NSX-T manager.

## Example Usage

```
data "nsxt_ns_service" "ns_service_dns" {
  display_name = "DNS"
}
```

## Argument Reference

* `id` - (Optional) The ID of NS service to retrieve

* `display_name` - (Optional) Display Name of the NS service to retrieve

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - Description of the NS service.
