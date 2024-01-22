---
subcategory: "DNS"
layout: "nsxt"
page_title: "NSXT: policy_gateway_dns_forwarder"
description: A policy gateway DNS forwarder data source.
---

# nsxt_policy_gateway_dns_forwarder

This data source provides information about policy gateways DNS forwarder.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_gateway_dns_forwarder" "my_dns_forwarder" {
  display_name = "dns-forwarder1"
  gateway_path = data.nsxt_policy_tier1_gateway.path
}
```

## Argument Reference

* `id` - (Optional) The ID of gateway DNS forwarder to retrieve.
* `display_name` - (Optional) The Display Name of the gateway DNS forwarder to retrieve.
* `gateway_path` - (Optional) Gateway Path for this Service.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
