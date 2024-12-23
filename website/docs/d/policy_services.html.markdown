---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: policy_services"
description: A policy services data source. This data source builds a list representation of the whole services table, containing `id`, `display_name` and `path` attributes.
---

# nsxt_policy_services

This data source builds a list of the whole policy Services table. Such list can be referenced in configuration to obtain object identifier attributes lookup by display name at a cost of single roundtrip to NSX, which improves apply and refresh
time at scale, compared to multiple instances of `nsxt_policy_service` data source.

## Example Usage

```hcl
data "nsxt_policy_services" "servicelist" {
}

resource "nsxt_policy_nat_rule" "dnat1" {
  display_name          = "dnat_rule1"
  action                = "DNAT"
  source_networks       = ["9.1.1.1", "9.2.1.1"]
  destination_networks  = ["11.1.1.1"]
  translated_networks   = ["10.1.1.1"]
  gateway_path          = nsxt_policy_tier1_gateway.t1gateway.path
  logging               = false
  firewall_match        = "MATCH_INTERNAL_ADDRESS"
  policy_based_vpn_mode = "BYPASS"
  service               = data.nsxt_policy_services.servicelist.items[index(data.nsxt_policy_services.servicelist.service_from_list.*.display_name, "DNS-UDP")].path

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `items` - List of policy service entries.
  * `id` - The ID of the service.
  * `display_name` - Display name of the service entry.
  * `path` - The NSX path of the policy resource.
