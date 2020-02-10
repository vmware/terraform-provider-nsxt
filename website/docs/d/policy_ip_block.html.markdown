---
layout: "nsxt"
page_title: "NSXT: policy_ip_block"
sidebar_current: "docs-nsxt-datasource-policy-ip-block"
description: Policy IP Block Config data source.
---

# nsxt_policy_ip_block

This data source provides information about policy IP Blocks configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_ip_block" "test" {
  display_name = "ipblock1"
}
```

## Argument Reference

* `id` - (Optional) The ID of IP Block Config to retrieve.
* `display_name` - (Optional) The Display Name prefix of the IP Block Config to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
