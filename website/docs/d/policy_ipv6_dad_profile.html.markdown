---
layout: "nsxt"
page_title: "NSXT: policy_ipv6_dad_profile"
sidebar_current: "docs-nsxt-datasource-policy-ipv6-dad-profile"
description: Policy Ipv6DadProfile data source.
---

# nsxt_policy_ipv6_dad_profile

This data source provides information about policy Ipv6DadProfile configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_ipv6_dad_profile" "test" {
  display_name = "ipv6-dad-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Ipv6DadProfile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Ipv6DadProfile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
