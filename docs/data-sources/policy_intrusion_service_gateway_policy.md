---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_gateway_policy"
description: A data source to retrieve an Intrusion Service Gateway Policy.
---

# nsxt_policy_intrusion_service_gateway_policy

This data source provides information about an existing Intrusion Service Gateway Policy configured on NSX.
This data source can be useful for fetching policy path to use in `nsxt_policy_intrusion_service_gateway_policy_rule` resource.

~> **NOTE:** This data source retrieves only the policy metadata (id, display_name, description, path, category). It does not retrieve the rules within the policy. To manage rules, use the `nsxt_policy_intrusion_service_gateway_policy_rule` resource with the policy path obtained from this data source.

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_intrusion_service_gateway_policy" "idps_gateway_policy" {
  display_name = "intrusion-service-gateway-policy"
}
```

## Argument Reference

* `id` - (Optional) The ID of the policy to retrieve.
* `display_name` - (Optional) The display name of the policy to retrieve.
* `domain` - (Optional) The domain of the policy. Defaults to `default`.
* `category` - (Optional) Category of the policy.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
