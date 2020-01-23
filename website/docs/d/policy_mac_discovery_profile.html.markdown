---
layout: "nsxt"
page_title: "NSXT: policy_mac_discovery_profile"
sidebar_current: "docs-nsxt-datasource-policy-mac-discovery-profile"
description: Policy MacDiscoveryProfile data source.
---

# nsxt_policy_mac_discovery_profile

This data source provides information about policy MacDiscoveryProfile configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_mac_discovery_profile" "test" {
  display_name = "mac-discovery-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of MacDiscoveryProfile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the MacDiscoveryProfile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
