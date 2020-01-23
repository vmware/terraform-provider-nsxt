---
layout: "nsxt"
page_title: "NSXT: policy_ip_discovery_profile"
sidebar_current: "docs-nsxt-datasource-policy-ip-discovery-profile"
description: Policy IpDiscoveryProfile data source.
---

# nsxt_policy_ip_discovery_profile

This data source provides information about policy IpDiscoveryProfile configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_ip_discovery_profile" "test" {
  display_name = "ip-discovery-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of IpDiscoveryProfile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the IpDiscoveryProfile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
