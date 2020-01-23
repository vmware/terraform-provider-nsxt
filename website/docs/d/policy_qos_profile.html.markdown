---
layout: "nsxt"
page_title: "NSXT: policy_qos_profile"
sidebar_current: "docs-nsxt-datasource-policy-qos-profile"
description: Policy QosProfile data source.
---

# nsxt_policy_qos_profile

This data source provides information about policy QosProfile configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_qos_profile" "test" {
  display_name = "qos-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of QosProfile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the QosProfile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
