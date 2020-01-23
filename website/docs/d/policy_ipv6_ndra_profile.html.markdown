---
layout: "nsxt"
page_title: "NSXT: policy_ipv6_ndra_profile"
sidebar_current: "docs-nsxt-datasource-policy-ipv6-ndra-profile"
description: Policy Ipv6NdraProfile data source.
---

# nsxt_policy_ipv6_ndra_profile

This data source provides information about policy Ipv6NdraProfile configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_ipv6_ndra_profile" "test" {
  display_name = "ipv6-ndra-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Ipv6NdraProfile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Ipv6NdraProfile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
