---
layout: "nsxt"
page_title: "NSXT: policy_service"
sidebar_current: "docs-nsxt-datasource-policy-service"
description: A policy service data source.
---

# nsxt_policy_service

This data source provides information about policy services configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_service" "dns_service" {
  display_name = "DNS"
}
```

## Argument Reference

* `id` - (Optional) The ID of service to retrieve.

* `display_name` - (Optional) The Display Name prefix of the service to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
