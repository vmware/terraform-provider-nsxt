---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_host_transport_node_profile"
description: A host transport node profile data source.
---

# nsxt_policy_host_transport_node_profile

This data source provides information about host transport node profiles configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_host_transport_node_profile" "host_transport_node_profile" {
  display_name = "host_transport_node_profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of host transport node profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the host transport node profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
