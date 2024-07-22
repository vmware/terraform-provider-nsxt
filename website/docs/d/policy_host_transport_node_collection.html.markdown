---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_host_transport_node_collection"
description: A host transport node collection data source.
---

# nsxt_policy_host_transport_node_collection

This data source provides information about host transport node collection configured on NSX.
This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_host_transport_node_collection" "host_transport_node_collection" {
  display_name = "host_transport_node_collection1"
  site_path    = data.nsxt_policy_site.paris.path
}
```

## Argument Reference

* `id` - (Optional) The ID of host transport node collection to retrieve.
* `display_name` - (Optional) The Display Name prefix of the host transport node collection to retrieve.
* `site_path` - (Optional) The path of the site which the Transport Node Collection belongs to.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `transport_node_profile_id` - Transport Node Profile Path.
* `unique_id` - A unique identifier assigned by the system.

